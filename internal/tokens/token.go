package tokens

import (
	"fmt"
	"strings"
)

const CurrentVersion = "v1"

type Token struct {
	Version string `db:"token_version"`
	Public  string `db:"token_public"`
	Secret  string `db:"token_secret"`
	Salt    string `db:"token_salt"`
}

// String converts a Token object to a string
func (t Token) String() string {
	return fmt.Sprintf("%s.%s.%s.%s", t.Version, t.Public, t.Secret, t.Salt)
}

// Parse converts a string containing a token into the Token type
func Parse(incomingToken string) Token {
	var tok Token
	toks := strings.Split(incomingToken, ".")

	tok.Version = toks[0]
	tok.Public = toks[1]
	tok.Secret = toks[2]

	if len(toks) == 4 {
		tok.Salt = toks[3]
	}

	return tok
}
