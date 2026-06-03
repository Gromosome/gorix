package interceptor

import "github.com/Gromosome/gorix/gorix"

type AuditInterceptor struct {
}

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
