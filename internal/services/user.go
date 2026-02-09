package services

import (
	"context"

	"github.com/davesaah/fatch/internal/database"
	"github.com/davesaah/fatch/types"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

func (s *Service) CreateUser(
	ctx context.Context, params database.RegisterUserParams,
) (int, *types.ErrorResponse, error) {
	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return 0, types.InternalServerErrorResponse(), err
	}
	defer tx.Rollback(ctx)

	qb := database.NewQueryBuilder(tx)
	otp, err := qb.CreateUser(ctx, params)
	if err != nil {
		pgErr := err.(*pgconn.PgError)
		return 0, types.ConflictErrorResponse(pgErr.Message), err
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, types.InternalServerErrorResponse(), err
	}

	return otp, nil, nil
}

func (s *Service) GetUserByID(
	ctx context.Context, userID pgtype.UUID,
) (*database.GetUserByIDRow, *types.ErrorResponse, error) {
	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return nil, types.InternalServerErrorResponse(), err
	}
	defer tx.Rollback(ctx)

	qb := database.NewQueryBuilder(tx)
	row, err := qb.GetUserByID(ctx, userID)
	if err != nil {
		return nil, types.BadRequestErrorResponse(err.Error()), err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, types.InternalServerErrorResponse(), err
	}

	return row, nil, nil
}

func (s *Service) VerifyUser(
	ctx context.Context, params database.VerifyUserParams,
) (*types.ErrorResponse, error) {
	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return types.InternalServerErrorResponse(), err
	}
	defer tx.Rollback(ctx)

	qb := database.NewQueryBuilder(tx)
	if err := qb.VerifyUser(ctx, params); err != nil {
		pgErr := err.(*pgconn.PgError)
		return types.ConflictErrorResponse(pgErr.Message), err
	}

	if err := tx.Commit(ctx); err != nil {
		return types.InternalServerErrorResponse(), err
	}

	return nil, nil
}
