package mapper

import (
	"context"
	"database/sql"
	"fmt"

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

func (m *Mapper) DB() (*sql.DB, error) {
	if m == nil || m.manager == nil {
		return nil, fmt.Errorf(
			"gorix mapper: database manager is unavailable",
		)
	}

	return m.manager.DB(m.connectionName)
}

func (m *Mapper) QueryOne(
	ctx context.Context,
	query string,
	args ...any,
) (*Row, error) {
	db, err := m.DB()
	if err != nil {
		return nil, err
	}

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return &Row{
		rows: rows,
	}, nil
}

func (m *Mapper) RegisterStatement(
	name string,
	query string,
) error {
	return m.statements.Register(name, query)
}

func (m *Mapper) Statement(name string) (string, error) {
	return m.statements.Get(name)
}

type Row struct {
	rows *sql.Rows
}

func (r *Row) Close() error {
	if r == nil || r.rows == nil {
		return nil
	}

	return r.rows.Close()
}
