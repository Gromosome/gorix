package user

import (
	"github.com/Gromosome/gorix/gorix"
	"github.com/Gromosome/gorix/impl-test/user/controller"
	"github.com/Gromosome/gorix/impl-test/user/repository"
	"github.com/Gromosome/gorix/impl-test/user/service"
)

type UserModule struct {
	userController controller.UserController
}

func NewUserModule() *UserModule {
	return &UserModule{}
}

func (m *UserModule) BasePath() gorix.BasePath {
	return "/user"
}

func (m *UserModule) APIVersion() gorix.APIVersion {
	return gorix.V1
}

func (m *UserModule) Providers() []any {
	return []any{
		gorix.Provider(
			repository.NewUserRepository,
			gorix.As((*repository.UserRepositoryPort)(nil)),
		),
		gorix.Provider(
			service.NewUserService,
			gorix.As((*service.UserServicePort)(nil)),
		),
	}
}

func (m *UserModule) Controllers() []any {
	return []any{
		gorix.Controller(
			controller.NewUserController,
			"/core",
		),
		gorix.Controller(
			controller.NewMediaController,
			"/media",
			gorix.V2,
		),
	}
}
