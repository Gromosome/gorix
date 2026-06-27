package user

import "github.com/Gromosome/gorix/gorix"

type UserRepositoryPort interface {
	FindByID(
		ctx *gorix.Context,
		id int64,
	) (*User, error)

	FindAll(
		ctx *gorix.Context,
	) ([]User, error)

	FindActiveUsers(
		ctx *gorix.Context,
		limit int,
		offset int,
	) ([]User, error)

	FindByEmail(
		ctx *gorix.Context,
		email string,
	) (*User, error)

	Summary(
		ctx *gorix.Context,
	) (*UserSummary, error)

	CreateWithAudit(
		ctx *gorix.Context,
		user *User,
	) error

	UpdateWithAudit(
		ctx *gorix.Context,
		user *User,
	) error

	DeleteByID(
		ctx *gorix.Context,
		id int64,
	) error
}
