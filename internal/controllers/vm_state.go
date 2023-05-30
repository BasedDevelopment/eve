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

package controllers

import (
	"github.com/BasedDevelopment/auto/pkg/models"
	"github.com/BasedDevelopment/eve/internal/auto"
)

func (hv *HV) GetVMState(vm *VM) (models.VMState, error) {
	vm.Mutex.Lock()
	defer vm.Mutex.Unlock()

	id := vm.ID.String()
	return hv.Auto.GetVMState(id)
}

func (hv *HV) SetVMState(vm *VM, state string) (models.VMState, error) {
	vm.Mutex.Lock()
	defer vm.Mutex.Unlock()

	id := vm.ID.String()

	var status uint8
	switch state {
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
	return hv.Auto.SetVMState(id, status)
}
