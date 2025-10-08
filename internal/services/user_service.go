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
func (s *UserService) CreateUser(ctx context.Context, params database.CreateUserParams) *types.ErrorResponse {
	tx, errResponse := initialiseDBTX(ctx)
	if errResponse != nil {
		return errResponse
	}
	defer tx.Rollback(ctx)

	qb := database.NewQueryBuilder(tx)
	if err := qb.CreateUser(ctx, params); err != nil {
		pgErr := err.(*pgconn.PgError)
		return types.BadRequestErrorResponse(pgErr.Message)
	}

	if err := tx.Commit(ctx); err != nil {
		return types.InternalServerErrorResponse()
	}

	return nil
}

// GetUserById retrieves a user by ID.
func (s *UserService) GetUserById(ctx context.Context, params database.GetUserByIdParams) (*database.GetUserByIdRow, *types.ErrorResponse) {
	tx, errResponse := initialiseDBTX(ctx)
	if errResponse != nil {
		return nil, errResponse
	}
	defer tx.Rollback(ctx)

	qb := database.NewQueryBuilder(tx)
	row, err := qb.GetUserById(ctx, params)
	if err != nil {
		// pgErr := err.(*pgconn.PgError)
		return nil, types.BadRequestErrorResponse(err.Error())
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, types.InternalServerErrorResponse()
	}

	return &row, nil
}
