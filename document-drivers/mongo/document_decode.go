package mongo

import (
	"context"
	"fmt"
	"reflect"

	"go.mongodb.org/mongo-driver/v2/bson"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
)

func decodeOneDocument(
	result *mongodriver.SingleResult,
	out any,
) error {
	if err := validateDecodeOut(out); err != nil {
		return err
	}

	var document bson.M
	if err := result.Decode(&document); err != nil {
		return err
	}

	normalizeIDForOutput(document, out)

	return decodeBSONMap(document, out)
}

func decodeManyDocuments(
	ctx context.Context,
	cursor *mongodriver.Cursor,
	out any,
) error {
	if err := validateDecodeOut(out); err != nil {
		return err
	}

	outValue := reflect.ValueOf(out)
	sliceValue := outValue.Elem()

	if sliceValue.Kind() != reflect.Slice {
		return fmt.Errorf(
			"gorix mongo: output must be a pointer to slice",
		)
	}

	elementType := sliceValue.Type().Elem()
	wantsPointerElement := elementType.Kind() == reflect.Pointer

	for cursor.Next(ctx) {
		var document bson.M

		if err := cursor.Decode(&document); err != nil {
			return err
		}

		element := reflect.New(elementType)
		target := element.Interface()

		if wantsPointerElement {
			element = reflect.New(elementType.Elem())
			target = element.Interface()
		}

		normalizeIDForOutput(document, target)

		if err := decodeBSONMap(document, target); err != nil {
			return err
		}

		if wantsPointerElement {
			sliceValue.Set(
				reflect.Append(
					sliceValue,
					element,
				),
			)
		} else {
			sliceValue.Set(
				reflect.Append(
					sliceValue,
					element.Elem(),
				),
			)
		}
	}

	if err := cursor.Err(); err != nil {
		return err
	}

	return nil
}

func decodeBSONMap(
	document bson.M,
	out any,
) error {
	raw, err := bson.Marshal(document)
	if err != nil {
		return err
	}

	return bson.Unmarshal(raw, out)
}

func normalizeIDForOutput(
	document bson.M,
	out any,
) {
	if document == nil {
		return
	}

	id, exists := document["_id"]
	if !exists {
		return
	}

	if !outputIDFieldIsString(out) {
		return
	}

	switch typedID := id.(type) {
	case bson.ObjectID:
		document["_id"] = typedID.Hex()
	}
}

func outputIDFieldIsString(out any) bool {
	if out == nil {
		return false
	}

	valueType := reflect.TypeOf(out)

	for valueType.Kind() == reflect.Pointer {
		valueType = valueType.Elem()
	}

	if valueType.Kind() == reflect.Slice {
		valueType = valueType.Elem()

		for valueType.Kind() == reflect.Pointer {
			valueType = valueType.Elem()
		}
	}

	if valueType.Kind() != reflect.Struct {
		return false
	}

	for i := 0; i < valueType.NumField(); i++ {
		field := valueType.Field(i)

		if !field.IsExported() {
			continue
		}

		if !isIDField(field) {
			continue
		}

		fieldType := field.Type

		for fieldType.Kind() == reflect.Pointer {
			fieldType = fieldType.Elem()
		}

		return fieldType.Kind() == reflect.String
	}

	return false
}

func validateDecodeOut(out any) error {
	if out == nil {
		return fmt.Errorf("gorix mongo: output cannot be nil")
	}

	value := reflect.ValueOf(out)
	if value.Kind() != reflect.Pointer || value.IsNil() {
		return fmt.Errorf(
			"gorix mongo: output must be a non-nil pointer",
		)
	}

	return nil
}
