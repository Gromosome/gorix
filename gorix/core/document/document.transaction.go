package document

import (
	"database/sql"
	"errors"
	"fmt"

	docdriver "github.com/Gromosome/gorix/document-driver-manager"
	"github.com/Gromosome/gorix/gorix/config"
	gorixcontext "github.com/Gromosome/gorix/gorix/core/context"
)

type Tx = docdriver.Tx
type TxOptions = docdriver.TxOptions

type TransactionFunc func(
	ctx *gorixcontext.Context,
	tx Tx,
) error

func WithTransaction(
	ctx *gorixcontext.Context,
	manager *Manager,
	connectionName string,
	options *TxOptions,
	fn TransactionFunc,
) error {
	if ctx == nil {
		return fmt.Errorf("gorix document: context cannot be nil")
	}

	if manager == nil {
		return fmt.Errorf("gorix document: manager cannot be nil")
	}

	if fn == nil {
		return fmt.Errorf("gorix document: transaction function cannot be nil")
	}

	if connectionName == "" {
		connectionName = config.DefaultConnectionName
	}

	connection, err := manager.Connection(connectionName)
	if err != nil {
		return err
	}

	database := connection.Database()
	if database == nil {
		return fmt.Errorf(
			"gorix document: database is unavailable for connection %q",
			connectionName,
		)
	}

	transactor, ok := database.(docdriver.Transactor)
	if !ok {
		return TransactionUnsupportedError(connectionName)
	}

	txOptions := docdriver.TxOptions{}
	if options != nil {
		txOptions = *options
	}

	tx, err := transactor.BeginTx(ctx, txOptions)
	if err != nil {
		return err
	}

	defer func() {
		if recovered := recover(); recovered != nil {
			_ = tx.Rollback(ctx)
			panic(recovered)
		}
	}()

	if err := fn(ctx, tx); err != nil {
		rollbackErr := tx.Rollback(ctx)
		if rollbackErr != nil &&
			!errors.Is(rollbackErr, sql.ErrTxDone) {
			return fmt.Errorf(
				"gorix document: transaction failed: %v; rollback failed: %w",
				err,
				rollbackErr,
			)
		}

		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("gorix document: commit failed: %w", err)
	}

	return nil
}
