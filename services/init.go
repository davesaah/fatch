// package services performs 2 responsibilities:
// 1. Fetch data from database
// 2. Return data/error
package services

import (
	"context"

	"github.com/jackc/pgx/v5"
	"gitlab.com/davesaah/fatch/database"
)

// initialiseDBTX creates a new DB transaction instance for querying
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
