package couchdb

import (
	"fmt"
	"reflect"
	"strings"
)

func documentID(document any) (string, bool) {
	return readTaggedString(document, "id")
}

func documentRevision(document any) (string, bool) {
	return readTaggedString(document, "revision")
}

func setDocumentID(document any, value string) {
	setTaggedString(document, "id", value)
}

func setDocumentRevision(document any, value string) {
	setTaggedString(document, "revision", value)
}

func readTaggedString(
	document any,
	target string,
) (string, bool) {
	value := reflect.ValueOf(document)
	if !value.IsValid() {
		return "", false
	}

	for value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return "", false
		}

		value = value.Elem()
	}

	if value.Kind() == reflect.Map &&
		value.Type().Key().Kind() == reflect.String {
		for _, key := range targetKeys(target) {
			mapValue := value.MapIndex(
				reflect.ValueOf(key),
			)
			if !mapValue.IsValid() {
				continue
			}

			text := strings.TrimSpace(
				fmt.Sprint(mapValue.Interface()),
			)
			if text != "" {
				return text, true
			}
		}

		return "", false
	}

	if value.Kind() != reflect.Struct {
		return "", false
	}

	valueType := value.Type()

	for i := 0; i < value.NumField(); i++ {
		structField := valueType.Field(i)
		fieldValue := value.Field(i)

		if !structField.IsExported() ||
			fieldValue.Kind() != reflect.String {
			continue
		}

		if !fieldMatchesTarget(structField, target) {
			continue
		}

		text := strings.TrimSpace(
			fieldValue.String(),
		)

		if text != "" {
			return text, true
		}
	}

	return "", false
}

func setTaggedString(
	document any,
	target string,
	text string,
) {
	value := reflect.ValueOf(document)
	if !value.IsValid() {
		return
	}

	if value.Kind() != reflect.Pointer ||
		value.IsNil() {
		return
	}

	value = value.Elem()

	if value.Kind() == reflect.Map &&
		value.Type().Key().Kind() == reflect.String {
		key := targetKeys(target)[0]

		value.SetMapIndex(
			reflect.ValueOf(key),
			reflect.ValueOf(text),
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
			!fieldValue.CanSet() ||
			fieldValue.Kind() != reflect.String {
			continue
		}

		if !fieldMatchesTarget(structField, target) {
			continue
		}

		fieldValue.SetString(text)
		return
	}
}

func fieldMatchesTarget(
	field reflect.StructField,
	target string,
) bool {
	documentTag := tagName(
		field.Tag.Get("document"),
	)

	if target == "id" && documentTag == "id" {
		return true
	}

	if target == "revision" &&
		(documentTag == "revision" || documentTag == "rev") {
		return true
	}

	names := []string{
		tagName(field.Tag.Get("json")),
		tagName(field.Tag.Get("couchdb")),
		tagName(field.Tag.Get("bson")),
		strings.ToLower(field.Name),
	}

	for _, name := range names {
		for _, key := range targetKeys(target) {
			if name == key {
				return true
			}
		}
	}

	return false
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

func targetKeys(target string) []string {
	switch target {
	case "id":
		return []string{"_id", "id"}

	case "revision":
		return []string{"_rev", "rev"}

	default:
		return []string{target}
	}
}
