<p align="center">
  <img src="https://www.cdn.gromosome.com/main/img/logos/gromo-framework/Gorix.png" alt="Gorix logo" width="180">
</p>

<p align="center">
  <a href="https://paypal.me/PressnaveKiruparaj?locale.x=en_US&amp;country.x=AE"><img src="https://img.shields.io/badge/PayPal-Support-blue?logo=paypal" alt="Support on PayPal"></a>
  <a href="https://opencollective.com/gorix"><img src="https://img.shields.io/badge/OpenCollective-Gorix-7FADF2?logo=opencollective" alt="Support Gorix on Open Collective"></a>
  <a href="https://pkg.go.dev/github.com/Gromosome/gorix/gorix"><img src="https://pkg.go.dev/badge/github.com/Gromosome/gorix/gorix.svg" alt="Go Reference"></a>
  <a href="https://github.com/Gromosome/gorix/actions"><img src="https://img.shields.io/badge/CI%2FCD-GitHub%20Actions-2088FF?logo=githubactions" alt="GitHub Actions CI/CD"></a>
</p>

# Gorix SQL Driver Manager

`sql-driver-manager` provides a small registry and wrapper around Go SQL drivers. It lets Gorix open configured databases by logical driver name and normalize driver-specific errors into a common `Error` type.

## Install

```bash
go get github.com/Gromosome/gorix/sql-driver-manager
```

Import one or more driver adapter modules in the application so they register themselves:

```go
import (
	sqldriver "github.com/Gromosome/gorix/sql-driver-manager"
	_ "github.com/Gromosome/gorix/sql-drivers/postgres"
)
```

## Open a Database

```go
manager, err := sqldriver.Open(ctx, sqldriver.Config{
	Driver:          "postgres",
	DSN:             "postgres://root:root@localhost:5432/root_db?sslmode=disable",
	MaxOpenConns:    25,
	MaxIdleConns:    10,
	ConnMaxLifetime: 30 * time.Minute,
	ConnMaxIdleTime: 5 * time.Minute,
	PingTimeout:     5 * time.Second,
})
if err != nil {
	return err
}
defer manager.Close()
```

## Driver Registry

Adapters implement:

```go
type Adapter interface {
	Name() string
	SQLDriverName() string
	Normalize(error) *Error
}
```

Use:

- `Register(adapter)`: registers an adapter, usually from an adapter module `init()`.
- `Lookup(name)`: returns a registered adapter by logical name.
- `RegisteredDrivers()`: lists registered logical driver names.

If a driver is not registered, blank-import its adapter module.

## Error Normalization

Normalized errors use:

```go
type Error struct {
	Kind       ErrorKind
	Driver     string
	Code       string
	Message    string
	Constraint string
	Table      string
	Column     string
	Cause      error
}
```

Common kinds include:

- `duplicate_key`
- `unique_violation`
- `foreign_key_violation`
- `not_null_violation`
- `check_constraint_violation`
- `serialization_failure`
- `deadlock`
- `connection_error`
- `timeout`
- `syntax_error`
- `permission_denied`
- `table_not_found`
- `column_not_found`

Use helpers:

```go
if dbErr, ok := sqldriver.AsError(err); ok {
	switch dbErr.Kind {
	case sqldriver.ErrorUniqueViolation:
		// return HTTP 409 or domain conflict
	}
}

if sqldriver.IsKind(err, sqldriver.ErrorTimeout) {
	// retry or return timeout response
}
```

## Query API

`Manager` wraps common `database/sql` operations and normalizes returned errors:

- `ExecContext`
- `QueryContext`
- `QueryRowContext`
- `PrepareContext`
- `BeginTx`
- `WithTx`
- `PingContext`
- `Stats`
- `Close`

Rows, row scans, statements, transactions, commits, rollbacks, and close operations are normalized through the same adapter.
