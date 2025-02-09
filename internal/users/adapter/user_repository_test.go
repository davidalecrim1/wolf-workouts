package adapter

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/davidalecrim1/wolf-workouts/internal/users/app"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	testcontainerPostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var db *sqlx.DB

func TestMain(m *testing.M) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pgContainer := createTestDatabase(ctx)
	defer pgContainer.Terminate(ctx)

	connectionString, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		log.Fatalf("failed to get connection string: %v", err)
	}

	db, err = sqlx.Connect("postgres", connectionString)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	loadDatabaseSchema(ctx)
	os.Exit(m.Run())
}

func createTestDatabase(ctx context.Context) (container *testcontainerPostgres.PostgresContainer) {
	pgContainer, err := testcontainerPostgres.Run(
		ctx,
		"docker.io/postgres:16.4-alpine3.20",
		testcontainerPostgres.WithDatabase("postgres"),
		testcontainerPostgres.WithUsername("postgres"),
		testcontainerPostgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		log.Fatalf("failed to create database for tests: %v", err)
	}
	return pgContainer
}

func loadDatabaseSchema(ctx context.Context) {
	schema, err := os.ReadFile("../../../scripts/database/user/schema.sql")
	if err != nil {
		log.Fatalf("failed to read schema file: %v", err)
	}

	_, err = db.ExecContext(ctx, string(schema))
	if err != nil {
		log.Fatalf("failed to load schema: %v", err)
	}
}

func TestPostgresUserRepository_CreateUser(t *testing.T) {
	repo := NewPostgresUserRepository(db)

	user, err := app.NewUser("John Doe", "john.doe@example.com", "password")
	require.NoError(t, err)

	err = repo.CreateUser(context.Background(), user)
	require.NoError(t, err)

	createdUser, err := repo.GetUserByEmail(context.Background(), "john.doe@example.com")
	require.NoError(t, err)
	require.Equal(t, user, createdUser)
}
