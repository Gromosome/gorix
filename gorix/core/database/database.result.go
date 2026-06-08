package database

import "fmt"

func (r Result) LastInsertID() (
	int64,
	error,
) {
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
	if r.native == nil {
		return 0, fmt.Errorf(
			"gorix database: result is unavailable",
		)
	}
	return r.native.RowsAffected()
}
