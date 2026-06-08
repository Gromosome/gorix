package user

import (
	"github.com/Gromosome/gorix/gorix/core/database"
	"github.com/Gromosome/gorix/gorix/core/database/mapper"
	"github.com/Gromosome/gorix/gorix/core/database/repository"
	"github.com/Gromosome/gorix/impl-test/user/entity"
)

type UserRepository struct {
	databases *database.Manager
	mapper    *mapper.Mapper
	repo      *repository.Repository[entity.User, int64]
}

func NewUserRepository(
	databases *database.Manager,
) *UserRepository {
	return &UserRepository{
		databases: databases,
		mapper:    mapper.New(databases),

		// This constructor must be lazy.
		// It must not request manager.DB() during module registration.
		repo: repository.NewRepository[entity.User, int64](
			databases,
		),
	}
}
