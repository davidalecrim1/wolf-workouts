package queries

import (
	"context"

	"github.com/davidalecrim1/wolf-workouts/internal/trainings/app"
)

type TrainingQueriesRepository interface {
	FindTrainingsForUser(ctx context.Context, userID string) ([]*app.Training, error)
}

type TrainingQueriesHandler struct {
	repo TrainingQueriesRepository
}

func NewTrainingQueriesHandler(repo TrainingQueriesRepository) *TrainingQueriesHandler {
	return &TrainingQueriesHandler{repo: repo}
}

func (h *TrainingQueriesHandler) FindTrainingsForUser(ctx context.Context, userID string) ([]*app.Training, error) {
	return h.repo.FindTrainingsForUser(ctx, userID)
}
