package global

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
	if ctx.Context.ResponseType() == gorix.ResponseTypeJSON {
		if validationErr, ok := errors.AsType[*gorix.ValidationError](ctx.Error); ok {
			_ = ctx.Context.Status(gorix.StatusBadRequest).JSONFault(ErrorDTO{
				Success: false,
				Error:   validationErr.ErrorAsList(),
				Message: validationErr.Error(),
			})
			return
		}
		if dbError, ok := gorix.AsDatabaseError(ctx.Error); ok {
			switch dbError.Kind {
			case gorix.DatabaseErrorDuplicateKey,
				gorix.DatabaseErrorUniqueViolation:

				_ = ctx.Context.
					Status(gorix.StatusConflict).
					JSONFault(
						ErrorDTO{
							Success: false,
							Error:   dbError.Constraint,
							Message: "resource already exists",
						})

			case gorix.DatabaseErrorForeignKey,
				gorix.DatabaseErrorNotNull,
				gorix.DatabaseErrorCheckConstraint:

				_ = ctx.Context.
					Status(gorix.StatusBadRequest).
					JSONFault(
						ErrorDTO{
							Success: false,
							Error:   dbError.Constraint,
							Message: "foreign key constraint failed / Database ErrorNotNull / Check Constraint",
						})

			case gorix.DatabaseErrorTimeout:
				_ = ctx.Context.
					Status(gorix.StatusGatewayTimeout).
					JSONFault(ErrorDTO{
						Success: false,
						Error:   "database operation timed out",
						Message: "database operation timed out",
					})

			case gorix.DatabaseErrorConnection:
				_ = ctx.Context.
					Status(gorix.StatusServiceUnavailable).
					JSONFault(ErrorDTO{
						Success: false,
						Error:   "database unavailable",
						Message: "database unavailable",
					})

			default:
				_ = ctx.Context.
					Status(gorix.StatusInternalServerError).
					JSONFault(ErrorDTO{
						Success: false,
						Error:   "database operation failed",
						Message: "database operation failed",
					})
			}
			return
		}

		_ = ctx.Context.Status(ctx.StatusCode).JSONFault(ErrorDTO{
			Success: false,
			Error:   ctx.Error.Error(),
			Message: ctx.Error.Error(),
		})

	}

}
