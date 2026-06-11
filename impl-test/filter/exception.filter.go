package filter

import (
	"errors"

	"github.com/Gromosome/gorix/gorix"
)

type ExceptionFilter struct {
}

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
		switch dbError.Kind {
		case gorix.DatabaseErrorDuplicateKey,
			gorix.DatabaseErrorUniqueViolation:

			_ = ctx.Context.
				Status(gorix.StatusConflict).
				JSON(map[string]any{
					"success":    false,
					"error":      "resource already exists",
					"constraint": dbError.Constraint,
				})

		case gorix.DatabaseErrorForeignKey,
			gorix.DatabaseErrorNotNull,
			gorix.DatabaseErrorCheckConstraint:

			_ = ctx.Context.
				Status(gorix.StatusBadRequest).
				JSON(map[string]any{
					"success": false,
					"error":   "invalid database operation",
				})

		case gorix.DatabaseErrorTimeout:
			_ = ctx.Context.
				Status(gorix.StatusGatewayTimeout).
				JSON(map[string]any{
					"success": false,
					"error":   "database operation timed out",
				})

		case gorix.DatabaseErrorConnection:
			_ = ctx.Context.
				Status(gorix.StatusServiceUnavailable).
				JSON(map[string]any{
					"success": false,
					"error":   "database unavailable",
				})

		default:
			_ = ctx.Context.
				Status(gorix.StatusInternalServerError).
				JSON(map[string]any{
					"success": false,
					"error":   "database operation failed",
				})
		}
		return
	}

	_ = ctx.Context.Status(ctx.StatusCode).JSON(map[string]any{
		"success": false,
		"message": ctx.Error.Error(),
	})
}
