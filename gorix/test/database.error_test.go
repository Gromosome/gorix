package test

import (
	"errors"
	"testing"

	"github.com/Gromosome/gorix/gorix"
	sqldriver "github.com/Gromosome/gorix/sql-driver-manager"
)

func TestDatabaseErrorAliases(t *testing.T) {
	cause := errors.New("duplicate key")
	err := &gorix.DatabaseError{
		Kind:   gorix.DatabaseErrorDuplicateKey,
		Driver: "test",
		Cause:  cause,
	}

	got, ok := gorix.AsDatabaseError(err)
	if !ok {
		t.Fatal("AsDatabaseError did not recognize DatabaseError")
	}

	if got.Kind != sqldriver.ErrorDuplicateKey {
		t.Fatalf("unexpected database error kind: %v", got.Kind)
	}

	if !gorix.IsDatabaseErrorKind(err, gorix.DatabaseErrorDuplicateKey) {
		t.Fatal("IsDatabaseErrorKind did not match duplicate key")
	}
}
