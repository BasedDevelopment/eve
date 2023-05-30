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
	"fmt"
	"net/http"

	"github.com/BasedDevelopment/eve/internal/controllers"
	"github.com/BasedDevelopment/eve/internal/util"
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

	hv.Mutex.Lock()
	defer hv.Mutex.Unlock()
	var vms []*controllers.VM
	for _, vm := range hv.VMs {
		vm.Mutex.Lock()
		defer vm.Mutex.Unlock()
		vms = append(vms, vm)
	}

	// Send response
	if err := eUtil.WriteResponse(vms, w, http.StatusOK); err != nil {
		eUtil.WriteError(w, r, err, http.StatusInternalServerError, "Failed to marshall/send response")
	}
}

func getVM(w http.ResponseWriter, r *http.Request) (*controllers.HV, *controllers.VM) {
	hvid, err := uuid.Parse(chi.URLParam(r, "hypervisor"))
	if err != nil {
		eUtil.WriteError(w, r, err, http.StatusBadRequest, "Invalid hypervisor ID")
		return nil, nil
	}

	vmid, err := uuid.Parse(chi.URLParam(r, "virtual_machine"))
	if err != nil {
		eUtil.WriteError(w, r, err, http.StatusBadRequest, "Invalid VM ID")
		return nil, nil
	}

	vm, ok := controllers.Cloud.HVs[hvid].VMs[vmid]
	if !ok {
		eUtil.WriteError(w, r, fmt.Errorf("VM not found"), http.StatusNotFound, "Invalid VM ID")
		return nil, nil
	}

	return controllers.Cloud.HVs[hvid], vm
}

func GetVM(w http.ResponseWriter, r *http.Request) {
	_, vm := getVM(w, r)
	if vm == nil {
		return
	}
	vm.Mutex.Lock()
	defer vm.Mutex.Unlock()

	// TODO: polish response json

	// Send response
	if err := eUtil.WriteResponse(vm, w, http.StatusOK); err != nil {
		eUtil.WriteError(w, r, err, http.StatusInternalServerError, "Failed to marshall/send response")
	}
}

func GetVMState(w http.ResponseWriter, r *http.Request) {
	hv, vm := getVM(w, r)
	if vm == nil {
		return
	}

	state, err := hv.GetVMState(vm)
	if err != nil {
		eUtil.WriteError(w, r, err, http.StatusInternalServerError, "Failed to get VM state")
		return
	}

	if err := eUtil.WriteResponse(state, w, http.StatusOK); err != nil {
		eUtil.WriteError(w, r, err, http.StatusInternalServerError, "Failed to marshall/send response")
	}
}

func SetVMState(w http.ResponseWriter, r *http.Request) {
	hv, vm := getVM(w, r)
	if vm == nil {
		return
	}

	req := new(util.SetStateRequest)
	if err := util.ParseRequest(r, req); err != nil {
		eUtil.WriteError(w, r, err, http.StatusBadRequest, "Failed to parse request")
		return
	}

	respState, err := hv.SetVMState(vm, req.State)
	if err != nil {
		eUtil.WriteError(w, r, err, http.StatusInternalServerError, "Failed to set VM state")
		return
	}

	if err := eUtil.WriteResponse(respState, w, http.StatusOK); err != nil {
		eUtil.WriteError(w, r, err, http.StatusInternalServerError, "Failed to marshall/send response")
	}
}

func CreateVM(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	hvid, err := uuid.Parse(chi.URLParam(r, "hypervisor"))
	if err != nil {
		eUtil.WriteError(w, r, err, http.StatusBadRequest, "Invalid hypervisor ID")
	}
	hv := controllers.Cloud.HVs[hvid]

	vm := new(util.VMCreateRequest)
	if err := util.ParseRequest(r, vm); err != nil {
		eUtil.WriteError(w, r, err, http.StatusBadRequest, "Failed to parse request")
		return
	}

	vmid, err := hv.CreateVM(ctx, vm, hvid)
	if err != nil {
		eUtil.WriteError(w, r, err, http.StatusInternalServerError, "Failed to create VM")
		return
	}

	if err := eUtil.WriteResponse(vmid, w, http.StatusCreated); err != nil {
		eUtil.WriteError(w, r, err, http.StatusInternalServerError, "Failed to marshall/send response")
	}
}

func GetVMConsole(w http.ResponseWriter, r *http.Request) {
	hv, vm := getVM(w, r)
	if vm == nil {
		eUtil.WriteError(w, r, nil, http.StatusNotFound, "virtual machine not found")
		return
	}

	hv.Auto.WsReq(w, r, vm.ID.String())
}

func DeleteVM(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	hv, vm := getVM(w, r)
	if vm == nil {
		return
	}

	vmid := vm.ID.String()

	// destroy and undefine in controller
	if err := hv.DeleteVM(ctx, vmid); err != nil {
		eUtil.WriteError(w, r, err, http.StatusInternalServerError, "Failed to delete VM")
		return
	}

	if err := eUtil.WriteResponse(vmid, w, http.StatusOK); err != nil {
		eUtil.WriteError(w, r, err, http.StatusInternalServerError, "Failed to marshall/send response")
	}
}
