package mongo

import (
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func setDocumentID(document any, id any) {
	if document == nil || id == nil {
		return
	}

	value := reflect.ValueOf(document)

	if value.Kind() != reflect.Pointer || value.IsNil() {
		return
	}

	value = value.Elem()

	if value.Kind() == reflect.Map &&
		value.Type().Key().Kind() == reflect.String {
		value.SetMapIndex(
			reflect.ValueOf("_id"),
			reflect.ValueOf(id),
		)
		return
	}

	if value.Kind() != reflect.Struct {
		return
	}

	valueType := value.Type()

	for i := 0; i < value.NumField(); i++ {
		structField := valueType.Field(i)
		fieldValue := value.Field(i)

		if !structField.IsExported() ||
			!fieldValue.CanSet() {
			continue
		}

		if !isIDField(structField) {
			continue
		}

		assignID(fieldValue, id)
		return
	}
}

func isIDField(field reflect.StructField) bool {
	documentTag := tagName(
		field.Tag.Get("document"),
	)

	if documentTag == "id" {
		return true
	}

	names := []string{
		tagName(field.Tag.Get("bson")),
		tagName(field.Tag.Get("json")),
		strings.ToLower(field.Name),
	}

	for _, name := range names {
		if name == "_id" || name == "id" {
			return true
		}
	}

	return false
}

func assignID(field reflect.Value, id any) {
	if !field.CanSet() {
		return
	}

	switch field.Kind() {
	case reflect.String:
		switch typedID := id.(type) {
		case bson.ObjectID:
			field.SetString(typedID.Hex())

		case string:
			field.SetString(typedID)
		}

	default:
		idValue := reflect.ValueOf(id)
		if idValue.IsValid() &&
			idValue.Type().AssignableTo(field.Type()) {
			field.Set(idValue)
		}
	}
}

func tagName(tag string) string {
	tag = strings.TrimSpace(tag)

	if tag == "" || tag == "-" {
		return ""
	}

	if index := strings.Index(tag, ","); index >= 0 {
		tag = tag[:index]
	}

	return strings.TrimSpace(tag)
}
