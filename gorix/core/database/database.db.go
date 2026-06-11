package database

import (
	"fmt"
	"strings"

	gorixcontext "github.com/Gromosome/gorix/gorix/core/context"
)

func (db *DB) Exec(
	ctx *gorixcontext.Context,
	query string,
	args ...any,
) Result {
	if err := validateDBOperation(
		db,
		ctx,
		query,
	); err != nil {
		return Result{err: err}
	}

	result, err := db.native.ExecContext(
		ctx,
		query,
		args...,
	)
	if err != nil {
		return Result{
			err: fmt.Errorf(
				"gorix database: statement execution failed: %w",
				err,
			),
		}
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
	if err := validateDBOperation(
		db,
		ctx,
		query,
	); err != nil {
		return &Rows{err: err}
	}

	rows, err := db.native.QueryContext(
		ctx,
		query,
		args...,
	)
	if err != nil {
		return &Rows{
			err: fmt.Errorf(
				"gorix database: query execution failed: %w",
				err,
			),
		}
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
	if err := validateDBOperation(
		db,
		ctx,
		query,
	); err != nil {
		return &Row{err: err}
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

	if ctx == nil {
		return fmt.Errorf(
			"gorix database: context cannot be nil",
		)
	}

	if err := ctx.Err(); err != nil {
		return fmt.Errorf(
			"gorix database: context is closed: %w",
			err,
		)
	}

	return db.native.PingContext(ctx)
}

func validateDBOperation(
	db *DB,
	ctx *gorixcontext.Context,
	query string,
) error {
	if db == nil || db.native == nil {
		return fmt.Errorf(
			"gorix database: database is unavailable",
		)
	}

	if ctx == nil {
		return fmt.Errorf(
			"gorix database: context cannot be nil",
		)
	}

	if err := ctx.Err(); err != nil {
		return fmt.Errorf(
			"gorix database: context is closed: %w",
			err,
		)
	}

	if strings.TrimSpace(query) == "" {
		return fmt.Errorf(
			"gorix database: query cannot be empty",
		)
	}

	return nil
}
