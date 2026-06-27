package couchdb

import (
	"fmt"
	"reflect"

	kivik "github.com/go-kivik/kivik/v4"
)

func scanResultSet(
	rows *kivik.ResultSet,
	out any,
) error {
	if rows == nil {
		return fmt.Errorf(
			"gorix couchdb: result set is nil",
		)
	}

	destination := reflect.ValueOf(out)
	if destination.Kind() != reflect.Pointer ||
		destination.IsNil() {
		return fmt.Errorf(
			"gorix couchdb: output must be a non-nil pointer",
		)
	}

	slice := destination.Elem()
	if slice.Kind() != reflect.Slice {
		return fmt.Errorf(
			"gorix couchdb: output must point to a slice",
		)
	}

	elementType := slice.Type().Elem()

	for rows.Next() {
		var element reflect.Value
		var scanTarget any

		if elementType.Kind() == reflect.Pointer {
			element = reflect.New(elementType.Elem())
			scanTarget = element.Interface()
		} else {
			element = reflect.New(elementType)
			scanTarget = element.Interface()
		}

		if err := rows.ScanDoc(scanTarget); err != nil {
			return err
		}

		if elementType.Kind() == reflect.Pointer {
			slice = reflect.Append(slice, element)
		} else {
			slice = reflect.Append(slice, element.Elem())
		}
	}

	destination.Elem().Set(slice)
	return nil
}
