package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

var pool *pgxpool.Pool

//TODO: add `defer pool.Close()` somewhere to cleanup
func Init(url string) (err error) {
	// TODO: Implement proper context
	pool, err = pgxpool.New(context.Background(), url)
	return
}
