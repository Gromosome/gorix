package gorix

import (
	"github.com/Gromosome/gorix/gorix/core/database"
	"github.com/Gromosome/gorix/gorix/core/database/mapper"
	"github.com/Gromosome/gorix/gorix/core/database/orm"
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

type Dialect = orm.Dialect
type QueryBuilder = orm.QueryBuilder

var NewQueryBuilder = orm.NewQueryBuilder

var ErrEntityNotFound = orm.ErrEntityNotFound
var ErrMissingID = orm.ErrMissingID
