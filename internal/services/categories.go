package services

import (
	"context"

	"github.com/davesaah/fatch/internal/database"
	"github.com/davesaah/fatch/types"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

func (s *Service) GetCategories(ctx context.Context, userID pgtype.UUID) ([]database.GetAllCategoriesRow, *types.ErrorResponse, error) {
	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return nil, types.InternalServerErrorResponse(), err
	}
	defer tx.Rollback(ctx)

	qb := database.NewQueryBuilder(tx)
	rows, err := qb.GetCategories(ctx, userID)
	if err != nil {
		pgErr := err.(*pgconn.PgError)
		return nil, types.NotFoundErrorResponse(pgErr.Message), err
	}

	return rows, nil, nil
}

func (s *Service) GetCategoryByID(ctx context.Context, params database.GetCategoryByIDParams) (*database.GetCategoryByIDRow, *types.ErrorResponse, error) {
	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return nil, types.InternalServerErrorResponse(), err
	}
	defer tx.Rollback(ctx)

	qb := database.NewQueryBuilder(tx)
	row, err := qb.GetCategoryByID(ctx, params)
	if err != nil {
		pgErr := err.(*pgconn.PgError)
		return nil, types.NotFoundErrorResponse(pgErr.Message), err
	}

	return row, nil, nil
}

func (s *Service) AddCategory(ctx context.Context, params database.CreateCategoryParams) (*types.ErrorResponse, error) {
	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return types.InternalServerErrorResponse(), err
	}
	defer tx.Rollback(ctx)

	qb := database.NewQueryBuilder(tx)
	if err := qb.AddCategory(ctx, params); err != nil {
		pgErr := err.(*pgconn.PgError)
		return types.BadRequestErrorResponse(pgErr.Message), err
	}

	if err := tx.Commit(ctx); err != nil {
		return types.InternalServerErrorResponse(), err
	}

	return nil, nil
}
