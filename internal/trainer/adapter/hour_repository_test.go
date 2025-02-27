package adapter

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/davidalecrim1/wolf-workouts/internal/trainer/domain"
	testHelpers "github.com/davidalecrim1/wolf-workouts/internal/trainer/test"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
)

var collectionHours *mongo.Collection

func TestMain(m *testing.M) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	col, closeFn := testHelpers.GetMongoTestDatabase(ctx, "hours")
	collectionHours = col
	defer closeFn()

	os.Exit(m.Run())
}

func TestHourRepository(t *testing.T) {
	t.Run("should create or update an hour", func(t *testing.T) {
		repo := NewHourMongoDbRepository(collectionHours)
		require.NotNil(t, repo)

		ctx := context.Background()

		timestamp := time.Now().Truncate(time.Hour).AddDate(0, 0, 2)
		h, err := domain.NewAvailableHour(timestamp)
		require.NoError(t, err)

		err = repo.UpdateHour(ctx, h, func(hour *domain.Hour) (*domain.Hour, error) {
			if err = hour.ScheduleTraining(); err != nil {
				return nil, err
			}

			return hour, nil
		})

		require.NoError(t, err)
	})
}
