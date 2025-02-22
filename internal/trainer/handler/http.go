package handler

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// func httpHealthCheck(db *mongo.Client) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		ctx := context.Background()

// 		err := db.Ping(ctx, nil)
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"message": "Database connection failed"})
// 			return
// 		}

// 		c.JSON(http.StatusOK, gin.H{"message": "OK"})
// 	}
// }

type TrainerHttpHandler struct {
	db *mongo.Client
}

func NewTrainerHttpHandler(db *mongo.Client) *TrainerHttpHandler {
	return &TrainerHttpHandler{
		db: db,
	}
}

func (h *TrainerHttpHandler) HealthCheck(c *gin.Context) {
	err := h.db.Ping(c.Request.Context(), nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database connection failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}

func (h *TrainerHttpHandler) RegisterRoutes(r *gin.Engine) {
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/healthz", h.HealthCheck)
}
