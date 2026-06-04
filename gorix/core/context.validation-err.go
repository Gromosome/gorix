package core

import "fmt"

type FieldError struct {
	Field   string `json:"field"`
	Rule    string `json:"rule"`
	Message string `json:"message"`
}

type ValidationError struct {
	Errors []FieldError `json:"errors"`
}

func (e *ValidationError) Error() string {
	msg := ""
	for i, fieldError := range e.Errors {
		msg += fieldError.Message
		if i < len(e.Errors)-1 {
			msg += ", "
		}
	}
	return msg
}
func (e *ValidationError) ErrorAsList() []FieldError {
	return e.Errors
}

func NewValidationError(errors []FieldError) error {
	return &ValidationError{
		Errors: errors,
	}
}

func NewFieldError(field string, rule string, message string) FieldError {
	return FieldError{
		Field:   field,
		Rule:    rule,
		Message: message,
	}
}

func RequiredError(field string) FieldError {
	return NewFieldError(field, "required", fmt.Sprintf("%s is required", field))
}
