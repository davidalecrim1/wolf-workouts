package command

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/davidalecrim1/wolf-workouts/internal/trainings/adapters"
	"github.com/davidalecrim1/wolf-workouts/internal/trainings/app"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestScheduleTraining(t *testing.T) {
	repo := NewFakeTrainingRepository()
	commandHandler := NewTrainingCommandHandler(repo)

	testCases := []struct {
		name           string
		commandFactory func(t *testing.T) *ScheduleTrainingCommand
		expectedError  error
	}{
		{
			name: "should schedule a training",
			commandFactory: func(t *testing.T) *ScheduleTrainingCommand {
				return &ScheduleTrainingCommand{
					UserID:           uuid.New().String(),
					Notes:            "This is a test training",
					TrainingDateTime: time.Now(),
				}
			},
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			command := tc.commandFactory(t)

			err := commandHandler.ScheduleTraining(context.Background(), command)
			require.Equal(t, tc.expectedError, err)
		})
	}
}

type FakeTrainingRepository struct {
	trainings map[string]*app.Training
	mu        sync.RWMutex
}

func NewFakeTrainingRepository() *FakeTrainingRepository {
	return &FakeTrainingRepository{
		trainings: make(map[string]*app.Training),
	}
}

func (r *FakeTrainingRepository) CreateTraining(ctx context.Context, training *app.Training) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.trainings[training.ID] = training
	return nil
}

func (r *FakeTrainingRepository) GetTrainingByID(ctx context.Context, userID string, trainingID string) (*app.Training, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	training, ok := r.trainings[trainingID]
	if !ok {
		return nil, adapters.ErrTrainingNotFound
	}

	return training, nil
}
