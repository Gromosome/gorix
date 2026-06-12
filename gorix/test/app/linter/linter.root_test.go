package linter

import (
	"os"
	"path/filepath"
	"testing"

	linter2 "github.com/Gromosome/gorix/gorix/app/linter"
)

func TestValidateRootAllowsOnlyMainGoAtRoot(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte("package main\n"), 0o600); err != nil {
		t.Fatalf("failed to write main.go: %v", err)
	}
	if err := linter2.ValidateRoot(dir); err != nil {
		t.Fatalf("ValidateRoot returned error: %v", err)
	}

	if err := os.WriteFile(filepath.Join(dir, "extra.go"), []byte("package main\n"), 0o600); err != nil {
		t.Fatalf("failed to write extra.go: %v", err)
	}
	if err := linter2.ValidateRoot(dir); err == nil {
		t.Fatal("expected root validation error")
	}
}
