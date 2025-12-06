// Package database defines the database connections and methods to manipulate
// and retrieve data
package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"gitlab.com/davesaah/fatch/internal/config"
)

func NewPool(ctx context.Context) (*pgxpool.Pool, error) {
	config, err := config.LoadDBConfig()
	if err != nil {
		return nil, err
	}

	dburl := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?search_path=%s",
		config.User, config.Password, config.Host, config.Port, config.DBName, config.Schema,
	)

	poolConfig, err := pgxpool.ParseConfig(dburl)
	if err != nil {
		return nil, err
	}

	// Pool tuning
	poolConfig.MaxConns = 20
	poolConfig.MinConns = 5
	poolConfig.MaxConnLifetime = time.Hour
	poolConfig.MaxConnIdleTime = 30 * time.Minute

	return pgxpool.NewWithConfig(ctx, poolConfig)
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
