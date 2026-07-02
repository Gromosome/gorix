package service

import (
	"github.com/Gromosome/gorix/gorix"
	"github.com/Gromosome/gorix/impl-test/user/dto"
	"github.com/Gromosome/gorix/impl-test/user/repository"
)

type UserServicePort interface {
	GetByID(
		ctx *gorix.Context,
		id int64,
	) (*repository.User, error)

	GetAll(
		ctx *gorix.Context,
	) ([]repository.User, error)

	GetActive(
		ctx *gorix.Context,
		query dto.UserQueryDto,
	) ([]repository.User, error)

	GetSummary(
		ctx *gorix.Context,
	) (*repository.UserSummary, error)

	Create(
		ctx *gorix.Context,
		request dto.CreateUserDto,
	) (*repository.User, error)

	Update(
		ctx *gorix.Context,
		id int64,
		request dto.UpdateUserDto,
	) (*repository.User, error)

	Delete(
		ctx *gorix.Context,
		id int64,
	) (any, error)
}
