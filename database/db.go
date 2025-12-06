// Package database defines the database connections and methods to manipulate
// and retrieve data
package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"gitlab.com/davesaah/fatch/config"
)

func NewConnection(ctx context.Context) (*pgx.Conn, error) {
	config, err := config.LoadDBConfig()
	if err != nil {
		return nil, err
	}

	dburl := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?search_path=%s",
		config.User, config.Password, config.Host, config.Port, config.DBName, config.Schema,
	)
	return pgx.Connect(ctx, dburl)
}

// DBTX is an interface that abstracts both *pgx.Conn and pgx.Tx for executing queries.
type DBTX interface {
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
	Query(context.Context, string, ...any) (pgx.Rows, error)
	QueryRow(context.Context, string, ...any) pgx.Row
}

// Queries provides methods to execute SQL queries and commands.
type Queries struct {
	db DBTX
}

// NewQueryBuilder creates a new Queries instance with the given DBTX (database transaction)
func NewQueryBuilder(db DBTX) *Queries {
	return &Queries{db: db}
}
