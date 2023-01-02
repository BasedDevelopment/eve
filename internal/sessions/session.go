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

package sessions

import (
	"context"
	"time"

	"github.com/BasedDevelopment/eve/internal/db"
	"github.com/BasedDevelopment/eve/internal/tokens"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
)

type Session struct {
	Owner   uuid.UUID `db:"owner"`
	Version string    `db:"token_version"`
	Public  string    `db:"token_public"`
	Secret  string    `db:"token_secret"`
	Salt    string    `db:"token_salt"`
	Created time.Time `db:"created"`
	Expires time.Time `db:"expires"`
}

// push pushes a Session to the database
func (s Session) push(ctx context.Context) error {
	// Object in database is same as Session type but with
	// Split up token into individual Version, Public, Secret, Salt cols

	_, err := db.Pool.Exec(
		ctx,
		"INSERT INTO sessions (owner, token_version, token_public, token_secret, token_salt, created, expires) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		s.Owner,   // owner
		s.Version, // token_version
		s.Public,  // token_public
		s.Secret,  // token_secret
		s.Salt,    // token_salt
		s.Created, // created_at
		s.Expires, // expires
	)

	return err
}

func GetSession(ctx context.Context, token tokens.Token) (Session, error) {
	var session Session

	// Query pgx.Rows from the database.
	rows, _ := db.Pool.Query(ctx, `SELECT * FROM sessions WHERE token_public=$1`, token.Public)

	// Scan rows into session
	if err := pgxscan.ScanOne(&session, rows); err != nil {
		return Session{}, err
	}

	return session, nil
}

func (s Session) isExpired() bool {
	return time.Now().After(s.Expires)
}
