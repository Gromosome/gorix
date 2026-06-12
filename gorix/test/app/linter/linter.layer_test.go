package linter

import (
	"os"
	"path/filepath"
	"testing"

	linter2 "github.com/Gromosome/gorix/gorix/app/linter"
)

func TestValidateControllerFileAcceptsConstructorAndRoutes(t *testing.T) {
	path := writeLintFile(t, "user.controller.go", `
package user
import "github.com/Gromosome/gorix/gorix"
type UserController struct {}
func NewUserController() *UserController { return &UserController{} }
func (UserController) Find() (gorix.Method, gorix.Path, gorix.RouteHandler) { return gorix.GET, "/", nil }
`)

	if err := linter2.ValidateControllerFile(path); err != nil {
		t.Fatalf("validateControllerFile returned error: %v", err)
	}
}

func TestValidateModuleFileAcceptsModuleContract(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "user")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatalf("failed to create dir: %v", err)
	}
	path := filepath.Join(dir, "user.module.go")
	content := `
package user
import "github.com/Gromosome/gorix/gorix"
type UserModule struct {}
func NewUserModule() *UserModule { return &UserModule{} }
func (UserModule) BasePath() gorix.BasePath { return "/users" }
func (UserModule) Controllers() []any { return nil }
func (UserModule) Providers() []any { return nil }
`
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("failed to write module file: %v", err)
	}

	if err := linter2.ValidateModuleFile(path, "user"); err != nil {
		t.Fatalf("validateModuleFile returned error: %v", err)
	}
}

func TestValidateServiceFileAcceptsConstructor(t *testing.T) {
	path := writeLintFile(t, "user.service.go", `
package user
type UserService struct {}
func NewUserService() *UserService { return &UserService{} }
func (UserService) Find() {}
`)

	if err := linter2.ValidateServiceFile(path); err != nil {
		t.Fatalf("validateServiceFile returned error: %v", err)
	}
}
