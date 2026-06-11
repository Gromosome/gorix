package sql_driver_manager

import (
	"errors"
	"testing"
)

type testAdapter struct{}

func (testAdapter) Name() string          { return "test" }
func (testAdapter) SQLDriverName() string { return "test-sql-driver" }
func (testAdapter) Normalize(err error) *Error {
	if err == nil {
		return nil
	}
	return &Error{
		Kind:    ErrorUnknown,
		Driver:  "test",
		Message: err.Error(),
		Cause:   err,
	}
}

func TestRegisterAndLookup(t *testing.T) {
	Register(testAdapter{})

	adapter, err := Lookup("test")
	if err != nil {
		t.Fatalf("Lookup returned error: %v", err)
	}

	if adapter.SQLDriverName() != "test-sql-driver" {
		t.Fatalf("unexpected SQL driver name: %s", adapter.SQLDriverName())
	}
}

func TestErrorUnwrap(t *testing.T) {
	cause := errors.New("driver error")
	wrapped := (&testAdapter{}).Normalize(cause)

	if !errors.Is(wrapped, cause) {
		t.Fatal("normalized error must preserve original cause")
	}
}
