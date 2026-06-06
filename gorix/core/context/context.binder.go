package context

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

func (c *Context) BindBody(target any) error {
	if target == nil {
		return NewValidationError([]FieldError{
			NewFieldError("body", "bind", "body target cannot be nil"),
		})
	}

	defer c.R.Body.Close()

	decoder := json.NewDecoder(c.R.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(target); err != nil {
		return NewValidationError([]FieldError{
			NewFieldError(
				"body",
				"json",
				fmt.Sprintf("invalid JSON body: %v", err),
			),
		})
	}

	return ValidateStruct(target)
}

func (c *Context) BindQuery(target any) error {
	if target == nil {
		return NewValidationError([]FieldError{
			NewFieldError("query", "bind", "query target cannot be nil"),
		})
	}

	if err := bindValues(target, c.R.URL.Query(), "query"); err != nil {
		return err
	}

	return ValidateStruct(target)
}

func (c *Context) BindParams(target any) error {
	if target == nil {
		return NewValidationError([]FieldError{
			NewFieldError("params", "bind", "path parameter target cannot be nil"),
		})
	}

	values := make(url.Values)

	for key, value := range c.params {
		values.Set(key, value)
	}

	if err := bindValues(target, values, "param"); err != nil {
		return err
	}

	return ValidateStruct(target)
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
