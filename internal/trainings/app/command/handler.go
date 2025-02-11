package command

import (
	"context"
	"fmt"

	"github.com/davidalecrim1/wolf-workouts/internal/trainings/app"
)

type TrainingRepository interface {
	CreateTraining(ctx context.Context, training *app.Training) error
	GetTrainingByID(ctx context.Context, userID string, trainingID string) (*app.Training, error)
}

type TrainingCommandHandler struct {
	repo TrainingRepository
}

func NewTrainingCommandHandler(repo TrainingRepository) *TrainingCommandHandler {
	return &TrainingCommandHandler{repo: repo}
}

func (h *TrainingCommandHandler) ScheduleTraining(ctx context.Context, cmd *ScheduleTrainingCommand) error {
	training := app.NewTraining(cmd.UserID, cmd.Username, cmd.TrainingDateTime, cmd.Notes)
	err := h.repo.CreateTraining(ctx, training)
	if err != nil {
		return fmt.Errorf("failed to create training: %w", err)
	}

	// TODO: Call the ScheduleTraining in trainer service

	return nil
}
