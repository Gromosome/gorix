package sqlite_modern

import (
	"database/sql"
	"errors"
	"testing"

	sqldriver "github.com/Gromosome/gorix/sql-driver-manager"
)

func TestAdapterRegistration(t *testing.T) {
	registered, err := sqldriver.Lookup(" sqlite-modern ")
	if err != nil {
		t.Fatalf("Lookup returned error: %v", err)
	}

	if registered.Name() != "sqlite-modern" {
		t.Fatalf("unexpected adapter name: %s", registered.Name())
	}
	if registered.SQLDriverName() != "sqlite" {
		t.Fatalf("unexpected SQL driver name: %s", registered.SQLDriverName())
	}
}

func TestNormalizeSQLiteModernRuntimeError(t *testing.T) {
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
	if normalized.Driver != "sqlite-modern" {
		t.Fatalf("unexpected driver: %s", normalized.Driver)
	}
	if !errors.Is(normalized, err) {
		t.Fatal("normalized error must wrap original cause")
	}
}

func TestNormalizeGenericError(t *testing.T) {
	cause := errors.New("bad connection")

	err := (adapter{}).Normalize(cause)
	if err.Kind != sqldriver.ErrorConnection {
		t.Fatalf("unexpected generic error kind: %s", err.Kind)
	}
	if err.Driver != "sqlite-modern" {
		t.Fatalf("unexpected driver: %s", err.Driver)
	}
	if !errors.Is(err, cause) {
		t.Fatal("generic error must wrap original cause")
	}
}
