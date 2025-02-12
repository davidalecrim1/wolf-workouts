package adapters

import (
	"context"
	"os"
	"strconv"
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
	commandsRepo := NewPostgresTrainingsCommandsRepository(db)
	queriesRepo := NewPostgresTrainingsQueriesRepository(db)

	t.Run("should create a training", func(t *testing.T) {
		userID := uuid.New().String()
		training := app.NewTraining(userID, time.Now(), "Notes")

		err := commandsRepo.CreateTraining(context.Background(), training)
		require.NoError(t, err)
		assertTrainingEqual(t, training, getTrainingByID(t, queriesRepo, training.UserID, training.ID))
	})
}

func TestPostgresTrainingsRepository_FindTrainingsForUser(t *testing.T) {
	t.Parallel()
	commandsRepo := NewPostgresTrainingsCommandsRepository(db)
	queriesRepo := NewPostgresTrainingsQueriesRepository(db)

	t.Run("should find trainings for user", func(t *testing.T) {
		userID := uuid.New().String()
		amount := 4
		trainings := make([]*app.Training, amount)

		for i := 0; i < amount; i++ {
			trainings[i] = app.NewTraining(
				userID,
				time.Now().Add(time.Duration(i)*time.Hour),
				"Notes "+strconv.Itoa(i),
			)

			err := commandsRepo.CreateTraining(context.Background(), trainings[i])
			require.NoError(t, err)
		}

		actualTrainings, err := queriesRepo.FindTrainingsForUser(context.Background(), userID)
		require.NoError(t, err)
		require.Equal(t, amount, len(actualTrainings))

		for i, expected := range trainings {
			assertTrainingEqual(t, expected, actualTrainings[i])
		}
	})

	t.Run("should return empty slice when no trainings exist", func(t *testing.T) {
		userID := uuid.New().String()

		trainings, err := queriesRepo.FindTrainingsForUser(context.Background(), userID)
		require.NoError(t, err)
		require.Empty(t, trainings)
	})
}

func getTrainingByID(t *testing.T, repo *PostgresTrainingsQueriesRepository, userID, trainingID string) *app.Training {
	t.Helper()
	training, err := repo.GetTrainingByID(context.Background(), userID, trainingID)
	require.NoError(t, err)
	return training
}

func assertTrainingEqual(t *testing.T, expected, actual *app.Training) {
	t.Helper()
	require.Equal(t, expected.ID, actual.ID)
	require.Equal(t, expected.UserID, actual.UserID)
	require.Equal(t, expected.Notes, actual.Notes)
	require.Equal(
		t,
		expected.TrainingDateTime.Format("2006-01-02 15:04:05.999999"),
		actual.TrainingDateTime.Format("2006-01-02 15:04:05.999999"),
		"TrainingDateTime mismatch",
	)
}
