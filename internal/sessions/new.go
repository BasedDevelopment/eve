/*
 * eve - management toolkit for libvirt servers
 * Copyright (C) 2022-2023  BNS Services LLC

 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.

 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.

 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package sessions

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/BasedDevelopment/eve/internal/profile"
	"github.com/BasedDevelopment/eve/internal/tokens"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/sha3"
)

const version = "v1"

var expirey = 24 * time.Hour

func prngString(size int) string {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to generate random string")
		return ""
	}

	return fmt.Sprintf("%x", b)
}

func generateStrings(bits []int) (a, b, c string, err error) {
	a = prngString(bits[0])
	b = prngString(bits[1])
	c = prngString(bits[2])

	return a, b, c, err
}

// New creates a new authentication session in the database (login)
func NewSession(ctx context.Context, user profile.Profile) (tokens.Token, error) {
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
		Owner:   user.ID,
		Version: version,
		Public:  public,
		Secret:  fmt.Sprintf("%x", saltedSecret),
		Salt:    salt,
		Created: time.Now(),
		Expires: time.Now().Add(expirey),
	}

	// Push the session to the database
	if err := session.push(ctx); err != nil {
		return tokens.Token{}, err
	}

	// This is the token we want to give to the user
	return tokens.Token{
		Version: version,
		Public:  public,
		Secret:  base64.URLEncoding.EncodeToString([]byte(secret)),
		Salt:    salt,
	}, nil
}
