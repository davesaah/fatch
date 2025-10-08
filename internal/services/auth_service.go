package services

import (
	"context"

	"github.com/davidreturns08/fatch/internal/database"
	"github.com/davidreturns08/fatch/internal/types"
	"github.com/jackc/pgx/v5/pgconn"
)

// AuthService provides authentication-related services.
type AuthService struct{}

// VerifyPassword verifies the password of a user.
func (s *AuthService) VerifyPassword(ctx context.Context, params database.VerifyPasswordParams) (*database.VerifyPasswordRow, *types.ErrorResponse) {

	tx, errResponse := initialiseDBTX(ctx)
	if errResponse != nil {
		return nil, errResponse
	}
	defer tx.Rollback(ctx)

	qb := database.NewQueryBuilder(tx)

	row, err := qb.VerifyPassword(ctx, params)
	if err != nil {
		pgErr := err.(*pgconn.PgError)
		return nil, types.BadRequestErrorResponse(pgErr.Message)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, types.InternalServerErrorResponse()
	}

	return &row, nil
}

// ChangePassword changes the password of a user.
func (s *AuthService) ChangePassword(ctx context.Context, params database.ChangePasswordParams) *types.ErrorResponse {

	tx, errResponse := initialiseDBTX(ctx)
	if errResponse != nil {
		return errResponse
	}
	defer tx.Rollback(ctx)

	qb := database.NewQueryBuilder(tx)

	err := qb.ChangePassword(ctx, params)
	if err != nil {
		pgErr := err.(*pgconn.PgError)
		return types.BadRequestErrorResponse(pgErr.Message)
	}

	if err := tx.Commit(ctx); err != nil {
		return types.InternalServerErrorResponse()
	}

	return nil
}
