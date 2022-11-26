package controllers

import (
	"context"
	"time"

	"github.com/ericzty/eve/internal/db"
	"github.com/google/uuid"
)

type Profile struct {
	ID        uuid.UUID
	Name      string
	Email     string
	Password  string
	Disabled  bool
	LastLogin time.Time `db:"last_login"`
	Created   time.Time
	Updated   time.Time
	Remarks   string
}

func (p *Profile) New()             {}
func (p *Profile) Get(id uuid.UUID) {}
func (p *Profile) Update()          {}
func (p *Profile) Delete()          {}
func (p *Profile) GetHash(ctx context.Context) (string, error) {
	var hash string

	return "", db.Pool.QueryRow(ctx, "SELECT password FROM profile WHERE email = $1", p.Email).Scan(&hash)

}

// func (p *Profile) AddToken(ctx context.Context, token) error {
// 	_, err := db.Pool.Exec(
// 		ctx,
// 		"INSERT INTO tokens (token_public, token_private, profile_id, expires) VALUES ($1, $2, $3)",
// 		p
// 	)

// 	return err
// }
