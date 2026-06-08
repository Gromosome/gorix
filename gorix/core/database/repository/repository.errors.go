package repository

import "errors"

var (
	ErrEntityNotFound = errors.New("gorix repository: entity not found")
	ErrMissingID      = errors.New("gorix repository: primary key value is missing")
)
