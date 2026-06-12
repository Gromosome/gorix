package context

import (
	"testing"

	context2 "github.com/Gromosome/gorix/gorix/core/context"
)

func TestValidationErrorFormatsMessages(t *testing.T) {
	err := context2.NewValidationError([]context2.FieldError{
		context2.NewFieldError("name", "required", "name is required"),
		context2.NewFieldError("email", "email", "email must be valid"),
	})

	if err.Error() != "name is required, email must be valid" {
		t.Fatalf("unexpected error message: %s", err.Error())
	}

	validationError := err.(*context2.ValidationError)
	if len(validationError.ErrorAsList()) != 2 {
		t.Fatalf("unexpected field error count: %d", len(validationError.ErrorAsList()))
	}
}
