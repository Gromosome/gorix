package linter

import (
	"os"
	"path/filepath"
	"testing"

	linter2 "github.com/Gromosome/gorix/gorix/app/linter"
)

func TestValidatePackageDirectoryDispatchesValidators(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "user")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatalf("failed to create dir: %v", err)
	}

	files := map[string]string{
		"user.dto.go": `
package user
type UserDto struct { Name string ` + "`json:\"name\"`" + ` }
`,
		"user.service.go": `
package user
type UserService struct {}
func NewUserService() *UserService { return &UserService{} }
`,
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o600); err != nil {
			t.Fatalf("failed to write %s: %v", name, err)
		}
	}

	if err := linter2.ValidatePackageDirectory(dir); err != nil {
		t.Fatalf("ValidatePackageDirectory returned error: %v", err)
	}
}
