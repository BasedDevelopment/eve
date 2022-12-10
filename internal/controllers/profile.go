package controllers

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ericzty/eve/internal/db"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
)

type Profile struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Email     string    `json:"email" db:"email"`
	Password  string    `db:"password"`
	Disabled  bool      `db:"disabled"`
	IsAdmin   bool      `db:"is_admin"`
	LastLogin time.Time `json:"last_login" db:"last_login"`
	Created   time.Time `json:"created" db:"created"`
	Updated   time.Time `json:"updated" db:"updated"`
	Remarks   string    `db:"remarks"`
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
	rows, err := db.Pool.Query(ctx, "SELECT * FROM profile WHERE id = $1", p.ID)

	if err != nil {
		return Profile{}, fmt.Errorf("%w %v", QueryErr, err)
	}

	// Scan rows into session
	if err := pgxscan.ScanOne(&profile, rows); err != nil {
		return Profile{}, err
	}

	return profile, nil
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
