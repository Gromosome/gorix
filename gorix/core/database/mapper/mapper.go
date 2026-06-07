package mapper

import (
	"fmt"

	gorixcontext "github.com/Gromosome/gorix/gorix/core/context"
	"github.com/Gromosome/gorix/gorix/core/database"
)

type Mapper struct {
	manager        *database.Manager
	connectionName string
	statements     *StatementRegistry
}

func New(
	manager *database.Manager,
	connectionNames ...string,
) *Mapper {
	connectionName := database.DefaultConnectionName

	if len(connectionNames) > 0 &&
		connectionNames[0] != "" {
		connectionName = connectionNames[0]
	}

	return &Mapper{
		manager:        manager,
		connectionName: connectionName,
		statements:     NewStatementRegistry(),
	}
}

func (m *Mapper) ConnectionName() string {
	if m == nil {
		return ""
	}

	return m.connectionName
}

func (m *Mapper) DB() (*database.DB, error) {
	if m == nil {
		return nil, fmt.Errorf(
			"gorix mapper: mapper cannot be nil",
		)
	}

	if m.manager == nil {
		return nil, fmt.Errorf(
			"gorix mapper: database manager is unavailable",
		)
	}

	connectionName := m.connectionName

	if connectionName == "" {
		connectionName = database.DefaultConnectionName
	}

	db, err := m.manager.DB(connectionName)
	if err != nil {
		return nil, fmt.Errorf(
			"gorix mapper: failed to resolve database connection %q: %w",
			connectionName,
			err,
		)
	}

	return db, nil
}

func (m *Mapper) QueryOne(
	ctx *gorixcontext.Context,
	target any,
	query string,
	args ...any,
) error {
	if err := validateMapperContext(ctx); err != nil {
		return err
	}

	db, err := m.DB()
	if err != nil {
		return err
	}

	return QueryOneInto(
		ctx,
		db,
		target,
		query,
		args...,
	)
}

func (m *Mapper) QueryMany(
	ctx *gorixcontext.Context,
	target any,
	query string,
	args ...any,
) error {
	if err := validateMapperContext(ctx); err != nil {
		return err
	}

	db, err := m.DB()
	if err != nil {
		return err
	}

	return QueryManyInto(
		ctx,
		db,
		target,
		query,
		args...,
	)
}

func (m *Mapper) Exec(
	ctx *gorixcontext.Context,
	query string,
	args ...any,
) (database.Result, error) {
	if err := validateMapperContext(ctx); err != nil {
		return database.Result{}, err
	}

	db, err := m.DB()
	if err != nil {
		return database.Result{}, err
	}

	return Exec(
		ctx,
		db,
		query,
		args...,
	)
}

func (m *Mapper) RegisterStatement(
	name string,
	query string,
) error {
	if m == nil {
		return fmt.Errorf(
			"gorix mapper: mapper cannot be nil",
		)
	}

	if m.statements == nil {
		m.statements = NewStatementRegistry()
	}

	return m.statements.Register(name, query)
}

func (m *Mapper) MustRegisterStatement(
	name string,
	query string,
) {
	if err := m.RegisterStatement(name, query); err != nil {
		panic(err)
	}
}

func (m *Mapper) Statement(
	name string,
) (string, error) {
	if m == nil {
		return "", fmt.Errorf(
			"gorix mapper: mapper cannot be nil",
		)
	}

	if m.statements == nil {
		return "", fmt.Errorf(
			"gorix mapper: statement registry is unavailable",
		)
	}

	return m.statements.Get(name)
}

func (m *Mapper) QueryNamedOne(
	ctx *gorixcontext.Context,
	target any,
	statementName string,
	args ...any,
) error {
	query, err := m.Statement(statementName)
	if err != nil {
		return err
	}

	return m.QueryOne(
		ctx,
		target,
		query,
		args...,
	)
}

func (m *Mapper) QueryNamedMany(
	ctx *gorixcontext.Context,
	target any,
	statementName string,
	args ...any,
) error {
	query, err := m.Statement(statementName)
	if err != nil {
		return err
	}

	return m.QueryMany(
		ctx,
		target,
		query,
		args...,
	)
}

func (m *Mapper) ExecNamed(
	ctx *gorixcontext.Context,
	statementName string,
	args ...any,
) (database.Result, error) {
	query, err := m.Statement(statementName)
	if err != nil {
		return database.Result{}, err
	}

	return m.Exec(
		ctx,
		query,
		args...,
	)
}

func (m *Mapper) WithExecutor(
	executor database.Executor,
) *ExecutorMapper {
	return &ExecutorMapper{
		executor: executor,
	}
}

type ExecutorMapper struct {
	executor database.Executor
}

func (m *ExecutorMapper) QueryOne(
	ctx *gorixcontext.Context,
	target any,
	query string,
	args ...any,
) error {
	return QueryOneInto(
		ctx,
		m.executor,
		target,
		query,
		args...,
	)
}

func (m *ExecutorMapper) QueryMany(
	ctx *gorixcontext.Context,
	target any,
	query string,
	args ...any,
) error {
	return QueryManyInto(
		ctx,
		m.executor,
		target,
		query,
		args...,
	)
}

func (m *ExecutorMapper) Exec(
	ctx *gorixcontext.Context,
	query string,
	args ...any,
) (database.Result, error) {
	return Exec(
		ctx,
		m.executor,
		query,
		args...,
	)
}
