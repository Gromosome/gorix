package database

import (
	"errors"
	"testing"

	database2 "github.com/Gromosome/gorix/gorix/core/database"
)

func TestNilRowsAndRowReturnErrors(t *testing.T) {
	var rows *database2.Rows
	if rows.Next() {
		t.Fatal("nil rows should not advance")
	}
	if _, err := rows.Columns(); err == nil {
		t.Fatal("nil rows Columns should return error")
	}
	if err := rows.Scan(); err == nil {
		t.Fatal("nil rows Scan should return error")
	}
	if err := rows.Err(); err == nil {
		t.Fatal("nil rows Err should return error")
	}
	if err := rows.Close(); err != nil {
		t.Fatalf("nil rows Close should be nil, got %v", err)
	}

	var row *database2.Row
	if err := row.Scan(); err == nil {
		t.Fatal("nil row Scan should return error")
	}
}

func TestRowsAndRowReturnStoredErrors(t *testing.T) {
	cause := errors.New("failed")
	rows := &database2.Rows{err: cause}
	if rows.Next() {
		t.Fatal("errored rows should not advance")
	}
	if err := rows.Scan(); !errors.Is(err, cause) {
		t.Fatalf("Scan did not return stored error: %v", err)
	}
	if _, err := rows.Columns(); !errors.Is(err, cause) {
		t.Fatalf("Columns did not return stored error: %v", err)
	}
	if err := rows.Err(); !errors.Is(err, cause) {
		t.Fatalf("Err did not return stored error: %v", err)
	}

	row := &database2.Row{err: cause}
	if err := row.Scan(); !errors.Is(err, cause) {
		t.Fatalf("row Scan did not return stored error: %v", err)
	}
}
