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
	eUtil "github.com/BasedDevelopment/eve/pkg/util"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func GetVMs(w http.ResponseWriter, r *http.Request) {
	// Get hv ID from request
	hvidStr := chi.URLParam(r, "hypervisor")
	hvid, err := uuid.Parse(hvidStr)
	if err != nil {
		eUtil.WriteError(w, r, err, http.StatusBadRequest, "Invalid hypervisor ID")
		return
	}
	hv := controllers.Cloud.HVs[hvid]

	var vms []*controllers.VM
	for _, vm := range hv.VMs {
		vms = append(vms, vm)
	}

	// Send response
	if err := eUtil.WriteResponse(vms, w, http.StatusOK); err != nil {
		eUtil.WriteError(w, r, err, http.StatusInternalServerError, "Failed to marshall/send response")
	}
}

func GetVM(w http.ResponseWriter, r *http.Request) {
	// Get hv ID from request
	hvidStr := chi.URLParam(r, "hypervisor")
	hvid, err := uuid.Parse(hvidStr)
	if err != nil {
		eUtil.WriteError(w, r, err, http.StatusBadRequest, "Invalid hypervisor ID")
		return
	}
	hv := controllers.Cloud.HVs[hvid]

	// Get vm ID from request
	vmidStr := chi.URLParam(r, "virtual_machine")
	vmid, err := uuid.Parse(vmidStr)
	if err != nil {
		eUtil.WriteError(w, r, err, http.StatusBadRequest, "Invalid VM ID")
		return
	}
	// TODO: polish response json

	// Send response
	if err := eUtil.WriteResponse(hv.VMs[vmid], w, http.StatusOK); err != nil {
		eUtil.WriteError(w, r, err, http.StatusInternalServerError, "Failed to marshall/send response")
	}
}
