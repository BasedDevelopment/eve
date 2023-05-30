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
	"github.com/rs/zerolog/log"
)

type VM struct {
	Mutex    sync.Mutex           `json:"-" db:"-"`
	ID       uuid.UUID            `json:"id"`
	HV       uuid.UUID            `json:"hv" db:"hv_id"`
	Hostname string               `json:"hostname"`
	UserID   uuid.UUID            `json:"user" db:"profile_id"`
	CPU      int                  `json:"cpu"`
	Memory   int64                `json:"memory"`
	Nics     map[string]VMNic     `json:"nics" db:"-"`
	Storages map[string]VMStorage `json:"storages" db:"-"`
	Created  time.Time            `json:"created"`
	Updated  time.Time            `json:"updated"`
	Remarks  string               `json:"remarks"`
	Domain   models.VM            `json:"-" db:"-"`
}

type VMNic struct {
	Mutex   sync.Mutex `db:"-" json:"-"`
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
	Mutex   sync.Mutex `db:"-" json:"-"`
	ID      uuid.UUID
	Size    int
	Created time.Time
	Updated time.Time
	Remarks string
}

func (hv *HV) getVMsFromDB() (vms []VM, err error) {
	// Fetch VMs from database for this hypervisor
	rows, queryErr := db.Pool.Query(context.Background(),
		"SELECT * FROM vm WHERE hv_id = $1", hv.ID)

	if queryErr != nil {
		return nil, fmt.Errorf("failed to query VMs: %w", queryErr)
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

// Fetch VMs from the DB and Libvirt, marshall them into the HV struct, and
// check for inconsistencies
func (hv *HV) InitVMs() error {
	// Use Auto to get the VMs from Libvirt
	vms, err := hv.Auto.GetLibvirtVMs()
	if err != nil {
		return err
	}

	// Get VMs from DB
	dbVMs, err := hv.getVMsFromDB()
	if err != nil {
		return err
	}

	// Check for mismatch in amount of VM's on libvirt vs in database
	if len(dbVMs) != len(vms) {
		log.Warn().
			Str("hv", hv.Hostname).
			Int("db", len(dbVMs)).
			Int("libvirt", len(vms)).
			Msg("VM count mismatch")
	}

	hv.Mutex.Lock()
	defer hv.Mutex.Unlock()
	// Load the VMs of a hypervisor into memory
	for i := range dbVMs {
		id := dbVMs[i].ID

		// Load from database
		hv.VMs[id] = &dbVMs[i]

		// Consistency check w/ libvirt
		for j := range vms {
			// Check if libvirt VM id == database ID
			if vms[j].ID == id {
				// Domain is the VM object, direct from libvirt
				hv.VMs[id].Domain = vms[j]

				// Check consistency between our database & libvirt/auto.
				go hv.checkVMConsistency(hv.VMs[id])
			}
		}
	}

	//go hv.checkUndefinedVMs()

	return nil
}

func (hv *HV) checkVMConsistency(dbvm *VM) {
	dom := dbvm.Domain

	// Check for CPU count
	if dom.CPU != dbvm.CPU {
		log.Warn().
			Str("hv", hv.Hostname).
			Str("vm", dbvm.Hostname).
			Int("db", dbvm.CPU).
			Int("libvirt", dom.CPU).
			Msg("CPU count mismatch")
	}

	// Check for memory size
	if dom.Memory != dbvm.Memory {
		log.Warn().
			Str("hv", hv.Hostname).
			Str("vm", dbvm.Hostname).
			Int64("db", dbvm.Memory).
			Int64("libvirt", dom.Memory).
			Msg("Memory size mismatch")
	}
}

func (hv *HV) DeleteVM(ctx context.Context, vmid string) error {
	if err := hv.Auto.DeleteVM(vmid); err != nil {
		return err
	}

	if _, err := db.Pool.Exec(
		ctx,
		"DELETE FROM vm WHERE id = $1",
		vmid,
	); err != nil {
		return err
	}

	_, err := hv.getVMsFromDB()
	return err
}
