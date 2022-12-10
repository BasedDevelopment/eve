package sessions

import (
	"context"
	"crypto/subtle"
	"encoding/base64"
	"fmt"

	"github.com/ericzty/eve/internal/tokens"
	"golang.org/x/crypto/sha3"
)

// ValidateSession takes a token and finds its session. Returns true if valid, false if anything else
func ValidateSession(ctx context.Context, incomingToken tokens.Token) bool {
	// Get the session from the database
	session, err := GetSession(ctx, incomingToken)

	if err != nil {
		fmt.Println(err)
		return false // Error fetching session, almost definitely unauthenticated
	}

	//# Prepare the incoming secret for comparison
	// Decode incoming token from base64
	decodedSecret, decodeErr := base64.URLEncoding.DecodeString(incomingToken.Secret)

	if decodeErr != nil {
		return false // Error while decoding secret from b64, assume unauthenticated
	}

	// Salt and Hash the token
	buf := []byte(string(decodedSecret) + incomingToken.Salt) // Append the salt to the secret
	saltedSecret := make([]byte, 64)
	sha3.ShakeSum256(saltedSecret, buf) // Hash the string with the combined secret and salt

	// Compare the two secrets
	if subtle.ConstantTimeCompare(
		[]byte(session.Secret),                  // secret from the database (already in hex)
		[]byte(fmt.Sprintf("%x", saltedSecret)), // secret from the request (now salted & hashed, and converted to hex)
	) != 1 {
		return false // Invalid Token, unauthenticated
	}

	// Check expiry
	if session.isExpired() {
		return false // Expired token, unauthenticated
	}

	return true // Passed all checks, authenticated
}
