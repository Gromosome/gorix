package database

import (
	"strings"
	"testing"

	gorixcontext "github.com/Gromosome/gorix/gorix/core/context"
	database2 "github.com/Gromosome/gorix/gorix/core/database"
)

func TestDBOperationsValidateInputs(t *testing.T) {
	var db *database2.DB
	if err := db.Exec(gorixcontext.Background(), "SELECT 1").Err(); err == nil {
		t.Fatal("nil DB Exec should return error")
	}

	db = &database2.DB{}
	if err := db.Exec(nil, "SELECT 1").Err(); err == nil {
		t.Fatal("nil context should return error")
	}

	if err := db.Exec(gorixcontext.Background(), "  ").Err(); err == nil || !strings.Contains(err.Error(), "database is unavailable") {
		t.Fatalf("unexpected validation ordering/error: %v", err)
	}
}
