package app

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/davidalecrim1/wolf-workouts/internal/users/config"

	"github.com/stretchr/testify/require"
)

var c *config.Config

func TestMain(m *testing.M) {
	c = config.NewConfig("secret")

	os.Exit(m.Run())
}

func createTestService(t *testing.T) (*UserService, *FakeUserRepository) {
	t.Helper()
	repo := NewFakeUserRepository()
	return NewUserService(repo, c), repo
}

func createTestUser(t *testing.T, svc *UserService, email, password string, roleString string) *User {
	t.Helper()
	user, err := NewUser("John Doe", email, password, roleString)
	require.NoError(t, err)
	require.NoError(t, svc.CreateUser(context.Background(), user))
	return user
}

func TestUserService_CreateUser(t *testing.T) {
	t.Parallel()

	t.Run("should create a user", func(t *testing.T) {
		svc, repo := createTestService(t)
		user := createTestUser(t, svc, "john.doe@example.com", "password", "trainee")

		createdUser, err := repo.GetUserByEmail(context.Background(), user.Email)
		require.NoError(t, err)
		require.Equal(t, user.ID, createdUser.ID)
	})

	t.Run("should return an error if the user does not exist", func(t *testing.T) {
		svc, _ := createTestService(t)
		user, err := NewUser("John Doe", "john.doe@example.com", "password", "trainee")
		require.NoError(t, err)

		invalidEmail := "invalid_email"

		user, err = svc.GetUserByEmail(context.Background(), invalidEmail)
		require.Nil(t, user)
		require.True(t, errors.Is(err, ErrInvalidEmailOrPassword))
	})

	t.Run("should not allow invalid roles", func(t *testing.T) {
		user, err := NewUser("John Doe", "john.doe@example.com", "password", "invalid_role")
		require.Error(t, err)
		require.ErrorIs(t, err, ErrInvalidRole)
		require.Nil(t, user)
	})

	t.Run("should be able to login a user", func(t *testing.T) {
		svc, _ := createTestService(t)
		_ = createTestUser(t, svc, "john.doe@example.com", "password", "trainee")

		token, err := svc.LoginUser(context.Background(), "john.doe@example.com", "password")
		require.NoError(t, err)
		require.NotEmpty(t, token)
	})

	t.Run("should return an error if the password is incorrect", func(t *testing.T) {
		svc, _ := createTestService(t)
		_ = createTestUser(t, svc, "john.doe@example.com", "password", "trainee")

		token, err := svc.LoginUser(context.Background(), "john.doe@example.com", "wrong_password")
		require.True(t, errors.Is(err, ErrInvalidEmailOrPassword))
		require.Empty(t, token)
	})

	t.Run("shouldn't be able to see the user password", func(t *testing.T) {
		svc, _ := createTestService(t)
		providedPassword := "password"
		_ = createTestUser(t, svc, "john.doe@example.com", providedPassword, "trainee")

		savedUser, err := svc.GetUserByEmail(context.Background(), "john.doe@example.com")
		require.NoError(t, err)
		require.NotNil(t, savedUser)
		require.NotEqual(t, savedUser.HashedPassword, providedPassword)
	})

	t.Run("should return an error is the user in the login is not found", func(t *testing.T) {
		svc, _ := createTestService(t)
		token, err := svc.LoginUser(context.Background(), "john.doe@example.com", "password")
		require.True(t, errors.Is(err, ErrInvalidEmailOrPassword))
		require.Empty(t, token)
	})

	t.Run("should detect existing user", func(t *testing.T) {
		svc, repo := createTestService(t)
		registeredUser := createTestUser(t, svc, "exists@test.com", "password", "trainee")

		foundUser, err := repo.GetUserByEmail(context.Background(), registeredUser.Email)

		require.NoError(t, err)
		require.Equal(t, registeredUser.ID, foundUser.ID)
	})
}

type FakeID string

type FakeUserRepository struct {
	users map[FakeID]*User
}

func NewFakeUserRepository() *FakeUserRepository {
	return &FakeUserRepository{
		users: make(map[FakeID]*User),
	}
}

func (r *FakeUserRepository) CreateUser(ctx context.Context, user *User) error {
	r.users[FakeID(user.ID)] = user
	return nil
}

func (r *FakeUserRepository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}

	return nil, errors.New("user not found")
}
