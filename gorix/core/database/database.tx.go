package database

import (
	"fmt"
	"strings"

	gorixcontext "github.com/Gromosome/gorix/gorix/core/context"
)

func (tx *Tx) Exec(
	ctx *gorixcontext.Context,
	query string,
	args ...any,
) Result {
	if err := validateTxOperation(tx, ctx, query); err != nil {
		return Result{err: err}
	}

	result, err := tx.native.ExecContext(
		ctx,
		query,
		args...,
	)
	if err != nil {
		return Result{
			err: fmt.Errorf(
				"gorix database: transaction statement failed: %w",
				err,
			),
		}
	}

	return Result{
		native: result,
	}
}

func (tx *Tx) Query(
	ctx *gorixcontext.Context,
	query string,
	args ...any,
) *Rows {
	if err := validateTxOperation(tx, ctx, query); err != nil {
		return &Rows{err: err}
	}

	rows, err := tx.native.QueryContext(
		ctx,
		query,
		args...,
	)
	if err != nil {
		return &Rows{
			err: fmt.Errorf(
				"gorix database: transaction query failed: %w",
				err,
			),
		}
	}

	return &Rows{
		native: rows,
	}
}

func (tx *Tx) QueryRow(
	ctx *gorixcontext.Context,
	query string,
	args ...any,
) *Row {
	if err := validateTxOperation(tx, ctx, query); err != nil {
		return &Row{err: err}
	}

	return &Row{
		native: tx.native.QueryRowContext(
			ctx,
			query,
			args...,
		),
	}
}

func validateTxOperation(
	tx *Tx,
	ctx *gorixcontext.Context,
	query string,
) error {
	if tx == nil || tx.native == nil {
		return fmt.Errorf(
			"gorix database: transaction is unavailable",
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
