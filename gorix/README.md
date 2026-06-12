<p align="center">
  <img src="https://www.cdn.gromosome.com/main/img/logos/gromo-framework/Gorix.png" alt="Gorix logo" width="180">
</p>

# Gorix User Guide

Gorix applications are structured around modules. A module owns one feature area and usually contains controllers, services, repositories, DTOs, entities, mappers, and any feature-specific hooks.

The `impl-test` module in this workspace is the reference implementation for the examples below.

## Recommended Project Structure

```text
my-api/
  application.yaml
  .env
  main.go
  user/
    user.module.go
    user.controller.go
    user.service.go
    user.repository.go
    user.dto.go
    entity/
      user.entity.go
      user-summary.entity.go
    mapper/
      user.mapper.go
  filter/
    exception.filter.go
  interceptor/
    audit.interceptor.go
  middleware/
    auth.middleware.go
```

Use this structure consistently:

- `main.go`: creates the app, registers global hooks, registers modules, and starts the server.
- `*.module.go`: defines module base path, providers, and controllers.
- `*.controller.go`: defines HTTP routes and binds request data.
- `*.service.go`: contains business logic and coordinates repositories.
- `*.repository.go`: contains persistence logic and transaction boundaries.
- `*.dto.go`: defines body, query, and path DTOs with validation tags.
- `entity/*.entity.go`: defines database-backed structs and table metadata.
- `mapper/*.go`: optional mapping helpers for entity/DTO conversions.
- `filter/*.filter.go`: converts errors into HTTP responses.
- `interceptor/*.interceptor.go`: runs before/after controller execution.
- `middleware/*.middleware.go`: wraps HTTP handlers before routing logic.

## `main.go`

```go
package main

import (
	"github.com/Gromosome/gorix/gorix"
	"github.com/acme/my-api/filter"
	"github.com/acme/my-api/interceptor"
	"github.com/acme/my-api/middleware"
	"github.com/acme/my-api/user"
	_ "github.com/Gromosome/gorix/sql-drivers/postgres"
)

func main() {
	app := gorix.NewApp()

	app.Use(
		gorix.Apply(middleware.AuthMiddleware()).Only("/user/*"),
	)
	app.UseInterceptors(
		interceptor.NewAuditInterceptor(),
	)
	app.UseFilters(
		filter.NewExceptionFilter(),
	)

	app.RegisterModules(
		user.NewUserModule(),
	)
	app.Listen()
}
```

Blank-import one SQL driver module for each configured database driver. For PostgreSQL, use:

```go
_ "github.com/Gromosome/gorix/sql-drivers/postgres"
```

## Module

Each module supplies its base route path and dependency constructors:

```go
package user

import "github.com/Gromosome/gorix/gorix"

type UserModule struct{}

func NewUserModule() *UserModule {
	return &UserModule{}
}

func (m *UserModule) BasePath() gorix.BasePath {
	return "/user"
}

func (m *UserModule) Providers() []any {
	return []any{
		NewUserService,
		NewUserRepository,
	}
}

func (m *UserModule) Controllers() []any {
	return []any{
		NewUserController,
	}
}
```

Gorix uses constructor injection. Register services, repositories, and other injectable providers in `Providers()`. Register controller constructors in `Controllers()`.

## Controller

Controller methods with no arguments are discovered as routes when they return:

```go
gorix.Method, gorix.Path, gorix.RouteHandler
```

Example:

```go
type UserController struct {
	userService *UserService
}

func NewUserController(userService *UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (c *UserController) FindByID() (gorix.Method, gorix.Path, gorix.RouteHandler) {
	return gorix.GET, "/:id", func(ctx *gorix.Context) (any, error) {
		var params UserPathDto
		if err := ctx.BindParams(&params); err != nil {
			return nil, err
		}

		return c.userService.GetByID(ctx, params.ID)
	}
}

func (c *UserController) Create() (gorix.Method, gorix.Path, gorix.RouteHandler) {
	return gorix.POST, "/", func(ctx *gorix.Context) (any, error) {
		var body CreateUserDto
		if err := ctx.BindBody(&body); err != nil {
			return nil, err
		}

		user, err := c.userService.Create(ctx, body)
		if err != nil {
			return nil, err
		}

		ctx.Status(gorix.StatusCreated)
		return user, nil
	}
}
```

Common context methods:

- `BindParams(&dto)`: binds route params using `param` tags.
- `BindQuery(&dto)`: binds query string values using `query` tags.
- `BindBody(&dto)`: binds JSON request bodies using `json` tags.
- `Status(code)`: sets the HTTP response status.
- `JSON(value)`: writes a JSON response directly.

## DTO

DTOs should keep transport concerns out of services and repositories:

```go
type UserPathDto struct {
	ID int64 `param:"id" validate:"required,min=1"`
}

type CreateUserDto struct {
	Name  string `json:"name" validate:"required,min=3,max=100"`
	Email string `json:"email" validate:"required,email"`
}

type UserQueryDto struct {
	Active *bool `query:"active"`
	Limit  int   `query:"limit" validate:"min=1,max=100"`
	Offset int   `query:"offset" validate:"min=0"`
}
```

## Service

Services contain business logic and depend on repositories:

```go
type UserService struct {
	userRepository *UserRepository
}

func NewUserService(userRepository *UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (s *UserService) GetByID(ctx *gorix.Context, id int64) (*User, error) {
	return s.userRepository.FindByID(ctx, id)
}

func (s *UserService) Create(ctx *gorix.Context, request CreateUserDto) (*User, error) {
	user := &User{
		Name:   request.Name,
		Email:  request.Email,
		Active: true,
	}

	if err := s.userRepository.CreateWithAudit(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}
```

## Entity

Entities use `db` tags for column mapping and optional `repository` tags for repository metadata:

```go
package entity

import "time"

type User struct {
	ID        int64     `json:"id" db:"id" repository:"primaryKey,autoIncrement"`
	Name      string    `json:"name" db:"name"`
	Email     string    `json:"email" db:"email"`
	Active    bool      `json:"active" db:"active"`
	CreatedAt time.Time `json:"createdAt" db:"created_at" repository:"readOnly"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at" repository:"readOnly"`
}

func (User) TableName() string {
	return "users"
}
```

## Repository

Repositories should isolate SQL access. Gorix provides:

- `gorix.NewSQLRepository[T, ID](databases)`: generic CRUD-style repository.
- `gorix.NewSQLMapper(databases)`: raw SQL mapper for custom queries.
- `gorix.WithTransaction(ctx, db, options, func(ctx, tx) error)`: transaction helper.

Example:

```go
type User = entity.User

type UserRepository struct {
	databases *gorix.DBManager
	mapper    *gorix.SQLMapper
	repo      *gorix.SQLRepository[entity.User, int64]
}

func NewUserRepository(databases *gorix.DBManager) *UserRepository {
	userRepo, err := gorix.NewSQLRepository[entity.User, int64](databases)
	if err != nil {
		panic(err)
	}

	return &UserRepository{
		databases: databases,
		mapper:    gorix.NewSQLMapper(databases),
		repo:      userRepo,
	}
}

func (r *UserRepository) FindByID(ctx *gorix.Context, id int64) (*User, error) {
	return r.repo.FindByID(ctx, id)
}

func (r *UserRepository) FindActiveUsers(ctx *gorix.Context, limit int, offset int) ([]User, error) {
	users := make([]User, 0)
	err := r.mapper.QueryMany(
		ctx,
		&users,
		`
			SELECT id, name, email, active, created_at, updated_at
			FROM users
			WHERE active = $1
			ORDER BY created_at DESC
			LIMIT $2
			OFFSET $3
		`,
		true,
		limit,
		offset,
	)
	return users, err
}
```

Transaction example:

```go
func (r *UserRepository) CreateWithAudit(ctx *gorix.Context, user *User) error {
	db, err := r.databases.DB()
	if err != nil {
		return err
	}

	return gorix.WithTransaction(ctx, db, nil, func(ctx *gorix.Context, tx *gorix.DBTx) error {
		transactionRepo := r.repo.WithExecutor(tx)
		transactionMapper := r.mapper.WithExecutor(tx)

		if err := transactionRepo.Insert(ctx, user); err != nil {
			return err
		}

		result := transactionMapper.Exec(
			ctx,
			`INSERT INTO user_audit (user_id, action) VALUES ($1, $2)`,
			user.ID,
			"USER_CREATED",
		)
		return result.Err()
	})
}
```

## Mapper

Use mapper files when transformations become non-trivial:

```go
package mapper

import (
	"github.com/acme/my-api/user"
	"github.com/acme/my-api/user/entity"
)

func NewUserFromCreateDTO(dto user.CreateUserDto) *entity.User {
	return &entity.User{
		Name:   dto.Name,
		Email:  dto.Email,
		Active: true,
	}
}
```

Keep mapping deterministic and side-effect free. Put validation in DTO tags or services, not mappers.

## Middleware

Middleware wraps handlers and can short-circuit the request:

```go
func AuthMiddleware() gorix.Middleware {
	return func(next gorix.Handler) gorix.Handler {
		return func(ctx *gorix.Context) error {
			token := ctx.R.Header.Get("Authorization")
			if token != "Bearer dev-token" {
				return ctx.Status(gorix.StatusUnauthorized).JSON(map[string]any{
					"success": false,
					"error":   "unauthorized",
				})
			}
			return next(ctx)
		}
	}
}
```

Register globally or with route rules:

```go
app.Use(gorix.Apply(AuthMiddleware()).Only("/user/*"))
```

## Interceptor

Interceptors run around controller execution. Use them for response envelopes, audit logging, timing, or route-level metadata:

```go
type AuditInterceptor struct{}

func NewAuditInterceptor() *AuditInterceptor {
	return &AuditInterceptor{}
}

func (i *AuditInterceptor) Before(ctx *gorix.ExecutionContext) error {
	return nil
}

func (i *AuditInterceptor) After(ctx *gorix.ExecutionContext) error {
	ctx.Response = map[string]any{
		"success": true,
		"data":    ctx.Response,
	}
	return nil
}
```

Register with:

```go
app.UseInterceptors(NewAuditInterceptor())
```

## Filter

Filters convert errors into HTTP responses:

```go
type ExceptionFilter struct{}

func NewExceptionFilter() *ExceptionFilter {
	return &ExceptionFilter{}
}

func (f *ExceptionFilter) Catch(ctx *gorix.ExceptionContext) {
	if validationErr, ok := errors.AsType[*gorix.ValidationError](ctx.Error); ok {
		_ = ctx.Context.Status(gorix.StatusBadRequest).JSON(map[string]any{
			"success": false,
			"message": validationErr.Error(),
			"errors":  validationErr.ErrorAsList(),
		})
		return
	}

	if dbError, ok := gorix.AsDatabaseError(ctx.Error); ok {
		_ = ctx.Context.Status(gorix.StatusInternalServerError).JSON(map[string]any{
			"success": false,
			"error":   dbError.Kind,
		})
		return
	}

	_ = ctx.Context.Status(ctx.StatusCode).JSON(map[string]any{
		"success": false,
		"message": ctx.Error.Error(),
	})
}
```

Register with:

```go
app.UseFilters(NewExceptionFilter())
```

## Context

`gorix.Context` wraps request and response behavior. It is passed to every route handler and should also be passed through service and repository calls so cancellation, request-scoped data, and transactions remain consistent.

Recommended usage:

- Bind request data in controllers.
- Pass `ctx` into services and repositories.
- Set HTTP status in controllers or filters.
- Write direct responses only when a hook or middleware must short-circuit.

## YAML and Environment Configuration

Gorix loads `application.yaml` from the working directory used to start the app.

Example:

```yaml
env: local
gorix:
  app:
    host: ${GORIX_HOST:-0.0.0.0}
    port: ${GORIX_PORT:-8070}
    prod: ${GORIX_PROD:-false}
  databases:
    default:
      driver: postgres
      dsn: ${DATABASE_DSN}
      max-open-connections: ${DB_MAX_OPEN:-25}
      max-idle-connections: ${DB_MAX_IDLE:-10}
      connection-max-lifetime: ${DB_CONN_MAX_LIFETIME:-30m}
      connection-max-idle-time: ${DB_CONN_MAX_IDLE_TIME:-5m}
      ping-timeout: ${DB_PING_TIMEOUT:-5s}
```

Environment file behavior:

- `env: ""` loads `.env` if it exists; missing `.env` is allowed.
- `env: local` loads `local.env`; missing `local.env` is an error.
- Environment names may contain letters, numbers, `_`, and `-`.
- Existing OS, Docker, Kubernetes, and CI variables are not overwritten by dotenv files.
- Dotenv values can reference earlier variables in the same file.
- Unquoted values may contain inline comments only when `#` is preceded by whitespace.
- Double-quoted values use Go string unquoting; single-quoted values are taken literally.

Example `.env`:

```dotenv
GORIX_PORT=8070
DB_HOST=localhost
DB_PORT=5432
DB_NAME=root_db
DATABASE_DSN=postgres://root:root@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable
```

Application defaults:

- `gorix.app.host`: `0.0.0.0`
- `gorix.app.port`: `8080`
- `gorix.app.prod`: `true` when omitted

When `prod` is `false`, Gorix runs project validation before starting the server.
