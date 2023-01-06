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

func GetVM(w http.ResponseWriter, r *http.Request) {
	// Fetch VM owned by the user
	ctx := r.Context()
	userID := ctx.Value("owner").(uuid.UUID)

	reqVmid := chi.URLParam(r, "virtual_machine")
	vmid, err := uuid.Parse(reqVmid)
	if err != nil {
		util.WriteError(w, r, nil, http.StatusBadRequest, "invalid virtual machine id")
		return
	}

	_, targetVM, err := getUserVM(userID, vmid)
	if err != nil {
		if err == errNotFound {
			util.WriteError(w, r, nil, http.StatusNotFound, "virtual machine not found")
			return
		} else if err == errForbidden {
			util.WriteError(w, r, nil, http.StatusForbidden, "virtual machine not owned by user")
			return
		} else {
			util.WriteError(w, r, nil, http.StatusInternalServerError, "failed to fetch virtual machine")
			return
		}
	}

	// Send response
	if err := util.WriteResponse(targetVM, w, http.StatusOK); err != nil {
		util.WriteError(w, r, err, http.StatusInternalServerError, "Failed to marshall/send response")
	}
}

func UpdateVM(w http.ResponseWriter, r *http.Request) {
	// Fetch VM owned by the user
	ctx := r.Context()

	// Decode request
	req := new(util.UserVMUpdateRequest)

	if err := util.ParseRequest(r, req); err != nil {
		util.WriteError(w, r, err, http.StatusBadRequest, "Failed to parse request")
		return
	}

	// Fetch VM
	userID := ctx.Value("owner").(uuid.UUID)
	reqVmid := chi.URLParam(r, "virtual_machine")
	vmid, err := uuid.Parse(reqVmid)
	if err != nil {
		util.WriteError(w, r, nil, http.StatusBadRequest, "invalid virtual machine id")
		return
	}

	targetHV, targetVM, err := getUserVM(userID, vmid)
	if err != nil {
		if err == errNotFound {
			util.WriteError(w, r, nil, http.StatusNotFound, "virtual machine not found")
			return
		} else if err == errForbidden {
			util.WriteError(w, r, nil, http.StatusForbidden, "virtual machine not owned by user")
			return
		} else {
			util.WriteError(w, r, nil, http.StatusInternalServerError, "failed to fetch virtual machine")
			return
		}
	}

	if req.Hostname != "" {
	}

	if req.State != "" {
		switch req.State {
		case "start":
			if err := targetHV.Libvirt.VMStart(targetVM.Domain); err != nil {
				util.WriteError(w, r, err, http.StatusInternalServerError, "failed to start virtual machine")
				return
			}
		case "reboot":
			if err := targetHV.Libvirt.VMReboot(targetVM.Domain); err != nil {
				util.WriteError(w, r, err, http.StatusInternalServerError, "failed to reboot virtual machine")
				return
			}
		case "poweroff":
			if err := targetHV.Libvirt.VMPowerOff(targetVM.Domain); err != nil {
				util.WriteError(w, r, err, http.StatusInternalServerError, "failed to power off virtual machine")
				return
			}
		case "stop":
			if err := targetHV.Libvirt.VMStop(targetVM.Domain); err != nil {
				util.WriteError(w, r, err, http.StatusInternalServerError, "failed to stop virtual machine")
				return
			}
		case "reset":
			if err := targetHV.Libvirt.VMReset(targetVM.Domain); err != nil {
				util.WriteError(w, r, err, http.StatusInternalServerError, "failed to reset virtual machine")
				return
			}
		}
	}

	response := map[string]interface{}{
		"vm": targetVM.ID,
	}

	// Send response
	if err := util.WriteResponse(response, w, http.StatusOK); err != nil {
		util.WriteError(w, r, err, http.StatusInternalServerError, "Failed to marshall/send response")
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
