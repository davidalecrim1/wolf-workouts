package app

import (
	"time"

	"github.com/google/uuid"
)

type Training struct {
	ID               string
	UserID           string
	Username         string
	TrainingDateTime time.Time
	Notes            string
}

func NewTraining(
	userID,
	username string,
	trainingDateTime time.Time,
	notes string,
) *Training {
	if userID == "" || username == "" || trainingDateTime.IsZero() {
		panic("invalid training data")
	}

	return &Training{
		ID:               uuid.New().String(),
		UserID:           userID,
		Username:         username,
		TrainingDateTime: trainingDateTime,
		Notes:            notes,
	}
}
