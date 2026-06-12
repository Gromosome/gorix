package database

import (
	"testing"

	database2 "github.com/Gromosome/gorix/gorix/core/database"
)

func TestManagerConnectionLookupAndClose(t *testing.T) {
	manager := database2.NewManager()

	if manager.Has("missing") {
		t.Fatal("manager should not have missing connection")
	}
	if _, err := manager.Connection("missing"); err == nil {
		t.Fatal("missing connection should return error")
	}
	if err := manager.Close(); err != nil {
		t.Fatalf("closing empty manager returned error: %v", err)
	}
}
