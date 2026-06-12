package mapper

import (
	"testing"

	mapper2 "github.com/Gromosome/gorix/gorix/core/database/mapper"
)

func TestStatementRegistryRegisterGetAndHas(t *testing.T) {
	registry := mapper2.NewStatementRegistry()

	if registry.Has("find") {
		t.Fatal("registry should not have missing statement")
	}
	if err := registry.Register("find", "SELECT 1"); err != nil {
		t.Fatalf("Register returned error: %v", err)
	}
	if !registry.Has("find") {
		t.Fatal("registry should have registered statement")
	}
	statement, err := registry.Get("find")
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
	if statement != "SELECT 1" {
		t.Fatalf("unexpected statement: %s", statement)
	}
}

func TestStatementRegistryRejectsInvalidStatements(t *testing.T) {
	registry := mapper2.NewStatementRegistry()

	if err := registry.Register("", "SELECT 1"); err == nil {
		t.Fatal("empty name should be rejected")
	}
	if err := registry.Register("empty", " "); err == nil {
		t.Fatal("empty statement should be rejected")
	}
	if err := registry.Register("find", "SELECT 1"); err != nil {
		t.Fatalf("Register returned error: %v", err)
	}
	if err := registry.Register("find", "SELECT 2"); err == nil {
		t.Fatal("duplicate statement should be rejected")
	}
}
