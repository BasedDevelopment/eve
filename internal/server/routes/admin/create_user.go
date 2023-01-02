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

package admin

import (
	"net/http"

	"github.com/BasedDevelopment/eve/internal/controllers"
	"github.com/BasedDevelopment/eve/internal/util"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Decode request
	req := new(util.CreateRequest)

	if err := util.ParseRequest(r, req); err != nil {
		util.WriteError(w, r, err, http.StatusBadRequest, "Failed to parse request")
		return
	}

	// New profile instance
	profile := controllers.Profile{
		Email:    req.Email,
		Name:     req.Name,
		Disabled: req.Disabled,
		IsAdmin:  req.IsAdmin,
		Remarks:  req.Remarks,
	}

	if profile, _ := profile.Get(ctx); profile.Name != "" {
		util.WriteError(w, r, nil, http.StatusBadRequest, "User already exists")
		return
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)

	if err != nil {
		util.WriteError(w, r, err, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	profile.Password = string(hash)
	uuid, err := profile.New(ctx)

	if err != nil {
		util.WriteError(w, r, err, http.StatusInternalServerError, "Failed to create user")
		return
	}

	resp := (map[string]interface{}{
		"uuid": uuid,
	})

	util.WriteResponse(resp, w, http.StatusCreated)
}
