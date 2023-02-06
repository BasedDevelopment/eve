//go:build !integration
// +build !integration

package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVerFromDec(t *testing.T) {
	assert.Equal(t, "1.0.0", VerFromDec(1_000_000))
	assert.Equal(t, "1.0.1", VerFromDec(1_000_001))
	assert.Equal(t, "1.1.0", VerFromDec(1_001_000))
	assert.Equal(t, "1.1.1", VerFromDec(1_001_001))
	assert.Equal(t, "999.999.999", VerFromDec(999_999_999))
}

func TestDivmod(t *testing.T) {
	func() {
		x, y := Divmod(1, 1)
		assert.Equal(t, 1, x)
		assert.Equal(t, 0, y)
	}()
	func() {
		x, y := Divmod(1, 2)
		assert.Equal(t, 0, x)
		assert.Equal(t, 1, y)
	}()
}
