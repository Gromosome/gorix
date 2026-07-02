package controller

import (
	"github.com/Gromosome/gorix/gorix"
	"github.com/Gromosome/gorix/impl-test/promotion"
	"github.com/Gromosome/gorix/impl-test/user/dto"
	"github.com/Gromosome/gorix/impl-test/user/service"
)

type UserController struct {
	userService      service.UserServicePort
	promotionService *promotion.PromotionService
}

func NewUserController(userService service.UserServicePort, promotionService *promotion.PromotionService) *UserController {
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
		params gorix.Params[dto.UserPathDto],
	) (any, error) {
		return ctx.
			Status(gorix.StatusOK).
			ResponseEntityJSON(func() (any, error) {
				return c.userService.GetByID(
					ctx,
					params.Value.ID,
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
		query gorix.Query[dto.UserQueryDto],
	) (any, error) {
		return ctx.
			Status(gorix.StatusOK).
			ResponseEntityJSON(func() (any, error) {
				return c.userService.GetActive(
					ctx,
					query.Value,
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
		body gorix.Body[dto.CreateUserDto],
	) (any, error) {
		return ctx.
			Status(gorix.StatusCreated).
			ResponseEntityJSON(func() (any, error) {
				return c.userService.Create(
					ctx,
					body.Value,
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
		params gorix.Params[dto.UserPathDto],
		body gorix.Body[dto.UpdateUserDto],
	) (any, error) {
		return ctx.
			Status(gorix.StatusOK).
			ResponseEntityJSON(func() (any, error) {
				return c.userService.Update(
					ctx,
					params.Value.ID,
					body.Value,
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
		params gorix.Params[dto.UserPathDto],
	) (any, error) {
		return ctx.
			Status(gorix.StatusNoContent).
			ResponseEntityJSON(func() (any, error) {
				return c.userService.Delete(
					ctx,
					params.Value.ID,
				)
			})
	}
}
