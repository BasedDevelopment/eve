//go:build !integration
// +build !integration

package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestService(t *testing.T) {
	assert := assert.New(t)
	assert.HTTPSuccess(Service().ServeHTTP, "GET", "/", nil)
}
