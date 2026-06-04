package core

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

func ValidateStruct(target any) error {
	if target == nil {
		return nil
	}

	value := reflect.ValueOf(target)

	if value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return nil
		}

		value = value.Elem()
	}

	if value.Kind() != reflect.Struct {
		return nil
	}

	valueType := value.Type()
	errors := make([]FieldError, 0)

	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		fieldInfo := valueType.Field(i)

		if !fieldInfo.IsExported() {
			continue
		}

		validateTag := fieldInfo.Tag.Get("validate")
		if validateTag == "" || validateTag == "-" {
			continue
		}

		fieldName := getDTOFieldName(fieldInfo)

		fieldErrors := validateField(fieldName, field, validateTag)
		errors = append(errors, fieldErrors...)
	}

	if len(errors) > 0 {
		return NewValidationError(errors)
	}

	return nil
}

func getDTOFieldName(field reflect.StructField) string {
	for _, tag := range []string{"json", "query", "param"} {
		value := field.Tag.Get(tag)
		if value == "" || value == "-" {
			continue
		}

		name := strings.Split(value, ",")[0]
		if name != "" {
			return name
		}
	}

	return field.Name
}

func validateField(fieldName string, field reflect.Value, validateTag string) []FieldError {
	rules := splitValidationRules(validateTag)
	errors := make([]FieldError, 0)

	for _, rule := range rules {
		rule = strings.TrimSpace(rule)

		if rule == "" {
			continue
		}

		switch {
		case rule == "required":
			if isZeroValue(field) {
				errors = append(errors, RequiredError(fieldName))
			}

		case rule == "email":
			if !isZeroValue(field) && !isEmail(field.String()) {
				errors = append(errors, NewFieldError(
					fieldName,
					"email",
					fmt.Sprintf("%s must be a valid email", fieldName),
				))
			}

		case strings.HasPrefix(rule, "min="):
			errors = append(errors, validateMin(fieldName, field, rule)...)

		case strings.HasPrefix(rule, "max="):
			errors = append(errors, validateMax(fieldName, field, rule)...)

		case strings.HasPrefix(rule, "oneof="):
			errors = append(errors, validateOneOf(fieldName, field, rule)...)

		case strings.HasPrefix(rule, "regex="):
			errors = append(errors, validateRegex(fieldName, field, rule)...)
		}
	}

	return errors
}
func splitValidationRules(tag string) []string {
	rules := make([]string, 0)

	start := 0
	inRegex := false
	braceDepth := 0

	for i, ch := range tag {
		if strings.HasPrefix(tag[i:], "regex=") {
			inRegex = true
		}

		switch ch {
		case '{', '[', '(':
			if inRegex {
				braceDepth++
			}

		case '}', ']', ')':
			if inRegex && braceDepth > 0 {
				braceDepth--
			}

		case ',':
			if !inRegex || braceDepth == 0 {
				rules = append(rules, strings.TrimSpace(tag[start:i]))
				start = i + 1
				inRegex = false
			}
		}
	}

	rules = append(rules, strings.TrimSpace(tag[start:]))

	return rules
}
func validateRegex(fieldName string, field reflect.Value, rule string) []FieldError {
	pattern := strings.TrimPrefix(rule, "regex=")

	if pattern == "" {
		return nil
	}

	value := fmt.Sprintf("%v", field.Interface())

	if isZeroValue(field) {
		return nil
	}

	matched, err := regexp.MatchString(pattern, value)
	if err != nil {
		return []FieldError{
			NewFieldError(
				fieldName,
				"regex",
				fmt.Sprintf("%s has invalid regex pattern", fieldName),
			),
		}
	}

	if matched {
		return nil
	}

	return []FieldError{
		NewFieldError(
			fieldName,
			"regex",
			fmt.Sprintf("%s format is invalid", fieldName),
		),
	}
}

func isZeroValue(field reflect.Value) bool {
	return field.IsZero()
}

func isEmail(value string) bool {
	pattern := `^[^\s@]+@[^\s@]+\.[^\s@]+$`
	return regexp.MustCompile(pattern).MatchString(value)
}

func validateMin(fieldName string, field reflect.Value, rule string) []FieldError {
	minRaw := strings.TrimPrefix(rule, "min=")
	min, err := strconv.ParseFloat(minRaw, 64)
	if err != nil {
		return nil
	}

	valid := true

	switch field.Kind() {
	case reflect.String:
		valid = float64(len(field.String())) >= min

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		valid = float64(field.Int()) >= min

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		valid = float64(field.Uint()) >= min

	case reflect.Float32, reflect.Float64:
		valid = field.Float() >= min

	case reflect.Slice, reflect.Array:
		valid = float64(field.Len()) >= min
	}

	if valid {
		return nil
	}

	return []FieldError{
		NewFieldError(
			fieldName,
			rule,
			fmt.Sprintf("%s must be at least %s", fieldName, minRaw),
		),
	}
}

func validateMax(fieldName string, field reflect.Value, rule string) []FieldError {
	maxRaw := strings.TrimPrefix(rule, "max=")
	max, err := strconv.ParseFloat(maxRaw, 64)
	if err != nil {
		return nil
	}

	valid := true

	switch field.Kind() {
	case reflect.String:
		valid = float64(len(field.String())) <= max

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		valid = float64(field.Int()) <= max

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		valid = float64(field.Uint()) <= max

	case reflect.Float32, reflect.Float64:
		valid = field.Float() <= max

	case reflect.Slice, reflect.Array:
		valid = float64(field.Len()) <= max
	}

	if valid {
		return nil
	}

	return []FieldError{
		NewFieldError(
			fieldName,
			rule,
			fmt.Sprintf("%s must be at most %s", fieldName, maxRaw),
		),
	}
}

func validateOneOf(fieldName string, field reflect.Value, rule string) []FieldError {
	raw := strings.TrimPrefix(rule, "oneof=")
	options := strings.Fields(raw)

	value := fmt.Sprintf("%v", field.Interface())

	for _, option := range options {
		if value == option {
			return nil
		}
	}

	return []FieldError{
		NewFieldError(
			fieldName,
			rule,
			fmt.Sprintf("%s must be one of: %s", fieldName, strings.Join(options, ", ")),
		),
	}
}
