package main

import (
	"github.com/Gromosome/gorix/gorix/engine"
	"github.com/Gromosome/gorix/gorix/hook"
	middlewares "github.com/Gromosome/gorix/impl-test/middleware"
	"github.com/Gromosome/gorix/impl-test/promotion"
	"github.com/Gromosome/gorix/impl-test/user"
)

func main() {
	app := engine.NewApp()
	app.Use(
		hook.Apply(middlewares.AuthMiddleware()).Only("/user/*"),
	)
	promotionModule := promotion.NewPromotionModule()
	userModule := user.NewUserModule()
	app.RegisterModules(promotionModule, userModule)
	app.Listen()
}
