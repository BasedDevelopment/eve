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
	"context"
	"net/http"

	"github.com/BasedDevelopment/eve/internal/db"
	"github.com/BasedDevelopment/eve/internal/profile"
	"github.com/BasedDevelopment/eve/internal/util"
	eUtil "github.com/BasedDevelopment/eve/pkg/util"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Decode request
	req := new(util.UserCreateRequest)

	if err := util.ParseRequest(r, req); err != nil {
		eUtil.WriteError(w, r, err, http.StatusBadRequest, "Failed to parse request")
		return
	}

	// New profile instance
	profile := profile.Profile{
		Email:    req.Email,
		Name:     req.Name,
		Disabled: req.Disabled,
		IsAdmin:  req.IsAdmin,
		Remarks:  req.Remarks,
	}

	if profile, _ := profile.Get(ctx); profile.Name != "" {
		eUtil.WriteError(w, r, nil, http.StatusBadRequest, "User already exists")
		return
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)

	if err != nil {
		eUtil.WriteError(w, r, err, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	profile.Password = string(hash)
	uuid, err := profile.New(ctx)

	if err != nil {
		eUtil.WriteError(w, r, err, http.StatusInternalServerError, "Failed to create user")
		return
	}

	resp := (map[string]interface{}{
		"uuid": uuid,
	})

	eUtil.WriteResponse(resp, w, http.StatusCreated)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Pool.Query(context.Background(), "SELECT * FROM profile")
	if err != nil {
		eUtil.WriteError(w, r, err, http.StatusInternalServerError, "Failed to get users")
		return
	}
	defer rows.Close()

	users, err := pgx.CollectRows(rows, pgx.RowToStructByName[profile.Profile])
	if err != nil {
		eUtil.WriteError(w, r, err, http.StatusInternalServerError, "Failed to collect users")
		return
	}

	eUtil.WriteResponse(users, w, http.StatusOK)
}
