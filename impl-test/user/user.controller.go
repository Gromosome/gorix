package user

import (
	"github.com/Gromosome/gorix/gorix"
	"github.com/Gromosome/gorix/impl-test/promotion"
)

type UserController struct {
	userService      *UserService
	promotionService *promotion.PromotionService
}

func NewUserController(userService *UserService, promotionService *promotion.PromotionService) *UserController {
	return &UserController{
		userService:      userService,
		promotionService: promotionService,
	}
}

func (c *UserController) GetUserList() (gorix.Method, gorix.Path, gorix.RouteHandler) {
	return gorix.GET, "/find", func(ctx *gorix.Context) (any, error) {
		return c.promotionService.GetPromotionList(), nil
	}
}

func (c *UserController) GetUserPromotions() (gorix.Method, gorix.Path, gorix.RouteHandler) {
	return gorix.GET, "/promotions", func(ctx *gorix.Context) (any, error) {
		return c.promotionService.GetPromotionList(), nil
	}
}
