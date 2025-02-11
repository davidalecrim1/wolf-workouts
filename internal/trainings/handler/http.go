package handler

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/davidalecrim1/wolf-workouts/internal/trainings/app/command"
	gen "github.com/davidalecrim1/wolf-workouts/internal/trainings/handler/generated"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// TODO: REST API FOR REACT

// TODO: Create auth middleware

// TODO: Create handler for schedule training

type HttpTrainingHandler struct {
	commandHandler *command.TrainingCommandHandler
}

func NewHttpTrainingHandler(commandHandler *command.TrainingCommandHandler) *HttpTrainingHandler {
	return &HttpTrainingHandler{commandHandler: commandHandler}
}

func (h *HttpTrainingHandler) ScheduleTraining(c *gin.Context) {
	var req gen.ScheduleTrainingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("ScheduleTraining - failed to bind request", "error", err)
		c.JSON(http.StatusBadRequest, gen.ResponseError{Message: err.Error()})
		return
	}

	userID, username, err := getUserDataFromContext(c)
	if err != nil {
		slog.Error("ScheduleTraining - failed to get user data", "error", err)
		c.JSON(http.StatusUnauthorized, gen.ResponseError{Message: "Unauthorized"})
		return
	}

	cmd := command.ScheduleTrainingCommand{
		UserID:           userID,
		Username:         username,
		TrainingDateTime: req.TrainingDatetime,
		Notes:            req.Notes,
	}

	if err := h.commandHandler.ScheduleTraining(c.Request.Context(), &cmd); err != nil {
		slog.Error("ScheduleTraining - failed to handle command", "error", err)
		c.JSON(http.StatusInternalServerError, gen.ResponseError{Message: err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func getUserDataFromContext(c *gin.Context) (userID string, username string, err error) {
	claims, ok := c.Get(JWT_CLAIMS_KEY)
	if !ok {
		return "", "", fmt.Errorf("failed to get JWT claims")
	}

	claimsMap, ok := claims.(jwt.MapClaims)
	if !ok {
		return "", "", fmt.Errorf("failed to get JWT claims")
	}

	userID, ok = claimsMap["sub"].(string)
	if !ok {
		return "", "", fmt.Errorf("failed to get JWT claims")
	}

	username, ok = claimsMap["user_name"].(string)
	if !ok {
		return "", "", fmt.Errorf("failed to get JWT claims")
	}

	return userID, username, nil
}

func (h *HttpTrainingHandler) RegisterRoutes(middleware gin.HandlerFunc, router *gin.Engine) {
	router.POST("/trainings", middleware, h.ScheduleTraining)
}
