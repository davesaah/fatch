package services

import (
	"context"

	"github.com/davidreturns08/fatch/internal/database"
	"github.com/davidreturns08/fatch/internal/types"
	"github.com/jackc/pgx/v5/pgconn"
)

type CurrencyService struct{}

func (cs *CurrencyService) GetCurrencyByID(ctx context.Context, id int) (*database.GetCurrencyByIdRow, *types.ErrorResponse, error) {
	tx, err := initialiseDBTX(ctx)
	if err != nil {
		return nil, types.InternalServerErrorResponse(), err
	}
	defer tx.Rollback(ctx)

	qb := database.NewQueryBuilder(tx)
	row, err := qb.GetCurrencyByID(ctx, id)
	if err != nil {
		pgErr := err.(*pgconn.PgError)
		return nil, types.BadRequestErrorResponse(pgErr.Message), err
	}

	return &row, nil, nil
}

func (cs *CurrencyService) GetAllCurrencies(ctx context.Context) ([]database.GetAllCurrenciesRow, *types.ErrorResponse, error) {
	tx, err := initialiseDBTX(ctx)
	if err != nil {
		return nil, types.InternalServerErrorResponse(), err
	}
	defer tx.Rollback(ctx)

	qb := database.NewQueryBuilder(tx)
	rows, err := qb.GetAllCurrencies(ctx)
	if err != nil {
		pgErr := err.(*pgconn.PgError)
		return nil, types.BadRequestErrorResponse(pgErr.Message), err
	}

	return rows, nil, nil
}
