package database

import (
	"database/sql"
	"errors"
	"fmt"
)

func (r *Rows) Next() bool {
	if r == nil || r.err != nil || r.native == nil {
		return false
	}

	return r.native.Next()
}

func (r *Rows) Scan(
	destinations ...any,
) error {
	if r == nil {
		return fmt.Errorf(
			"gorix database: rows are unavailable",
		)
	}

	if r.err != nil {
		return r.err
	}

	if r.native == nil {
		return fmt.Errorf(
			"gorix database: native rows are unavailable",
		)
	}

	return r.native.Scan(destinations...)
}

func (r *Rows) Columns() (
	[]string,
	error,
) {
	if r == nil {
		return nil, fmt.Errorf(
			"gorix database: rows are unavailable",
		)
	}

	if r.err != nil {
		return nil, r.err
	}

	if r.native == nil {
		return nil, fmt.Errorf(
			"gorix database: native rows are unavailable",
		)
	}

	return r.native.Columns()
}

func (r *Rows) Err() error {
	if r == nil {
		return fmt.Errorf(
			"gorix database: rows are unavailable",
		)
	}

	if r.err != nil {
		return r.err
	}

	if r.native == nil {
		return fmt.Errorf(
			"gorix database: native rows are unavailable",
		)
	}

	return r.native.Err()
}

func (r *Rows) Close() error {
	if r == nil {
		return nil
	}

	if r.native == nil {
		return nil
	}

	return r.native.Close()
}

func (r *Row) Scan(
	destinations ...any,
) error {
	if r == nil {
		return fmt.Errorf(
			"gorix database: row is unavailable",
		)
	}

	if r.err != nil {
		return r.err
	}

	if r.native == nil {
		return fmt.Errorf(
			"gorix database: native row is unavailable",
		)
	}

	err := r.native.Scan(destinations...)

	if errors.Is(err, sql.ErrNoRows) {
		return ErrNoRows
	}

	return err
}
