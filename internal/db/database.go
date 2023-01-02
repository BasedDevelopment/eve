/*
 * eve - management toolkit for libvirt servers
 * Copyright (C) 2022-2023  BNS Services LLC

 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.

 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.

 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

// The database connection pool that will be used by the rest of the app
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
