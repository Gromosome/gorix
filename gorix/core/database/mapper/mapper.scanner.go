package mapper

import (
	"fmt"
	"reflect"
	"strings"
)

type RowsScanner interface {
	Columns() ([]string, error)
	Scan(dest ...any) error
}

type structMetadata struct {
	fields map[string][]int
}

func buildStructMetadata(structType reflect.Type) structMetadata {
	fields := make(map[string][]int)

	collectStructFields(structType, nil, fields)

	return structMetadata{
		fields: fields,
	}
}

func collectStructFields(
	structType reflect.Type,
	parentIndex []int,
	fields map[string][]int,
) {
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)

		if !field.IsExported() {
			continue
		}

		index := append(
			append([]int(nil), parentIndex...),
			i,
		)

		fieldType := field.Type

		if field.Anonymous {
			if fieldType.Kind() == reflect.Pointer {
				fieldType = fieldType.Elem()
			}

			if fieldType.Kind() == reflect.Struct {
				collectStructFields(fieldType, index, fields)
				continue
			}
		}

		columnName := field.Tag.Get("db")
		if columnName == "-" {
			continue
		}

		if columnName == "" {
			columnName = toSnakeCase(field.Name)
		}

		fields[strings.ToLower(columnName)] = index
	}
}

func ScanStruct[T any](row RowsScanner) (T, error) {
	var result T

	value := reflect.ValueOf(&result).Elem()
	if value.Kind() != reflect.Struct {
		return result, fmt.Errorf(
			"gorix mapper: target type must be a struct, got %s",
			value.Kind(),
		)
	}

	columns, err := row.Columns()
	if err != nil {
		return result, fmt.Errorf(
			"gorix mapper: failed to read columns: %w",
			err,
		)
	}

	metadata := buildStructMetadata(value.Type())
	destinations := make([]any, len(columns))

	for i, column := range columns {
		index, found := metadata.fields[strings.ToLower(column)]

		if !found {
			var ignored any
			destinations[i] = &ignored
			continue
		}

		field := value.FieldByIndex(index)

		if !field.CanAddr() {
			var ignored any
			destinations[i] = &ignored
			continue
		}

		destinations[i] = field.Addr().Interface()
	}

	if err := row.Scan(destinations...); err != nil {
		return result, fmt.Errorf(
			"gorix mapper: failed to scan row: %w",
			err,
		)
	}

	return result, nil
}

func toSnakeCase(value string) string {
	if value == "" {
		return ""
	}

	var builder strings.Builder

	for i, character := range value {
		if character >= 'A' && character <= 'Z' {
			if i > 0 {
				builder.WriteByte('_')
			}

			builder.WriteRune(character + ('a' - 'A'))
			continue
		}

		builder.WriteRune(character)
	}

	return builder.String()
}
