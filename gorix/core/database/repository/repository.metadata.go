package repository

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

type TableNamer interface {
	TableName() string
}

type FieldMetadata struct {
	Index         int
	FieldName     string
	ColumnName    string
	Type          reflect.Type
	PrimaryKey    bool
	AutoIncrement bool
	ReadOnly      bool
	OmitEmpty     bool
}

type EntityMetadata struct {
	Type       reflect.Type
	TableName  string
	Fields     []FieldMetadata
	PrimaryKey *FieldMetadata
}

var metadataCache sync.Map

func MetadataOf[T any]() (*EntityMetadata, error) {
	var entity T

	entityType := reflect.TypeOf(entity)

	if entityType == nil {
		return nil, fmt.Errorf(
			"gorix repository: entity type cannot be nil",
		)
	}

	if entityType.Kind() == reflect.Pointer {
		entityType = entityType.Elem()
	}

	if entityType.Kind() != reflect.Struct {
		return nil, fmt.Errorf(
			"gorix repository: entity must be a struct, got %s",
			entityType.Kind(),
		)
	}

	if cached, exists := metadataCache.Load(entityType); exists {
		return cached.(*EntityMetadata), nil
	}

	metadata, err := buildEntityMetadata(entityType)
	if err != nil {
		return nil, err
	}

	metadataCache.Store(entityType, metadata)
	return metadata, nil
}

func buildEntityMetadata(
	entityType reflect.Type,
) (*EntityMetadata, error) {
	tableName := defaultTableName(entityType.Name())

	pointerValue := reflect.New(entityType)
	if namer, ok := pointerValue.Interface().(TableNamer); ok {
		tableName = namer.TableName()
	} else if namer, ok := pointerValue.Elem().Interface().(TableNamer); ok {
		tableName = namer.TableName()
	}

	metadata := &EntityMetadata{
		Type:      entityType,
		TableName: tableName,
		Fields:    make([]FieldMetadata, 0),
	}

	for i := 0; i < entityType.NumField(); i++ {
		field := entityType.Field(i)

		if !field.IsExported() {
			continue
		}

		columnName := field.Tag.Get("db")
		if columnName == "-" {
			continue
		}

		if columnName == "" {
			columnName = toSnakeCase(field.Name)
		}

		options := parseORMOptions(field.Tag.Get("repository"))

		fieldMetadata := FieldMetadata{
			Index:         i,
			FieldName:     field.Name,
			ColumnName:    columnName,
			Type:          field.Type,
			PrimaryKey:    options["primaryKey"],
			AutoIncrement: options["autoIncrement"],
			ReadOnly:      options["readOnly"],
			OmitEmpty:     options["omitEmpty"],
		}

		metadata.Fields = append(
			metadata.Fields,
			fieldMetadata,
		)

		if fieldMetadata.PrimaryKey {
			if metadata.PrimaryKey != nil {
				return nil, fmt.Errorf(
					"gorix repository: entity %s has multiple primary keys; composite keys are not supported yet",
					entityType.Name(),
				)
			}

			copy := fieldMetadata
			metadata.PrimaryKey = &copy
		}
	}

	if metadata.PrimaryKey == nil {
		return nil, fmt.Errorf(
			"gorix repository: entity %s must declare one primary key using repository:\"primaryKey\"",
			entityType.Name(),
		)
	}

	return metadata, nil
}

func parseORMOptions(tag string) map[string]bool {
	options := make(map[string]bool)

	for _, item := range strings.Split(tag, ",") {
		item = strings.TrimSpace(item)

		if item != "" {
			options[item] = true
		}
	}

	return options
}

func defaultTableName(typeName string) string {
	name := toSnakeCase(typeName)

	if strings.HasSuffix(name, "s") {
		return name
	}

	return name + "s"
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
