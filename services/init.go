// package services performs 2 responsibilities:
// 1. Fetch data from database
// 2. Return data/error
package services

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"gitlab.com/davesaah/fatch/database"
)

// initialiseDBTX creates a new DB transaction instance for querying
func initialiseDBTX(ctx context.Context) (pgx.Tx, *pgxpool.Pool, error) {
	conn, err := database.NewConnection(ctx)
	if err != nil {
		return nil, nil, err
	}

	tx, err := conn.Begin(ctx)
	if err != nil {
		return nil, nil, err
	}

	return tx, conn, nil
}
