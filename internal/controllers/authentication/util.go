package authentication

import (
	"crypto/rand"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func hashAndSalt(pwd []byte) string {
	// Use GenerateFromPassword to hash & salt pwd.
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Fatalln(err)
	} // GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash)
}

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
