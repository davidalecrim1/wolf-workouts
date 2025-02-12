package app

import (
	"time"

	"github.com/google/uuid"
)

type Training struct {
	ID               string
	UserID           string
	TrainingDateTime time.Time
	Notes            string
}

func NewTraining(
	userID string,
	trainingDateTime time.Time,
	notes string,
) *Training {
	if userID == "" || trainingDateTime.IsZero() {
		panic("invalid training data")
	}

	return &Training{
		ID:               uuid.New().String(),
		UserID:           userID,
		TrainingDateTime: trainingDateTime,
		Notes:            notes,
	}
}
