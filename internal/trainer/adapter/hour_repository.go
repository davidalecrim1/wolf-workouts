package adapter

import (
	"context"

	"github.com/davidalecrim1/wolf-workouts/internal/trainer/domain"
	"go.mongodb.org/mongo-driver/mongo"
)

type HourMongoDbRepository struct {
	collection *mongo.Collection
}

func NewHourMongoDbRepository(c *mongo.Collection) *HourMongoDbRepository {
	return &HourMongoDbRepository{
		collection: c,
	}
}

func (r *HourMongoDbRepository) UpdateHour(
	ctx context.Context,
	h *domain.Hour,
	updateFn func(h *domain.Hour) (*domain.Hour, error),
) error {
	// TODO: Save in mongo. Create the DB model.
	return nil
}
