package services

import (
	"context"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"gitlab.com/davesaah/fatch/database"
	"gitlab.com/davesaah/fatch/types"
)

type AuthService struct{}

func (s *AuthService) Login(
	ctx context.Context, params database.LoginParams,
) (*pgtype.UUID, *types.ErrorResponse, error) {
	tx, err := initialiseDBTX(ctx)
	if err != nil {
		return nil, types.InternalServerErrorResponse(), err
	}
	defer tx.Rollback(ctx)

	qb := database.NewQueryBuilder(tx)

	userID, err := qb.VerifyPassword(ctx, params)
	if err != nil {
		pgErr := err.(*pgconn.PgError)
		return nil, types.BadRequestErrorResponse(pgErr.Message), err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, types.InternalServerErrorResponse(), err
	}

	return userID, nil, nil
}

func (s *AuthService) ChangePassword(
	ctx context.Context, params database.ChangePasswordParams,
) (*types.ErrorResponse, error) {

	tx, err := initialiseDBTX(ctx)
	if err != nil {
		return types.InternalServerErrorResponse(), err
	}
	defer tx.Rollback(ctx)

	qb := database.NewQueryBuilder(tx)

	err = qb.ChangePassword(ctx, params)
	if err != nil {
		pgErr := err.(*pgconn.PgError)
		return types.BadRequestErrorResponse(pgErr.Message), err
	}

	if err := tx.Commit(ctx); err != nil {
		return types.InternalServerErrorResponse(), err
	}

	return nil, nil
}
