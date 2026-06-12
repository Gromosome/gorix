package context

import (
	"strings"
	"testing"

	context2 "github.com/Gromosome/gorix/gorix/core/context"
)

type validationDTO struct {
	Name  string   `json:"name" validate:"required,min=3,max=5"`
	Email string   `json:"email" validate:"email"`
	Role  string   `json:"role" validate:"oneof=admin user"`
	Code  string   `json:"code" validate:"regex=^[A-Z]{2},[0-9]{2}$"`
	Tags  []string `json:"tags" validate:"min=1,max=2"`
}

func TestValidateStructAcceptsValidDTO(t *testing.T) {
	dto := validationDTO{
		Name:  "bob",
		Email: "bob@example.com",
		Role:  "admin",
		Code:  "AB,12",
		Tags:  []string{"a"},
	}

	if err := context2.ValidateStruct(dto); err != nil {
		t.Fatalf("ValidateStruct returned error: %v", err)
	}
}

func TestValidateStructReturnsFieldErrors(t *testing.T) {
	dto := validationDTO{
		Name:  "bo",
		Email: "invalid",
		Role:  "guest",
		Code:  "bad",
		Tags:  []string{"a", "b", "c"},
	}

	err := context2.ValidateStruct(dto)
	if err == nil {
		t.Fatal("expected validation error")
	}

	message := err.Error()
	for _, part := range []string{"name must be at least 3", "email must be a valid email", "role must be one of", "code format is invalid", "tags must be at most 2"} {
		if !strings.Contains(message, part) {
			t.Fatalf("expected message %q in %q", part, message)
		}
	}
}

func TestValidateStructIgnoresNilAndNonStruct(t *testing.T) {
	if err := context2.ValidateStruct(nil); err != nil {
		t.Fatalf("nil target should be valid: %v", err)
	}
	if err := context2.ValidateStruct("not struct"); err != nil {
		t.Fatalf("non-struct target should be valid: %v", err)
	}
}
