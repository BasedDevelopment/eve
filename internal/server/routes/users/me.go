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

package users

import (
	"net/http"

	"github.com/BasedDevelopment/eve/internal/controllers"
	"github.com/BasedDevelopment/eve/internal/util"
	"github.com/google/uuid"
)

func GetSelf(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	profile := controllers.Profile{ID: ctx.Value("owner").(uuid.UUID)}
	profile, err := profile.Get(ctx)

	if err != nil {
		util.WriteError(w, r, err, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	response := map[string]interface{}{
		"id":         profile.ID,
		"name":       profile.Name,
		"email":      profile.Email,
		"last_login": profile.LastLogin,
		"created":    profile.Created,
		"updated":    profile.Updated,
	}

	err = util.WriteResponse(response, w, http.StatusOK)

	if err != nil {
		util.WriteError(w, r, err, http.StatusInternalServerError, "Internal Server Error")
		return
	}
}
