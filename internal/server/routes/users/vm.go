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
	"errors"
	"net/http"

	"github.com/BasedDevelopment/eve/internal/controllers"
	"github.com/BasedDevelopment/eve/internal/util"
	eUtil "github.com/BasedDevelopment/eve/pkg/util"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

var (
	errNotFound  = errors.New("not found")
	errForbidden = errors.New("forbidden")
)

func GetVMs(w http.ResponseWriter, r *http.Request) {
	// Fetch a list of VMs owner by the user across all HVs
	ctx := r.Context()
	userID := ctx.Value("owner").(uuid.UUID)

	cloud := controllers.Cloud
	var response []interface{}

	for _, hv := range cloud.HVs {
		for _, vm := range hv.VMs {
			if vm.UserID == userID {
				response = append(response, map[string]interface{}{
					"hypervisor": hv,
					"name":       vm.Hostname,
				})
			}
		}
	}

	// Send response
	if err := eUtil.WriteResponse(response, w, http.StatusOK); err != nil {
		eUtil.WriteError(w, r, err, http.StatusInternalServerError, "Failed to marshall/send response")
	}

}

func GetVM(w http.ResponseWriter, r *http.Request) {
	// Fetch VM owned by the user
	ctx := r.Context()
	userID := ctx.Value("owner").(uuid.UUID)

	reqVmid := chi.URLParam(r, "virtual_machine")
	vmid, err := uuid.Parse(reqVmid)
	if err != nil {
		eUtil.WriteError(w, r, nil, http.StatusBadRequest, "invalid virtual machine id")
		return
	}

	targetHV, targetVM, err := getUserVM(userID, vmid)
	if err != nil {
		if err == errNotFound {
			eUtil.WriteError(w, r, nil, http.StatusNotFound, "virtual machine not found")
			return
		} else if err == errForbidden {
			eUtil.WriteError(w, r, nil, http.StatusForbidden, "virtual machine not owned by user")
			return
		} else {
			eUtil.WriteError(w, r, nil, http.StatusInternalServerError, "failed to fetch virtual machine")
			return
		}
	}
	//TODO
	_ = targetHV

	// Send response
	if err := eUtil.WriteResponse(targetVM, w, http.StatusOK); err != nil {
		eUtil.WriteError(w, r, err, http.StatusInternalServerError, "Failed to marshall/send response")
	}
}

func UpdateVM(w http.ResponseWriter, r *http.Request) {
	// Fetch VM owned by the user
	ctx := r.Context()

	// Decode request
	req := new(util.UserVMUpdateRequest)

	if err := util.ParseRequest(r, req); err != nil {
		eUtil.WriteError(w, r, err, http.StatusBadRequest, "Failed to parse request")
		return
	}

	// Fetch VM
	userID := ctx.Value("owner").(uuid.UUID)
	reqVmid := chi.URLParam(r, "virtual_machine")
	vmid, err := uuid.Parse(reqVmid)
	if err != nil {
		eUtil.WriteError(w, r, nil, http.StatusBadRequest, "invalid virtual machine id")
		return
	}

	targetHV, targetVM, err := getUserVM(userID, vmid)
	if err != nil {
		if err == errNotFound {
			eUtil.WriteError(w, r, nil, http.StatusNotFound, "virtual machine not found")
			return
		} else if err == errForbidden {
			eUtil.WriteError(w, r, nil, http.StatusForbidden, "virtual machine not owned by user")
			return
		} else {
			eUtil.WriteError(w, r, nil, http.StatusInternalServerError, "failed to fetch virtual machine")
			return
		}
	}

	if req.Hostname != "" {
	}

	if req.State != "" {
		switch req.State {
		case "start":
			//		if err := targetHV.Libvirt.VMStart(targetVM.Domain); err != nil {
			//			eUtil.WriteError(w, r, err, http.StatusInternalServerError, "failed to start virtual machine")
			//			return
			//		}
		case "reboot":
		case "poweroff":
		case "stop":
		case "reset":
		}
	}

	//if err := targetHV.fetchVMState(targetVM); err != nil {
	//	eUtil.WriteError(w, r, err, http.StatusInternalServerError, "failed to fetch virtual machine state")
	//	return
	//}
	_ = targetHV
	_ = targetVM

	response := map[string]interface{}{
		//	"state": targetVM.StateStr,
	}

	// Send response
	if err := eUtil.WriteResponse(response, w, http.StatusOK); err != nil {
		eUtil.WriteError(w, r, err, http.StatusInternalServerError, "Failed to marshall/send response")
	}
}

func getUserVM(userid uuid.UUID, vmid uuid.UUID) (*controllers.HV, *controllers.VM, error) {
	// Fetch VM owned by the user
	cloud := controllers.Cloud

	for _, hv := range cloud.HVs {
		for _, vm := range hv.VMs {
			if vm.ID == vmid {
				if vm.UserID != userid {
					return nil, nil, errForbidden
				} else {
					return hv, vm, nil
				}
			}
		}
	}
	return nil, nil, errNotFound
}
