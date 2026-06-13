<p align="center">
  <img src="https://www.cdn.gromosome.com/main/img/logos/gromo-framework/Gorix.png" alt="Gorix logo" width="180">
</p>

<p align="center">
  <a href="https://paypal.me/PressnaveKiruparaj?locale.x=en_US&amp;country.x=AE"><img src="https://img.shields.io/badge/PayPal-Support-blue?logo=paypal" alt="Support on PayPal"></a>
  <a href="https://opencollective.com/gorix"><img src="https://img.shields.io/badge/OpenCollective-Gorix-7FADF2?logo=opencollective" alt="Support Gorix on Open Collective"></a>
  <a href="https://pkg.go.dev/github.com/Gromosome/gorix/gorix"><img src="https://pkg.go.dev/badge/github.com/Gromosome/gorix/gorix.svg" alt="Go Reference"></a>
  <a href="https://github.com/Gromosome/gorix/actions"><img src="https://img.shields.io/badge/CI%2FCD-GitHub%20Actions-2088FF?logo=githubactions" alt="GitHub Actions CI/CD"></a>
</p>

# Gorix SQL Driver Adapters

`sql-drivers` contains adapter modules that register database drivers with `github.com/Gromosome/gorix/sql-driver-manager`.

Import adapters for side effects in your application:

```go
import (
	_ "github.com/Gromosome/gorix/sql-drivers/postgres"
)
```

Then configure the matching logical driver name in `application.yaml`.

## Available Drivers

| Adapter module | Logical driver | SQL driver name | Mapper module | Mapper version | Underlying driver | Driver version |
| --- | --- | --- | --- | --- | --- | --- |
| `github.com/Gromosome/gorix/sql-drivers/postgres` | `postgres` | `pgx` | `github.com/Gromosome/gorix/sql-driver-manager` | `v1.0.0` | `github.com/jackc/pgx/v5` | `v5.7.6` |
| `github.com/Gromosome/gorix/sql-drivers/mysql` | `mysql` | `mysql` | `github.com/Gromosome/gorix/sql-driver-manager` | `v1.0.0` | `github.com/go-sql-driver/mysql` | `v1.9.3` |
| `github.com/Gromosome/gorix/sql-drivers/sqlite3` | `sqlite3` | `sqlite3` | `github.com/Gromosome/gorix/sql-driver-manager` | `v1.0.0` | `github.com/mattn/go-sqlite3` | `v1.14.32` |
| `github.com/Gromosome/gorix/sql-drivers/sqlite-modern` | `sqlite-modern` | `sqlite` | `github.com/Gromosome/gorix/sql-driver-manager` | `v1.0.0` | `modernc.org/sqlite` | `v1.39.1` |
| `github.com/Gromosome/gorix/sql-drivers/mssql` | `mssql` | `sqlserver` | `github.com/Gromosome/gorix/sql-driver-manager` | `v1.0.0` | `github.com/microsoft/go-mssqldb` | `v1.9.3` |
| `github.com/Gromosome/gorix/sql-drivers/oracle` | `oracle` | `godror` | `github.com/Gromosome/gorix/sql-driver-manager` | `v1.0.0` | `github.com/godror/godror` | `v0.49.4` |

## Configuration Examples

PostgreSQL:

```yaml
gorix:
  databases:
    default:
      driver: postgres
      dsn: postgres://user:password@localhost:5432/app?sslmode=disable
```

MySQL:

```yaml
gorix:
  databases:
    default:
      driver: mysql
      dsn: user:password@tcp(localhost:3306)/app?parseTime=true
```

SQLite using `mattn/go-sqlite3`:

```yaml
gorix:
  databases:
    default:
      driver: sqlite3
      dsn: file:app.db?_foreign_keys=on
```

SQLite using `modernc.org/sqlite`:

```yaml
gorix:
  databases:
    default:
      driver: sqlite-modern
      dsn: file:app.db?_pragma=foreign_keys(1)
```

Microsoft SQL Server:

```yaml
gorix:
  databases:
    default:
      driver: mssql
      dsn: sqlserver://user:password@localhost:1433?database=app
```

Oracle:

```yaml
gorix:
  databases:
    default:
      driver: oracle
      dsn: user/password@localhost:1521/FREEPDB1
```

## Error Mapping

Each adapter maps native driver errors to normalized `sql-driver-manager` kinds such as:

- unique or duplicate key violations
- foreign key violations
- not-null and check constraint violations
- serialization failures and deadlocks
- connection errors and timeouts
- syntax, permission, table, and column errors

Gorix consumes these normalized errors so application filters can return consistent HTTP responses across database engines.

## Usage with Gorix

```go
package main

import (
	"github.com/Gromosome/gorix/gorix"
	_ "github.com/Gromosome/gorix/sql-drivers/postgres"
)

func main() {
	app := gorix.NewApp()
	app.RegisterModules(NewUserModule())
	app.Listen()
}
```

The logical `driver` in `application.yaml` must match the imported adapter. If it does not, startup fails with a message telling you to blank-import the wrapper module.
