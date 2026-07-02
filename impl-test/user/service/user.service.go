package service

import (
	gorixcontext "github.com/Gromosome/gorix/gorix/core/context"
	"github.com/Gromosome/gorix/impl-test/user/dto"
	"github.com/Gromosome/gorix/impl-test/user/repository"
)

type UserService struct {
	userRepository repository.UserRepositoryPort
}

func NewUserService(userRepository repository.UserRepositoryPort) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (s *UserService) GetByID(
	ctx *gorixcontext.Context,
	id int64,
) (*repository.User, error) {
	return s.userRepository.FindByID(ctx, id)
}

func (s *UserService) GetAll(
	ctx *gorixcontext.Context,
) ([]repository.User, error) {
	return s.userRepository.FindAll(ctx)
}

func (s *UserService) GetActive(
	ctx *gorixcontext.Context,
	query dto.UserQueryDto,
) ([]repository.User, error) {
	return s.userRepository.FindActiveUsers(
		ctx,
		query.Limit,
		query.Offset,
	)
}

func (s *UserService) GetSummary(
	ctx *gorixcontext.Context,
) (*repository.UserSummary, error) {
	return s.userRepository.Summary(ctx)
}

func (s *UserService) Create(
	ctx *gorixcontext.Context,
	request dto.CreateUserDto,
) (*repository.User, error) {

	user := &repository.User{
		Name:   request.Name,
		Email:  request.Email,
		Active: true,
	}

	if err := s.userRepository.CreateWithAudit(
		ctx,
		user,
	); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Update(
	ctx *gorixcontext.Context,
	id int64,
	request dto.UpdateUserDto,
) (*repository.User, error) {
	user, err := s.userRepository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	user.Name = request.Name
	user.Active = request.Active

	if err := s.userRepository.UpdateWithAudit(
		ctx,
		user,
	); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Delete(
	ctx *gorixcontext.Context,
	id int64,
) (any, error) {
	return nil, s.userRepository.DeleteByID(ctx, id)
}
