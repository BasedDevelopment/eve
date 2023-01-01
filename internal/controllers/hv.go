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
	mutex    sync.Mutex               `db:"-"`
	ID       uuid.UUID                `json:"id"`
	Hostname string                   `json:"hostname"`
	IP       net.IP                   `json:"ip"`
	Port     int                      `json:"port"`
	Site     string                   `json:"site"`
	Nics     map[uuid.UUID]*HVNic     `json:"nics" db:"-"`
	Storages map[uuid.UUID]*HVStorage `json:"storages" db:"-"`
	VMs      map[uuid.UUID]*VM        `json:"vms" db:"-"`
	Created  time.Time                `json:"created"`
	Updated  time.Time                `json:"updated"`
	Remarks  string                   `json:"remarks"`
	Status   util.Status              `json:"status" db:"-"`
	Version  string                   `json:"version" db:"-"`
	Libvirt  *libvirt.Libvirt         `json:"-" db:"-"`
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

	// Collect the rows into the HV struct
	HVs, collectErr := pgx.CollectRows(rows, pgx.RowToStructByName[HV])

	if collectErr != nil {
		return fmt.Errorf("Error collecting hv: %w", collectErr)
	}

	cloud.mutex.Lock()
	defer cloud.mutex.Unlock()

	cloud.HVs = make(map[string]*HV)
	for i := range HVs {
		hvId := HVs[i].ID.String()
		cloud.HVs[hvId] = &HVs[i]
		HVs[i].Nics = make(map[uuid.UUID]*HVNic)
		HVs[i].Storages = make(map[uuid.UUID]*HVStorage)
		HVs[i].VMs = make(map[uuid.UUID]*VM)
	}

	return
}

func (hv *HV) ensureConn() error {
	if !hv.Libvirt.IsConnected() {
		return hv.connect()
	}
	return nil
}

func (hv *HV) connect() error {
	hv.mutex.Lock()
	defer hv.mutex.Unlock()

	v, err := hv.Libvirt.Connect()
	if err != nil {
		hv.Status = util.STATUS_OFFLINE
		return err
	} else {
		hv.Status = util.STATUS_ONLINE
		hv.Version = v
		return nil
	}
}

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

func (hv *HV) InitVMs() error {
	if err := hv.ensureConn(); err != nil {
		return err
	}

	//TODO: Get VMs from db
	uuids, err := hv.Libvirt.GetVMs()
	if err != nil {
		return err
	}
	fmt.Println(uuids)

	//TODO: Compare DB's list of UUIDs and the list from libvirt

	//TODO: Init VMs
	return nil
}
