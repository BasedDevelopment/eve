//go:build integration
// +build integration

package test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
}

const (
	dbUrl = "postgres://postgres:password@localhost:5432/postgres"
	host  = "http://localhost:3000"

	adminName         = "Admin Test"
	adminEmail        = "admin@testing.com"
	adminPassword     = "adminPasswordTest"
	adminPasswordHash = "$2a$11$MCwHsWkdVATPJ0URdUFg9uvY6UdskKO.Mwc3Y2e9LKi.5GQFOhTCq"

	userName     = "User Test"
	userEmail    = "eric@testing.com"
	userPassword = "userPasswordTest"
)

var (
	pool *pgxpool.Pool

	adminId uuid.UUID
	userId  uuid.UUID
)

func (ts *TestSuite) SetupSuite() {
	adminId = uuid.New()
	userId = uuid.New()

	ctx := context.Background()

	config, err := pgxpool.ParseConfig(dbUrl)
	assert.NoError(ts.T(), err)

	pool, err = pgxpool.NewWithConfig(ctx, config)
	assert.NoError(ts.T(), err)

	// Insert test user and hypervisor
	_, err = pool.Exec(
		ctx,
		"INSERT INTO profile (id, name, email, password, is_admin) VALUES ($1, $2, $3, $4, $5)",
		adminId, adminName, adminEmail, adminPasswordHash, true,
	)
	assert.NoError(ts.T(), err)
}

func (ts *TestSuite) TearDownSuite() {
	ctx := context.Background()
	_, err := pool.Exec(ctx, "DELETE FROM sessions WHERE owner = $1", adminId)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to delete test admin sessions")
	}
	_, err = pool.Exec(ctx, "DELETE FROM profile WHERE id = $1", adminId)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to delete test admin")
	}
	_, err = pool.Exec(ctx, "DELETE FROM sessions WHERE owner = $1", userId)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to delete test user sessions")
	}
	_, err = pool.Exec(ctx, "DELETE FROM profile WHERE id = $1", userId)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to delete test user")
	}
	pool.Close()
}

func TestInit(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
