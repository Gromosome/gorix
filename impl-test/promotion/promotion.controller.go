// Package promotion module
package promotion

import (
	"github.com/Gromosome/gorix/gorix"
)

type PromotionController struct {
	promotionService PromotionService
}

func NewPromotionController(promotionService PromotionService) PromotionController {
	return PromotionController{
		promotionService: promotionService,
	}
}

func (c *PromotionController) GetPromotionList() (gorix.Method, gorix.Path, []string) {
	return gorix.GET, "/find", c.promotionService.GetPromotionList()
}
