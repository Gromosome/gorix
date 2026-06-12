package sqlite3

import (
	"database/sql"
	"errors"
	"testing"

	sqldriver "github.com/Gromosome/gorix/sql-driver-manager"
	sqlite3 "github.com/mattn/go-sqlite3"
)

func TestAdapterRegistration(t *testing.T) {
	registered, err := sqldriver.Lookup(" sqlite3 ")
	if err != nil {
		t.Fatalf("Lookup returned error: %v", err)
	}

	if registered.Name() != "sqlite3" {
		t.Fatalf("unexpected adapter name: %s", registered.Name())
	}
	if registered.SQLDriverName() != "sqlite3" {
		t.Fatalf("unexpected SQL driver name: %s", registered.SQLDriverName())
	}
}

func TestNormalizeSQLite3Error(t *testing.T) {
	cause := sqlite3.Error{
		Code:         sqlite3.ErrConstraint,
		ExtendedCode: sqlite3.ErrConstraintUnique,
	}

	err := (adapter{}).Normalize(cause)
	if err.Kind != sqldriver.ErrorUniqueViolation {
		t.Fatalf("unexpected error kind: %s", err.Kind)
	}
	if err.Driver != "sqlite3" {
		t.Fatalf("unexpected driver: %s", err.Driver)
	}
	if err.Code != "2067" {
		t.Fatalf("unexpected code: %s", err.Code)
	}
	if !errors.Is(err, cause) {
		t.Fatal("normalized error must wrap original cause")
	}
}

func TestNormalizeSQLite3RuntimeError(t *testing.T) {
	db, err := sql.Open((adapter{}).SQLDriverName(), ":memory:")
	if err != nil {
		t.Fatalf("sql.Open returned error: %v", err)
	}
	defer db.Close()

	if _, err = db.Exec("create table users (email text unique)"); err != nil {
		t.Fatalf("create table returned error: %v", err)
	}
	if _, err = db.Exec("insert into users (email) values (?)", "a@example.test"); err != nil {
		t.Fatalf("first insert returned error: %v", err)
	}
	_, err = db.Exec("insert into users (email) values (?)", "a@example.test")
	if err == nil {
		t.Fatal("expected duplicate insert error")
	}

	normalized := (adapter{}).Normalize(err)
	if normalized.Kind != sqldriver.ErrorUniqueViolation {
		t.Fatalf("unexpected runtime error kind: %s", normalized.Kind)
	}
}
