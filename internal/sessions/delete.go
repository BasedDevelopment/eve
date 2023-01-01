package sessions

import (
	"context"

	"github.com/BasedDevelopment/eve/internal/db"
	"github.com/BasedDevelopment/eve/internal/tokens"
)

// Delete removes a session from the database (logout)
func Delete(ctx context.Context, token tokens.Token) error {
	_, err := db.Pool.Exec(ctx, "DELETE FROM sessions WHERE token_public = $1", token.Public)

	return err
}
