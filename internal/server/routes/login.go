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

package routes

import (
	"net/http"

	"github.com/BasedDevelopment/eve/internal/profile"
	"github.com/BasedDevelopment/eve/internal/sessions"
	"github.com/BasedDevelopment/eve/internal/util"
	eUtil "github.com/BasedDevelopment/eve/pkg/util"
	"golang.org/x/crypto/bcrypt"
)

func Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse request body
	req := new(util.LoginRequest)

	if err := util.ParseRequest(r, req); err != nil {
		eUtil.WriteError(w, r, nil, http.StatusBadRequest, "Failed to parse login request")
		return
	}

	// New profile instance
	profile := profile.Profile{Email: req.Email}
	profile, err := profile.Get(ctx)

	if err != nil {
		eUtil.WriteError(w, r, nil, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Validate password
	if err := bcrypt.CompareHashAndPassword([]byte(profile.Password), []byte(req.Password)); err != nil {
		eUtil.WriteError(w, r, nil, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Issue token
	userToken, err := sessions.NewSession(ctx, profile)

	if err != nil {
		eUtil.WriteError(w, r, err, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	// Send token to client
	eUtil.WriteResponse(map[string]string{
		"token": userToken.String(),
	}, w, http.StatusOK)
}
