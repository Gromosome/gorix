package database

import (
	"database/sql"
	"fmt"
)

func (r Result) Err() error {
	return r.err
}

func (r Result) Native() sql.Result {
	return r.native
}

func ErrResult(err error) Result {
	return Result{err: err}
}

func (r Result) LastInsertID() (
	int64,
	error,
) {
	if r.err != nil {
		return 0, r.err
	}

	if r.native == nil {
		return 0, fmt.Errorf(
			"gorix database: result is unavailable",
		)
	}

	return r.native.LastInsertId()
}

func (r Result) RowsAffected() (
	int64,
	error,
) {
	if r.err != nil {
		return 0, r.err
	}

	if r.native == nil {
		return 0, fmt.Errorf(
			"gorix database: result is unavailable",
		)
	}

	return r.native.RowsAffected()
}
