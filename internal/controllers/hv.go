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
	"sync"
	"time"

	"github.com/BasedDevelopment/auto/pkg/models"
	"github.com/BasedDevelopment/eve/internal/auto"
	"github.com/BasedDevelopment/eve/internal/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type HV struct {
	ID         uuid.UUID         `json:"id"`
	Hostname   string            `json:"hostname"`
	AutoUrl    string            `json:"auto_url" db:"auto_url"`
	AutoSerial string            `json:"auto_serial" db:"auto_serial"`
	Site       string            `json:"site"`
	Created    time.Time         `json:"created"`
	Updated    time.Time         `json:"updated"`
	Remarks    string            `json:"remarks"`
	Mutex      sync.Mutex        `json:"-" db:"-"`
	VMs        map[uuid.UUID]*VM `json:"-" db:"-"`
	Auto       *auto.Auto        `json:"-" db:"-"`
	Libvirt    *models.HV        `json:"-" db:"-"`
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

	cloud.Mutex.Lock()
	defer cloud.Mutex.Unlock()

	// Populate the cloud struct with the HVs with a map of uuids to HVs
	cloud.HVs = make(map[uuid.UUID]*HV)
	for i := range HVs {
		cloud.HVs[HVs[i].ID] = &HVs[i]
		HVs[i].Auto = &auto.Auto{
			Url:    "https://" + HVs[i].AutoUrl,
			Serial: HVs[i].AutoSerial,
		}
		HVs[i].VMs = make(map[uuid.UUID]*VM)
	}

	return
}

// Initialize the HV libvirt connection
func (hv *HV) Init() error {
	if err := hv.Refresh(); err != nil {
		return err
	}

	if err := hv.InitVMs(); err != nil {
		return err
	}

	return nil
}

func (hv *HV) Refresh() error {
	hv.Mutex.Lock()
	defer hv.Mutex.Unlock()

	if libvirt, err := hv.Auto.GetLibvirt(); err != nil {
		return err
	} else {
		hv.Libvirt = &libvirt
	}
	return nil
}
