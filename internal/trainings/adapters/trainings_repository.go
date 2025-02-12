package adapters

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/davidalecrim1/wolf-workouts/internal/trainings/app"
	"github.com/jmoiron/sqlx"
)

var ErrTrainingNotFound = fmt.Errorf("training not found")

// Commands Section

type PostgresTrainingsCommandsRepository struct {
	db *sqlx.DB
}

func NewPostgresTrainingsCommandsRepository(db *sqlx.DB) *PostgresTrainingsCommandsRepository {
	return &PostgresTrainingsCommandsRepository{db: db}
}

func (r *PostgresTrainingsCommandsRepository) CreateTraining(ctx context.Context, training *app.Training) error {
	query := `
		INSERT INTO trainings (uuid, user_id, training_datetime, notes)
		VALUES (:uuid, :user_id, :training_datetime, :notes)
	`
	trainingModel := marshalTraining(training)

	_, err := r.db.NamedExecContext(ctx, query, trainingModel)
	if err != nil {
		slog.Error("failed to create training", "error", err)
		return fmt.Errorf("failed to create training: %w", err)
	}

	return nil
}

// Queries Section

type PostgresTrainingsQueriesRepository struct {
	db *sqlx.DB
}

func NewPostgresTrainingsQueriesRepository(db *sqlx.DB) *PostgresTrainingsQueriesRepository {
	return &PostgresTrainingsQueriesRepository{db: db}
}

func (r *PostgresTrainingsQueriesRepository) FindTrainingsForUser(ctx context.Context, userID string) ([]*app.Training, error) {
	query := `
		SELECT uuid, user_id, training_datetime, notes, created_at, updated_at
		FROM trainings
		WHERE user_id = $1
	`

	var trainings []*trainingModel
	err := r.db.SelectContext(ctx, &trainings, query, userID)
	if err == sql.ErrNoRows {
		slog.Debug("no trainings found for user", "user_id", userID)
		return nil, nil
	}

	if err != nil {
		slog.Error("failed to find trainings for user", "error", err)
		return nil, fmt.Errorf("failed to find trainings for user: %w", err)
	}

	return unmarshalTrainings(trainings), nil
}

func (r *PostgresTrainingsQueriesRepository) GetTrainingByID(ctx context.Context, userID string, trainingID string) (*app.Training, error) {
	query := `
		SELECT uuid, user_id,  training_datetime, notes, created_at, updated_at
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

	return unmarshalTraining(&trainingModel), nil
}
