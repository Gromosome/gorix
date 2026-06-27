package user

import (
	"github.com/Gromosome/gorix/gorix"
)

type UserModule struct {
	userController UserController
}

func NewUserModule() *UserModule {
	return &UserModule{}
}

func (m *UserModule) BasePath() gorix.BasePath {
	return "/user"
}

func (m *UserModule) Providers() []any {
	return []any{
		gorix.Provider(
			NewUserRepository,
			gorix.As((*UserRepositoryPort)(nil)),
		),
		gorix.Provider(
			NewUserService,
			gorix.As((*UserServicePort)(nil)),
		),
	}
}

func (m *UserModule) Controllers() []any {
	return []any{
		NewUserController,
	}
}
