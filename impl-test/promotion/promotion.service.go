package promotion

import "fmt"

type PromotionService struct {
}

func NewPromotionService() *PromotionService {
	return &PromotionService{}
}
func (s *PromotionService) GetPromotionList() []string {
	panic(fmt.Errorf("promotion service error"))
}
