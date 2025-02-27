package command

import (
	"context"
	"time"

	"github.com/davidalecrim1/wolf-workouts/internal/trainer/domain"
)

type ScheduleTrainingHandler struct {
	writeRepo WriteRepository
}

type WriteRepository interface {
	UpdateHour(
		ctx context.Context,
		h *domain.Hour,
		updateFn func(h *domain.Hour) (*domain.Hour, error)) error
}

func NewScheduleTrainingHandler(repo WriteRepository) *ScheduleTrainingHandler {
	return &ScheduleTrainingHandler{
		writeRepo: repo,
	}
}

type ScheduleTrainingCommand struct {
	Timestamp time.Time
}

func (h *ScheduleTrainingHandler) Handle(ctx context.Context, cmd *ScheduleTrainingCommand) error {
	hour, err := domain.NewAvailableHour(
		cmd.Timestamp,
	)
	if err != nil {
		return err
	}

	return h.writeRepo.UpdateHour(ctx, hour, func(h *domain.Hour) (*domain.Hour, error) {
		if err = hour.ScheduleTraining(); err != nil {
			return nil, err
		}

		return hour, nil
	})
}
