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
		params gorix.Params[PostPathDto],
	) (any, error) {
		return ctx.
			Status(gorix.StatusOK).
			ResponseEntityJSON(func() (any, error) {
				return c.postService.FindByID(
					ctx,
					params.Value.ID,
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
		query gorix.Query[PostQueryDto],
	) (any, error) {
		return ctx.
			Status(gorix.StatusOK).
			ResponseEntityJSON(func() (any, error) {
				return c.postService.Find(
					ctx,
					query.Value,
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
		query gorix.Query[PostQueryDto],
	) (any, error) {
		return ctx.
			Status(gorix.StatusOK).
			ResponseEntityJSON(func() (any, error) {
				return c.postService.Count(
					ctx,
					query.Value,
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
		body gorix.Body[CreatePostDto],
	) (any, error) {
		return ctx.
			Status(gorix.StatusCreated).
			ResponseEntityJSON(func() (any, error) {
				return c.postService.Create(
					ctx,
					body.Value,
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
		body gorix.Body[CreatePostDto],
	) (any, error) {
		return ctx.
			Status(gorix.StatusCreated).
			ResponseEntityJSON(func() (any, error) {
				return c.postService.CreateTx(
					ctx,
					body.Value,
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
		params gorix.Params[PostPathDto],
		body gorix.Body[UpdatePostDto],
	) (any, error) {
		return ctx.
			Status(gorix.StatusOK).
			ResponseEntityJSON(func() (any, error) {
				return c.postService.Update(
					ctx,
					params.Value.ID,
					body.Value,
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
		params gorix.Params[PostPathDto],
	) (any, error) {
		return ctx.
			Status(gorix.StatusNoContent).
			ResponseEntityJSON(func() (any, error) {
				return c.postService.DeleteByID(
					ctx,
					params.Value.ID,
				)
			})
	}
}
