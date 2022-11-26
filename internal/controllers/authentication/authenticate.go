package authentication

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"

	"github.com/ericzty/eve/internal/controllers"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/sha3"
)

func Authenticate(email, password string) (public, internal Token, err error) {
	profile := controllers.Profile{}
	hash, err := profile.GetHash(context.Background())

	if err != nil {
		return Token{}, Token{}, err
	}

	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil {
		p, s, salt, err := generateStrings([]int{64, 64, 32})

		if err != nil {
			log.Fatal(err)
		}

		buf := []byte(s + salt)
		secret := make([]byte, 64)
		sha3.ShakeSum256(secret, buf)

		userToken := Token{
			Version: "v1",
			Public:  p,
			Secret:  base64.URLEncoding.EncodeToString([]byte(s)),
			Salt:    salt,
		}

		serverToken := Token{
			Version: "v1",
			Public:  p,
			Secret:  fmt.Sprintf("%x", secret),
			Salt:    salt,
		}

		return userToken, serverToken, nil
	}

	return Token{}, Token{}, nil
}
