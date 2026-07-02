package user

import (
	"github.com/Gromosome/gorix/gorix"
	"github.com/Gromosome/gorix/impl-test/user/repository"
)

type MockUserRepository struct{}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{}
}

func (m *MockUserRepository) FindByID(
	ctx *gorix.Context,
	id int64,
) (*repository.User, error) {
	return &repository.User{
		ID:     id,
		Name:   "Mock User",
		Email:  "mock@gorix.dev",
		Active: true,
	}, nil
}

func (m *MockUserRepository) FindAll(
	ctx *gorix.Context,
) ([]repository.User, error) {
	return []repository.User{
		{
			ID:     1,
			Name:   "Mock User",
			Email:  "mock@gorix.dev",
			Active: true,
		},
	}, nil
}

func (m *MockUserRepository) FindActiveUsers(
	ctx *gorix.Context,
	limit int,
	offset int,
) ([]repository.User, error) {
	return m.FindAll(ctx)
}

func (m *MockUserRepository) FindByEmail(
	ctx *gorix.Context,
	email string,
) (*repository.User, error) {
	return &repository.User{
		ID:     1,
		Name:   "Mock User",
		Email:  email,
		Active: true,
	}, nil
}

func (m *MockUserRepository) Summary(
	ctx *gorix.Context,
) (*repository.UserSummary, error) {
	return &repository.UserSummary{
		TotalUsers:  1,
		ActiveUsers: 1,
	}, nil
}

func (m *MockUserRepository) CreateWithAudit(
	ctx *gorix.Context,
	u *repository.User,
) error {
	u.ID = 100
	return nil
}

func (m *MockUserRepository) UpdateWithAudit(
	ctx *gorix.Context,
	u *repository.User,
) error {
	return nil
}

func (m *MockUserRepository) DeleteByID(
	ctx *gorix.Context,
	id int64,
) error {
	return nil
}
