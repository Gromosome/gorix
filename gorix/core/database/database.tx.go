package database

import gorixcontext "github.com/Gromosome/gorix/gorix/core/context"

func (tx *Tx) Exec(
	ctx *gorixcontext.Context,
	query string,
	args ...any,
) Result {
	result, err := tx.native.ExecContext(
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

func (tx *Tx) Query(
	ctx *gorixcontext.Context,
	query string,
	args ...any,
) *Rows {
	rows, err := tx.native.QueryContext(
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

func (tx *Tx) QueryRow(
	ctx *gorixcontext.Context,
	query string,
	args ...any,
) *Row {
	return &Row{
		native: tx.native.QueryRowContext(
			ctx,
			query,
			args...,
		),
	}
}
