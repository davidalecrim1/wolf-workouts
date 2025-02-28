package adapter

import (
	"context"
	"time"

	"github.com/davidalecrim1/wolf-workouts/internal/trainer/domain"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type HourMongoDbRepository struct {
	collection *mongo.Collection
}

func NewHourMongoDbRepository(c *mongo.Collection) *HourMongoDbRepository {
	return &HourMongoDbRepository{
		collection: c,
	}
}

type HourModel struct {
	HourID       string `bson:"_id" json:"id"`
	Availability int    `bson:"availability" json:"availability"`
}

func NewHourModel(h time.Time, availability int) *HourModel {
	return &HourModel{
		HourID:       h.Format("2006-01-02T15:04"),
		Availability: availability,
	}
}

func (r *HourMongoDbRepository) UpdateHour(
	ctx context.Context,
	h *domain.Hour,
	updateFn func(h *domain.Hour) (*domain.Hour, error),
) error {
	updatedHour, err := updateFn(h)
	if err != nil {
		zap.S().Error("failed to insert in the database", "error", err, "ctx", ctx)
		return err
	}

	hm := NewHourModel(updatedHour.Hour, int(h.GetAvailability()))
	_, err = r.collection.InsertOne(ctx, hm)
	if err != nil {
		zap.S().Error("failed to while inserting in the database", "error", err, "ctx", ctx)
	}

	return nil
}
