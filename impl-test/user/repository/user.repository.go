package user

import (
	"fmt"

	"github.com/Gromosome/gorix/gorix"
	"github.com/Gromosome/gorix/impl-test/user/entity"
)

type User = entity.User
type UserSummary = entity.UserSummary
type UserRepository struct {
	databases *gorix.DBManager
	mapper    *gorix.SQLMapper
	repo      *gorix.SQLRepository[entity.User, int64]
}

func NewUserRepository(
	databases *gorix.DBManager,
) *UserRepository {
	userRepo, err := gorix.NewSQLRepository[entity.User, int64](
		databases,
	)
	if err != nil {
		panic(err)
	}
	return &UserRepository{
		databases: databases,
		mapper:    gorix.NewSQLMapper(databases),
		repo:      userRepo,
	}
}

func (r *UserRepository) FindByID(
	ctx *gorix.Context,
	id int64,
) (*User, error) {
	user, err := r.repo.FindByID(ctx, id)
	if err != nil {
		if gorix.IsEntityNotFound(err) {
			return nil, fmt.Errorf(
				"user %d not found",
				id,
			)
		}
		return nil, err
	}
	return user, nil
}
func (r *UserRepository) FindAll(
	ctx *gorix.Context,
) ([]User, error) {
	return r.repo.FindAll(ctx)
}

func (r *UserRepository) Save(
	ctx *gorix.Context,
	user *User,
) error {
	if user == nil {
		return fmt.Errorf(
			"user repository: user cannot be nil",
		)
	}

	return r.repo.Save(ctx, user)
}

func (r *UserRepository) Update(
	ctx *gorix.Context,
	user *User,
) error {
	if user == nil {
		return fmt.Errorf(
			"user repository: user cannot be nil",
		)
	}

	return r.repo.Update(ctx, user)
}

func (r *UserRepository) DeleteByID(
	ctx *gorix.Context,
	id int64,
) error {
	return r.repo.DeleteByID(ctx, id)
}

func (r *UserRepository) FindActiveUsers(
	ctx *gorix.Context,
	limit int,
	offset int,
) ([]User, error) {
	if limit <= 0 {
		limit = 20
	}

	users := make([]User, 0)

	err := r.mapper.QueryMany(
		ctx,
		&users,
		`
			SELECT
				id,
				name,
				email,
				active,
				created_at,
				updated_at
			FROM users
			WHERE active = $1
			ORDER BY created_at DESC
			LIMIT $2
			OFFSET $3
		`,
		true,
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) FindByEmail(
	ctx *gorix.Context,
	email string,
) (*User, error) {
	user := &User{}

	err := r.mapper.QueryOne(
		ctx,
		user,
		`
			SELECT
				id,
				name,
				email,
				active,
				created_at,
				updated_at
			FROM users
			WHERE email = $1
		`,
		email,
	)
	if err != nil {
		if gorix.DBIsNoRows(err) {
			return nil, fmt.Errorf(
				"user with email %s not found",
				email,
			)
		}

		return nil, err
	}

	return user, nil
}

func (r *UserRepository) Summary(
	ctx *gorix.Context,
) (*UserSummary, error) {
	summary := &UserSummary{}

	err := r.mapper.QueryOne(
		ctx,
		summary,
		`
			SELECT
				COUNT(*) AS total_users,
				COUNT(*) FILTER (
					WHERE active = TRUE
				) AS active_users
			FROM users
		`,
	)
	if err != nil {
		return nil, err
	}

	return summary, nil
}

func (r *UserRepository) CreateWithAudit(
	ctx *gorix.Context,
	user *User,
) error {
	if user == nil {
		return fmt.Errorf(
			"user repository: user cannot be nil",
		)
	}

	db, err := r.databases.DB()
	if err != nil {
		return err
	}

	return gorix.WithTransaction(
		ctx,
		db,
		&gorix.DBTxOptions{
			Isolation: gorix.DBIsolationReadCommitted,
			ReadOnly:  false,
		},
		func(
			ctx *gorix.Context,
			tx *gorix.DBTx,
		) error {
			transactionORM := r.repo.WithExecutor(tx)

			transactionMapper :=
				r.mapper.WithExecutor(tx)

			if err := transactionORM.Insert(
				ctx,
				user,
			); err != nil {
				return err
			}

			result := transactionMapper.Exec(
				ctx,
				`
					INSERT INTO user_audit (
						user_id,
						action
					)
					VALUES ($1, $2)
				`,
				user.ID,
				"USER_CREATED",
			)

			return result.Err()
		},
	)
}

func (r *UserRepository) UpdateWithAudit(
	ctx *gorix.Context,
	user *User,
) error {
	if user == nil {
		return fmt.Errorf(
			"user repository: user cannot be nil",
		)
	}

	db, err := r.databases.DB()
	if err != nil {
		return err
	}

	return gorix.WithTransaction(
		ctx,
		db,
		nil,
		func(
			ctx *gorix.Context,
			tx *gorix.DBTx,
		) error {
			transactionORM := r.repo.WithExecutor(tx)

			transactionMapper :=
				r.mapper.WithExecutor(tx)

			if err := transactionORM.Update(
				ctx,
				user,
			); err != nil {
				return err
			}

			result := transactionMapper.Exec(
				ctx,
				`
					INSERT INTO user_audit (
						user_id,
						action
					)
					VALUES ($1, $2)
				`,
				user.ID,
				"USER_UPDATED",
			)

			return result.Err()
		},
	)
}
