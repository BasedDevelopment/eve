package controllers

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ericzty/eve/internal/db"
	"github.com/ericzty/eve/internal/tokens"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Profile struct {
	ID        uuid.UUID
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string
	Disabled  bool
	IsAdmin   bool      `db:"is_admin"`
	LastLogin time.Time `json:"lastLogin" db:"last_login"`
	Created   time.Time `json:"created"`
	Updated   time.Time `json:"updated"`
	Remarks   string
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

func (p *Profile) Get(ctx context.Context) (profile Profile, err error) {
	row, err := db.Pool.Query(ctx, "SELECT * FROM profile WHERE id = $1", p.ID)
	if err != nil {
		return Profile{}, fmt.Errorf("%w %v", QueryErr, err)
	}
	profilerow, err := pgx.CollectRows(row, pgx.RowToStructByName[Profile])
	profile = profilerow[0]
	if err != nil {
		return Profile{}, fmt.Errorf("%w %v", CollectErr, err)
	}
	return

}
func (p *Profile) Update() {}
func (p *Profile) Delete() {}

func (p *Profile) GetHash(ctx context.Context) (hash string, err error) {
	err = db.Pool.QueryRow(ctx, "SELECT password FROM profile WHERE email = $1", p.Email).Scan(&hash)
	return
}

func IsAdmin(ctx context.Context, id string) (isAdmin bool, err error) {
	err = db.Pool.QueryRow(ctx, "SELECT is_admin FROM profile WHERE id = $1", id).Scan(&isAdmin)
	return
}

func Logout(ctx context.Context, reqToken string) (err error) {
	publicPart := tokens.Parse(reqToken)
	_, err = db.Pool.Exec(ctx, "DELETE FROM token WHERE token_public = $1", publicPart)

	return
}
