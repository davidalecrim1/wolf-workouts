package app

import (
	"context"
	"errors"
	"time"

	"github.com/davidalecrim1/wolf-workouts/internal/users/config"
	"github.com/dgrijalva/jwt-go"
)

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
		return nil, errors.Join(err, ErrInvalidEmailOrPassword)
	}

	return user, nil
}

func (s *UserService) LoginUser(ctx context.Context, email string, password string) (string, error) {
	u, err := s.userRepository.GetUserByEmail(ctx, email)
	if err != nil {
		return "", errors.Join(err, ErrInvalidEmailOrPassword)
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
			"user_id": u.ID,
			"exp":     time.Now().Add(time.Hour * 24).Unix(),
		},
	)

	tokenString, err := token.SignedString(s.config.GetJWTSecret())
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
