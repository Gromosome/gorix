package user

import (
	"github.com/Gromosome/gorix/gorix"
	"github.com/Gromosome/gorix/impl-test/user/dto"
	"github.com/Gromosome/gorix/impl-test/user/repository"
)

type MockUserService struct {
	GetByIDCalled bool
	GetByIDValue  int64

	CreateCalled bool
	CreateBody   dto.CreateUserDto
}

func NewMockUserService() *MockUserService {
	return &MockUserService{}
}

func (m *MockUserService) GetByID(
	ctx *gorix.Context,
	id int64,
) (*repository.User, error) {
	m.GetByIDCalled = true
	m.GetByIDValue = id

	return &repository.User{
		ID:     id,
		Name:   "Mock Controller User",
		Email:  "controller@gorix.dev",
		Active: true,
	}, nil
}

func (m *MockUserService) GetAll(
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

func (m *MockUserService) GetActive(
	ctx *gorix.Context,
	query dto.UserQueryDto,
) ([]repository.User, error) {
	return m.GetAll(ctx)
}

func (m *MockUserService) GetSummary(
	ctx *gorix.Context,
) (*repository.UserSummary, error) {
	return &repository.UserSummary{
		TotalUsers:  1,
		ActiveUsers: 1,
	}, nil
}

func (m *MockUserService) Create(
	ctx *gorix.Context,
	request dto.CreateUserDto,
) (*repository.User, error) {
	m.CreateCalled = true
	m.CreateBody = request

	return &repository.User{
		ID:     100,
		Name:   request.Name,
		Email:  request.Email,
		Active: true,
	}, nil
}

func (m *MockUserService) Update(
	ctx *gorix.Context,
	id int64,
	request dto.UpdateUserDto,
) (*repository.User, error) {
	return &repository.User{
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
