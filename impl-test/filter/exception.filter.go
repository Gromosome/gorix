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

	_ = ctx.Context.Status(ctx.StatusCode).JSON(map[string]any{
		"success": false,
		"message": ctx.Error.Error(),
	})
}
