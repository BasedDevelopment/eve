package controllers

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ericzty/eve/internal/controllers/authentication"
	"github.com/ericzty/eve/internal/db"
	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v5"
)

type Profile struct {
	ID        uuid.UUID
	Name      string
	Email     string
	Password  string
	Disabled  bool
	IsAdmin   bool      `db:"is_admin"`
	LastLogin time.Time `db:"last_login"`
	Created   time.Time
	Updated   time.Time
	Remarks   pgtype.Text
}

func (p *Profile) New(ctx context.Context) (id string, err error) {
	// Generate UUID
	p.ID = uuid.New()
	id = p.ID.String()

	_, err = db.Pool.Exec(ctx, "INSERT INTO profile (id, name, email, password, disabled, is_admin, remarks) VALUES ($1, $2, $3, $4, $5, $6, $7)", id, p.Name, p.Email, p.Password, p.Disabled, p.IsAdmin, p.Remarks)
	return
}

var QueryErr = errors.New("Query error:")
var CollectErr = errors.New("Collect error:")

func (p *Profile) Get(ctx context.Context) (err error) {
	row, err := db.Pool.Query(ctx, "SELECT * FROM profile WHERE id = $1", p.ID)
	if err != nil {
		return fmt.Errorf("%w %v", QueryErr, err)
	}
	profile, err := pgx.CollectRows(row, pgx.RowToStructByName[Profile])
	if err != nil {
		return fmt.Errorf("%w %v", CollectErr, err)
	}
	fmt.Println(profile)
	return nil

}
func (p *Profile) Update() {}
func (p *Profile) Delete() {}

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

	// Set expiry
	expiry := time.Now().Add(24 * time.Hour)

	// Get user ID
	var id uuid.UUID
	if err := db.Pool.QueryRow(ctx, "SELECT id FROM profile WHERE email = $1", p.Email).Scan(&id); err != nil {
		return "", err
	}

	// Set Last Login
	time := time.Now()
	_, err = db.Pool.Exec(ctx, "UPDATE profile SET last_login = $1 WHERE id = $2", time, id)
	if err != nil {
		return "", err
	}

	// Store token
	_, err = db.Pool.Exec(ctx, "INSERT INTO token (token_public, token_private, profile_id, expires) VALUES ($1, $2, $3, $4)", publicPart, serverToken, id.String(), expiry)
	return
}

func IsAdmin(ctx context.Context, id string) (isAdmin bool, err error) {
	err = db.Pool.QueryRow(ctx, "SELECT is_admin FROM profile WHERE id = $1", id).Scan(&isAdmin)
	return
}

func Logout(ctx context.Context, reqToken string) (err error) {
	publicPart, err := authentication.GetPublicPart(reqToken)
	if err != nil {
		return
	}
	_, err = db.Pool.Exec(ctx, "DELETE FROM token WHERE token_public = $1", publicPart)
	return
}
