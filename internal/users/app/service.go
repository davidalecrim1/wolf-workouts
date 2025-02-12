package app

import (
	"context"
	"errors"
	"time"

	"github.com/davidalecrim1/wolf-workouts/internal/users/config"
	"github.com/dgrijalva/jwt-go"
)

var ErrUserAlreadyExists = errors.New("user already exists")

type UserService struct {
	userRepository UserRepository
	config         *config.Config
}

func NewUserService(
	userRepository UserRepository,
	config *config.Config,
) *UserService {
	return &UserService{
		userRepository: userRepository,
		config:         config,
	}
}

type UserRepository interface {
	CreateUser(ctx context.Context, user *User) error
	GetUserByEmail(ctx context.Context, email string) (*User, error)
}

func (s *UserService) CreateUser(ctx context.Context, user *User) error {
	return s.userRepository.CreateUser(ctx, user)
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	user, err := s.userRepository.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, ErrInvalidEmailOrPassword
	}

	return user, nil
}

func (s *UserService) LoginUser(ctx context.Context, email string, password string) (string, error) {
	u, err := s.userRepository.GetUserByEmail(ctx, email)
	if err != nil {
		return "", ErrInvalidEmailOrPassword
	}

	if !u.IsPasswordCorrect(password) {
		return "", ErrInvalidEmailOrPassword
	}

	return s.generateToken(u)
}

func (s *UserService) generateToken(u *User) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"sub":       u.ID,
			"user_name": u.Name,
			"user_role": u.Role.String(),
			"exp":       time.Now().Add(time.Hour * 24).Unix(),
		},
	)

	token.Header = map[string]interface{}{
		"alg": "HS256",
		"typ": "JWT",
	}

	tokenString, err := token.SignedString(s.config.GetJWTSecret())
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
