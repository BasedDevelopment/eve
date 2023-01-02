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

package middlewares

import (
	"context"
	"net/http"

	"github.com/BasedDevelopment/eve/internal/sessions"
	"github.com/BasedDevelopment/eve/internal/util"
)

// UserContext fetches the owner of the request from the current session and
// and appends it to the request context. Requires Auth, required MustBeAdmin.
func UserContext(next http.Handler) http.Handler {
	// This function doesn't check whether a user is authenticated
	// and as such should only be used after Auth has been called.

	// It is required for the MustBeAdmin middleware though, since
	// that middleware uses the profile in the request context.
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Get session
		requestToken := getToken(w, r) // function from auth middleware; ges token from authorization header
		session, err := sessions.GetSession(ctx, requestToken)

		if err != nil {
			util.WriteError(w, r, err, http.StatusInternalServerError, "Internal Server Error")

			return
		}

		// Add the owner of the session to r.Context
		ctx = context.WithValue(ctx, "owner", session.Owner)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
