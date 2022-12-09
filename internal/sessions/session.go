package sessions

import (
	"context"
	"time"

	"github.com/ericzty/eve/internal/db"
	"github.com/ericzty/eve/internal/tokens"
	"github.com/google/uuid"
)

type Session struct {
	Owner   uuid.UUID // TODO: Make this a reference?
	Token   tokens.Token
	Created time.Time
	Expires time.Time
}

// push pushes a Session to the database
func (s Session) push(ctx context.Context) error {
	// Object in database is same as Session type but with
	// Split up token into individual Version, Public, Secret, Salt cols

	_, err := db.Pool.Exec(
		ctx,
		"INSERT INTO token (owner, token_version, token_public, token_secret, token_salt, created_at, expires) VALUES ($1, $2, $3, $4)",
		s.Owner,         // owner
		s.Token.Version, // token_version
		s.Token.Public,  // token_public
		s.Token.Secret,  // token_secret
		s.Token.Salt,    // token_salt
		s.Created,       // created_at
		s.Expires,       // expires
	)

	return err
}

func getSession(ctx context.Context, token tokens.Token) (Session, error) {
	var session Session

	err := db.Pool.QueryRow(ctx, "SELECT * FROM token WHERE token_public = $1", token.Public).Scan(&session)

	if err != nil {
		return Session{}, nil
	}

	return session, nil
}

func (s Session) isExpired() bool {
	return s.Expires.After(time.Now())
}
