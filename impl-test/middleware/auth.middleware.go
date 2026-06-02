package middlewares

import (
	"github.com/Gromosome/gorix/gorix"
	"github.com/Gromosome/gorix/gorix/core"
)

func AuthMiddleware() gorix.Middleware {
	return func(next gorix.Handler) gorix.Handler {
		return func(c *gorix.Context) error {
			token := c.R.Header.Get("Authorization")

			if token != "Bearer dev-token" {
				return c.Status(core.StatusUnauthorized).JSON(map[string]any{
					"success": false,
					"error":   "unauthorized",
				})
			}

			return next(c)
		}
	}
}
