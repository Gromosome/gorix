package middlewares

import (
	"github.com/Gromosome/gorix/gorix"
)

func AuthMiddleware() gorix.Middleware {
	return func(next gorix.Handler) gorix.Handler {
		return func(c *gorix.Context) error {
			token := c.R.Header.Get("Authorization")

			if token != "Bearer dev-token" {
				return c.Status(gorix.StatusUnauthorized).JSON(map[string]any{
					"success": false,
					"error":   "unauthorized",
				})
			}
			return next(c)
		}
	}
}
