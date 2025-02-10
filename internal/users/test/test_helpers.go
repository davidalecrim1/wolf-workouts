package test

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	testcontainerPostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func GetTestDatabase(ctx context.Context) (conn *sqlx.DB, close func()) {
	pgContainer := CreateTestDatabase(ctx)

	connectionString, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		log.Fatalf("failed to get connection string: %v", err)
	}

	db, err := sqlx.Connect("postgres", connectionString)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	loadDatabaseSchema(ctx, db)

	return db, func() {
		err := pgContainer.Terminate(ctx)
		if err != nil {
			log.Fatalf("failed to terminate container: %v", err)
		}
	}
}

func CreateTestDatabase(ctx context.Context) (container *testcontainerPostgres.PostgresContainer) {
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

func loadDatabaseSchema(ctx context.Context, db *sqlx.DB) {
	schema, err := os.ReadFile("../../../scripts/database/users/schema.sql")
	if err != nil {
		log.Fatalf("failed to read schema file: %v", err)
	}

	_, err = db.ExecContext(ctx, string(schema))
	if err != nil {
		log.Fatalf("failed to load schema: %v", err)
	}
}
