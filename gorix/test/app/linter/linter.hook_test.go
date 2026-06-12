package linter

import (
	"testing"

	linter2 "github.com/Gromosome/gorix/gorix/app/linter"
)

func TestValidateMiddlewareFileAcceptsFactory(t *testing.T) {
	path := writeLintFile(t, "auth.middleware.go", `
package middleware
import "github.com/Gromosome/gorix/gorix"
func AuthMiddleware() gorix.Middleware { return nil }
`)

	if err := linter2.ValidateMiddlewareFile(path); err != nil {
		t.Fatalf("validateMiddlewareFile returned error: %v", err)
	}
}

func TestValidateInterceptorFileRequiresBeforeAndAfter(t *testing.T) {
	path := writeLintFile(t, "audit.interceptor.go", `
package interceptor
import "github.com/Gromosome/gorix/gorix"
type AuditInterceptor struct {}
func (AuditInterceptor) Before(ctx *gorix.ExecutionContext) error { return nil }
func (AuditInterceptor) After(ctx *gorix.ExecutionContext) error { return nil }
`)

	if err := linter2.ValidateInterceptorFile(path); err != nil {
		t.Fatalf("validateInterceptorFile returned error: %v", err)
	}

	invalidPath := writeLintFile(t, "audit.interceptor.go", `
package interceptor
type AuditInterceptor struct {}
`)
	if err := linter2.ValidateInterceptorFile(invalidPath); err == nil {
		t.Fatal("expected interceptor validation error")
	}
}

func TestValidateFilterFileRequiresCatch(t *testing.T) {
	path := writeLintFile(t, "exception.filter.go", `
package filter
import "github.com/Gromosome/gorix/gorix"
type ExceptionFilter struct {}
func (ExceptionFilter) Catch(ctx *gorix.ExceptionContext) {}
`)

	if err := linter2.ValidateFilterFile(path); err != nil {
		t.Fatalf("validateFilterFile returned error: %v", err)
	}
}
