package repository

import (
	"fmt"
)

type Dialect interface {
	Name() string
	Placeholder(position int) string
	QuoteIdentifier(identifier string) string
	SupportsReturning() bool
}

type PostgresDialect struct{}

func (PostgresDialect) Name() string {
	return "postgres"
}

func (PostgresDialect) Placeholder(position int) string {
	return fmt.Sprintf("$%d", position)
}

func (PostgresDialect) QuoteIdentifier(identifier string) string {
	return `"` + identifier + `"`
}

func (PostgresDialect) SupportsReturning() bool {
	return true
}

type MySQLDialect struct{}

func (MySQLDialect) Name() string {
	return "mysql"
}

func (MySQLDialect) Placeholder(_ int) string {
	return "?"
}

func (MySQLDialect) QuoteIdentifier(identifier string) string {
	return "`" + identifier + "`"
}

func (MySQLDialect) SupportsReturning() bool {
	return false
}

type SQLiteDialect struct{}

func (SQLiteDialect) Name() string {
	return "sqlite3"
}

func (SQLiteDialect) Placeholder(_ int) string {
	return "?"
}

func (SQLiteDialect) QuoteIdentifier(identifier string) string {
	return `"` + identifier + `"`
}

func (SQLiteDialect) SupportsReturning() bool {
	return true
}

func ResolveDialect(driver string) (Dialect, error) {
	switch driver {
	case "pgx", "postgres", "postgresql":
		return PostgresDialect{}, nil

	case "mysql":
		return MySQLDialect{}, nil

	case "sqlite3", "sqlite3":
		return SQLiteDialect{}, nil

	default:
		return nil, fmt.Errorf(
			"gorix repository: unsupported driver dialect %q",
			driver,
		)
	}
}
