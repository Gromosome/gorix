package interceptor

import (
	"fmt"

	"github.com/Gromosome/gorix/gorix"
)

type AuditInterceptor struct {
}

func NewAuditInterceptor() *AuditInterceptor {
	return &AuditInterceptor{}
}

func (i *AuditInterceptor) Before(ctx *gorix.ExecutionContext) error {

	fmt.Println(fmt.Sprintf("Reached : Controller - %s : Handler - %s ", ctx.Controller, ctx.Handler))

	return nil
}

func (i *AuditInterceptor) After(ctx *gorix.ExecutionContext) error {
	responseType := ctx.Context.ResponseType()
	fmt.Println(responseType)
	if responseType == gorix.ResponseTypeAuto || responseType == gorix.ResponseTypeJSON {
		ctx.Response = map[string]any{
			"success": true,
			"data":    ctx.Response,
		}
	}

	return nil
}
