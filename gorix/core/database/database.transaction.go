package database

import (
	"database/sql"
	"errors"
	"fmt"

	gorixcontext "github.com/Gromosome/gorix/gorix/core/context"
)

type TransactionFunc func(
	ctx *gorixcontext.Context,
	tx *Tx,
) error

func WithTransaction(
	ctx *gorixcontext.Context,
	db *DB,
	options *TxOptions,
	fn TransactionFunc,
) error {
	if ctx == nil {
		return fmt.Errorf(
			"gorix database: context cannot be nil",
		)
	}

	if db == nil || db.native == nil {
		return fmt.Errorf(
			"gorix database: database cannot be nil",
		)
	}

	if fn == nil {
		return fmt.Errorf(
			"gorix database: transaction function cannot be nil",
		)
	}

	nativeTx, err := db.native.BeginTx(
		ctx,
		toNativeTxOptions(options),
	)
	if err != nil {
		return fmt.Errorf(
			"gorix database: begin transaction failed: %w",
			err,
		)
	}

	tx := &Tx{
		native: nativeTx,
	}

	defer func() {
		if recovered := recover(); recovered != nil {
			_ = nativeTx.Rollback()
			panic(recovered)
		}
	}()

	if err := fn(ctx, tx); err != nil {
		rollbackErr := nativeTx.Rollback()

		if rollbackErr != nil &&
			!errors.Is(
				rollbackErr,
				sql.ErrTxDone,
			) {
			return fmt.Errorf(
				"gorix database: transaction failed: %v; rollback failed: %w",
				err,
				rollbackErr,
			)
		}

		return err
	}

	if err := nativeTx.Commit(); err != nil {
		return fmt.Errorf(
			"gorix database: commit failed: %w",
			err,
		)
	}

	return nil
}

func toNativeTxOptions(
	options *TxOptions,
) *sql.TxOptions {
	if options == nil {
		return nil
	}

	return &sql.TxOptions{
		Isolation: toNativeIsolation(
			options.Isolation,
		),
		ReadOnly: options.ReadOnly,
	}
}

func toNativeIsolation(
	level IsolationLevel,
) sql.IsolationLevel {
	switch level {
	case IsolationReadUncommitted:
		return sql.LevelReadUncommitted
	case IsolationReadCommitted:
		return sql.LevelReadCommitted
	case IsolationWriteCommitted:
		return sql.LevelWriteCommitted
	case IsolationRepeatableRead:
		return sql.LevelRepeatableRead
	case IsolationSnapshot:
		return sql.LevelSnapshot
	case IsolationSerializable:
		return sql.LevelSerializable
	case IsolationLinearizable:
		return sql.LevelLinearizable
	default:
		return sql.LevelDefault
	}
}
