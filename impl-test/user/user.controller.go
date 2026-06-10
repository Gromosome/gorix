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

		if err := ctx.BindParams(&params); err != nil {
			return nil, err
		}

		return c.userService.GetByID(
			ctx,
			params.ID,
		)
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
		return c.userService.GetAll(ctx)
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

		if err := ctx.BindQuery(&query); err != nil {
			return nil, err
		}

		return c.userService.GetActive(
			ctx,
			query,
		)
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

		if err := ctx.BindBody(&body); err != nil {
			return nil, err
		}

		user, err := c.userService.Create(
			ctx,
			body,
		)
		if err != nil {
			return nil, err
		}

		ctx.Status(gorix.StatusCreated)

		return user, nil
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

		if err := ctx.BindParams(&params); err != nil {
			return nil, err
		}

		var body UpdateUserDto

		if err := ctx.BindBody(&body); err != nil {
			return nil, err
		}

		return c.userService.Update(
			ctx,
			params.ID,
			body,
		)
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

		if err := ctx.BindParams(&params); err != nil {
			return nil, err
		}

		if err := c.userService.Delete(
			ctx,
			params.ID,
		); err != nil {
			return nil, err
		}

		ctx.Status(gorix.StatusNoContent)

		return nil, nil
	}
}
