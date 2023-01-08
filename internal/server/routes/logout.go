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
	"strings"

	"github.com/BasedDevelopment/eve/internal/sessions"
	"github.com/BasedDevelopment/eve/internal/tokens"
	eUtil "github.com/BasedDevelopment/eve/pkg/util"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	header := strings.Split(r.Header.Get("Authorization"), "Bearer ")
	reqToken, err := tokens.Parse(header[1])
	if err != nil {
		eUtil.WriteError(w, r, nil, http.StatusBadRequest, "Invalid token")
		return
	}

	sessions.Delete(ctx, reqToken)

	if err := sessions.Delete(ctx, reqToken); err != nil {
		eUtil.WriteError(w, r, nil, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	eUtil.WriteResponse(map[string]interface{}{
		"message": "logout success",
	}, w, http.StatusOK)
}
