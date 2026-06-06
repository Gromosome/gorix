package main

import (
	"github.com/Gromosome/gorix/gorix"
	"github.com/Gromosome/gorix/impl-test/filter"
	"github.com/Gromosome/gorix/impl-test/interceptor"
	middlewares "github.com/Gromosome/gorix/impl-test/middleware"
	"github.com/Gromosome/gorix/impl-test/promotion"
	"github.com/Gromosome/gorix/impl-test/user"
)

func main() {
	app := gorix.NewApp()
	app.Use(
		gorix.Apply(middlewares.AuthMiddleware()).Only("/promotion/*"),
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
