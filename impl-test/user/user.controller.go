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

func (c *UserController) GetUserPromotions() (gorix.Method, gorix.Path, gorix.RouteHandler) {
	return gorix.GET, "/promotions", func(ctx *gorix.Context) (any, error) {
		return c.promotionService.GetPromotionList(), nil
	}
}

func (c *UserController) FindByID() (
	gorix.Method,
	gorix.Path,
	gorix.RouteHandler,
) {
	return gorix.GET, "/:id", func(
		ctx *gorix.Context,
	) (any, error) {
		var params UserPathDto
		return ctx.
			Status(gorix.StatusOK).
			BindParams(&params).
			ResponseEntityJSON(func() (any, error) {
				return c.userService.GetByID(
					ctx,
					params.ID,
				)
			})
	}
}

func (c *UserController) FindAll() (
	gorix.Method,
	gorix.Path,
	gorix.RouteHandler,
) {
	return gorix.GET, "/", func(
		ctx *gorix.Context,
	) (any, error) {

		return ctx.Status(gorix.StatusOK).ResponseEntityXML(func() (any, error) {
			return c.userService.GetAll(ctx)
		})
	}
}

func (c *UserController) FindActive() (
	gorix.Method,
	gorix.Path,
	gorix.RouteHandler,
) {
	return gorix.GET, "/active", func(
		ctx *gorix.Context,
	) (any, error) {
		query := UserQueryDto{
			Limit: 20,
		}
		return ctx.
			Status(gorix.StatusOK).
			BindQuery(&query).
			ResponseEntityJSON(func() (any, error) {
				return c.userService.GetActive(
					ctx,
					query,
				)
			})
	}
}

func (c *UserController) Summary() (
	gorix.Method,
	gorix.Path,
	gorix.RouteHandler,
) {
	return gorix.GET, "/summary", func(
		ctx *gorix.Context,
	) (any, error) {
		return c.userService.GetSummary(ctx)
	}
}

func (c *UserController) Create() (
	gorix.Method,
	gorix.Path,
	gorix.RouteHandler,
) {
	return gorix.POST, "/", func(
		ctx *gorix.Context,
	) (any, error) {
		var body CreateUserDto
		return ctx.
			BindBody(&body).
			Status(gorix.StatusCreated).
			ResponseEntityJSON(func() (any, error) {
				return c.userService.Create(
					ctx,
					body,
				)
			})
	}
}

func (c *UserController) Update() (
	gorix.Method,
	gorix.Path,
	gorix.RouteHandler,
) {
	return gorix.PUT, "/:id", func(
		ctx *gorix.Context,
	) (any, error) {
		var params UserPathDto
		var body UpdateUserDto
		return ctx.
			Status(gorix.StatusOK).
			BindBody(&body).
			BindParams(&params).
			ResponseEntityJSON(func() (any, error) {
				return c.userService.Update(
					ctx,
					params.ID,
					body,
				)
			})
	}
}

func (c *UserController) Delete() (
	gorix.Method,
	gorix.Path,
	gorix.RouteHandler,
) {
	return gorix.DELETE, "/:id", func(
		ctx *gorix.Context,
	) (any, error) {
		var params UserPathDto
		return ctx.
			Status(gorix.StatusNoContent).
			BindParams(&params).
			ResponseEntityJSON(func() (any, error) {
				return c.userService.Delete(
					ctx,
					params.ID,
				)
			})
	}
}
