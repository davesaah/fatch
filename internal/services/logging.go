package services

import (
	"context"

	"github.com/davesaah/fatch/internal/database"
)

func (s *Service) Log(ctx context.Context, arg *database.Log) error {
	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qb := database.NewQueryBuilder(tx)
	err = qb.InsertLog(ctx, arg)
	if err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}
