package post

import "github.com/Gromosome/gorix/gorix"

type PostController struct {
	postService *PostService
}

func NewPostController(
	postService *PostService,
) *PostController {
	return &PostController{
		postService: postService,
	}
}

func (c *PostController) EnsureSchema() (
	gorix.Method,
	gorix.Path,
	gorix.RouteHandler,
) {
	return gorix.POST, "/schema", func(
		ctx *gorix.Context,
	) (any, error) {
		return ctx.
			Status(gorix.StatusOK).
			ResponseEntityJSON(func() (any, error) {
				return c.postService.EnsureSchema(ctx)
			})
	}
}

func (c *PostController) FindByID() (
	gorix.Method,
	gorix.Path,
	gorix.RouteHandler,
) {
	return gorix.GET, "/:id", func(
		ctx *gorix.Context,
	) (any, error) {
		var params PostPathDto

		return ctx.
			BindParams(&params).
			Status(gorix.StatusOK).
			ResponseEntityJSON(func() (any, error) {
				return c.postService.FindByID(
					ctx,
					params.ID,
				)
			})
	}
}

func (c *PostController) Find() (
	gorix.Method,
	gorix.Path,
	gorix.RouteHandler,
) {
	return gorix.GET, "/", func(
		ctx *gorix.Context,
	) (any, error) {
		var query PostQueryDto

		return ctx.
			BindQuery(&query).
			Status(gorix.StatusOK).
			ResponseEntityJSON(func() (any, error) {
				return c.postService.Find(
					ctx,
					query,
				)
			})
	}
}

func (c *PostController) Count() (
	gorix.Method,
	gorix.Path,
	gorix.RouteHandler,
) {
	return gorix.GET, "/count", func(
		ctx *gorix.Context,
	) (any, error) {
		var query PostQueryDto

		return ctx.
			BindQuery(&query).
			Status(gorix.StatusOK).
			ResponseEntityJSON(func() (any, error) {
				return c.postService.Count(
					ctx,
					query,
				)
			})
	}
}

func (c *PostController) Create() (
	gorix.Method,
	gorix.Path,
	gorix.RouteHandler,
) {
	return gorix.POST, "/", func(
		ctx *gorix.Context,
	) (any, error) {
		var body CreatePostDto

		return ctx.
			BindBody(&body).
			Status(gorix.StatusCreated).
			ResponseEntityJSON(func() (any, error) {
				return c.postService.Create(
					ctx,
					body,
				)
			})
	}
}

func (c *PostController) CreateTx() (
	gorix.Method,
	gorix.Path,
	gorix.RouteHandler,
) {
	return gorix.POST, "/tx", func(
		ctx *gorix.Context,
	) (any, error) {
		var body CreatePostDto

		return ctx.
			BindBody(&body).
			Status(gorix.StatusCreated).
			ResponseEntityJSON(func() (any, error) {
				return c.postService.CreateTx(
					ctx,
					body,
				)
			})
	}
}

func (c *PostController) Update() (
	gorix.Method,
	gorix.Path,
	gorix.RouteHandler,
) {
	return gorix.PUT, "/:id", func(
		ctx *gorix.Context,
	) (any, error) {
		var params PostPathDto
		var body UpdatePostDto

		return ctx.
			BindParams(&params).
			BindBody(&body).
			Status(gorix.StatusOK).
			ResponseEntityJSON(func() (any, error) {
				return c.postService.Update(
					ctx,
					params.ID,
					body,
				)
			})
	}
}

func (c *PostController) Delete() (
	gorix.Method,
	gorix.Path,
	gorix.RouteHandler,
) {
	return gorix.DELETE, "/:id", func(
		ctx *gorix.Context,
	) (any, error) {
		var params PostPathDto

		return ctx.
			BindParams(&params).
			Status(gorix.StatusNoContent).
			ResponseEntityJSON(func() (any, error) {
				return c.postService.DeleteByID(
					ctx,
					params.ID,
				)
			})
	}
}
