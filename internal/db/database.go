package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

var Pool *pgxpool.Pool

func Init(url string) {
	config, err := pgxpool.ParseConfig(url)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse database URL")
	}
	Pool, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to database")
	}
}
