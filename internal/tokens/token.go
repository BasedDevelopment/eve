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
func Parse(incomingToken string) (Token, error) {
	var tok Token
	toks := strings.Split(incomingToken, ".")
	if len(toks) != 4 {
		return tok, fmt.Errorf("invalid token format")
	}

	tok.Version = toks[0]
	tok.Public = toks[1]
	tok.Secret = toks[2]
	tok.Salt = toks[3]

	return tok, nil
}
