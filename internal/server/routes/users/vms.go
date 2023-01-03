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

func GetVMs(w http.ResponseWriter, r *http.Request) {
	// Fetch a list of VMs owner by the user across all HVs
	ctx := r.Context()
	userID := ctx.Value("owner").(uuid.UUID)

	cloud := controllers.Cloud
	var hvs []map[string]interface{}

	for _, hv := range cloud.HVs {
		for _, vm := range hv.VMs {
			if vm.UserID == userID {
				hvs = append(hvs, map[string]interface{}{
					"hypervisor": hv.Hostname,
					"vm":         vm,
				})
			}
		}
	}

	// Send response
	if err := util.WriteResponse(hvs, w, http.StatusOK); err != nil {
		util.WriteError(w, r, err, http.StatusInternalServerError, "Failed to marshall/send response")
	}

}
