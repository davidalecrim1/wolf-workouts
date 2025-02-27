package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"strings"
	"time"

	"github.com/davidalecrim1/wolf-workouts/internal/trainer/adapter"
	"github.com/davidalecrim1/wolf-workouts/internal/trainer/app/command"
	"github.com/davidalecrim1/wolf-workouts/internal/trainer/handler"
	gen "github.com/davidalecrim1/wolf-workouts/internal/trainer/handler/generated"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

func main() {
	ctx := context.Background()

	mongoEndpoint := os.Getenv("TRAINER_MONGODB_ENDPOINT")

	opts := options.Client().
		ApplyURI(mongoEndpoint).
		SetServerSelectionTimeout(5 * time.Second)

	mongoDbClient, err := mongo.Connect(ctx, opts)
	if err != nil {
		slog.Error("failed to create MongoDB client", "error", err)
		panic(err)
	}
	defer mongoDbClient.Disconnect(ctx)

	err = mongoDbClient.Ping(ctx, nil)
	if err != nil {
		slog.Error("failed to ping MongoDB", "error", err)
		panic(err)
	}

	trainerDatabase := mongoDbClient.Database(os.Getenv("TRAINER_MONGODB_DATABASE"))
	hoursCollection := trainerDatabase.Collection(os.Getenv("TRAINER_MONGODB_COLLECTION_HOURS"))
	_ = trainerDatabase.Collection(os.Getenv("TRAINER_MONGODB_COLLECTION_DATES"))
	slog.Info("successfully connected to MongoDB")

	serverType := strings.ToLower(os.Getenv("TRAINER_SERVER_TYPE"))
	port := os.Getenv("TRAINER_API_PORT")

	switch serverType {
	case "http":
		router := gin.Default()

		trainerHandler := handler.NewTrainerHttpHandler(mongoDbClient)
		trainerHandler.RegisterRoutes(router)

		slog.Info(fmt.Sprintf("starting to listen for %v on port %v", serverType, port))
		err := router.Run(":" + port)
		if err != nil {
			slog.Error("failed to start http server", "error", err)
			os.Exit(1)
		}

	case "grpc":
		s := grpc.NewServer()

		lis, err := net.Listen("tcp", ":"+port)
		if err != nil {
			slog.Error("failed to start grpc server", "error", err)
			os.Exit(1)
		}

		hoursRepo := adapter.NewHourMongoDbRepository(hoursCollection)
		cmdSth := command.NewScheduleTrainingHandler(
			hoursRepo,
		)

		trainerHandler := handler.NewTrainerGrpcHandler(mongoDbClient, cmdSth)
		gen.RegisterTrainerServiceServer(s, trainerHandler)

		slog.Info(fmt.Sprintf("starting to listen for %v on port %v", serverType, port))
		err = s.Serve(lis)
		if err != nil {
			slog.Error("failed to start grpc server", "error", err)
			os.Exit(1)
		}

	default:
		panic("invalid server type")
	}
}
