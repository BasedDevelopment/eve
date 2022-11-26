package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

// TODO: add `defer pool.Close()` somewhere to cleanup
func Init(url string) (err error) {
	// TODO: Implement proper context
	Pool, err = pgxpool.New(context.Background(), url)
	return
}
