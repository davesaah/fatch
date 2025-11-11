package services

import (
	"context"

	"github.com/jackc/pgx/v5/pgconn"
	"gitlab.com/davesaah/fatch/database"
	"gitlab.com/davesaah/fatch/types"
)

type CurrencyService struct{}

func (cs *CurrencyService) GetCurrencyByID(ctx context.Context, id int) (*database.GetCurrencyByIDRow, *types.ErrorResponse, error) {
	tx, err := initialiseDBTX(ctx)
	if err != nil {
		return nil, types.InternalServerErrorResponse(), err
	}
	defer tx.Rollback(ctx)

	qb := database.NewQueryBuilder(tx)
	row, err := qb.GetCurrencyByID(ctx, id)
	if err != nil {
		pgErr := err.(*pgconn.PgError)
		return nil, types.NotFoundErrorResponse(pgErr.Message), err
	}

	return row, nil, nil
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
		return nil, types.NotFoundErrorResponse(pgErr.Message), err
	}

	return rows, nil, nil
}
