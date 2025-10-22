package services

import (
	"context"

	"github.com/davidreturns08/fatch/internal/database"
	"github.com/davidreturns08/fatch/internal/types"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

type AccountService struct{}

func (a *AccountService) CreateAccount(ctx context.Context, params database.CreateAccountParams) (*database.GetAllUserAccountsRow, *types.ErrorResponse, error) {
	tx, err := initialiseDBTX(ctx)
	if err != nil {
		return nil, types.InternalServerErrorResponse(), err
	}
	defer tx.Rollback(ctx)

	qb := database.NewQueryBuilder(tx)
	row, err := qb.CreateAccount(ctx, params)
	if err != nil {
		pgErr := err.(*pgconn.PgError)
		return nil, types.ConflictErrorResponse(pgErr.Message), err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, types.InternalServerErrorResponse(), err
	}

	return row, nil, nil
}

func (a *AccountService) GetAccountDetails(ctx context.Context, params database.GetAccountDetailsParams) (*database.GetAccountDetailsRow, *types.ErrorResponse, error) {
	tx, err := initialiseDBTX(ctx)
	if err != nil {
		return nil, types.InternalServerErrorResponse(), err
	}
	defer tx.Rollback(ctx)

	qb := database.NewQueryBuilder(tx)
	row, err := qb.GetAccountDetails(ctx, params)
	if err != nil {
		pgErr := err.(*pgconn.PgError)
		return nil, types.NotFoundErrorResponse(pgErr.Message), err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, types.InternalServerErrorResponse(), err
	}

	return row, nil, nil
}

func (a *AccountService) GetAllUserAccounts(ctx context.Context, userID pgtype.UUID) ([]database.GetAllUserAccountsRow, *types.ErrorResponse, error) {
	tx, err := initialiseDBTX(ctx)
	if err != nil {
		return nil, types.InternalServerErrorResponse(), err
	}
	defer tx.Rollback(ctx)

	qb := database.NewQueryBuilder(tx)
	rows, err := qb.GetAllUserAccounts(ctx, userID)
	if err != nil {
		pgErr := err.(*pgconn.PgError)
		return nil, types.NotFoundErrorResponse(pgErr.Message), err
	}

	return rows, nil, nil
}

func (a *AccountService) ArchiveAccount(ctx context.Context, params database.ArchiveAccountByIDParams) (*types.ErrorResponse, error) {
	tx, err := initialiseDBTX(ctx)
	if err != nil {
		return types.InternalServerErrorResponse(), err
	}
	defer tx.Rollback(ctx)

	qb := database.NewQueryBuilder(tx)
	err = qb.ArchiveAccountByID(ctx, params)
	if err != nil {
		pgErr := err.(*pgconn.PgError)
		return types.BadRequestErrorResponse(pgErr.Message), err
	}

	if err := tx.Commit(ctx); err != nil {
		return types.InternalServerErrorResponse(), err
	}

	return nil, nil
}
