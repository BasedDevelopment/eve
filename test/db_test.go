//go:build integration
// +build integration

package test

import (
	"context"

	"github.com/stretchr/testify/assert"
)

func (ts *TestSuite) TestDBPing() {
	assert.NoError(ts.T(), pool.Ping(context.Background()))
}
