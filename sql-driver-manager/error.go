package sql_driver_manager

import (
	"errors"
	"fmt"
)

type ErrorKind string

const (
	ErrorUnknown         ErrorKind = "unknown"
	ErrorDuplicateKey    ErrorKind = "duplicate_key"
	ErrorUniqueViolation ErrorKind = "unique_violation"
	ErrorForeignKey      ErrorKind = "foreign_key_violation"
	ErrorNotNull         ErrorKind = "not_null_violation"
	ErrorCheckConstraint ErrorKind = "check_constraint_violation"
	ErrorSerialization   ErrorKind = "serialization_failure"
	ErrorDeadlock        ErrorKind = "deadlock"
	ErrorConnection      ErrorKind = "connection_error"
	ErrorTimeout         ErrorKind = "timeout"
	ErrorSyntax          ErrorKind = "syntax_error"
	ErrorPermission      ErrorKind = "permission_denied"
	ErrorTableNotFound   ErrorKind = "table_not_found"
	ErrorColumnNotFound  ErrorKind = "column_not_found"
)

type Error struct {
	Kind       ErrorKind
	Driver     string
	Code       string
	Message    string
	Constraint string
	Table      string
	Column     string
	Cause      error
}

func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.Message != "" {
		return e.Message
	}
	if e.Cause != nil {
		return e.Cause.Error()
	}
	return fmt.Sprintf("sql driver error: kind=%s driver=%s code=%s", e.Kind, e.Driver, e.Code)
}

func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Cause
}

func AsError(err error) (*Error, bool) {
	var target *Error
	ok := errors.As(err, &target)
	return target, ok
}

func IsKind(err error, kind ErrorKind) bool {
	target, ok := AsError(err)
	return ok && target.Kind == kind
}
