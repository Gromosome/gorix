package filter

import "github.com/Gromosome/gorix/gorix"

type ExceptionFilter struct {
}

func NewExceptionFilter() *ExceptionFilter {
	return &ExceptionFilter{}
}

func (f *ExceptionFilter) Catch(ctx *gorix.ExceptionContext) {
	_ = ctx.Context.Status(gorix.StatusBadRequest).JSON(map[string]any{
		"success": false,
		"error":   ctx.Error.Error(),
	})
}
