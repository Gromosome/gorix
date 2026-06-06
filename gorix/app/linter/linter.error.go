package linter

import "fmt"

type ValidationError struct {
	File    string
	Line    int
	Column  int
	Message string
}

func (e *ValidationError) Error() string {
	if e.Line > 0 && e.Column > 0 {
		return fmt.Sprintf(
			"%s:%d:%d: gorix validation error: %s",
			e.File,
			e.Line,
			e.Column,
			e.Message,
		)
	}

	return fmt.Sprintf(
		"%s: gorix validation error: %s",
		e.File,
		e.Message,
	)
}

func newValidationError(file string, line int, column int, message string) error {
	return &ValidationError{
		File:    file,
		Line:    line,
		Column:  column,
		Message: message,
	}
}
