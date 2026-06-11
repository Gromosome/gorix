package repository

import (
	"fmt"
	"strings"
)

type GeneratedKeyStrategy uint8

const (
	GeneratedKeyUnsupported GeneratedKeyStrategy = iota

	// MySQL and similar drivers expose the generated key through Result.
	GeneratedKeyLastInsertID

	// PostgreSQL and modern SQLite use:
	// INSERT ... RETURNING id
	GeneratedKeyReturning

	// SQL Server uses:
	// INSERT ... OUTPUT INSERTED.id VALUES ...
	GeneratedKeyOutputInserted

	// Oracle uses:
	// INSERT ... RETURNING id INTO :n
	GeneratedKeyReturningInto
)

type Dialect interface {
	// Name returns the canonical Gorix wrapper name.
	Name() string

	// Placeholder returns the SQL parameter placeholder for the position.
	Placeholder(position int) string

	// QuoteIdentifier safely quotes table and column identifiers.
	//
	// It also supports qualified identifiers such as:
	// public.users
	// users.id
	QuoteIdentifier(identifier string) string

	// SupportsReturning reports whether the dialect supports the literal
	// trailing "RETURNING column" SQL clause.
	//
	// Keep this temporarily for compatibility with the current repository.
	SupportsReturning() bool

	// GeneratedKeyStrategy describes how generated primary keys are obtained.
	GeneratedKeyStrategy() GeneratedKeyStrategy
}

// -----------------------------------------------------------------------------
// PostgreSQL
// -----------------------------------------------------------------------------

type PostgresDialect struct{}

func (PostgresDialect) Name() string {
	return "postgres"
}

func (PostgresDialect) Placeholder(position int) string {
	return fmt.Sprintf("$%d", normalizePosition(position))
}

func (PostgresDialect) QuoteIdentifier(identifier string) string {
	return quoteQualifiedIdentifier(
		identifier,
		`"`,
		`"`,
	)
}

func (PostgresDialect) SupportsReturning() bool {
	return true
}

func (PostgresDialect) GeneratedKeyStrategy() GeneratedKeyStrategy {
	return GeneratedKeyReturning
}

// -----------------------------------------------------------------------------
// MySQL
// -----------------------------------------------------------------------------

type MySQLDialect struct{}

func (MySQLDialect) Name() string {
	return "mysql"
}

func (MySQLDialect) Placeholder(_ int) string {
	return "?"
}

func (MySQLDialect) QuoteIdentifier(identifier string) string {
	return quoteQualifiedIdentifier(
		identifier,
		"`",
		"`",
	)
}

func (MySQLDialect) SupportsReturning() bool {
	return false
}

func (MySQLDialect) GeneratedKeyStrategy() GeneratedKeyStrategy {
	return GeneratedKeyLastInsertID
}

// -----------------------------------------------------------------------------
// Microsoft SQL Server
// -----------------------------------------------------------------------------

type MSSQLDialect struct{}

func (MSSQLDialect) Name() string {
	return "mssql"
}

func (MSSQLDialect) Placeholder(position int) string {
	return fmt.Sprintf(
		"@p%d",
		normalizePosition(position),
	)
}

func (MSSQLDialect) QuoteIdentifier(identifier string) string {
	return quoteQualifiedIdentifier(
		identifier,
		"[",
		"]",
	)
}

func (MSSQLDialect) SupportsReturning() bool {
	return false
}

func (MSSQLDialect) GeneratedKeyStrategy() GeneratedKeyStrategy {
	return GeneratedKeyOutputInserted
}

// -----------------------------------------------------------------------------
// Oracle
// -----------------------------------------------------------------------------

type OracleDialect struct{}

func (OracleDialect) Name() string {
	return "oracle"
}

func (OracleDialect) Placeholder(position int) string {
	return fmt.Sprintf(
		":%d",
		normalizePosition(position),
	)
}

func (OracleDialect) QuoteIdentifier(identifier string) string {
	return quoteQualifiedIdentifier(
		identifier,
		`"`,
		`"`,
	)
}

func (OracleDialect) SupportsReturning() bool {
	return false
}

func (OracleDialect) GeneratedKeyStrategy() GeneratedKeyStrategy {
	return GeneratedKeyReturningInto
}

// -----------------------------------------------------------------------------
// SQLite
// -----------------------------------------------------------------------------

type SQLiteDialect struct {
	name string
}

func NewSQLiteDialect(name string) SQLiteDialect {
	return SQLiteDialect{
		name: name,
	}
}

func (d SQLiteDialect) Name() string {
	if d.name == "" {
		return "sqlite3"
	}

	return d.name
}

func (SQLiteDialect) Placeholder(_ int) string {
	return "?"
}

func (SQLiteDialect) QuoteIdentifier(identifier string) string {
	return quoteQualifiedIdentifier(
		identifier,
		`"`,
		`"`,
	)
}

func (SQLiteDialect) SupportsReturning() bool {
	return true
}

func (SQLiteDialect) GeneratedKeyStrategy() GeneratedKeyStrategy {
	return GeneratedKeyReturning
}

// -----------------------------------------------------------------------------
// Resolution
// -----------------------------------------------------------------------------

func ResolveDialect(driver string) (Dialect, error) {
	name := strings.ToLower(
		strings.TrimSpace(driver),
	)

	switch name {
	// Canonical Gorix wrapper name.
	case "postgres":
		return PostgresDialect{}, nil

	// Backward-compatible native aliases.
	case "pgx", "postgresql":
		return PostgresDialect{}, nil

	case "mysql":
		return MySQLDialect{}, nil

	case "mssql":
		return MSSQLDialect{}, nil

	case "sqlserver":
		return MSSQLDialect{}, nil

	case "oracle":
		return OracleDialect{}, nil

	case "godror":
		return OracleDialect{}, nil

	case "sqlite3":
		return NewSQLiteDialect("sqlite3"), nil

	case "sqlite-modern":
		return NewSQLiteDialect("sqlite-modern"), nil

	case "sqlite":
		return NewSQLiteDialect("sqlite-modern"), nil

	default:
		return nil, fmt.Errorf(
			"gorix repository: unsupported database dialect %q",
			driver,
		)
	}
}

// -----------------------------------------------------------------------------
// Helpers
// -----------------------------------------------------------------------------

func normalizePosition(position int) int {
	if position < 1 {
		return 1
	}

	return position
}

func quoteQualifiedIdentifier(
	identifier string,
	openQuote string,
	closeQuote string,
) string {
	identifier = strings.TrimSpace(identifier)

	if identifier == "" {
		return ""
	}

	parts := strings.Split(identifier, ".")

	for index, part := range parts {
		part = strings.TrimSpace(part)

		// Keep wildcard identifiers unquoted:
		// users.*
		if part == "*" {
			parts[index] = part
			continue
		}

		// Escape the closing quote inside the identifier.
		//
		// PostgreSQL/Oracle/SQLite:
		// user"name -> "user""name"
		//
		// MySQL:
		// user`name -> `user``name`
		//
		// SQL Server:
		// user]name -> [user]]name]
		part = strings.ReplaceAll(
			part,
			closeQuote,
			closeQuote+closeQuote,
		)

		parts[index] =
			openQuote + part + closeQuote
	}

	return strings.Join(parts, ".")
}
