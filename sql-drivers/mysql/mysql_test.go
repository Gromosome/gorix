package mysql

import (
	"errors"
	"testing"

	sqldriver "github.com/Gromosome/gorix/sql-driver-manager"
	mysql "github.com/go-sql-driver/mysql"
)

func TestAdapterRegistration(t *testing.T) {
	registered, err := sqldriver.Lookup(" mysql ")
	if err != nil {
		t.Fatalf("Lookup returned error: %v", err)
	}

	if registered.Name() != "mysql" {
		t.Fatalf("unexpected adapter name: %s", registered.Name())
	}
	if registered.SQLDriverName() != "mysql" {
		t.Fatalf("unexpected SQL driver name: %s", registered.SQLDriverName())
	}
}

func TestNormalizeMySQLError(t *testing.T) {
	cause := &mysql.MySQLError{
		Number:  1062,
		Message: "duplicate entry",
	}

	err := (adapter{}).Normalize(cause)
	if err.Kind != sqldriver.ErrorDuplicateKey {
		t.Fatalf("unexpected error kind: %s", err.Kind)
	}
	if err.Driver != "mysql" {
		t.Fatalf("unexpected driver: %s", err.Driver)
	}
	if err.Code != "1062" {
		t.Fatalf("unexpected code: %s", err.Code)
	}
	if !errors.Is(err, cause) {
		t.Fatal("normalized error must wrap original cause")
	}
}

func TestNormalizeGenericError(t *testing.T) {
	cause := errors.New("connection refused")

	err := (adapter{}).Normalize(cause)
	if err.Kind != sqldriver.ErrorConnection {
		t.Fatalf("unexpected generic error kind: %s", err.Kind)
	}
	if err.Driver != "mysql" {
		t.Fatalf("unexpected driver: %s", err.Driver)
	}
	if !errors.Is(err, cause) {
		t.Fatal("generic error must wrap original cause")
	}
}
