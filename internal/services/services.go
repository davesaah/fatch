package services

import (
	"context"

	"github.com/davidreturns08/fatch/internal/database"
	"github.com/davidreturns08/fatch/internal/types"
	"github.com/jackc/pgx/v5"
)

func initialiseDBTX(ctx context.Context) (pgx.Tx, *types.ErrorResponse) {
	conn, err := database.NewConnection(ctx)
	if err != nil {
		return nil, types.InternalServerErrorResponse()
	}

	tx, err := conn.Begin(ctx)
	if err != nil {
		return nil, types.InternalServerErrorResponse()
	}

	return tx, nil
}
