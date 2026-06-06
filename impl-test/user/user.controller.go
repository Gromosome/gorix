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
		return c.userService.GetUserList(), nil
	}
}

func (c *UserController) GetUserPromotions() (gorix.Method, gorix.Path, gorix.RouteHandler) {
	return gorix.GET, "/promotions", func(ctx *gorix.Context) (any, error) {
		return c.promotionService.GetPromotionList(), nil
	}
}
func (c *UserController) CreateUser() (gorix.Method, gorix.Path, gorix.RouteHandler) {
	return gorix.POST, "/create", func(ctx *gorix.Context) (any, error) {
		var body CreateUserBodyDto

		if err := ctx.BindBody(&body); err != nil {
			return nil, err
		}

		return c.userService.CreateUser(body), nil
	}
}

func (c *UserController) GetUserByID() (gorix.Method, gorix.Path, gorix.RouteHandler) {
	return gorix.GET, "page/:id", func(ctx *gorix.Context) (any, error) {
		var params UserPathDto

		if err := ctx.BindParams(&params); err != nil {
			return nil, err
		}

		return params, nil
	}
}

func (c *UserController) SearchUsers() (gorix.Method, gorix.Path, gorix.RouteHandler) {
	return gorix.GET, "/search", func(ctx *gorix.Context) (any, error) {
		var query UserSearchQueryDto

		if err := ctx.BindQuery(&query); err != nil {
			return nil, err
		}

		return query, nil
	}
}
