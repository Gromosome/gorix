package database

import (
	"testing"

	gorixcontext "github.com/Gromosome/gorix/gorix/core/context"
	database2 "github.com/Gromosome/gorix/gorix/core/database"
)

func TestTxOperationsValidateInputs(t *testing.T) {
	var tx *database2.Tx
	if err := tx.Exec(gorixcontext.Background(), "SELECT 1").Err(); err == nil {
		t.Fatal("nil Tx Exec should return error")
	}

	tx = &database2.Tx{}
	if err := tx.Exec(nil, "SELECT 1").Err(); err == nil {
		t.Fatal("nil context should return error")
	}
}
