package core

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
		return fmt.Errorf("gorix: BindBody target cannot be nil")
	}

	defer c.R.Body.Close()

	if err := json.NewDecoder(c.R.Body).Decode(target); err != nil {
		return err
	}

	return ValidateStruct(target)
}

func (c *Context) BindQuery(target any) error {
	if target == nil {
		return fmt.Errorf("gorix: BindQuery target cannot be nil")
	}

	values := c.R.URL.Query()

	if err := bindValues(target, values, "query"); err != nil {
		return err
	}

	return ValidateStruct(target)
}

func (c *Context) BindParams(target any) error {
	if target == nil {
		return fmt.Errorf("gorix: BindParams target cannot be nil")
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
		return fmt.Errorf("gorix: bind target must be non-nil pointer")
	}

	elem := targetValue.Elem()

	if elem.Kind() != reflect.Struct {
		return fmt.Errorf("gorix: bind target must point to struct")
	}

	elemType := elem.Type()

	for i := 0; i < elem.NumField(); i++ {
		field := elem.Field(i)
		fieldType := elemType.Field(i)

		if !field.CanSet() {
			continue
		}

		key := fieldType.Tag.Get(tagName)
		if key == "" || key == "-" {
			key = fieldType.Name
		}

		raw := values.Get(key)
		if raw == "" {
			continue
		}

		if err := setFieldValue(field, raw); err != nil {
			return fmt.Errorf("gorix: failed to bind %s field %s: %w", tagName, fieldType.Name, err)
		}
	}

	return nil
}

func setFieldValue(field reflect.Value, raw string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(raw)
		return nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		value, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			return err
		}
		field.SetInt(value)
		return nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		value, err := strconv.ParseUint(raw, 10, 64)
		if err != nil {
			return err
		}
		field.SetUint(value)
		return nil

	case reflect.Float32, reflect.Float64:
		value, err := strconv.ParseFloat(raw, 64)
		if err != nil {
			return err
		}
		field.SetFloat(value)
		return nil

	case reflect.Bool:
		value, err := strconv.ParseBool(raw)
		if err != nil {
			return err
		}
		field.SetBool(value)
		return nil

	case reflect.Slice:
		return setSliceValue(field, raw)

	default:
		return fmt.Errorf("unsupported field type %s", field.Kind())
	}
}

func setSliceValue(field reflect.Value, raw string) error {
	parts := strings.Split(raw, ",")
	slice := reflect.MakeSlice(field.Type(), 0, len(parts))

	for _, part := range parts {
		part = strings.TrimSpace(part)
		elem := reflect.New(field.Type().Elem()).Elem()

		if err := setFieldValue(elem, part); err != nil {
			return err
		}

		slice = reflect.Append(slice, elem)
	}

	field.Set(slice)
	return nil
}
