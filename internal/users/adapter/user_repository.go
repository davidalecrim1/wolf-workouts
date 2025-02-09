package adapter

import (
	"context"
	"time"

	"github.com/davidalecrim1/wolf-workouts/internal/users/app"
	"github.com/jmoiron/sqlx"
)

type userModel struct {
	ID             string    `db:"uuid"`
	Name           string    `db:"name"`
	Email          string    `db:"email"`
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
		INSERT INTO users (uuid, name, email, hashed_password)
		VALUES (:uuid, :name, :email, :hashed_password)
	`
	_, err := r.db.NamedExecContext(ctx, query, r.mapUserToUserModel(user))
	return err
}

func (r *PostgresUserRepository) GetUserByEmail(ctx context.Context, email string) (*app.User, error) {
	query := `SELECT * FROM users WHERE email = $1`
	var userModel userModel
	err := r.db.GetContext(ctx, &userModel, query, email)
	if err != nil {
		return nil, err
	}
	return r.mapUserModelToUser(userModel), nil
}

func (r *PostgresUserRepository) mapUserModelToUser(userModel userModel) *app.User {
	return &app.User{
		ID:             userModel.ID,
		Name:           userModel.Name,
		Email:          userModel.Email,
		HashedPassword: userModel.HashedPassword,
	}
}

func (r *PostgresUserRepository) mapUserToUserModel(user *app.User) userModel {
	return userModel{
		ID:             user.ID,
		Name:           user.Name,
		Email:          user.Email,
		HashedPassword: user.HashedPassword,
	}
}
