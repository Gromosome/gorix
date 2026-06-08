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

type Mapper = mapper.Mapper
type StatementRegistry = mapper.StatementRegistry

var NewMapper = mapper.New
var MapperExec = mapper.Exec

type Dialect = repository.Dialect
type QueryBuilder = repository.QueryBuilder

var NewQueryBuilder = repository.NewQueryBuilder

var ErrEntityNotFound = repository.ErrEntityNotFound
var ErrMissingID = repository.ErrMissingID
