package test

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetMongoTestDatabase(ctx context.Context, collectionName string) (collection *mongo.Collection, closeFn func()) {
	mgContainer := CreateTestDatabase(ctx)

	hostWithPort, err := mgContainer.PortEndpoint(ctx, nat.Port("27017"), "")
	if err != nil {
		log.Fatalf("failed to get container port: %s", err)
	}

	connectionString := fmt.Sprintf("mongodb://root:example@%s", hostWithPort)

	opts := options.Client().
		ApplyURI(connectionString).
		SetServerSelectionTimeout(5 * time.Second)

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		log.Fatalf("failed to create mongodb test container client %v", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("failed to ping mongodb test container %v", err)
	}

	trainerDatabase := client.Database("trainer")

	return trainerDatabase.Collection(collectionName), func() {
		if err := mgContainer.Terminate(ctx); err != nil {
			log.Printf("failed to terminate container: %s", err)
		}
	}
}

func CreateTestDatabase(ctx context.Context) testcontainers.Container {
	initScriptPath := "../../../scripts/database/trainer/init.js"

	req := testcontainers.ContainerRequest{
		Image:        "mongo:7.0",
		ExposedPorts: []string{"27017/tcp"},
		Env: map[string]string{
			"MONGO_INITDB_ROOT_USERNAME": "root",
			"MONGO_INITDB_ROOT_PASSWORD": "example",
		},
		Files: []testcontainers.ContainerFile{
			{
				HostFilePath:      initScriptPath,
				ContainerFilePath: "/docker-entrypoint-initdb.d/init.js",
				FileMode:          0o755,
			},
		},
	}

	mgContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatalf("failed to create mongodb database for tests: %v", err)
	}

	return mgContainer
}
