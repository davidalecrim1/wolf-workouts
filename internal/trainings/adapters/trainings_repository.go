package adapters

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/davidalecrim1/wolf-workouts/internal/trainings/app"
	"github.com/jmoiron/sqlx"
)

var ErrTrainingNotFound = fmt.Errorf("training not found")

type PostgresTrainingsRepository struct {
	db *sqlx.DB
}

func NewPostgresTrainingsRepository(db *sqlx.DB) *PostgresTrainingsRepository {
	return &PostgresTrainingsRepository{db: db}
}

type trainingModel struct {
	ID               string    `db:"uuid"`
	UserID           string    `db:"user_id"`
	Username         string    `db:"username"`
	TrainingDateTime time.Time `db:"training_datetime"`
	Notes            string    `db:"notes"`
	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`
}

func (r *PostgresTrainingsRepository) CreateTraining(ctx context.Context, training *app.Training) error {
	query := `
		INSERT INTO trainings (uuid, user_id, username, training_datetime, notes)
		VALUES (:uuid, :user_id, :username, :training_datetime, :notes)
	`
	trainingModel, err := r.marshalTraining(training)
	if err != nil {
		return err
	}

	_, err = r.db.NamedExecContext(ctx, query, trainingModel)
	if err != nil {
		slog.Error("failed to create training", "error", err)
		return fmt.Errorf("failed to create training: %w", err)
	}

	return nil
}

func (r *PostgresTrainingsRepository) marshalTraining(training *app.Training) (trainingModel, error) {
	return trainingModel{
		ID:               training.ID,
		UserID:           training.UserID,
		Username:         training.Username,
		TrainingDateTime: training.TrainingDateTime,
		Notes:            training.Notes,
	}, nil
}

func (r *PostgresTrainingsRepository) GetTrainingByID(ctx context.Context, userID string, trainingID string) (*app.Training, error) {
	query := `
		SELECT uuid, user_id, username, training_datetime, notes, created_at, updated_at
		FROM trainings
		WHERE user_id = $1 AND uuid = $2
	`

	var trainingModel trainingModel
	err := r.db.GetContext(ctx, &trainingModel, query, userID, trainingID)
	if err == sql.ErrNoRows {
		return nil, ErrTrainingNotFound
	}

	if err != nil {
		slog.Error("failed to get training by id", "error", err)
		return nil, fmt.Errorf("failed to get training by id: %w", err)
	}

	return r.unmarshalTraining(&trainingModel), nil
}

func (r *PostgresTrainingsRepository) unmarshalTraining(trainingModel *trainingModel) *app.Training {
	return &app.Training{
		ID:               trainingModel.ID,
		UserID:           trainingModel.UserID,
		Username:         trainingModel.Username,
		TrainingDateTime: trainingModel.TrainingDateTime,
		Notes:            trainingModel.Notes,
	}
}
