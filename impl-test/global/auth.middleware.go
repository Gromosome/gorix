package global

import (
	"github.com/Gromosome/gorix/gorix"
	"github.com/Gromosome/gorix/gorix/core/context"
)

func AuthMiddleware() gorix.Middleware {
	return func(next gorix.Handler) gorix.Handler {
		return func(c *gorix.Context) error {
			token := c.R.Header.Get("Authorization")

			if token != "Bearer dev-token" {
				return c.Status(gorix.StatusUnauthorized).SOAPFault11(context.SOAP11FaultMustUnderstand, "UnAuthorized", "UnAuthorized")
			}
			return next(c)
		}
	}
}
