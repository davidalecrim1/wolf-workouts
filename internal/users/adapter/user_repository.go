package adapter

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"github.com/davidalecrim1/wolf-workouts/internal/users/app"
	"github.com/jmoiron/sqlx"
)

type userModel struct {
	ID             string    `db:"uuid"`
	Name           string    `db:"name"`
	Email          string    `db:"email"`
	Role           string    `db:"role"`
	HashedPassword string    `db:"hashed_password"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}

type PostgresUserRepository struct {
	db *sqlx.DB
}

func NewPostgresUserRepository(db *sqlx.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) CreateUser(ctx context.Context, user *app.User) error {
	query := `
		INSERT INTO public.users (uuid, name, email, role, hashed_password)
		VALUES (:uuid, :name, :email, :role, :hashed_password)
	`
	userModel, err := r.mapUserToUserModel(user)
	if err != nil {
		slog.Error("failed to create user", "error", err)
		return err
	}

	_, err = r.db.NamedExecContext(ctx, query, userModel)
	if err != nil {
		slog.Error("failed to create user", "error", err)
		return err
	}

	return nil
}

func (r *PostgresUserRepository) GetUserByEmail(ctx context.Context, email string) (*app.User, error) {
	query := `SELECT uuid, name, email, role, hashed_password, created_at, updated_at FROM users WHERE email = $1`
	var userModel userModel
	err := r.db.GetContext(ctx, &userModel, query, email)

	if errors.Is(err, sql.ErrNoRows) {
		slog.Debug("user not found", "email", email)
		return nil, errors.New("user not found")
	}

	if err != nil {
		slog.Error("failed to get user by email", "error", err)
		return nil, err
	}
	return r.mapUserModelToUser(userModel)
}

func (r *PostgresUserRepository) mapUserModelToUser(userModel userModel) (*app.User, error) {
	role, err := app.ParseRole(userModel.Role)
	if err != nil {
		return nil, err
	}

	return &app.User{
		ID:             userModel.ID,
		Name:           userModel.Name,
		Email:          userModel.Email,
		HashedPassword: userModel.HashedPassword,
		Role:           role,
	}, nil
}

func (r *PostgresUserRepository) mapUserToUserModel(user *app.User) (userModel, error) {
	role, err := app.ParseRole(user.Role.String())
	if err != nil {
		return userModel{}, err
	}

	return userModel{
		ID:             user.ID,
		Name:           user.Name,
		Email:          user.Email,
		Role:           role.String(),
		HashedPassword: user.HashedPassword,
	}, nil
}
