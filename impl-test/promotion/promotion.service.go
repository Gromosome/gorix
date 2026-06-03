package promotion

type PromotionService struct {
}

func NewPromotionService() *PromotionService {
	return &PromotionService{}
}
func (s *PromotionService) GetPromotionList() []string {
	return []string{"promo1", "promo2", "promo3"}
}
