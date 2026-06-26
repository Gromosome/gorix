// Package promotion module
package promotion

import (
	"github.com/Gromosome/gorix/gorix"
)

type PromotionController struct {
	promotionService *PromotionService
}

func NewPromotionController(promotionService *PromotionService) *PromotionController {
	return &PromotionController{
		promotionService: promotionService,
	}
}

func (c *PromotionController) GetPromotionList() (gorix.Method, gorix.Path, gorix.RouteHandler) {
	return gorix.GET, "/find", func(ctx *gorix.Context) (any, error) {
		return ctx.ResponseEntityXML(func() (any, error) {
			return c.promotionService.GetPromotionList(), nil
		})

	}
}
