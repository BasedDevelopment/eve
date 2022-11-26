package db

import "context"

func GetUserHash(ctx context.Context, email string) (hash string, err error) {
	var hash string
	err := pool.QueryRow(ctx, "SELECT password FROM profile WHERE email = $1", email).Scan(&hash)
	return
}

func AddToken(ctx context.Context, email string, token_p string, token_p string) (err error) {
	_, err := pool.Exec(ctx, "INSERT INTO tokens (token_public, token_private, profile_id, expires) VALUES ($1, $2, $3)", email, token_p, token_r)
	return
}
