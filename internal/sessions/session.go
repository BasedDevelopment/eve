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
