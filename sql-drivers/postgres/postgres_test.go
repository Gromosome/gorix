package postgres

import (
	"errors"
	"testing"

	sqldriver "github.com/Gromosome/gorix/sql-driver-manager"
	"github.com/jackc/pgx/v5/pgconn"
)

func TestAdapterRegistration(t *testing.T) {
	registered, err := sqldriver.Lookup(" postgres ")
	if err != nil {
		t.Fatalf("Lookup returned error: %v", err)
	}

	if registered.Name() != "postgres" {
		t.Fatalf("unexpected adapter name: %s", registered.Name())
	}
	if registered.SQLDriverName() != "pgx" {
		t.Fatalf("unexpected SQL driver name: %s", registered.SQLDriverName())
	}
}

func TestNormalizePostgresError(t *testing.T) {
	cause := &pgconn.PgError{
		Code:           "23505",
		Message:        "duplicate key",
		ConstraintName: "users_email_key",
		TableName:      "users",
		ColumnName:     "email",
	}

	err := (adapter{}).Normalize(cause)
	if err.Kind != sqldriver.ErrorUniqueViolation {
		t.Fatalf("unexpected error kind: %s", err.Kind)
	}
	if err.Driver != "postgres" {
		t.Fatalf("unexpected driver: %s", err.Driver)
	}
	if err.Code != "23505" {
		t.Fatalf("unexpected code: %s", err.Code)
	}
	if err.Constraint != "users_email_key" || err.Table != "users" || err.Column != "email" {
		t.Fatalf("metadata was not preserved: %+v", err)
	}
	if !errors.Is(err, cause) {
		t.Fatal("normalized error must wrap original cause")
	}
}

func TestNormalizeGenericError(t *testing.T) {
	cause := errors.New("broken pipe")

	err := (adapter{}).Normalize(cause)
	if err.Kind != sqldriver.ErrorConnection {
		t.Fatalf("unexpected generic error kind: %s", err.Kind)
	}
	if err.Driver != "postgres" {
		t.Fatalf("unexpected driver: %s", err.Driver)
	}
	if !errors.Is(err, cause) {
		t.Fatal("generic error must wrap original cause")
	}
}
