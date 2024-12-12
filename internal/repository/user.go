package repository

import (
	"context"
	"fmt"
	db "mpc/internal/db/sqlc"
	"mpc/internal/model"
	"mpc/pkg/utils"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	queries *db.Queries
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{queries: db.New(pool)}
}

// CreateUser creates a new user
func (r *UserRepository) CreateUser(ctx context.Context, email, passwordHash string) (model.User, error) {
	user, err := r.queries.CreateUser(ctx, db.CreateUserParams{
		ID:           utils.ToPgUUID(uuid.New()),
		Email:        email,
		PasswordHash: passwordHash,
		Status:       "active",
		CreatedAt:    utils.CurrentPgTimestamp(),
		UpdatedAt:    utils.CurrentPgTimestamp(),
	})
	if err != nil {
		return model.User{}, fmt.Errorf("failed to create user: %w", err)
	}
	return toUserModel(user), nil
}

// GetUserByEmail retrieves a user by their email
func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	user, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return model.User{}, fmt.Errorf("failed to get user by email: %w", err)
	}
	return toUserModel(user), nil
}

// GetUserByID retrieves a user by their ID
func (r *UserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (model.User, error) {
	user, err := r.queries.GetUserByID(ctx, utils.ToPgUUID(id))
	if err != nil {
		return model.User{}, fmt.Errorf("failed to get user by ID: %w", err)
	}
	return toUserModel(user), nil
}

// UpdateUser updates a user
func (r *UserRepository) UpdateUser(ctx context.Context, user model.User) (model.User, error) {
	updatedUser, err := r.queries.UpdateUser(ctx, db.UpdateUserParams{
		ID:           utils.ToPgUUID(user.ID),
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		Status:       user.Status,
		UpdatedAt:    utils.CurrentPgTimestamp(),
	})
	if err != nil {
		return model.User{}, fmt.Errorf("failed to update user: %w", err)
	}
	return toUserModel(updatedUser), nil
}

// toUserModel converts a sqlc user to a model user
func toUserModel(sqlcUser db.User) model.User {
	return model.User{
		ID:           utils.ToUUID(sqlcUser.ID),
		Email:        sqlcUser.Email,
		PasswordHash: sqlcUser.PasswordHash,
		Status:       sqlcUser.Status,
		CreatedAt:    sqlcUser.CreatedAt.Time,
		UpdatedAt:    sqlcUser.UpdatedAt.Time,
	}
}
