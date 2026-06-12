package database

import (
	"database/sql"
	"fmt"
)

func (r Result) Err() error {
	return r.err
}

func NewErrResult(err error) Result {
	return Result{err: err}
}

func NewResult(native sql.Result) Result {
	return Result{native: native}
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
