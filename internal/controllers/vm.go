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
	"sync"
	"time"

	"github.com/BasedDevelopment/auto/pkg/models"
	"github.com/BasedDevelopment/eve/internal/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// TODO: Delete a bunch of this and start using auto method
// TODO: don't forget to consistency check
type VM struct {
	mutex    sync.Mutex           `db:"-" json:"-"`
	ID       uuid.UUID            `json:"id"`
	HV       uuid.UUID            `db:"hv_id" json:"hv"`
	Hostname string               `json:"hostname"`
	UserID   uuid.UUID            `db:"profile_id" json:"user"`
	CPU      int                  `json:"cpu"`
	Memory   int64                `json:"memory"`
	Nics     map[string]VMNic     `db:"-" json:"nics"`
	Storages map[string]VMStorage `db:"-" json:"storages"`
	Created  time.Time            `json:"created"`
	Updated  time.Time            `json:"updated"`
	Remarks  string               `json:"remarks"`
	Domain   models.VM            `json:"-"`
}

type VMNic struct {
	mutex   sync.Mutex `db:"-" json:"-"`
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
	mutex   sync.Mutex `db:"-" json:"-"`
	ID      uuid.UUID
	Size    int
	Created time.Time
	Updated time.Time
	Remarks string
}

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

// Fetch VMs from the DB and Libvirt, marshall them into the HV struct,
// and check for consistency
func (hv *HV) InitVMs() error {
	// Fetch VMs from HV
	libvirtVMs, err := hv.Auto.GetLibvirtVMs()
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
		vm := hv.VMs[dbVMs[i].ID]
		vm = &dbVMs[i]
		for i := range libvirtVMs {
			if libvirtVMs[i].ID == vm.ID {
				vm.Domain = libvirtVMs[i]
			}
		}
	}

	//go consistencyCheck(libvirtVMs, hv)
	//go hv.checkUndefinedVMs()

	return nil
}

// Fetch VM state and state reason
func (hv *HV) fetchVMState(vm *VM) (models.VMState, error) {
	//TODO
	return models.VMState{}, nil
}

func (hv *HV) checkVMConsistency(domain map[uuid.UUID]models.VM, db map[uuid.UUID]*VM) error {
	for uuid, dom := range domain {

		_ = dom
		_ = db
		_ = uuid
		// Get the VM specs from libvirt

		// Check if the VM is in the DB

		// Check for CPU count

		// Check for memory size
	}
	return nil
}

func (hv *HV) checkUndefinedVMs() {
	//doms, err := hv.Libvirt.GetUndefinedVMs()
	//if err != nil {
	//	log.Error().Err(err).Msg("failed to get undefined VMs")
	//}
	//if len(doms) > 0 {
	//	log.Warn().Str("hv", hv.Hostname).Int("count", len(doms)).Msg("undefined VMs found")
	//} else {
	//	log.Info().Str("hv", hv.Hostname).Msg("no undefined VMs found")
	//}
}
