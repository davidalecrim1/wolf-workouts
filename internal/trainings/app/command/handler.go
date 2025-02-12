package command

import (
	"context"
	"fmt"

	"github.com/davidalecrim1/wolf-workouts/internal/trainings/app"
)

type TrainingCommandsRepository interface {
	CreateTraining(ctx context.Context, training *app.Training) error
}

type TrainingCommandHandler struct {
	repo TrainingCommandsRepository
}

func NewTrainingCommandHandler(repo TrainingCommandsRepository) *TrainingCommandHandler {
	return &TrainingCommandHandler{repo: repo}
}

func (h *TrainingCommandHandler) ScheduleTraining(ctx context.Context, cmd *ScheduleTrainingCommand) error {
	training := app.NewTraining(cmd.UserID, cmd.TrainingDateTime, cmd.Notes)
	err := h.repo.CreateTraining(ctx, training)
	if err != nil {
		return fmt.Errorf("failed to create training: %w", err)
	}

	// TODO: Call the ScheduleTraining in trainer service

	return nil
}
