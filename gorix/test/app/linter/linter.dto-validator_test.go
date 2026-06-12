package linter

import (
	"os"
	"path/filepath"
	"testing"

	linter2 "github.com/Gromosome/gorix/gorix/app/linter"
)

func TestValidateDTOFileAcceptsValidDTO(t *testing.T) {
	path := writeLintFile(t, "user.dto.go", `
package user
type CreateUserDto struct {
	Name string `+"`json:\"name\"`"+`
}`)

	if err := linter2.ValidateDTOFile(path); err != nil {
		t.Fatalf("validateDTOFile returned error: %v", err)
	}
}

func TestValidateDTOFileRejectsMissingTags(t *testing.T) {
	path := writeLintFile(t, "user.dto.go", `
package user
type CreateUserDto struct {
	Name string
}`)

	if err := linter2.ValidateDTOFile(path); err == nil {
		t.Fatal("expected DTO validation error")
	}
}

func writeLintFile(t *testing.T, name string, content string) string {
	t.Helper()

	dir := t.TempDir()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("failed to write lint file: %v", err)
	}
	return path
}
