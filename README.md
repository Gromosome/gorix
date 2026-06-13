<p align="center">
  <img src="https://www.cdn.gromosome.com/main/img/logos/gromo-framework/Gorix.png" alt="Gorix logo" width="180">
</p>

<p align="center">
  <a href="https://paypal.me/PressnaveKiruparaj?locale.x=en_US&amp;country.x=AE"><img src="https://img.shields.io/badge/PayPal-Support-blue?logo=paypal" alt="Support on PayPal"></a>
  <a href="https://opencollective.com/gorix"><img src="https://img.shields.io/badge/OpenCollective-Gorix-7FADF2?logo=opencollective" alt="Support Gorix on Open Collective"></a>
  <a href="https://pkg.go.dev/github.com/Gromosome/gorix/gorix"><img src="https://pkg.go.dev/badge/github.com/Gromosome/gorix/gorix.svg" alt="Go Reference"></a>
  <a href="https://github.com/Gromosome/gorix/actions"><img src="https://img.shields.io/badge/CI%2FCD-GitHub%20Actions-2088FF?logo=githubactions" alt="GitHub Actions CI/CD"></a>
</p>

# Gorix Backend Framework

Gorix is a Go backend framework for building modular HTTP APIs with dependency injection, controller-based routing, request/response context helpers, middleware, interceptors, exception filters, YAML configuration, environment expansion, and SQL database access.

This repository is organized as a Go workspace with separate modules:

- `gorix`: the framework module and main public API.
- `impl-test`: a runnable reference application that demonstrates project structure and framework usage.
- `sql-driver-manager`: a normalized SQL driver registry and database wrapper.
- `sql-drivers`: driver adapter modules for PostgreSQL, MySQL, SQLite, Microsoft SQL Server, and Oracle.

## Core Concepts

- **Application**: `gorix.NewApp()` creates the HTTP application, loads `application.yaml`, validates the project in non-production mode, connects databases, registers modules, and starts the server.
- **Modules**: each feature module exposes `BasePath()`, `Providers()`, and `Controllers()` so Gorix can register dependencies and routes.
- **Controllers**: controller methods return `gorix.Method`, `gorix.Path`, and `gorix.RouteHandler`.
- **Services**: services contain business logic and receive dependencies through constructor injection.
- **Repositories**: repositories isolate persistence using `gorix.SQLRepository`, `gorix.SQLMapper`, and transactions.
- **DTOs**: DTO structs bind and validate route params, query params, and request bodies.
- **Hooks**: middleware, interceptors, and filters handle cross-cutting request behavior.
- **Configuration**: `application.yaml` supports `${ENV}` and `${ENV:-fallback}` placeholders, plus optional `.env` files.

## Quick Start

Use the `impl-test` module as the reference application:

```bash
cd impl-test
go mod tidy
go run .
```

The example app:

- Registers auth middleware for `/promotion/*`.
- Registers an audit interceptor that wraps successful responses.
- Registers an exception filter for validation and database errors.
- Registers `promotion` and `user` modules.
- Uses PostgreSQL through a blank import of `github.com/Gromosome/gorix/sql-drivers/postgres`.

## Minimal Application

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

## Configuration

Create `application.yaml` at the application root:

```yaml
env: ""
gorix:
  app:
    host: 0.0.0.0
    port: ${GORIX_PORT:-8070}
    prod: false
  databases:
    default:
      driver: postgres
      dsn: ${DATABASE_DSN:-postgres://root:root@localhost:5432/root_db?sslmode=disable}
      max-open-connections: 25
      max-idle-connections: 10
      connection-max-lifetime: 30m
      connection-max-idle-time: 5m
```

If `env` is empty, Gorix optionally loads `.env`. If `env: local`, Gorix requires `local.env`.

## Documentation

- Framework user guide: `gorix/README.md`
- SQL driver manager: `sql-driver-manager/README.md`
- SQL driver adapters: `sql-drivers/README.md`
- Reference implementation: `impl-test`
