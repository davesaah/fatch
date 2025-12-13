package database

import (
	"context"
	"time"
)

const insertLog = "SELECT insert_log($1::timestamptz, $2::text, $3::text, $4::jsonb)"

type Log struct {
	Timestamp time.Time
	Level     string
	Service   string
	LogData   map[string]any
}

func (q *Queries) InsertLog(ctx context.Context, arg *Log) error {
	_, err := q.db.Exec(ctx, insertLog, arg.Timestamp, arg.Level, arg.Service, arg.LogData)
	return err
}
