<p align="center">
  <img src="https://www.cdn.gromosome.com/main/img/logos/gromo-framework/Gorix.png" alt="Gorix logo" width="180">
</p>

<p align="center">
  <a href="https://paypal.me/PressnaveKiruparaj?locale.x=en_US&amp;country.x=AE"><img src="https://img.shields.io/badge/PayPal-Support-blue?logo=paypal" alt="Support on PayPal"></a>
  <a href="https://opencollective.com/gorix"><img src="https://img.shields.io/badge/OpenCollective-Gorix-7FADF2?logo=opencollective" alt="Support Gorix on Open Collective"></a>
  <a href="https://pkg.go.dev/github.com/Gromosome/gorix/gorix"><img src="https://pkg.go.dev/badge/github.com/Gromosome/gorix/gorix.svg" alt="Go Reference"></a>
  <a href="https://github.com/Gromosome/gorix/actions"><img src="https://img.shields.io/badge/CI%2FCD-GitHub%20Actions-2088FF?logo=githubactions" alt="GitHub Actions CI/CD"></a>
</p>

# Gorix PostgreSQL Driver Adapter

Registers the PostgreSQL adapter with `github.com/Gromosome/gorix/sql-driver-manager`.

## Driver Map

| Adapter module | Logical driver | SQL driver name | Mapper module | Mapper version | Underlying driver | Driver version |
| --- | --- | --- | --- | --- | --- | --- |
| `github.com/Gromosome/gorix/sql-drivers/postgres` | `postgres` | `pgx` | `github.com/Gromosome/gorix/sql-driver-manager` | `v1.0.0` | `github.com/jackc/pgx/v5` | `v5.7.6` |

## Usage

```go
import (
	_ "github.com/Gromosome/gorix/sql-drivers/postgres"
)
```

```yaml
gorix:
  databases:
    default:
      driver: postgres
      dsn: postgres://user:password@localhost:5432/app?sslmode=disable
```
