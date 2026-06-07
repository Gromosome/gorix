package database

import (
	"fmt"

	gorixcontext "github.com/Gromosome/gorix/gorix/core/context"
)

func wrapDB(
	nativeDB DB,
) *DB {
	return &DB{
		native: nativeDB.native,
	}
}

func (db *DB) Exec(
	ctx *gorixcontext.Context,
	query string,
	args ...any,
) (Result, error) {
	if db == nil || db.native == nil {
		return Result{}, fmt.Errorf(
			"gorix database: database is unavailable",
		)
	}

	if ctx == nil {
		return Result{}, fmt.Errorf(
			"gorix database: context cannot be nil",
		)
	}

	result, err := db.native.ExecContext(
		ctx,
		query,
		args...,
	)
	if err != nil {
		return Result{}, err
	}

	return Result{
		native: result,
	}, nil
}

func (db *DB) Query(
	ctx *gorixcontext.Context,
	query string,
	args ...any,
) (*Rows, error) {
	if db == nil || db.native == nil {
		return nil, fmt.Errorf(
			"gorix database: database is unavailable",
		)
	}

	if ctx == nil {
		return nil, fmt.Errorf(
			"gorix database: context cannot be nil",
		)
	}

	rows, err := db.native.QueryContext(
		ctx,
		query,
		args...,
	)
	if err != nil {
		return nil, err
	}

	return &Rows{
		native: rows,
	}, nil
}

func (db *DB) QueryRow(
	ctx *gorixcontext.Context,
	query string,
	args ...any,
) *Row {
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
