package user

type UserService struct {
}

func (s *UserService) GetUserList() []string {
	return []string{"user1", "user2", "user3"}
}
