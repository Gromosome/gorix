package hook

import (
	"github.com/Gromosome/gorix/gorix/core"
)

type Handler func(c *core.Context) error

type Middleware func(handler Handler) Handler

type MiddlewareConfig struct {
	Middleware Middleware
	Rule       RouteRule
}

type MiddlewareBuilder struct {
	config MiddlewareConfig
}

func Apply(middleware Middleware) MiddlewareBuilder {
	return MiddlewareBuilder{
		config: MiddlewareConfig{
			Middleware: middleware,
			Rule:       RouteRule{},
		},
	}
}

func (b MiddlewareBuilder) Only(paths ...string) MiddlewareConfig {
	b.config.Rule.OnlyPaths = paths
	return b.config
}

func (b MiddlewareBuilder) Except(paths ...string) MiddlewareConfig {
	b.config.Rule.ExceptPaths = paths
	return b.config
}

func GlobalMiddleware(middleware Middleware) MiddlewareConfig {
	return MiddlewareConfig{
		Middleware: middleware,
		Rule:       RouteRule{},
	}
}

func ChainMiddlewares(handler Handler, middlewares ...Middleware) Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}
