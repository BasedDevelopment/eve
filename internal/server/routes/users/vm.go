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
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func GetVM(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value("owner").(uuid.UUID)

	cloud := controllers.Cloud

	reqVmid := chi.URLParam(r, "virtual_machine")
	vmid, err := uuid.Parse(reqVmid)
	if err != nil {
		util.WriteError(w, r, nil, http.StatusBadRequest, "invalid virtual machine id")
		return
	}

	response := new(controllers.VM)

	for _, hv := range cloud.HVs {
		for _, vm := range hv.VMs {
			if vm.ID == vmid {
				if vm.UserID != userID {
					util.WriteError(w, r, nil, http.StatusForbidden, "forbidden")
					return
				} else {
					response = cloud.HVs[hv.ID].VMs[vmid]
					break
				}
			}
		}
	}

	if response.ID == uuid.Nil {
		util.WriteError(w, r, nil, http.StatusNotFound, "virtual machine not found")
		return
	}

	// Send response
	if err := util.WriteResponse(response, w, http.StatusOK); err != nil {
		util.WriteError(w, r, err, http.StatusInternalServerError, "Failed to marshall/send response")
	}
}
