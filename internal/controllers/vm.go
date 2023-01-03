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
	"context"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/BasedDevelopment/eve/internal/db"
	"github.com/BasedDevelopment/eve/internal/libvirt"
	"github.com/BasedDevelopment/eve/internal/util"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type VM struct {
	ID          uuid.UUID
	HV          uuid.UUID `db:"hv_id"`
	Hostname    string
	UserID      uuid.UUID `db:"profile_id"`
	User        *Profile  `db:"-"`
	CPU         int
	Memory      int
	Nics        map[string]VMNic     `db:"-"`
	Storages    map[string]VMStorage `db:"-"`
	Created     time.Time
	Updated     time.Time
	Remarks     string
	Domain      libvirt.Dom `db:"-" json:"-"`
	State       util.Status `db:"-" json:"-"`
	StateStr    string      `db:"-"`
	StateReason string      `db:"-"`
}

type VMNic struct {
	ID      uuid.UUID
	name    string
	MAC     string
	IP      []net.IP `db:"-"`
	Created time.Time
	Updated time.Time
	Remarks string
	State   string `db:"-"`
}

type VMStorage struct {
	ID      uuid.UUID
	Size    int
	Created time.Time
	Updated time.Time
	Remarks string
}

// Fetch VMs from the DB and Libvirt, marshall them into the HV struct,
// and check for consistency
// TODO: concurrency
func (hv *HV) InitVMs() error {
	if err := hv.ensureConn(); err != nil {
		return err
	}

	// Get VMs from libvirt
	libvirtVMs, err := hv.getVMsFromLibvirt()
	if err != nil {
		return err
	}

	// Get VMs from DB
	dbVMs, err := hv.getVMsFromDB()
	if err != nil {
		return err
	}

	hv.mutex.Lock()
	defer hv.mutex.Unlock()

	// Marshall the HV.VMs struct in
	for i := range dbVMs {
		hv.VMs[dbVMs[i].ID] = &dbVMs[i]
		hv.VMs[dbVMs[i].ID].Domain = libvirtVMs[dbVMs[i].ID]
	}

	// Check if the VMs are consistent
	if err := hv.checkVMConsistency(libvirtVMs, hv.VMs); err != nil {
		return err
	}

	// Fetch VM state and state reason
	for uuid := range hv.VMs {
		if err := hv.GetVMState(hv.VMs[uuid]); err != nil {
			return err
		}
	}

	return nil
}

// Get the list of VMs from the DB
// Will be used to check consistency
func (hv *HV) getVMsFromDB() (vms []VM, err error) {
	// Query DB
	rows, queryErr := db.Pool.Query(context.Background(), "SELECT * FROM vm WHERE hv_id = $1", hv.ID)
	if queryErr != nil {
		err := fmt.Errorf("failed to query VMs: %w", queryErr)
		return nil, err
	}
	defer rows.Close()

	// Collect the rows into the VM struct
	vms, collectErr := pgx.CollectRows(rows, pgx.RowToStructByName[VM])

	if collectErr != nil {
		err := fmt.Errorf("error collecting VMs: %w", collectErr)
		return nil, err
	}
	return
}

// Get the list of VMs from libvirt
// Will be used to check consistency
func (hv *HV) getVMsFromLibvirt() (doms map[uuid.UUID]libvirt.Dom, err error) {
	if err := hv.ensureConn(); err != nil {
		return nil, err
	}

	doms, err = hv.Libvirt.GetVMs()
	if err != nil {
		return nil, err
	}
	return
}

func (hv *HV) checkVMConsistency(libvirt map[uuid.UUID]libvirt.Dom, db map[uuid.UUID]*VM) error {
	if err := hv.ensureConn(); err != nil {
		return err
	}

	for uuid, dom := range libvirt {

		// Get the VM specs from libvirt
		domSpec, err := hv.Libvirt.GetVMSpecs(dom)
		if err != nil {
			return err
		}

		// Check if the VM is in the DB
		if _, ok := db[uuid]; !ok {
			return fmt.Errorf("VM %s is not in the DB", uuid)
		}

		// Check for CPU count
		if domSpec.Vcpu.Text != strconv.Itoa(db[uuid].CPU) {
			return fmt.Errorf("CPU count mismatch for VM %s", uuid)
		}
	}
	return nil
}

func (hv *HV) GetVMState(vm *VM) (err error) {
	if err := hv.ensureConn(); err != nil {
		return err
	}

	stateInt, stateStr, reasonStr, err := hv.Libvirt.GetVMState(vm.Domain)
	if err != nil {
		return err
	}

	vm.State = stateInt
	vm.StateStr = stateStr
	vm.StateReason = reasonStr
	return nil
}
