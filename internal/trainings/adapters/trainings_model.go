package adapters

import (
	"time"

	"github.com/davidalecrim1/wolf-workouts/internal/trainings/app"
)

type trainingModel struct {
	ID               string    `db:"uuid"`
	UserID           string    `db:"user_id"`
	TrainingDateTime time.Time `db:"training_datetime"`
	Notes            string    `db:"notes"`
	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`
}

func unmarshalTraining(trainingModel *trainingModel) *app.Training {
	return &app.Training{
		ID:               trainingModel.ID,
		UserID:           trainingModel.UserID,
		TrainingDateTime: trainingModel.TrainingDateTime,
		Notes:            trainingModel.Notes,
	}
}

func marshalTraining(training *app.Training) trainingModel {
	return trainingModel{
		ID:               training.ID,
		UserID:           training.UserID,
		TrainingDateTime: training.TrainingDateTime,
		Notes:            training.Notes,
	}
}

func unmarshalTrainings(
	trainings []*trainingModel,
) []*app.Training {
	unmarshalledTrainings := make([]*app.Training, len(trainings))

	for i, training := range trainings {
		unmarshalledTrainings[i] = unmarshalTraining(training)
	}

	return unmarshalledTrainings
}
