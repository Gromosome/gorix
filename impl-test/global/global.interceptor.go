package global

import (
	"fmt"

	"github.com/Gromosome/gorix/gorix"
)

type GlobalInterceptor struct {
}

func NewGlobalInterceptor() *GlobalInterceptor {
	return &GlobalInterceptor{}
}

func (i *GlobalInterceptor) Before(ctx *gorix.ExecutionContext) error {
	fmt.Println(fmt.Sprintf("Reached : Controller - %s : Handler - %s ", ctx.Controller, ctx.Handler))
	return nil
}

func (i *GlobalInterceptor) After(ctx *gorix.ExecutionContext) error {
	responseType := ctx.Context.ResponseType()
	if responseType == gorix.ResponseTypeJSON {
		ctx.Response = ResponseDTO[any]{
			Success: true,
			Data:    ctx.Response,
		}
	}
	return nil
}
