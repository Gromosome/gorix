package postgres

import (
	"errors"

	sqldriver "github.com/Gromosome/gorix/sql-driver-manager"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type adapter struct{}

func (adapter) Name() string          { return "postgres" }
func (adapter) SQLDriverName() string { return "pgx" }

func (adapter) Normalize(err error) *sqldriver.Error {
	if err == nil {
		return nil
	}

	var target *pgconn.PgError
	if !errors.As(err, &target) {
		return sqldriver.GenericError("postgres", err)
	}

	kinds := map[string]sqldriver.ErrorKind{
		"23505": sqldriver.ErrorUniqueViolation,
		"23503": sqldriver.ErrorForeignKey,
		"23502": sqldriver.ErrorNotNull,
		"23514": sqldriver.ErrorCheckConstraint,
		"40001": sqldriver.ErrorSerialization,
		"40P01": sqldriver.ErrorDeadlock,
		"42601": sqldriver.ErrorSyntax,
		"42501": sqldriver.ErrorPermission,
		"42P01": sqldriver.ErrorTableNotFound,
		"42703": sqldriver.ErrorColumnNotFound,
		"08000": sqldriver.ErrorConnection,
		"08003": sqldriver.ErrorConnection,
		"08006": sqldriver.ErrorConnection,
	}

	kind := kinds[target.Code]
	if kind == "" {
		kind = sqldriver.ErrorUnknown
	}

	return &sqldriver.Error{
		Kind:       kind,
		Driver:     "postgres",
		Code:       target.Code,
		Message:    target.Message,
		Constraint: target.ConstraintName,
		Table:      target.TableName,
		Column:     target.ColumnName,
		Cause:      err,
	}
}

func init() {
	sqldriver.Register(adapter{})
}
