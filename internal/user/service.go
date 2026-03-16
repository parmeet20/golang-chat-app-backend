package user

import (
	"github.com/parmeet20/golang-chatapp/internal/auth"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserService struct {
	userRepo    *UserRepo
	authService *auth.AuthService
}

func NewUserService(userRepo *UserRepo, authService *auth.AuthService) *UserService {
	return &UserService{userRepo: userRepo, authService: authService}
}

func (s *UserService) Register(user *User) error {
	return s.userRepo.Create(user)
}

func (s *UserService) FindById(id primitive.ObjectID) (*User, error) {
	return s.userRepo.FindById(id)
}

func (s *UserService) Login(username, password string) (string, error) {

	user, err := s.userRepo.Login(username, password)
	if err != nil {
		return "", err
	}

	token, err := s.authService.GenerateToken(user.ID.Hex(), user.Email)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *UserService) FindByEmail(email string) (*User, error) {
	return s.userRepo.FindByEmail(email)
}
