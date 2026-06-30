package document

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"
)

type CollectionNamer interface {
	CollectionName() string
}

type DocumentField struct {
	GoName      string
	Name        string
	Index       []int
	Type        reflect.Type
	IsID        bool
	IsRev       bool
	ReadOnly    bool
	IndexField  bool
	UniqueIndex bool
}

type DocumentIndex struct {
	Name   string
	Fields []string
	Unique bool
}

type DocumentSchema struct {
	Collection string
	Indexes    []DocumentIndex
}

type DocumentMetadata struct {
	Type           reflect.Type
	CollectionName string
	IDField        *DocumentField
	RevField       *DocumentField
	Fields         []DocumentField
	Indexes        []DocumentIndex
}

func ParseMetadata[T any]() (*DocumentMetadata, error) {
	var entity T
	return ParseMetadataOf(entity)
}

func ParseMetadataOf(entity any) (*DocumentMetadata, error) {
	entityType := reflect.TypeOf(entity)
	if entityType == nil {
		return nil, fmt.Errorf("gorix document: entity type cannot be nil")
	}

	for entityType.Kind() == reflect.Pointer {
		entityType = entityType.Elem()
	}

	if entityType.Kind() != reflect.Struct {
		return nil, fmt.Errorf(
			"gorix document: entity must be a struct, got %s",
			entityType.Kind(),
		)
	}

	metadata := &DocumentMetadata{
		Type:           entityType,
		CollectionName: collectionName(entityType),
		Fields:         make([]DocumentField, 0),
		Indexes:        make([]DocumentIndex, 0),
	}

	for i := 0; i < entityType.NumField(); i++ {
		structField := entityType.Field(i)

		if structField.PkgPath != "" {
			continue
		}

		documentTag := strings.TrimSpace(structField.Tag.Get("document"))
		if documentTag == "-" {
			continue
		}

		options := parseDocumentOptions(documentTag)
		name := storageName(structField)

		field := DocumentField{
			GoName:      structField.Name,
			Name:        name,
			Index:       structField.Index,
			Type:        structField.Type,
			IsID:        options["id"] || isIDField(structField, name),
			IsRev:       options["rev"] || isRevField(structField, name),
			ReadOnly:    options["readonly"] || options["readOnly"],
			IndexField:  options["index"],
			UniqueIndex: options["uniqueIndex"] || options["unique"],
		}

		if field.IsID {
			copied := field
			metadata.IDField = &copied
		}

		if field.IsRev {
			copied := field
			metadata.RevField = &copied
		}

		if field.IndexField || field.UniqueIndex {
			metadata.Indexes = append(
				metadata.Indexes,
				DocumentIndex{
					Name:   name,
					Fields: []string{name},
					Unique: field.UniqueIndex,
				},
			)
		}

		metadata.Fields = append(metadata.Fields, field)
	}

	if metadata.IDField == nil {
		return nil, fmt.Errorf(
			"gorix document: %s must have a document id field; add document:\"id\"",
			entityType.Name(),
		)
	}

	if strings.TrimSpace(metadata.CollectionName) == "" {
		return nil, fmt.Errorf(
			"gorix document: collection name cannot be empty for %s",
			entityType.Name(),
		)
	}

	return metadata, nil
}

func (m *DocumentMetadata) Schema() DocumentSchema {
	if m == nil {
		return DocumentSchema{}
	}

	return DocumentSchema{
		Collection: m.CollectionName,
		Indexes:    m.Indexes,
	}
}

func (m *DocumentMetadata) IDValue(entity any) (any, error) {
	if m == nil || m.IDField == nil {
		return nil, fmt.Errorf("gorix document: id metadata is missing")
	}

	value, err := structValue(entity)
	if err != nil {
		return nil, err
	}

	field := value.FieldByIndex(m.IDField.Index)
	if !field.IsValid() {
		return nil, fmt.Errorf("gorix document: id field is invalid")
	}

	return field.Interface(), nil
}

func (m *DocumentMetadata) IsIDZero(entity any) bool {
	if m == nil || m.IDField == nil {
		return true
	}

	value, err := structValue(entity)
	if err != nil {
		return true
	}

	field := value.FieldByIndex(m.IDField.Index)
	if !field.IsValid() {
		return true
	}

	return field.IsZero()
}

func (m *DocumentMetadata) SetID(entity any, id any) error {
	if m == nil || m.IDField == nil || id == nil {
		return nil
	}

	return m.setField(entity, m.IDField, id)
}

func (m *DocumentMetadata) RevValue(entity any) string {
	if m == nil || m.RevField == nil {
		return ""
	}

	value, err := structValue(entity)
	if err != nil {
		return ""
	}

	field := value.FieldByIndex(m.RevField.Index)
	if !field.IsValid() || field.IsZero() {
		return ""
	}

	return fmt.Sprint(field.Interface())
}

func (m *DocumentMetadata) SetRev(entity any, rev string) error {
	if m == nil || m.RevField == nil || strings.TrimSpace(rev) == "" {
		return nil
	}

	return m.setField(entity, m.RevField, rev)
}

func (m *DocumentMetadata) setField(
	entity any,
	field *DocumentField,
	value any,
) error {
	entityValue, err := structValue(entity)
	if err != nil {
		return err
	}

	target := entityValue.FieldByIndex(field.Index)
	if !target.IsValid() {
		return fmt.Errorf(
			"gorix document: field %q is invalid",
			field.GoName,
		)
	}

	if !target.CanSet() {
		return fmt.Errorf(
			"gorix document: field %q cannot be set",
			field.GoName,
		)
	}

	source := reflect.ValueOf(value)
	if !source.IsValid() {
		return nil
	}

	if source.Type().AssignableTo(target.Type()) {
		target.Set(source)
		return nil
	}

	if source.Type().ConvertibleTo(target.Type()) {
		target.Set(source.Convert(target.Type()))
		return nil
	}

	if target.Kind() == reflect.String {
		target.SetString(fmt.Sprint(value))
		return nil
	}

	return fmt.Errorf(
		"gorix document: cannot assign %T to field %q of type %s",
		value,
		field.GoName,
		target.Type(),
	)
}

func structValue(entity any) (reflect.Value, error) {
	if entity == nil {
		return reflect.Value{}, fmt.Errorf("gorix document: entity cannot be nil")
	}

	value := reflect.ValueOf(entity)
	if value.Kind() != reflect.Pointer || value.IsNil() {
		return reflect.Value{}, fmt.Errorf(
			"gorix document: entity must be a non-nil pointer",
		)
	}

	value = value.Elem()

	for value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return reflect.Value{}, fmt.Errorf(
				"gorix document: nested entity pointer cannot be nil",
			)
		}
		value = value.Elem()
	}

	if value.Kind() != reflect.Struct {
		return reflect.Value{}, fmt.Errorf(
			"gorix document: entity must point to a struct",
		)
	}

	return value, nil
}

func collectionName(entityType reflect.Type) string {
	ptr := reflect.New(entityType)
	if namer, ok := ptr.Interface().(CollectionNamer); ok {
		return strings.TrimSpace(namer.CollectionName())
	}

	value := reflect.Zero(entityType)
	if value.CanInterface() {
		if namer, ok := value.Interface().(CollectionNamer); ok {
			return strings.TrimSpace(namer.CollectionName())
		}
	}

	return toSnakeCase(entityType.Name())
}

func parseDocumentOptions(tag string) map[string]bool {
	options := make(map[string]bool)

	for _, part := range strings.Split(tag, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		options[part] = true
	}

	return options
}

func storageName(field reflect.StructField) string {
	for _, tagName := range []string{"bson", "couchdb", "json"} {
		tagValue := strings.TrimSpace(field.Tag.Get(tagName))
		if tagValue == "" {
			continue
		}

		name := strings.TrimSpace(strings.Split(tagValue, ",")[0])
		if name == "" || name == "-" {
			continue
		}

		return name
	}

	return toSnakeCase(field.Name)
}

func isIDField(field reflect.StructField, name string) bool {
	return field.Name == "ID" ||
		field.Name == "Id" ||
		name == "_id" ||
		strings.EqualFold(name, "id")
}

func isRevField(field reflect.StructField, name string) bool {
	return field.Name == "Rev" ||
		field.Name == "Revision" ||
		name == "_rev" ||
		strings.EqualFold(name, "rev")
}

func toSnakeCase(value string) string {
	if value == "" {
		return ""
	}

	var builder strings.Builder

	for index, char := range value {
		if unicode.IsUpper(char) {
			if index > 0 {
				builder.WriteRune('_')
			}
			builder.WriteRune(unicode.ToLower(char))
			continue
		}

		builder.WriteRune(char)
	}

	return builder.String()
}
