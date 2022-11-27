package authentication

import (
	"context"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

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

var TokenErr = errors.New("Token Err:")
var TokenExpiredErr = errors.New("Token Expired Err:")
var ServerTokenErr = errors.New("Parsing Server Token Err:")

func VerifyToken(ctx context.Context, token string) (string, error) {

	userToken, parseTokenErr := parseToken(token)
	if parseTokenErr != nil {
		return "", fmt.Errorf("%w Error parsing user token: %v", TokenErr, parseTokenErr) // Invalid Token
	}

	id, serverTokenDB, expirey, dbErr := getToken(ctx, userToken.Public)
	if dbErr != nil {
		return "", fmt.Errorf("%w Error getting token from database: %v", TokenErr, dbErr) // Database Error
	}

	serverToken, parseServerTokenErr := parseToken(serverTokenDB)
	if parseServerTokenErr != nil {
		return "", fmt.Errorf("%w Error parsing server token: %v", ServerTokenErr) // Internal Server Err
	}

	// Decode token
	unb64dSecret, decodeErr := base64.URLEncoding.DecodeString(userToken.Secret)
	if decodeErr != nil {
		return "", fmt.Errorf("%w Error decoding token: %v", TokenErr, decodeErr) //Invalid Token
	}

	buf := []byte(string(unb64dSecret) + userToken.Salt)
	secret := make([]byte, 64)
	sha3.ShakeSum256(secret, buf)

	// Verify token
	if subtle.ConstantTimeCompare(
		[]byte(serverToken.Secret),        // secret from the database (already in hex)
		[]byte(fmt.Sprintf("%x", secret)), // secret from the request (now salted & hashed, and converted to hex)
	) != 1 {
		return "", fmt.Errorf("%w Token incorrect", TokenErr) // Invalid Token
	}

	// Check expirey
	if expirey.Before(time.Now()) {
		return "", fmt.Errorf("%w Token expired", TokenExpiredErr) // Invalid Token
	}

	return id, nil
}

func GetPublicPart(token string) (publicPart string, err error) {
	userToken, err := parseToken(token)
	return userToken.Public, err
}
