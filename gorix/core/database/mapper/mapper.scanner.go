package mapper

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

type RowScanner interface {
	Columns() ([]string, error)
	Scan(dest ...any) error
}

type structMetadata struct {
	fields map[string][]int
}

var structMetadataCache sync.Map

func ScanInto(
	row RowScanner,
	target any,
) error {
	if row == nil {
		return fmt.Errorf(
			"gorix mapper: row scanner cannot be nil",
		)
	}

	if target == nil {
		return fmt.Errorf(
			"gorix mapper: scan target cannot be nil",
		)
	}

	targetValue := reflect.ValueOf(target)

	if targetValue.Kind() != reflect.Pointer ||
		targetValue.IsNil() {
		return fmt.Errorf(
			"gorix mapper: scan target must be a non-nil pointer",
		)
	}

	structValue := targetValue.Elem()

	if structValue.Kind() != reflect.Struct {
		return fmt.Errorf(
			"gorix mapper: scan target must point to a struct, got %s",
			structValue.Kind(),
		)
	}

	columns, err := row.Columns()
	if err != nil {
		return fmt.Errorf(
			"gorix mapper: failed to read result columns: %w",
			err,
		)
	}

	metadata := getStructMetadata(
		structValue.Type(),
	)

	destinations := make([]any, len(columns))

	for index, column := range columns {
		columnKey := normalizeColumnName(column)

		fieldIndex, exists := metadata.fields[columnKey]
		if !exists {
			var ignored any
			destinations[index] = &ignored
			continue
		}

		field, err := resolveWritableField(
			structValue,
			fieldIndex,
		)
		if err != nil {
			return fmt.Errorf(
				"gorix mapper: failed to resolve column %q: %w",
				column,
				err,
			)
		}

		if !field.CanAddr() {
			var ignored any
			destinations[index] = &ignored
			continue
		}

		destinations[index] = field.Addr().Interface()
	}

	if err := row.Scan(destinations...); err != nil {
		return fmt.Errorf(
			"gorix mapper: failed to scan row: %w",
			err,
		)
	}

	return nil
}

func getStructMetadata(
	structType reflect.Type,
) structMetadata {
	if cached, exists := structMetadataCache.Load(
		structType,
	); exists {
		return cached.(structMetadata)
	}

	metadata := buildStructMetadata(structType)

	structMetadataCache.Store(
		structType,
		metadata,
	)

	return metadata
}

func buildStructMetadata(
	structType reflect.Type,
) structMetadata {
	fields := make(map[string][]int)

	collectStructFields(
		structType,
		nil,
		fields,
	)

	return structMetadata{
		fields: fields,
	}
}

func collectStructFields(
	structType reflect.Type,
	parentIndex []int,
	fields map[string][]int,
) {
	for index := 0; index < structType.NumField(); index++ {
		field := structType.Field(index)

		if !field.IsExported() {
			continue
		}

		fieldIndex := append(
			append([]int(nil), parentIndex...),
			index,
		)

		fieldType := field.Type

		if field.Anonymous {
			if fieldType.Kind() == reflect.Pointer {
				fieldType = fieldType.Elem()
			}

			if fieldType.Kind() == reflect.Struct {
				collectStructFields(
					fieldType,
					fieldIndex,
					fields,
				)
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

		fields[normalizeColumnName(columnName)] = fieldIndex
	}
}

func resolveWritableField(
	root reflect.Value,
	indexPath []int,
) (reflect.Value, error) {
	current := root

	for position, fieldIndex := range indexPath {
		if current.Kind() == reflect.Pointer {
			if current.IsNil() {
				if !current.CanSet() {
					return reflect.Value{}, fmt.Errorf(
						"pointer field cannot be initialized",
					)
				}

				current.Set(
					reflect.New(current.Type().Elem()),
				)
			}

			current = current.Elem()
		}

		if current.Kind() != reflect.Struct {
			return reflect.Value{}, fmt.Errorf(
				"expected struct while resolving field path, got %s",
				current.Kind(),
			)
		}

		field := current.Field(fieldIndex)

		if position == len(indexPath)-1 {
			if !field.CanSet() {
				return reflect.Value{}, fmt.Errorf(
					"field cannot be set",
				)
			}

			return field, nil
		}

		if field.Kind() == reflect.Pointer {
			if field.IsNil() {
				if !field.CanSet() {
					return reflect.Value{}, fmt.Errorf(
						"nested pointer field cannot be initialized",
					)
				}

				field.Set(
					reflect.New(field.Type().Elem()),
				)
			}

			field = field.Elem()
		}

		current = field
	}

	return reflect.Value{}, fmt.Errorf(
		"invalid field index path",
	)
}

func normalizeColumnName(
	value string,
) string {
	return strings.ToLower(
		strings.TrimSpace(value),
	)
}

func toSnakeCase(value string) string {
	if value == "" {
		return ""
	}

	var builder strings.Builder
	runes := []rune(value)

	for index, current := range runes {
		if index > 0 &&
			isUpper(current) &&
			(isLower(runes[index-1]) ||
				isDigit(runes[index-1]) ||
				index+1 < len(runes) && isLower(runes[index+1])) {
			builder.WriteRune('_')
		}

		builder.WriteRune(toLower(current))
	}

	return builder.String()
}

func isUpper(value rune) bool {
	return value >= 'A' && value <= 'Z'
}

func isLower(value rune) bool {
	return value >= 'a' && value <= 'z'
}

func isDigit(value rune) bool {
	return value >= '0' && value <= '9'
}

func toLower(value rune) rune {
	if isUpper(value) {
		return value + ('a' - 'A')
	}

	return value
}
