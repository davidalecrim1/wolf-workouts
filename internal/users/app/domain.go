package app

import (
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidEmailOrPassword = errors.New("invalid email or password")

type User struct {
	ID             string
	Name           string
	Email          string
	HashedPassword string
}

func NewUser(
	name string,
	email string,
	providedPassword string,
) (*User, error) {
	u := &User{
		Name:  name,
		Email: email,
	}

	u.ID = uuid.New().String()

	err := u.addHashedPassword(providedPassword)
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
