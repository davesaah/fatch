package services

import (
	"context"

	"github.com/davidreturns08/fatch/internal/database"
	"github.com/davidreturns08/fatch/internal/types"
	"github.com/jackc/pgx/v5/pgconn"
)

// UserService provides user-related services.
type UserService struct{}

// CreateUser creates a new user.
func (s *UserService) CreateUser(ctx context.Context, params database.CreateUserParams) (*types.ErrorResponse, error) {
	tx, err := initialiseDBTX(ctx)
	if err != nil {
		return types.InternalServerErrorResponse(), err
	}
	defer tx.Rollback(ctx)

	qb := database.NewQueryBuilder(tx)
	if err := qb.CreateUser(ctx, params); err != nil {
		pgErr := err.(*pgconn.PgError)
		return types.ConflictErrorResponse(pgErr.Message), err
	}

	if err := tx.Commit(ctx); err != nil {
		return types.InternalServerErrorResponse(), err
	}

	return nil, nil
}

// GetUserById retrieves a user by ID.
func (s *UserService) GetUserById(ctx context.Context, params database.GetUserByIdParams) (*database.GetUserByIdRow, *types.ErrorResponse, error) {
	tx, err := initialiseDBTX(ctx)
	if err != nil {
		return nil, types.InternalServerErrorResponse(), err
	}
	defer tx.Rollback(ctx)

	qb := database.NewQueryBuilder(tx)
	row, err := qb.GetUserById(ctx, params)
	if err != nil {
		return nil, types.BadRequestErrorResponse(err.Error()), err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, types.InternalServerErrorResponse(), err
	}

	return &row, nil, nil
}
