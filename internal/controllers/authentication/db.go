package authentication

import (
	"context"
	"time"

	"github.com/ericzty/eve/internal/db"
)

func getToken(ctx context.Context, publicPart string) (id string, serverToken string, expirey time.Time, err error) {
	err = db.Pool.QueryRow(ctx, "SELECT profile_id, token_private, expires FROM token WHERE token_public = $1", publicPart).Scan(&id, &serverToken, &expirey)
	return
}
