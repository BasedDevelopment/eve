package routes

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ericzty/eve/internal/db"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/sha3"
)

var (
	ctx = context.Background()
)

func Login(w http.ResponseWriter, r *http.Request) {
	var body credentials

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Print(err)
	}

	if body.Username == "" || body.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
	}

	hash := db.GetUserHash(ctx, body.Username)

	if bcrypt.CompareHashAndPassword([]byte(users[body.Username]), []byte(body.Password)) == nil {
		p, s, salt, err := generateStrings([]int{64, 64, 32})

		if err != nil {
			log.Fatal(err)
		}

		buf := []byte(s + salt)
		secret := make([]byte, 64)
		sha3.ShakeSum256(secret, buf)

		userToken := makeToken(Token{
			Version: "v1",
			Public:  p,
			Secret:  base64.URLEncoding.EncodeToString([]byte(s)),
			Salt:    salt,
		})

		serverToken := makeToken(Token{
			Version: "v1",
			Public:  p,
			Secret:  fmt.Sprintf("%x", secret),
			Salt:    salt,
		})

		if err := rdb.Set(ctx, p, serverToken, 0).Err(); err != nil {
			log.Fatal(err)
		}

		w.Write([]byte(userToken))
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Invalid credentials"))
	}
}
