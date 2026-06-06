package user

type UserService struct {
}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) GetUserList() []string {
	return []string{"user1", "user2", "user3"}
}

func (s *UserService) CreateUser(userDto CreateUserBodyDto) string {
	return "User created: " + userDto.Name
}
