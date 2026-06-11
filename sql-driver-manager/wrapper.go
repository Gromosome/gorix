package sql_driver_manager

import (
	"context"
	"database/sql"
)

type normalizer func(error) error

type Row struct {
	row       *sql.Row
	normalize normalizer
}

func (r *Row) Scan(destinations ...any) error {
	return r.normalize(r.row.Scan(destinations...))
}

type Rows struct {
	*sql.Rows
	normalize normalizer
}

func (r *Rows) Scan(destinations ...any) error {
	return r.normalize(r.Rows.Scan(destinations...))
}

func (r *Rows) Err() error {
	return r.normalize(r.Rows.Err())
}

func (r *Rows) Close() error {
	return r.normalize(r.Rows.Close())
}

type Stmt struct {
	stmt      *sql.Stmt
	normalize normalizer
}

func (s *Stmt) Close() error {
	return s.normalize(s.stmt.Close())
}

func (s *Stmt) ExecContext(
	ctx context.Context,
	args ...any,
) (sql.Result, error) {
	result, err := s.stmt.ExecContext(ctx, args...)
	return result, s.normalize(err)
}

func (s *Stmt) QueryContext(
	ctx context.Context,
	args ...any,
) (*Rows, error) {
	rows, err := s.stmt.QueryContext(ctx, args...)
	if err != nil {
		return nil, s.normalize(err)
	}

	return &Rows{
		Rows:      rows,
		normalize: s.normalize,
	}, nil
}

func (s *Stmt) QueryRowContext(
	ctx context.Context,
	args ...any,
) *Row {
	return &Row{
		row:       s.stmt.QueryRowContext(ctx, args...),
		normalize: s.normalize,
	}
}

type Tx struct {
	tx        *sql.Tx
	normalize normalizer
}

func (t *Tx) Commit() error {
	return t.normalize(t.tx.Commit())
}

func (t *Tx) Rollback() error {
	return t.normalize(t.tx.Rollback())
}

func (t *Tx) ExecContext(
	ctx context.Context,
	query string,
	args ...any,
) (sql.Result, error) {
	result, err := t.tx.ExecContext(ctx, query, args...)
	return result, t.normalize(err)
}

func (t *Tx) QueryContext(
	ctx context.Context,
	query string,
	args ...any,
) (*Rows, error) {
	rows, err := t.tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, t.normalize(err)
	}

	return &Rows{
		Rows:      rows,
		normalize: t.normalize,
	}, nil
}

func (t *Tx) QueryRowContext(
	ctx context.Context,
	query string,
	args ...any,
) *Row {
	return &Row{
		row:       t.tx.QueryRowContext(ctx, query, args...),
		normalize: t.normalize,
	}
}

func (t *Tx) PrepareContext(
	ctx context.Context,
	query string,
) (*Stmt, error) {
	statement, err := t.tx.PrepareContext(ctx, query)
	if err != nil {
		return nil, t.normalize(err)
	}

	return &Stmt{
		stmt:      statement,
		normalize: t.normalize,
	}, nil
}
