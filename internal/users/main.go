package main

import (
	"log"
	"os"

	"github.com/davidalecrim1/wolf-workouts/internal/users/adapter"
	"github.com/davidalecrim1/wolf-workouts/internal/users/app"
	"github.com/davidalecrim1/wolf-workouts/internal/users/config"
	"github.com/davidalecrim1/wolf-workouts/internal/users/handler"
	"github.com/davidalecrim1/wolf-workouts/internal/users/server"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sqlx.Connect("postgres", os.Getenv("USERS_DATABASE_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	config := config.NewConfig(os.Getenv("USERS_JWT_SECRET"))

	userRepository := adapter.NewPostgresUserRepository(db)
	userService := app.NewUserService(userRepository, config)
	userHandler := handler.NewUserHandler(userService)

	router := gin.Default()
	server := server.NewServer(db, os.Getenv("USERS_API_PORT"), router)

	server.RegisterRoutes(router)
	userHandler.RegisterRoutes(router)

	server.Start()
}
