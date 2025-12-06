// package services performs 2 responsibilities:
// 1. Fetch data from database
// 2. Return data/error
package services

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	DB *pgxpool.Pool
}

func NewService(db *pgxpool.Pool) *Service {
	return &Service{DB: db}
}
