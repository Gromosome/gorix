package database

import (
	"errors"
	"testing"

	database2 "github.com/Gromosome/gorix/gorix/core/database"
)

func TestIsNoRows(t *testing.T) {
	if !database2.IsNoRows(database2.ErrNoRows) {
		t.Fatal("ErrNoRows should be recognized")
	}
	if database2.IsNoRows(errors.New("other")) {
		t.Fatal("unrelated error should not be recognized as no rows")
	}
}

func TestNewErrorResult(t *testing.T) {
	cause := errors.New("failed")
	result := database2.NewErrorResult(cause)
	if result.Err() != cause {
		t.Fatal("NewErrorResult did not preserve error")
	}
}
