package controllers

import (
	"context"
	"time"

	"github.com/ericzty/eve/internal/controllers/authentication"
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

func (p *Profile) GetHash(ctx context.Context) (hash string, err error) {
	err = db.Pool.QueryRow(ctx, "SELECT password FROM profile WHERE email = $1", p.Email).Scan(&hash)
	return
}

func (p *Profile) IssueToken(ctx context.Context) (userToken string, err error) {
	// Generate Token
	userToken, serverToken, publicPart, err := authentication.GenerateToken()
	if err != nil {
		return "", err
	}

	// Set expirey
	expirey := time.Now().Add(24 * time.Hour)

	// Get user ID
	var id uuid.UUID
	if err := db.Pool.QueryRow(ctx, "SELECT id FROM profile WHERE email = $1", p.Email).Scan(&id); err != nil {
		return "", err
	}

	// Store token
	_, err = db.Pool.Exec(ctx, "INSERT INTO token (token_public, token_private, profile_id, expires) VALUES ($1, $2, $3, $4)", publicPart, serverToken, id.String(), expirey)
	return
}

func IsAdmin(ctx context.Context, id string) (isAdmin bool, err error) {
	err = db.Pool.QueryRow(ctx, "SELECT is_admin FROM profile WHERE id = $1", id).Scan(&isAdmin)
	return
}
