package services

import (
	"context"

	"gitlab.com/davesaah/fatch/internal/database"
	"github.com/jackc/pgx/v5"
)

func initialiseDBTX(ctx context.Context) (pgx.Tx, error) {
	conn, err := database.NewConnection(ctx)
	if err != nil {
		return nil, err
	}

	tx, err := conn.Begin(ctx)
	if err != nil {
		return nil, err
	}

	return tx, nil
}
