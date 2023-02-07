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

func GetHVs(w http.ResponseWriter, r *http.Request) {
	cloud := controllers.Cloud
	cloud.Mutex.Lock()
	defer cloud.Mutex.Unlock()
	var hvs []*controllers.HV
	for _, hv := range cloud.HVs {
		hv.Mutex.Lock()
		defer hv.Mutex.Unlock()
		hvs = append(hvs, hv)
	}

	// Send response
	if err := eUtil.WriteResponse(hvs, w, http.StatusOK); err != nil {
		eUtil.WriteError(w, r, err, http.StatusInternalServerError, "Failed to marshall/send response")
	}
}

func getHV(w http.ResponseWriter, r *http.Request) *controllers.HV {
	hvidStr := chi.URLParam(r, "hypervisor")
	hvid, err := uuid.Parse(hvidStr)
	if err != nil {
		eUtil.WriteError(w, r, err, http.StatusBadRequest, "Invalid hypervisor ID")
		return nil
	}
	hv, ok := controllers.Cloud.HVs[hvid]
	if !ok {
		eUtil.WriteError(w, r, err, http.StatusNotFound, "Hypervisor not found")
		return nil
	}
	return hv
}

func GetHV(w http.ResponseWriter, r *http.Request) {
	hv := getHV(w, r)
	if hv == nil {
		return
	}

	if err := hv.Refresh(); err != nil {
		eUtil.WriteError(w, r, err, http.StatusInternalServerError, "Failed to refresh hypervisor")
		return
	}

	hv.Mutex.Lock()
	defer hv.Mutex.Unlock()

	// Send response
	if err := eUtil.WriteResponse(hv.Specs, w, http.StatusOK); err != nil {
		eUtil.WriteError(w, r, err, http.StatusInternalServerError, "Failed to marshall/send response")
	}
}

func GetHVState(w http.ResponseWriter, r *http.Request) {
	hv := getHV(w, r)
	if hv == nil {
		return
	}

	if err := hv.Refresh(); err != nil {
		eUtil.WriteError(w, r, err, http.StatusInternalServerError, "Failed to refresh hypervisor state")
	}

	hv.Mutex.Lock()
	defer hv.Mutex.Unlock()

	response := map[string]interface{}{
		"state":     hv.Specs.Status,
		"state_str": hv.Specs.Status.String(),
		"reason":    hv.Specs.StatusReason,
	}

	if err := eUtil.WriteResponse(response, w, http.StatusOK); err != nil {
		eUtil.WriteError(w, r, err, http.StatusInternalServerError, "Failed to marshall/send response")
	}
}
