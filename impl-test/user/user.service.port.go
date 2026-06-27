package user

import "github.com/Gromosome/gorix/gorix"

type UserServicePort interface {
	GetByID(
		ctx *gorix.Context,
		id int64,
	) (*User, error)

	GetAll(
		ctx *gorix.Context,
	) ([]User, error)

	GetActive(
		ctx *gorix.Context,
		query UserQueryDto,
	) ([]User, error)

	GetSummary(
		ctx *gorix.Context,
	) (*UserSummary, error)

	Create(
		ctx *gorix.Context,
		request CreateUserDto,
	) (*User, error)

	Update(
		ctx *gorix.Context,
		id int64,
		request UpdateUserDto,
	) (*User, error)

	Delete(
		ctx *gorix.Context,
		id int64,
	) (any, error)
}
