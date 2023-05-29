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

	"github.com/BasedDevelopment/eve/internal/auto"
	"github.com/BasedDevelopment/eve/internal/controllers"
	"github.com/BasedDevelopment/eve/internal/util"
	eUtil "github.com/BasedDevelopment/eve/pkg/util"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func GetVMs(w http.ResponseWriter, r *http.Request) {
	// Fetch a list of VMs owner by the user across all HVs
	ctx := r.Context()
	userID := ctx.Value("owner").(uuid.UUID)

	cloud := controllers.Cloud
	var response []interface{}
	cloud.Mutex.Lock()
	defer cloud.Mutex.Unlock()

	for _, hv := range cloud.HVs {
		hv.Mutex.Lock()
		defer hv.Mutex.Unlock()
		for _, vm := range hv.VMs {
			if vm.UserID == userID {
				vm.Mutex.Lock()
				defer vm.Mutex.Unlock()
				response = append(response, map[string]interface{}{
					"hypervisor": hv.Hostname,
					"name":       vm.Hostname,
					"id":         vm.ID,
				})
			}
		}
	}

	// Send response
	if err := eUtil.WriteResponse(response, w, http.StatusOK); err != nil {
		eUtil.WriteError(w, r, err, http.StatusInternalServerError, "Failed to marshall/send response")
	}

}

func getUserVM(w http.ResponseWriter, r *http.Request) (*controllers.HV, *controllers.VM) {
	ctx := r.Context()
	cloud := controllers.Cloud
	userID := ctx.Value("owner").(uuid.UUID)
	reqVmid := chi.URLParam(r, "virtual_machine")
	vmid, err := uuid.Parse(reqVmid)
	if err != nil {
		eUtil.WriteError(w, r, nil, http.StatusBadRequest, "invalid virtual machine id")
		return nil, nil
	}

	for _, hv := range cloud.HVs {
		for _, vm := range hv.VMs {
			if vm.ID == vmid {
				if vm.UserID != userID {
					eUtil.WriteError(w, r, nil, http.StatusForbidden, "virtual machine not owned by user")
					return nil, nil
				} else {
					return hv, vm
				}
			}
		}
	}
	eUtil.WriteError(w, r, nil, http.StatusNotFound, "virtual machine not found")
	return nil, nil
}

func GetVM(w http.ResponseWriter, r *http.Request) {
	_, vm := getUserVM(w, r)
	if vm == nil {
		return
	}
	vm.Mutex.Lock()
	defer vm.Mutex.Unlock()

	// Send response
	if err := eUtil.WriteResponse(vm, w, http.StatusOK); err != nil {
		eUtil.WriteError(w, r, err, http.StatusInternalServerError, "Failed to marshall/send response")
	}
}

func GetVMState(w http.ResponseWriter, r *http.Request) {
	hv, vm := getUserVM(w, r)
	if vm == nil {
		return
	}
	vm.Mutex.Lock()
	defer vm.Mutex.Unlock()

	state, err := hv.Auto.GetVMState(vm.ID.String())
	if err != nil {
		eUtil.WriteError(w, r, err, http.StatusInternalServerError, "Failed to get VM state")
	}

	if err := eUtil.WriteResponse(state, w, http.StatusOK); err != nil {
		eUtil.WriteError(w, r, err, http.StatusInternalServerError, "Failed to marshall/send response")
	}
}

func SetVMState(w http.ResponseWriter, r *http.Request) {
	hv, vm := getUserVM(w, r)
	if vm == nil {
		eUtil.WriteError(w, r, nil, http.StatusNotFound, "virtual machine not found")
		return
	}
	vm.Mutex.Lock()
	defer vm.Mutex.Unlock()

	req := new(util.SetStateRequest)
	if err := util.ParseRequest(r, req); err != nil {
		eUtil.WriteError(w, r, err, http.StatusBadRequest, "Failed to parse request")
		return
	}

	var status uint8
	switch req.State {
	case "start":
		status = auto.Start
	case "reboot":
		status = auto.Reboot
	case "poweroff":
		status = auto.Poweroff
	case "stop":
		status = auto.Stop
	case "reset":
		status = auto.Reset
	}

	respState, err := hv.Auto.SetVMState(vm.ID.String(), status)
	if err != nil {
		eUtil.WriteError(w, r, err, http.StatusInternalServerError, "Failed to set VM state")
		return
	}

	if err := eUtil.WriteResponse(respState, w, http.StatusOK); err != nil {
		eUtil.WriteError(w, r, err, http.StatusInternalServerError, "Failed to marshall/send response")
	}
}

func GetVMConsole(w http.ResponseWriter, r *http.Request) {
	hv, vm := getUserVM(w, r)
	if vm == nil {
		eUtil.WriteError(w, r, nil, http.StatusNotFound, "virtual machine not found")
		return
	}
	hv.Auto.WsReq(w, r, vm.ID.String())
}
