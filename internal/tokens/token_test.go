//go:build !integration
// +build !integration

package tokens

import (
	"crypto/rand"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/sha3"
)

const (
	tokenVer = "v1"
	pubSize  = 64
	secSize  = 64
	saltSize = 32
)

var (
	tokenPubBytes  = make([]byte, pubSize/2)
	tokenSecBytes  = make([]byte, secSize/2)
	tokenSaltBytes = make([]byte, saltSize/2)

	tokenPub     string
	tokenSec     string
	tokenSecHash string
	tokenSalt    string

	tokenStruct Token
	tokenString string
)

func TestParse(t *testing.T) {
	_, err := rand.Read(tokenPubBytes)
	assert.NoError(t, err)
	_, err = rand.Read(tokenSecBytes)
	assert.NoError(t, err)
	_, err = rand.Read(tokenSaltBytes)
	assert.NoError(t, err)

	tokenPub = fmt.Sprintf("%x", tokenPubBytes)
	tokenSec = fmt.Sprintf("%x", tokenSecBytes)
	tokenSalt = fmt.Sprintf("%x", tokenSaltBytes)

	buf := []byte(tokenSec + tokenSalt)
	tokenSecHashBytes := make([]byte, 64+32)
	sha3.ShakeSum256(tokenSecHashBytes, buf)
	tokenSecHash = fmt.Sprintf("%x", tokenSecHashBytes)

	tokenString = tokenVer + "." + tokenPub + "." + tokenSecHash + "." + tokenSalt

	tokenStruct, err = Parse(tokenString)
	assert.NoError(t, err)
	assert.Equal(t, tokenVer, tokenStruct.Version)
	assert.Equal(t, tokenPub, tokenStruct.Public)
	assert.Equal(t, tokenSecHash, tokenStruct.Secret)
	assert.Equal(t, tokenSalt, tokenStruct.Salt)
}

func TestString(t *testing.T) {
	str := tokenStruct.String()
	assert.Equal(t, tokenString, str)
}
