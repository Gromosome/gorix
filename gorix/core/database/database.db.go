package database

import (
	"fmt"

	gorixcontext "github.com/Gromosome/gorix/gorix/core/context"
)

func (db *DB) Exec(
	ctx *gorixcontext.Context,
	query string,
	args ...any,
) Result {
	if db == nil || db.native == nil {
		return Result{err: fmt.Errorf(
			"gorix database: database is unavailable",
		)}
	}

	if ctx == nil {
		return Result{err: fmt.Errorf(
			"gorix database: context cannot be nil",
		)}
	}

	result, err := db.native.ExecContext(
		ctx,
		query,
		args...,
	)
	if err != nil {
		return Result{err: err}
	}

	return Result{
		native: result,
	}
}

func (db *DB) Query(
	ctx *gorixcontext.Context,
	query string,
	args ...any,
) *Rows {
	if db == nil || db.native == nil {
		return &Rows{err: fmt.Errorf(
			"gorix database: database is unavailable",
		)}
	}

	if ctx == nil {
		return &Rows{err: fmt.Errorf(
			"gorix database: context cannot be nil",
		)}
	}

	rows, err := db.native.QueryContext(
		ctx,
		query,
		args...,
	)
	if err != nil {
		return &Rows{err: err}
	}

	return &Rows{
		native: rows,
	}
}

func (db *DB) QueryRow(
	ctx *gorixcontext.Context,
	query string,
	args ...any,
) *Row {
	if db == nil || db.native == nil {
		return &Row{err: fmt.Errorf(
			"gorix database: database is unavailable",
		),
		}
	}
	return &Row{
		native: db.native.QueryRowContext(
			ctx,
			query,
			args...,
		),
	}
}

func (db *DB) Ping(
	ctx *gorixcontext.Context,
) error {
	if db == nil || db.native == nil {
		return fmt.Errorf(
			"gorix database: database is unavailable",
		)
	}

	return db.native.PingContext(ctx)
}
