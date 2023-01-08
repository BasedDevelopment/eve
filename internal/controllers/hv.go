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
	"github.com/rs/zerolog/log"
)

type HV struct {
	mutex          sync.Mutex            `db:"-" json:"-"`
	ID             uuid.UUID             `json:"id"`
	Hostname       string                `json:"hostname"`
	CPUModel       string                `json:"cpu_model" db:"-"`
	Arch           string                `json:"arch" db:"-"`
	RAMTotal       uint64                `json:"total_ram" db:"-"`
	RAMFree        uint64                `json:"free_ram" db:"-"`
	CPUCount       int32                 `json:"cpu_count" db:"-"`
	CPUFrequency   int32                 `json:"cpu_frequency_mhz" db:"-"`
	NUMANodes      int32                 `json:"numa_nodes" db:"-"`
	CPUSockets     int32                 `json:"cpu_sockets" db:"-"`
	CPUCores       int32                 `json:"cpu_cores" db:"-"`
	CPUThreads     int32                 `json:"cpu_threads" db:"-"`
	IP             net.IP                `json:"ip"`
	Port           int                   `json:"port"`
	Site           string                `json:"site"`
	Brs            map[string]*HVBr      `json:"-" db:"-"`
	Storages       map[string]*HVStorage `json:"-" db:"-"`
	VMs            map[uuid.UUID]*VM     `json:"-" db:"-"`
	Created        time.Time             `json:"created"`
	Updated        time.Time             `json:"updated"`
	Remarks        string                `json:"remarks"`
	Status         util.Status           `json:"status" db:"-"`
	StatusReason   string                `json:"status_reason" db:"-"`
	QemuVersion    string                `json:"qemu_version" db:"-"`
	LibvirtVersion string                `json:"libvirt_version" db:"-"`
	Libvirt        *libvirt.Libvirt      `json:"-" db:"-"`
}

type HVBr struct {
	Name string
}

// WIP
type HVStorage struct {
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
		HVs[i].Brs = make(map[string]*HVBr)
		HVs[i].Storages = make(map[string]*HVStorage)
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

// Get the HV specs
func (hv *HV) getHVSpecs() error {
	if err := hv.ensureConn(); err != nil {
		return err
	}

	qemuVer, err := hv.Libvirt.GetHVQemuVersion()
	if err != nil {
		return err
	}
	hv.QemuVersion = qemuVer

	libvirtVer, err := hv.Libvirt.GetHVLibvirtVersion()
	if err != nil {
		return err
	}
	hv.LibvirtVersion = libvirtVer

	specs, err := hv.Libvirt.GetHVSpecs()
	if err != nil {
		return err
	}

	for _, spec := range specs.Processor.Entry {
		if spec.Name == "version" {
			hv.CPUModel = spec.Text
		}
	}
	if hv.CPUModel == "" {
		log.Warn().
			Str("hv", hv.Hostname).
			Msg("Could not find CPU version")
	}
	return nil
}

// Get HV stats
func (hv *HV) getHVStats() error {
	if err := hv.ensureConn(); err != nil {
		return err
	}

	arch, memoryTotal, memoryFree, cpus, mhz, nodes, sockets, cores, threads, err := hv.Libvirt.GetHVStats()
	if err != nil {
		return err
	}

	hv.Arch = arch
	hv.RAMTotal = memoryTotal
	hv.RAMFree = memoryFree
	hv.CPUCount = cpus
	hv.CPUFrequency = mhz
	hv.NUMANodes = nodes
	hv.CPUSockets = sockets
	hv.CPUCores = cores
	hv.CPUThreads = threads

	return nil
}

func (hv *HV) getHVBrs() error {
	if err := hv.ensureConn(); err != nil {
		return err
	}

	brs, err := hv.Libvirt.GetHVBrs()
	if err != nil {
		return err
	}

	for _, br := range brs {
		hv.Brs[br.Name] = &HVBr{
			Name: br.Name,
		}
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

	err := hv.Libvirt.Connect()
	if err != nil {
		hv.Status = util.StatusUnknown
		hv.StatusReason = err.Error()
		return err
	} else {
		hv.Status = util.StatusRunning
		hv.StatusReason = "Connected to libvirt"
	}
	if err := hv.getHVSpecs(); err != nil {
		return err
	}
	if err := hv.getHVStats(); err != nil {
		return err
	}
	if err := hv.getHVBrs(); err != nil {
		return err
	}
	return nil
}
