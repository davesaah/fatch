package services

import (
	"context"

	"github.com/davesaah/fatch/internal/database"
	"github.com/davesaah/fatch/types"
	"github.com/jackc/pgx/v5/pgconn"
)

func (s *Service) GetCurrencyByID(
	ctx context.Context, id int,
) (*database.GetCurrencyByIDRow, *types.ErrorResponse, error) {
	tx, err := s.DB.Begin(ctx)
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

func (s *Service) GetAllCurrencies(
	ctx context.Context,
) ([]database.GetAllCurrenciesRow, *types.ErrorResponse, error) {
	tx, err := s.DB.Begin(ctx)
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
