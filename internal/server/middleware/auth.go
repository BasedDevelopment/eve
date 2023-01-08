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

package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/BasedDevelopment/eve/internal/controllers"
	"github.com/BasedDevelopment/eve/internal/sessions"
	"github.com/BasedDevelopment/eve/internal/tokens"
	eUtil "github.com/BasedDevelopment/eve/pkg/util"
	"github.com/google/uuid"
)

var (
	ErrMissingHeader = errors.New("missing Authorization header")
	ErrBadHeader     = errors.New("invalid Authorization header")
	ErrBadToken      = errors.New("invalid token")
)

func getToken(w http.ResponseWriter, r *http.Request) (tokens.Token, error) {
	authorizationHeader := r.Header.Get("Authorization")

	if authorizationHeader == "" {
		return tokens.Token{}, ErrMissingHeader
	}

	splitHeader := strings.Split(authorizationHeader, "Bearer ")
	if len(splitHeader) != 2 {
		return tokens.Token{}, ErrBadHeader
	}

	token, err := tokens.Parse(splitHeader[1])
	if err != nil {
		return tokens.Token{}, ErrBadToken
	}

	return token, nil
}

// Auth forces a user to be authenticated before continuing to the route
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		requestToken, err := getToken(w, r)

		if err != nil {
			switch err {
			case ErrBadHeader:
				eUtil.WriteError(w, r, nil, http.StatusBadRequest, ErrBadHeader.Error())
				return
			case ErrBadToken:
				eUtil.WriteError(w, r, nil, http.StatusUnauthorized, ErrBadToken.Error())
				return
			case ErrMissingHeader:
				eUtil.WriteError(w, r, nil, http.StatusBadRequest, ErrMissingHeader.Error())
				return
			}
		}

		if !sessions.ValidateSession(ctx, requestToken) {
			eUtil.WriteError(w, r, nil, http.StatusUnauthorized, "Unauthorized")
			return
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// MustBeAdmin verifies is an **already authenticated** user is an admin
func MustBeAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		profile := controllers.Profile{ID: ctx.Value("owner").(uuid.UUID)}
		profile, err := profile.Get(ctx)

		if err != nil {
			eUtil.WriteError(w, r, err, http.StatusInternalServerError, "Internal Server Error")

			return
		}

		if !profile.IsAdmin {
			eUtil.WriteError(w, r, nil, http.StatusUnauthorized, "Unauthorized")

			return
		}

		next.ServeHTTP(w, r)
	})
}
