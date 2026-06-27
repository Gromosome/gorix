package document_driver_manager

import "errors"

type ErrorKind string

const (
	ErrorUnknown            ErrorKind = "unknown"
	ErrorNotFound           ErrorKind = "not_found"
	ErrorDuplicateKey       ErrorKind = "duplicate_key"
	ErrorConflict           ErrorKind = "conflict"
	ErrorValidation         ErrorKind = "validation"
	ErrorConnection         ErrorKind = "connection_error"
	ErrorTimeout            ErrorKind = "timeout"
	ErrorPermissionDenied   ErrorKind = "permission_denied"
	ErrorDatabaseNotFound   ErrorKind = "database_not_found"
	ErrorCollectionNotFound ErrorKind = "collection_not_found"
	ErrorInvalidQuery       ErrorKind = "invalid_query"
)

type Error struct {
	Kind       ErrorKind
	Driver     string
	Code       string
	Message    string
	Database   string
	Collection string
	Cause      error
}

func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	if e.Message != "" {
		return e.Message
	}
	return string(e.Kind)
}

func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Cause
}

func AsError(err error) (*Error, bool) {
	var dbErr *Error
	if errors.As(err, &dbErr) {
		return dbErr, true
	}
	return nil, false
}

func IsKind(err error, kind ErrorKind) bool {
	dbErr, ok := AsError(err)
	return ok && dbErr.Kind == kind
}
