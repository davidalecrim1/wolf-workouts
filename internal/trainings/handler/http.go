package handler

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/davidalecrim1/wolf-workouts/internal/trainings/app/command"
	"github.com/davidalecrim1/wolf-workouts/internal/trainings/app/queries"
	gen "github.com/davidalecrim1/wolf-workouts/internal/trainings/handler/generated"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type HttpTrainingHandler struct {
	commandHandler *command.TrainingCommandHandler
	queriesHandler *queries.TrainingQueriesHandler
}

func NewHttpTrainingHandler(c *command.TrainingCommandHandler, q *queries.TrainingQueriesHandler) *HttpTrainingHandler {
	return &HttpTrainingHandler{
		commandHandler: c,
		queriesHandler: q,
	}
}

func (h *HttpTrainingHandler) ScheduleTraining(c *gin.Context) {
	var req gen.ScheduleTrainingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("ScheduleTraining - failed to bind request", "error", err)
		c.JSON(http.StatusBadRequest, gen.ResponseError{Message: err.Error()})
		return
	}

	userID, err := getUserIDFromContext(c)
	if err != nil {
		slog.Error("ScheduleTraining - failed to get user data", "error", err)
		c.JSON(http.StatusUnauthorized, gen.ResponseError{Message: "Unauthorized"})
		return
	}

	cmd := command.ScheduleTrainingCommand{
		UserID:           userID,
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

func (h *HttpTrainingHandler) GetTrainings(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		slog.Error("GetTrainings - failed to get user data", "error", err)
		c.JSON(http.StatusUnauthorized, gen.ResponseError{Message: "Unauthorized"})
		return
	}

	trainings, err := h.queriesHandler.FindTrainingsForUser(c.Request.Context(), userID)
	if err != nil {
		slog.Error("GetTrainings - failed to get trainings", "error", err)
		c.JSON(http.StatusInternalServerError, gen.ResponseError{Message: err.Error()})
		return
	}

	responseTrainings := make([]gen.Training, len(trainings))
	for i, training := range trainings {
		responseTrainings[i] = gen.Training{
			Id:               &training.ID,
			Notes:            &training.Notes,
			TrainingDatetime: &training.TrainingDateTime,
		}
	}

	c.JSON(http.StatusOK, responseTrainings)
}

func getUserIDFromContext(c *gin.Context) (userID string, err error) {
	claims, ok := c.Get(JWT_CLAIMS_KEY)
	if !ok {
		return "", fmt.Errorf("failed to get JWT claims")
	}

	claimsMap, ok := claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("failed to get JWT claims")
	}

	userID, ok = claimsMap["sub"].(string)
	if !ok {
		return "", fmt.Errorf("failed to get JWT claims")
	}

	return userID, nil
}

func (h *HttpTrainingHandler) RegisterRoutes(middleware gin.HandlerFunc, router *gin.Engine) {
	router.POST("/trainings", middleware, h.ScheduleTraining)
	router.GET("/trainings", middleware, h.GetTrainings)
}
