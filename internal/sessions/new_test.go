//go:build !integration
// +build !integration

package sessions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrngString(t *testing.T) {
	str1 := prngString(32)
	str2 := prngString(32)
	str3 := prngString(64)
	str4 := prngString(64)

	assert.NotEqual(t, str1, str2)
	assert.NotEqual(t, str3, str4)

	assert.Len(t, str1, 32)
	assert.Len(t, str2, 32)
	assert.Len(t, str3, 64)
	assert.Len(t, str4, 64)
}
