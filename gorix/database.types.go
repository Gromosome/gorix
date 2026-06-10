package gorix

import (
	"github.com/Gromosome/gorix/gorix/core/database"
	"github.com/Gromosome/gorix/gorix/core/database/mapper"
	"github.com/Gromosome/gorix/gorix/core/database/repository"
)

type DBManager = database.Manager
type DBConnection = database.Connection
type DBConfig = database.Config
type DBExecutor = database.Executor
type TransactionFunc = database.TransactionFunc

var WithTransaction = database.WithTransaction

type SQLMapper = mapper.Mapper
type SQLRepository[T any, ID comparable] = repository.Repository[T, ID]
type StatementRegistry = mapper.StatementRegistry

var NewSQLMapper = mapper.New
var SQLMapperExec = mapper.Exec

func NewSQLRepository[T any, ID comparable](
	manager *database.Manager,
	connectionNames ...string,
) (*SQLRepository[T, ID], error) {
	return repository.NewRepository[T, ID](
		manager,
		connectionNames...,
	)
}

type Dialect = repository.Dialect
type QueryBuilder = repository.QueryBuilder

var NewQueryBuilder = repository.NewQueryBuilder

var ErrEntityNotFound = repository.ErrEntityNotFound
var ErrMissingID = repository.ErrMissingID
var IsEntityNotFound = repository.IsEntityNotFound
var IsMissingID = repository.IsMissingID

var DBIsNoRows = database.IsNoRows

const (
	DBIsolationDefault         = database.IsolationDefault
	DBIsolationReadUncommitted = database.IsolationReadUncommitted
	DBIsolationReadCommitted   = database.IsolationReadCommitted
	DBIsolationWriteCommitted  = database.IsolationWriteCommitted
	DBIsolationRepeatableRead  = database.IsolationRepeatableRead
	DBIsolationSnapshot        = database.IsolationSnapshot
	DBIsolationSerializable    = database.IsolationSerializable
	DBIsolationLinearizable    = database.IsolationLinearizable
)

type DBTxOptions = database.TxOptions
type DBTx = database.Tx
