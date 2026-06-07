package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type TransactionFunc func(tx *sql.Tx) error

func WithTransaction(
	ctx context.Context,
	db *sql.DB,
	options *sql.TxOptions,
	fn TransactionFunc,
) error {
	if db == nil {
		return fmt.Errorf("gorix database: database cannot be nil")
	}

	if fn == nil {
		return fmt.Errorf(
			"gorix database: transaction function cannot be nil",
		)
	}

	tx, err := db.BeginTx(ctx, options)
	if err != nil {
		return fmt.Errorf(
			"gorix database: failed to begin transaction: %w",
			err,
		)
	}

	committed := false

	defer func() {
		if !committed {
			_ = tx.Rollback()
		}
	}()

	defer func() {
		if recovered := recover(); recovered != nil {
			_ = tx.Rollback()
			panic(recovered)
		}
	}()

	if err := fn(tx); err != nil {
		rollbackErr := tx.Rollback()

		if rollbackErr != nil &&
			!errors.Is(rollbackErr, sql.ErrTxDone) {
			return fmt.Errorf(
				"gorix database: operation failed: %v; rollback failed: %w",
				err,
				rollbackErr,
			)
		}

		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf(
			"gorix database: failed to commit transaction: %w",
			err,
		)
	}

	committed = true
	return nil
}
