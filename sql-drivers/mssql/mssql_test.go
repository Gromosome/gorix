package mssql

import (
	"errors"
	"testing"

	sqldriver "github.com/Gromosome/gorix/sql-driver-manager"
	mssql "github.com/microsoft/go-mssqldb"
)

func TestAdapterRegistration(t *testing.T) {
	registered, err := sqldriver.Lookup(" mssql ")
	if err != nil {
		t.Fatalf("Lookup returned error: %v", err)
	}

	if registered.Name() != "mssql" {
		t.Fatalf("unexpected adapter name: %s", registered.Name())
	}
	if registered.SQLDriverName() != "sqlserver" {
		t.Fatalf("unexpected SQL driver name: %s", registered.SQLDriverName())
	}
}

func TestNormalizeMSSQLError(t *testing.T) {
	cause := mssql.Error{
		Number:  2627,
		Message: "duplicate key",
	}

	err := (adapter{}).Normalize(cause)
	if err.Kind != sqldriver.ErrorDuplicateKey {
		t.Fatalf("unexpected error kind: %s", err.Kind)
	}
	if err.Driver != "mssql" {
		t.Fatalf("unexpected driver: %s", err.Driver)
	}
	if err.Code != "2627" {
		t.Fatalf("unexpected code: %s", err.Code)
	}
	var target mssql.Error
	if !errors.As(err, &target) {
		t.Fatal("normalized error must expose original cause")
	}
	if target.Number != cause.Number {
		t.Fatalf("unexpected wrapped error number: %d", target.Number)
	}
}

func TestNormalizeGenericError(t *testing.T) {
	cause := errors.New("connection reset by peer")

	err := (adapter{}).Normalize(cause)
	if err.Kind != sqldriver.ErrorConnection {
		t.Fatalf("unexpected generic error kind: %s", err.Kind)
	}
	if err.Driver != "mssql" {
		t.Fatalf("unexpected driver: %s", err.Driver)
	}
	if !errors.Is(err, cause) {
		t.Fatal("generic error must wrap original cause")
	}
}
