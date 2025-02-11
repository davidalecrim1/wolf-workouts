package main

import (
	"log/slog"
	"os"

	"github.com/davidalecrim1/wolf-workouts/internal/trainings/adapters"
	"github.com/davidalecrim1/wolf-workouts/internal/trainings/app/command"
	"github.com/davidalecrim1/wolf-workouts/internal/trainings/config"
	"github.com/davidalecrim1/wolf-workouts/internal/trainings/handler"
	"github.com/davidalecrim1/wolf-workouts/internal/trainings/server"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	logLevel := adapters.GetLogLevel(os.Getenv("LOG_LEVEL"))
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	}))
	slog.SetDefault(logger)

	db, err := sqlx.Connect("postgres", os.Getenv("TRAININGS_DATABASE_URL"))
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	config := config.NewConfig(os.Getenv("TRAININGS_JWT_SECRET"))

	trainingRepository := adapters.NewPostgresTrainingsRepository(db)
	commandTrainingHandler := command.NewTrainingCommandHandler(trainingRepository)
	// TODO: Create query training handler
	httpTrainingHandler := handler.NewHttpTrainingHandler(commandTrainingHandler)

	router := gin.Default()
	authMiddleware := handler.AuthMiddleware(config)
	httpTrainingHandler.RegisterRoutes(authMiddleware, router)

	server := server.NewServer(db, os.Getenv("TRAININGS_API_PORT"), router)
	server.RegisterRoutes(router)
	server.Start()
}
