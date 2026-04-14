package service

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}
func (s *UserService) GetUser() User {
	return User{
		ID:   1,
		Name: "Вася Залупкин",
	}
}
