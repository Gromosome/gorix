package user

import (
	"github.com/Gromosome/gorix/gorix"
	"github.com/Gromosome/gorix/impl-test/promotion"
)

type UserController struct {
	userService      UserService
	promotionService promotion.PromotionService
}

func NewUserController(userService UserService, promotionService promotion.PromotionService) UserController {
	return UserController{
		userService:      userService,
		promotionService: promotionService,
	}
}

func (c *UserController) GetUserList() (gorix.Method, gorix.Path, []string) {
	return gorix.GET, "/find", c.userService.GetUserList()
}

func (c *UserController) GetUserPromotions() (gorix.Method, gorix.Path, []string) {
	return gorix.GET, "/promotions", c.promotionService.GetPromotionList()
}
