package app

import (
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidEmailOrPassword = errors.New("invalid email or password")
	ErrInvalidRole            = errors.New("invalid role")
)

type Role uint8

const (
	RoleTrainer Role = iota + 1
	RoleTrainee
)

func (r Role) String() string {
	switch r {
	case RoleTrainer:
		return "trainer"
	case RoleTrainee:
		return "trainee"
	default:
		return "unknown"
	}
}

func ParseRole(s string) (Role, error) {
	switch s {
	case "trainer":
		return RoleTrainer, nil
	case "trainee":
		return RoleTrainee, nil
	default:
		return 0, ErrInvalidRole
	}
}

type User struct {
	ID             string
	Name           string
	Email          string
	HashedPassword string
	Role           Role
}

func NewUser(
	name string,
	email string,
	providedPassword string,
	roleString string,
) (*User, error) {
	role, err := ParseRole(roleString)
	if err != nil {
		return nil, err
	}

	u := &User{
		Name:  name,
		Email: email,
		Role:  role,
	}

	u.ID = uuid.New().String()

	err = u.addHashedPassword(providedPassword)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (u *User) addHashedPassword(providedPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(providedPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.HashedPassword = string(hash)
	return nil
}

func (u *User) IsPasswordCorrect(providedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.HashedPassword), []byte(providedPassword))
	return err == nil
}
