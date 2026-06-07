package database

import (
	"database/sql"
	"errors"
)

var ErrNoRows = errors.New(
	"gorix database: no rows found",
)

func IsNoRows(err error) bool {
	return errors.Is(err, ErrNoRows) ||
		errors.Is(err, sql.ErrNoRows)
}
