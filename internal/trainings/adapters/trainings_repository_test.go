package adapters

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/davidalecrim1/wolf-workouts/internal/trainings/app"
	testHelpers "github.com/davidalecrim1/wolf-workouts/internal/trainings/test"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

var db *sqlx.DB

func TestMain(m *testing.M) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conn, closeDB := testHelpers.GetTestDatabase(ctx)
	db = conn
	defer closeDB()

	os.Exit(m.Run())
}

func TestPostgresTrainingsRepository_CreateTraining(t *testing.T) {
	t.Parallel()
	repo := NewPostgresTrainingsRepository(db)

	testCases := []struct {
		name            string
		trainingFactory func(t *testing.T) *app.Training
	}{
		{
			name: "should create a training",
			trainingFactory: func(t *testing.T) *app.Training {
				userID := uuid.New().String()
				return app.NewTraining(userID, "John Doe", time.Now(), "Notes")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			training := tc.trainingFactory(t)
			err := repo.CreateTraining(context.Background(), training)
			require.NoError(t, err)
			assertCreatedTraining(t, repo, training)
		})
	}
}

func assertCreatedTraining(t *testing.T, repo *PostgresTrainingsRepository, training *app.Training) {
	createdTraining, err := repo.GetTrainingByID(context.Background(), training.UserID, training.ID)
	require.NoError(t, err)

	require.Equal(t, training.ID, createdTraining.ID)
	require.Equal(t, training.UserID, createdTraining.UserID)
	require.Equal(t, training.Notes, createdTraining.Notes)

	require.Equal(
		t,
		training.TrainingDateTime.Format("2006-01-02 15:04:05.999999"),
		createdTraining.TrainingDateTime.Format("2006-01-02 15:04:05.999999"),
		"TrainingDateTime mismatch",
	)
}
