package promotion

import (
	"github.com/Gromosome/gorix/gorix"
)

type PromotionModule struct {
	promotionController PromotionController
}

func NewPromotionModule() *PromotionModule {
	promotionService := PromotionService{}
	promotionController := NewPromotionController(promotionService)

	return &PromotionModule{
		promotionController: promotionController,
	}
}

func (m *PromotionModule) GetPromotionController() (gorix.BasePath, PromotionController) {
	return "/promotion", m.promotionController
}
