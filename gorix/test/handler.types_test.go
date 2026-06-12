package test

import (
	"testing"

	"github.com/Gromosome/gorix/gorix"
	"github.com/Gromosome/gorix/gorix/core/context"
)

func TestHandlerTypeAliasesExposeRoutePrimitives(t *testing.T) {
	if gorix.GET != context.GET || gorix.POST != context.POST || gorix.DELETE != context.DELETE {
		t.Fatal("HTTP method aliases do not match context package")
	}

	if gorix.StatusOK.Int() != 200 {
		t.Fatalf("unexpected StatusOK value: %d", gorix.StatusOK.Int())
	}

	fieldError := gorix.NewFieldError("email", "required", "email is required")
	err := gorix.NewValidationError([]gorix.FieldError{fieldError})
	validationError, ok := err.(*gorix.ValidationError)
	if !ok {
		t.Fatalf("unexpected validation error type %T", err)
	}
	if validationError.Error() != "email is required" {
		t.Fatalf("unexpected validation message: %s", validationError.Error())
	}
}
