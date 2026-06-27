package user

import (
	"github.com/Gromosome/gorix/gorix"
	"github.com/Gromosome/gorix/impl-test/user"
)

type MockUserService struct {
	GetByIDCalled bool
	GetByIDValue  int64

	CreateCalled bool
	CreateBody   user.CreateUserDto
}

func NewMockUserService() *MockUserService {
	return &MockUserService{}
}

func (m *MockUserService) GetByID(
	ctx *gorix.Context,
	id int64,
) (*user.User, error) {
	m.GetByIDCalled = true
	m.GetByIDValue = id

	return &user.User{
		ID:     id,
		Name:   "Mock Controller User",
		Email:  "controller@gorix.dev",
		Active: true,
	}, nil
}

func (m *MockUserService) GetAll(
	ctx *gorix.Context,
) ([]user.User, error) {
	return []user.User{
		{
			ID:     1,
			Name:   "Mock User",
			Email:  "mock@gorix.dev",
			Active: true,
		},
	}, nil
}

func (m *MockUserService) GetActive(
	ctx *gorix.Context,
	query user.UserQueryDto,
) ([]user.User, error) {
	return m.GetAll(ctx)
}

func (m *MockUserService) GetSummary(
	ctx *gorix.Context,
) (*user.UserSummary, error) {
	return &user.UserSummary{
		TotalUsers:  1,
		ActiveUsers: 1,
	}, nil
}

func (m *MockUserService) Create(
	ctx *gorix.Context,
	request user.CreateUserDto,
) (*user.User, error) {
	m.CreateCalled = true
	m.CreateBody = request

	return &user.User{
		ID:     100,
		Name:   request.Name,
		Email:  request.Email,
		Active: true,
	}, nil
}

func (m *MockUserService) Update(
	ctx *gorix.Context,
	id int64,
	request user.UpdateUserDto,
) (*user.User, error) {
	return &user.User{
		ID:     id,
		Name:   request.Name,
		Active: request.Active,
	}, nil
}

func (m *MockUserService) Delete(
	ctx *gorix.Context,
	id int64,
) (any, error) {
	return nil, nil
}
