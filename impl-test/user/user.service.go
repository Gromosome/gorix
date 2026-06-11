package user

import (
	gorixcontext "github.com/Gromosome/gorix/gorix/core/context"
)

type UserService struct {
	userRepository *UserRepository
}

func NewUserService(userRepository *UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (s *UserService) GetByID(
	ctx *gorixcontext.Context,
	id int64,
) (*User, error) {
	return s.userRepository.FindByID(ctx, id)
}

func (s *UserService) GetAll(
	ctx *gorixcontext.Context,
) ([]User, error) {
	return s.userRepository.FindAll(ctx)
}

func (s *UserService) GetActive(
	ctx *gorixcontext.Context,
	query UserQueryDto,
) ([]User, error) {
	return s.userRepository.FindActiveUsers(
		ctx,
		query.Limit,
		query.Offset,
	)
}

func (s *UserService) GetSummary(
	ctx *gorixcontext.Context,
) (*UserSummary, error) {
	return s.userRepository.Summary(ctx)
}

func (s *UserService) Create(
	ctx *gorixcontext.Context,
	request CreateUserDto,
) (*User, error) {

	user := &User{
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
	request UpdateUserDto,
) (*User, error) {
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
) error {
	return s.userRepository.DeleteByID(ctx, id)
}
