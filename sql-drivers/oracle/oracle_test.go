package oracle

import (
	"errors"
	"testing"

	sqldriver "github.com/Gromosome/gorix/sql-driver-manager"
)

type oracleTestError struct {
	code    int
	message string
}

func (e oracleTestError) Error() string {
	return e.message
}

func (e oracleTestError) Code() int {
	return e.code
}

func TestAdapterRegistration(t *testing.T) {
	registered, err := sqldriver.Lookup(" oracle ")
	if err != nil {
		t.Fatalf("Lookup returned error: %v", err)
	}

	if registered.Name() != "oracle" {
		t.Fatalf("unexpected adapter name: %s", registered.Name())
	}
	if registered.SQLDriverName() != "godror" {
		t.Fatalf("unexpected SQL driver name: %s", registered.SQLDriverName())
	}
}

func TestNormalizeOracleError(t *testing.T) {
	cause := oracleTestError{
		code:    1,
		message: "unique constraint violated",
	}

	err := (adapter{}).Normalize(cause)
	if err.Kind != sqldriver.ErrorUniqueViolation {
		t.Fatalf("unexpected error kind: %s", err.Kind)
	}
	if err.Driver != "oracle" {
		t.Fatalf("unexpected driver: %s", err.Driver)
	}
	if err.Code != "1" {
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
	if err.Driver != "oracle" {
		t.Fatalf("unexpected driver: %s", err.Driver)
	}
	if !errors.Is(err, cause) {
		t.Fatal("generic error must wrap original cause")
	}
}
