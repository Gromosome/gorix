package database

import gorixcontext "github.com/Gromosome/gorix/gorix/core/context"

func (tx *Tx) Exec(
	ctx *gorixcontext.Context,
	query string,
	args ...any,
) (Result, error) {
	result, err := tx.native.ExecContext(
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

func (tx *Tx) Query(
	ctx *gorixcontext.Context,
	query string,
	args ...any,
) (*Rows, error) {
	rows, err := tx.native.QueryContext(
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
