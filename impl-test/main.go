package main

import (
	_ "github.com/Gromosome/gorix/document-drivers/couchdb"
	_ "github.com/Gromosome/gorix/document-drivers/mongo"
	"github.com/Gromosome/gorix/impl-test/post"

	"github.com/Gromosome/gorix/gorix"
	"github.com/Gromosome/gorix/impl-test/global"
	"github.com/Gromosome/gorix/impl-test/promotion"
	"github.com/Gromosome/gorix/impl-test/user"
	_ "github.com/Gromosome/gorix/sql-drivers/postgres"
)

func main() {
	app := gorix.NewApp().APIPrefix("/api")
	app.Use(
		gorix.Apply(global.AuthMiddleware()).Only("/promotion/*"),
	)
	app.UseInterceptors(
		global.NewGlobalInterceptor(),
	)
	app.UseFilters(
		global.NewExceptionFilter(),
	)
	promotionModule := promotion.NewPromotionModule()
	userModule := user.NewUserModule()
	postModule := post.NewPostModule()
	app.RegisterModules(promotionModule, userModule, postModule)
	app.Listen()
}
