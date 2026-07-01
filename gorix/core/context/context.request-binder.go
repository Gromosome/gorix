package context

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

func (c *Context) BindBody(target any) *Context {
	if c == nil {
		return c
	}

	contentType := c.contentType()

	switch {
	case isJSONContentType(contentType):
		return c.BindJSONBody(target)

	case isXMLContentType(contentType):
		return c.BindXMLBody(target)

	default:
		c.setBindingError(NewValidationError([]FieldError{
			NewFieldError(
				"headers.Content-Type",
				"unsupported",
				"unsupported Content-Type",
			),
		}))
		return c
	}
}

func (c *Context) BindQuery(target any) *Context {
	if target == nil {
		c.setBindingError(NewValidationError([]FieldError{
			NewFieldError("query", "bind", "query target cannot be nil"),
		}))
		return c
	}

	if err := bindValues(target, c.R.URL.Query(), "query"); err != nil {
		c.setBindingError(err)
		return c
	}
	if err := ValidateStruct(target); err != nil {
		c.setBindingError(err)
		return c
	}
	return c
}

func (c *Context) BindParams(target any) *Context {
	if target == nil {
		c.setBindingError(NewValidationError([]FieldError{
			NewFieldError("params", "bind", "path parameter target cannot be nil"),
		}))
	}

	values := make(url.Values)

	for key, value := range c.params {
		values.Set(key, value)
	}

	if err := bindValues(target, values, "param"); err != nil {
		c.setBindingError(err)
		return c
	}

	if err := ValidateStruct(target); err != nil {
		c.setBindingError(err)
		return c
	}
	return c
}
func (c *Context) BindHeaders(target any) *Context {
	if target == nil {
		c.setBindingError(NewValidationError([]FieldError{
			NewFieldError("headers", "bind", "header target cannot be nil"),
		}))
		return c
	}
	if c.R == nil {
		c.setBindingError(NewValidationError([]FieldError{
			NewFieldError("headers", "bind", "request cannot be nil"),
		}))
		return c
	}
	values := make(url.Values)
	for key, rawValues := range c.R.Header {
		values[key] = rawValues
		// Allows lowercase header tags like `header:"authorization"`
		values[strings.ToLower(key)] = rawValues
	}
	if err := bindValues(target, values, "header"); err != nil {
		c.setBindingError(err)
		return c
	}
	if err := ValidateStruct(target); err != nil {
		c.setBindingError(err)
		return c
	}

	return c
}

func bindValues(target any, values url.Values, tagName string) error {
	targetValue := reflect.ValueOf(target)

	if targetValue.Kind() != reflect.Pointer || targetValue.IsNil() {
		return NewValidationError([]FieldError{
			NewFieldError(
				tagName,
				"bind",
				"bind target must be a non-nil pointer",
			),
		})
	}

	elem := targetValue.Elem()

	if elem.Kind() != reflect.Struct {
		return NewValidationError([]FieldError{
			NewFieldError(
				tagName,
				"bind",
				"bind target must point to a struct",
			),
		})
	}

	elemType := elem.Type()
	fieldErrors := make([]FieldError, 0)

	for i := 0; i < elem.NumField(); i++ {
		field := elem.Field(i)
		fieldType := elemType.Field(i)

		if !field.CanSet() {
			continue
		}

		key := fieldType.Tag.Get(tagName)
		if key == "-" {
			continue
		}

		if key == "" {
			key = fieldType.Name
		}

		rawValues, exists := values[key]
		if !exists || len(rawValues) == 0 {
			continue
		}

		if err := setFieldValues(field, rawValues); err != nil {
			fieldErrors = append(
				fieldErrors,
				NewBindFieldError(key, tagName, err),
			)
		}
	}

	if len(fieldErrors) > 0 {
		return NewValidationError(fieldErrors)
	}

	return nil
}

func setFieldValues(field reflect.Value, rawValues []string) error {
	if field.Kind() == reflect.Pointer {
		if len(rawValues) == 0 || strings.TrimSpace(rawValues[0]) == "" {
			return nil
		}

		elem := reflect.New(field.Type().Elem())

		if err := setFieldValues(elem.Elem(), rawValues); err != nil {
			return err
		}

		field.Set(elem)
		return nil
	}

	if field.Kind() == reflect.Slice {
		return setSliceValues(field, rawValues)
	}

	return setFieldValue(field, rawValues[0])
}

func setSliceValues(field reflect.Value, rawValues []string) error {
	parts := make([]string, 0)

	for _, raw := range rawValues {
		for _, part := range strings.Split(raw, ",") {
			part = strings.TrimSpace(part)

			if part != "" {
				parts = append(parts, part)
			}
		}
	}

	slice := reflect.MakeSlice(field.Type(), 0, len(parts))

	for _, part := range parts {
		elem := reflect.New(field.Type().Elem()).Elem()

		if err := setFieldValue(elem, part); err != nil {
			return err
		}

		slice = reflect.Append(slice, elem)
	}

	field.Set(slice)

	return nil
}
func setFieldValue(field reflect.Value, raw string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(raw)
		return nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		value, err := strconv.ParseInt(raw, 10, field.Type().Bits())
		if err != nil {
			return fmt.Errorf("must be a valid integer")
		}

		field.SetInt(value)
		return nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		value, err := strconv.ParseUint(raw, 10, field.Type().Bits())
		if err != nil {
			return fmt.Errorf("must be a valid unsigned integer")
		}

		field.SetUint(value)
		return nil

	case reflect.Float32, reflect.Float64:
		value, err := strconv.ParseFloat(raw, field.Type().Bits())
		if err != nil {
			return fmt.Errorf("must be a valid number")
		}

		field.SetFloat(value)
		return nil

	case reflect.Bool:
		value, err := strconv.ParseBool(raw)
		if err != nil {
			return fmt.Errorf("must be true or false")
		}

		field.SetBool(value)
		return nil

	default:
		return fmt.Errorf("unsupported field type %s", field.Kind())
	}
}
func (c *Context) BindJSONBody(target any) *Context {
	if c == nil {
		return c
	}
	if target == nil {
		c.setBindingError(NewValidationError([]FieldError{
			NewFieldError("body", "bind", "body target cannot be nil"),
		}))
		return c
	}
	contentType := c.contentType()
	if !isJSONContentType(contentType) {
		c.setBindingError(NewValidationError([]FieldError{
			NewFieldError(
				"headers.Content-Type",
				"content_type",
				"Content-Type must be application/json",
			),
		}))
		return c
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(c.R.Body)
	decoder := json.NewDecoder(c.R.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		c.setBindingError(NewValidationError([]FieldError{
			NewFieldError(
				"body",
				"json",
				fmt.Sprintf("invalid JSON body: %v", err),
			),
		}))
		return c
	}
	if err := ValidateStruct(target); err != nil {
		c.setBindingError(err)
		return c
	}
	return c
}
func (c *Context) BindXMLBody(target any) *Context {
	if c == nil {
		return c
	}

	if target == nil {
		c.setBindingError(NewValidationError([]FieldError{
			NewFieldError("body", "bind", "body target cannot be nil"),
		}))
		return c
	}

	contentType := c.contentType()

	if !isXMLContentType(contentType) {
		c.setBindingError(NewValidationError([]FieldError{
			NewFieldError(
				"headers.Content-Type",
				"content_type",
				"Content-Type must be application/xml or text/xml",
			),
		}))
		return c
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(c.R.Body)

	decoder := xml.NewDecoder(c.R.Body)

	if err := decoder.Decode(target); err != nil {
		c.setBindingError(NewValidationError([]FieldError{
			NewFieldError(
				"body",
				"xml",
				fmt.Sprintf("invalid XML body: %v", err),
			),
		}))
		return c
	}

	if err := ValidateStruct(target); err != nil {
		c.setBindingError(err)
		return c
	}

	return c
}
func isJSONContentType(contentType string) bool {
	return contentType == "application/json" ||
		strings.HasSuffix(contentType, "+json")
}

func isXMLContentType(contentType string) bool {
	return contentType == "application/xml" ||
		contentType == "text/xml" ||
		strings.HasSuffix(contentType, "+xml")
}
