package database

import (
	"errors"
	"testing"

	database2 "github.com/Gromosome/gorix/gorix/core/database"
)

type fakeSQLResult struct {
	lastInsertID int64
	rowsAffected int64
}

func (r fakeSQLResult) LastInsertId() (int64, error) { return r.lastInsertID, nil }
func (r fakeSQLResult) RowsAffected() (int64, error) { return r.rowsAffected, nil }

func TestResultReturnsNativeValues(t *testing.T) {
	result := database2.Result{native: fakeSQLResult{lastInsertID: 10, rowsAffected: 2}}

	id, err := result.LastInsertID()
	if err != nil {
		t.Fatalf("LastInsertID returned error: %v", err)
	}
	if id != 10 {
		t.Fatalf("unexpected last insert id: %d", id)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("RowsAffected returned error: %v", err)
	}
	if rows != 2 {
		t.Fatalf("unexpected rows affected: %d", rows)
	}
}

func TestResultReturnsStoredError(t *testing.T) {
	cause := errors.New("failed")
	result := database2.ErrResult(cause)

	if _, err := result.LastInsertID(); !errors.Is(err, cause) {
		t.Fatalf("LastInsertID did not return stored error: %v", err)
	}
	if _, err := result.RowsAffected(); !errors.Is(err, cause) {
		t.Fatalf("RowsAffected did not return stored error: %v", err)
	}
}
