package user

import (
	"github.com/Gromosome/gorix/gorix"
	"github.com/Gromosome/gorix/impl-test/promotion"
)

type UserModule struct {
	userController UserController
}

func NewUserModule() *UserModule {
	userService := UserService{}
	promotionService := promotion.PromotionService{}
	userController := NewUserController(userService, promotionService)
	return &UserModule{
		userController: userController,
	}
}

func (m *UserModule) GetUserController() (gorix.BasePath, UserController) {
	return "/user", m.userController
}
