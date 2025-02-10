package adapter

import (
	"context"
	"os"
	"testing"

	"github.com/davidalecrim1/wolf-workouts/internal/users/app"
	testHelpers "github.com/davidalecrim1/wolf-workouts/internal/users/test"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

var db *sqlx.DB

func TestMain(m *testing.M) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conn, closeDB := testHelpers.GetTestDatabase(ctx)
	db = conn
	defer closeDB()

	os.Exit(m.Run())
}

func TestPostgresUserRepository_CreateUser(t *testing.T) {
	repo := NewPostgresUserRepository(db)

	user, err := app.NewUser("John Doe", "john.doe@example.com", "password", "trainee")
	require.NoError(t, err)

	err = repo.CreateUser(context.Background(), user)
	require.NoError(t, err)

	createdUser, err := repo.GetUserByEmail(context.Background(), "john.doe@example.com")
	require.NoError(t, err)
	require.Equal(t, user, createdUser)
}
