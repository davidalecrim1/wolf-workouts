package main

import (
	"context"
	"fmt"
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
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	ctx := context.Background()

	logger := zap.Must(zap.NewDevelopment())
	defer logger.Sync()
	sugar := logger.Sugar()
	zap.ReplaceGlobals(sugar.Desugar())

	mongoEndpoint := os.Getenv("TRAINER_MONGODB_ENDPOINT")
	opts := options.Client().
		ApplyURI(mongoEndpoint).
		SetServerSelectionTimeout(5 * time.Second)
	mongoDbClient, err := mongo.Connect(ctx, opts)
	if err != nil {
		zap.S().Fatal("failed to create mongo database client", "error", err)
	}
	defer mongoDbClient.Disconnect(ctx)

	err = mongoDbClient.Ping(ctx, nil)
	if err != nil {
		zap.S().Fatal("failed to ping mongo database", "error", err)
	}

	trainerDatabase := mongoDbClient.Database(os.Getenv("TRAINER_MONGODB_DATABASE"))
	hoursCollection := trainerDatabase.Collection(os.Getenv("TRAINER_MONGODB_COLLECTION_HOURS"))
	_ = trainerDatabase.Collection(os.Getenv("TRAINER_MONGODB_COLLECTION_DATES"))
	zap.S().Info("successfully connected to mongo database")

	serverType := strings.ToLower(os.Getenv("TRAINER_SERVER_TYPE"))
	port := os.Getenv("TRAINER_API_PORT")

	switch serverType {
	case "http":
		router := gin.Default()

		trainerHandler := handler.NewTrainerHttpHandler(mongoDbClient)
		trainerHandler.RegisterRoutes(router)

		zap.S().Info(fmt.Sprintf("starting to listen for %v on port %v", serverType, port))
		err := router.Run(":" + port)
		if err != nil {
			zap.S().Fatal("failed to start http server", "error", err)
		}

	case "grpc":
		s := grpc.NewServer()

		lis, err := net.Listen("tcp", ":"+port)
		if err != nil {
			zap.S().Fatal("failed to start grpc server", "error", err)
		}

		hoursRepo := adapter.NewHourMongoDbRepository(hoursCollection)
		cmdSth := command.NewScheduleTrainingHandler(
			hoursRepo,
		)

		trainerHandler := handler.NewTrainerGrpcHandler(mongoDbClient, cmdSth)
		gen.RegisterTrainerServiceServer(s, trainerHandler)

		zap.S().Info(fmt.Sprintf("starting to listen for %v on port %v", serverType, port))
		err = s.Serve(lis)
		if err != nil {
			zap.S().Fatal("failed to start grpc server", "error", err)
		}

	default:
		zap.S().Fatal("the server type wasn't provided.")
	}
}
