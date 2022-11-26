package authentication

import (
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/sha3"
)

func GenerateToken() (userToken string, serverToken string, p string, err error) {
	p, s, salt, err := generateStrings([]int{64, 64, 32})

	buf := []byte(s + salt)
	secret := make([]byte, 64)
	sha3.ShakeSum256(secret, buf)

	userToken = makeToken(Token{
		Version: "v1",
		Public:  p,
		Secret:  base64.URLEncoding.EncodeToString([]byte(s)),
		Salt:    salt,
	})

	serverToken = makeToken(Token{
		Version: "v1",
		Public:  p,
		Secret:  fmt.Sprintf("%x", secret),
		Salt:    salt,
	})

	return
}
