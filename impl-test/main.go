package main

import (
	"github.com/Gromosome/gorix/gorix/engine"
	"github.com/Gromosome/gorix/gorix/hook"
	"github.com/Gromosome/gorix/impl-test/filter"
	"github.com/Gromosome/gorix/impl-test/interceptor"
	middlewares "github.com/Gromosome/gorix/impl-test/middleware"
	"github.com/Gromosome/gorix/impl-test/promotion"
	"github.com/Gromosome/gorix/impl-test/user"
)

func main() {
	app := engine.NewApp()
	app.Use(
		hook.Apply(middlewares.AuthMiddleware()).Only("/user/*"),
	)
	app.UseInterceptors(
		interceptor.NewAuditInterceptor(),
	)
	app.UseFilters(
		filter.NewExceptionFilter(),
	)
	promotionModule := promotion.NewPromotionModule()
	userModule := user.NewUserModule()
	app.RegisterModules(promotionModule, userModule)
	app.Listen()
}
