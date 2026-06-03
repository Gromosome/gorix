package promotion

import "github.com/Gromosome/gorix/gorix"

type PromotionModule struct {
}

func NewPromotionModule() *PromotionModule {
	return &PromotionModule{}
}

func (m *PromotionModule) BasePath() gorix.BasePath {
	return "/promotion"
}

func (m *PromotionModule) Providers() []any {
	return []any{
		NewPromotionService,
	}
}

func (m *PromotionModule) Controllers() []any {
	return []any{
		NewPromotionController,
	}
}
