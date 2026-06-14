# Contributing to Gorix

Thank you for helping improve Gorix. This project welcomes bug reports, documentation fixes, tests, driver adapters, and framework improvements.

## Ways to Contribute

- Report bugs and security concerns responsibly.
- Improve documentation, examples, and setup instructions.
- Add or improve tests for framework behavior and SQL drivers.
- Propose focused enhancements before investing in large changes.
- Support the project through Open Collective: <https://opencollective.com/gorix>.

## Before You Start

For non-trivial changes, open an issue first and describe:

- The problem or use case.
- The proposed behavior or API change.
- Any compatibility or migration impact.

Small documentation fixes, typo fixes, and narrowly scoped bug fixes can be submitted directly.

## Development Setup

This repository is a Go workspace with multiple modules.

```bash
go work sync
go test ./...
```

You can also test individual modules:

```bash
cd gorix
go test ./...

cd ../sql-driver-manager
go test ./...
```

Some SQL driver tests may require local database services or environment variables. If a test requires external infrastructure, document the requirement in the pull request.

## Repository Structure

Gorix is organized as a Go workspace. Keep changes in the module that owns the behavior:

- `gorix`: framework public API and core framework implementation.
- `gorix/app`: application startup, route registration, and framework validation.
- `gorix/app/linter`: project-structure validation rules for Gorix applications.
- `gorix/config`: application and database configuration types.
- `gorix/config/yaml`: YAML parsing, normalization, and environment expansion.
- `gorix/core/context`: request context, binding, validation, and response helpers.
- `gorix/core/database`: framework database abstraction and SQL helpers.
- `gorix/core/database/mapper`: row scanning and mapper registration.
- `gorix/core/database/repository`: repository metadata, dialects, and query builders.
- `gorix/di`: dependency injection container.
- `gorix/hook`: middleware, interceptor, filter, and route rule types.
- `sql-driver-manager`: SQL driver registry, manager, normalized errors, and wrappers.
- `sql-drivers/*`: individual database driver adapters.
- `impl-test`: reference application for framework usage.

## File Naming Standards

Follow the existing package and file naming pattern. Framework files use lowercase dot-separated names grouped by area:

- Public re-export files live at the module root, for example `app.go`, `database.error.go`, and `handler.types.go`.
- Internal framework implementation files use `<area>.<topic>.go`, for example `context.request.go`, `database.manager.go`, and `repository.query-builder.go`.
- Tests mirror the implementation name with `_test.go`, for example `context.request_test.go` and `repository.query-builder_test.go`.
- Keep package directories focused; add new framework code beside the related package instead of creating broad utility packages.
- Use `gofmt` for all Go files before submitting a pull request.

## Gorix Application Structure

When contributing examples, reference applications, or linter behavior, follow the Gorix application structure enforced by the framework:

```text
my-api/
  application.yaml
  main.go
  user/
    user.module.go
    user.controller.go
    user.service.go
    user.repository.go
    user.dto.go
    entity/
      user.entity.go
    mapper/
      user.mapper.go
  filter/
    exception.filter.go
  interceptor/
    audit.interceptor.go
  middleware/
    auth.middleware.go
```

Use these layer rules consistently:

- Root application directory allows `main.go` as the only Go file.
- Feature package names should match their directory names.
- Module files must be named `<package>.module.go`.
- Controller files should use `.controller.go`.
- Service files should use `.service.go`.
- Repository files should use `.repository.go`.
- DTO files should use `.dto.go`.
- Middleware files should use `.middleware.go`.
- Interceptor files should use `.interceptor.go`.
- Filter files should use `.filter.go`.

## Framework Layer Contracts

The framework linter validates application structure. Contributions must preserve these contracts unless the change intentionally updates the linter and documentation together:

- Module files contain exactly one module struct, exactly one constructor returning `*StructName`, `BasePath() gorix.BasePath`, and `Controllers() []any`; `Providers() []any` is allowed.
- Controller files contain exactly one controller struct, exactly one constructor returning the struct or pointer, and route methods with no parameters returning `gorix.Method`, `gorix.Path`, and `gorix.RouteHandler`.
- Service files contain exactly one service struct and exactly one constructor named `NewStructName` returning `*StructName`.
- DTO files contain at least one DTO struct; DTO struct names end with `Dto` or `DTO`, fields are exported, and fields include `json`, `query`, or `param` tags.
- Middleware files contain no receiver methods and at least one function returning `gorix.Middleware`.
- Interceptor files contain exactly one struct with `Before(ctx *gorix.ExecutionContext) error` and `After(ctx *gorix.ExecutionContext) error`.
- Filter files contain exactly one struct with `Catch(ctx *gorix.ExceptionContext)` and no return value.

## Test Standards

Tests in this repository intentionally live under `test/` directories that mirror the implementation package tree.

- Framework tests go under `gorix/test/...` using the same subdirectory path as the implementation, for example `gorix/core/context` is tested in `gorix/test/core/context`.
- SQL driver manager tests go under `sql-driver-manager/test`.
- Individual SQL driver adapter tests live beside the adapter module, for example `sql-drivers/postgres/postgres_test.go`.
- Test file names should mirror the implementation file name with `_test.go`.
- Test packages should match the tested area name when using mirrored framework tests, for example `package context`, `package repository`, and `package linter`.
- Use Go's standard `testing` package and direct assertions with `t.Fatal` or `t.Fatalf`.
- Name tests as `Test<Behavior>` and make the tested behavior explicit, for example `TestQueryBuilderBuildSelect`.
- Prefer deterministic unit tests. Use `t.TempDir()` for filesystem tests and avoid shared mutable state unless a package-level fake driver or adapter is required.
- Keep external service requirements out of default unit tests where practical. If a database service or environment variable is required, document it in the test and pull request.

## Pull Request Guidelines

- Keep pull requests focused on one logical change.
- Follow the existing Go style and run `gofmt` on changed Go files.
- Add or update tests when behavior changes.
- Update relevant README files when public APIs, configuration, or behavior changes.
- Do not include secrets, credentials, private keys, or production data.
- Explain verification steps in the pull request description.

## Reporting Issues

Public bug reports and feature requests can be opened in GitHub Issues. If the report includes private information, credentials, vulnerability details, or sensitive operational data, email the project directly instead:

**pressnavek@gromosome.com**

## Funding

Gorix is open source. Organizations and individuals can support maintenance through Open Collective:

<https://opencollective.com/gorix>
