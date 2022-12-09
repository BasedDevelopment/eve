package sessions

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"time"

	"github.com/ericzty/eve/internal/controllers"
	"github.com/ericzty/eve/internal/tokens"
	"golang.org/x/crypto/sha3"
)

func prngString() string {
	b := make([]byte, 10)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatalln("error:", err)
		return ""
	}

	return fmt.Sprintf("%x", b)
}

func generateStrings(bits []int) (a, b, c string, err error) {
	if a = prngString(); err != nil {
		return "", "", "", err
	}

	if b = prngString(); err != nil {
		return "", "", "", err
	}

	if c = prngString(); err != nil {
		return "", "", "", err
	}

	return a, b, c, err
}

// New creates a new authentication session in the database (login)
func NewSession(ctx context.Context, user controllers.Profile) (tokens.Token, error) {
	// Generate three pseudo-random numbers (in a string)
	public, secret, salt, err := generateStrings([]int{64, 64, 32})

	if err != nil {
		return tokens.Token{}, err
	}

	// Salt the token
	buf := []byte(secret + salt) // Append the salt to the secret
	saltedSecret := make([]byte, 64)
	sha3.ShakeSum256(saltedSecret, buf) // Hash the string with the combined secret and salt

	// This is what we store in the database
	session := Session{
		Owner: user.ID,
		Token: tokens.Token{
			Version: "v1",
			Public:  public,
			Secret:  fmt.Sprintf("%x", saltedSecret),
			Salt:    salt,
		},
		Created: time.Now(),
		Expires: time.Now().Add(24 * time.Hour), // expires in 1 day
	}

	// Push the session to the database
	if err := session.push(ctx); err != nil {
		return tokens.Token{}, err
	}

	// This is the token we want to give to the user
	return tokens.Token{
		Version: "v1",
		Public:  public,
		Secret:  base64.URLEncoding.EncodeToString([]byte(secret)),
		Salt:    salt,
	}, nil
}
