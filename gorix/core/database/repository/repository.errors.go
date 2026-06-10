package repository

import "errors"

var (
	ErrEntityNotFound = errors.New("gorix repository: entity not found")
	ErrMissingID      = errors.New("gorix repository: primary key value is missing")
)

func IsEntityNotFound(err error) bool {
	return errors.Is(err, ErrEntityNotFound)
}
func IsMissingID(err error) bool {
	return errors.Is(err, ErrMissingID)
}
