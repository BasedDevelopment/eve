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
	"net/http"
	"strings"

	"github.com/BasedDevelopment/eve/internal/controllers"
	"github.com/BasedDevelopment/eve/internal/sessions"
	"github.com/BasedDevelopment/eve/internal/tokens"
	"github.com/BasedDevelopment/eve/internal/util"
	"github.com/google/uuid"
)

func getToken(w http.ResponseWriter, r *http.Request) (token tokens.Token) {
	authorizationHeader := r.Header.Get("Authorization")

	if authorizationHeader == "" {
		util.WriteError(w, r, nil, http.StatusBadRequest, "Missing Authorization header")
		return tokens.Token{}
	}

	splitHeader := strings.Split(authorizationHeader, "Bearer ")
	token, err := tokens.Parse(splitHeader[1])
	if err != nil {
		util.WriteError(w, r, nil, http.StatusBadRequest, "Invalid token")
		return
	}
	return
}

// Auth forces a user to be authenticated before continuing to the route
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		requestToken := getToken(w, r)

		if !sessions.ValidateSession(ctx, requestToken) {
			util.WriteError(w, r, nil, http.StatusUnauthorized, "Unauthorized")

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
			util.WriteError(w, r, err, http.StatusInternalServerError, "Internal Server Error")

			return
		}

		if !profile.IsAdmin {
			util.WriteError(w, r, nil, http.StatusUnauthorized, "Unauthorized")

			return
		}

		next.ServeHTTP(w, r)
	})
}
