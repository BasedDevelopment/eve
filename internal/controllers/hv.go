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

	"github.com/BasedDevelopment/eve/internal/db"
	"github.com/BasedDevelopment/eve/internal/libvirt"
	"github.com/BasedDevelopment/eve/internal/util"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type HV struct {
	mutex        sync.Mutex               `db:"-" json:"-"`
	ID           uuid.UUID                `json:"id"`
	Hostname     string                   `json:"hostname"`
	IP           net.IP                   `json:"ip"`
	Port         int                      `json:"port"`
	Site         string                   `json:"site"`
	Nics         map[uuid.UUID]*HVNic     `json:"nics" db:"-"`
	Storages     map[uuid.UUID]*HVStorage `json:"storages" db:"-"`
	VMs          map[uuid.UUID]*VM        `json:"vms" db:"-"`
	Created      time.Time                `json:"created"`
	Updated      time.Time                `json:"updated"`
	Remarks      string                   `json:"remarks"`
	Status       util.Status              `json:"status" db:"-"`
	StatusReason string                   `json:"status_reason" db:"-"`
	Version      string                   `json:"version" db:"-"`
	Libvirt      *libvirt.Libvirt         `json:"-" db:"-"`
}

type HVNic struct {
	mutex   sync.Mutex `db:"-"`
	ID      uuid.UUID
	Name    string
	Mac     net.HardwareAddr
	IP      []net.IP
	Created time.Time
	Updated time.Time
	Remarks string
}

// TODO: Impl bridges

type HVStorage struct {
	mutex     sync.Mutex `db:"-"`
	ID        uuid.UUID
	Size      int
	Used      int `db:"-"`
	Available int `db:"-"`
	Created   time.Time
	Updated   time.Time
	Remarks   string
}

func getHVs(cloud *HVList) (err error) {
	// Reading HVs
	rows, queryErr := db.Pool.Query(context.Background(), "SELECT * FROM hv")

	if queryErr != nil {
		return fmt.Errorf("Error reading hv: %w", queryErr)
	}

	defer rows.Close()

	// Collect the rows into the HV struct
	HVs, collectErr := pgx.CollectRows(rows, pgx.RowToStructByName[HV])

	if collectErr != nil {
		return fmt.Errorf("Error collecting hv: %w", collectErr)
	}

	cloud.mutex.Lock()
	defer cloud.mutex.Unlock()

	// Populate the cloud struct with the HVs with a map of uuids to HVs
	cloud.HVs = make(map[uuid.UUID]*HV)
	for i := range HVs {
		cloud.HVs[HVs[i].ID] = &HVs[i]
		HVs[i].Nics = make(map[uuid.UUID]*HVNic)
		HVs[i].Storages = make(map[uuid.UUID]*HVStorage)
		HVs[i].VMs = make(map[uuid.UUID]*VM)
	}

	return
}

// Initialize the HV libvirt connection
func (hv *HV) Init() error {
	hv.Libvirt = libvirt.Init(hv.IP, hv.Port)
	if err := hv.ensureConn(); err != nil {
		return err
	}
	if err := hv.InitVMs(); err != nil {
		return err
	}
	return nil
}

// Ensure the HV libvirt connection is alive
func (hv *HV) ensureConn() error {
	if !hv.Libvirt.IsConnected() {
		return hv.connect()
	}
	return nil
}

// Connect to the HV libvirt
func (hv *HV) connect() error {
	hv.mutex.Lock()
	defer hv.mutex.Unlock()

	v, err := hv.Libvirt.Connect()
	if err != nil {
		hv.Status = util.StatusUnknown
		hv.StatusReason = err.Error()
		return err
	} else {
		hv.Status = util.StatusRunning
		hv.StatusReason = "Connected to libvirt"
		hv.Version = v
		return nil
	}
}
