package document_driver_manager

import "context"

const (
	ErrorTransactionUnsupported ErrorKind = "transaction_unsupported"
	ErrorTransactionAborted     ErrorKind = "transaction_aborted"
)

type Executor interface {
	Collection(name string) Collection
}

type TxOptions struct {
	ReadOnly     bool
	ReadConcern  string
	WriteConcern string
}

type Tx interface {
	Executor
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type Transactor interface {
	BeginTx(ctx context.Context, options TxOptions) (Tx, error)
}
