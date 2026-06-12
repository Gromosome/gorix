package linter

import (
	"testing"

	linter2 "github.com/Gromosome/gorix/gorix/app/linter"
)

func TestValidationErrorFormatsWithAndWithoutPosition(t *testing.T) {
	withPosition := (&linter2.ValidationError{
		File:    "user.go",
		Line:    2,
		Column:  3,
		Message: "invalid",
	}).Error()
	if withPosition != "user.go:2:3: gorix validation error: invalid" {
		t.Fatalf("unexpected positioned error: %s", withPosition)
	}

	withoutPosition := (&linter2.ValidationError{
		File:    "user.go",
		Message: "invalid",
	}).Error()
	if withoutPosition != "user.go: gorix validation error: invalid" {
		t.Fatalf("unexpected unpositioned error: %s", withoutPosition)
	}
}
