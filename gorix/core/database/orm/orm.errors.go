package orm

import "errors"

var (
	ErrEntityNotFound = errors.New("gorix orm: entity not found")
	ErrMissingID      = errors.New("gorix orm: primary key value is missing")
)
