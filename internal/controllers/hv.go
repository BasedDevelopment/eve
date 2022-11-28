package controllers

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/ericzty/eve/internal/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type HV struct {
	mutex    sync.Mutex            `db:"-"`
	ID       uuid.UUID             `json:"id"`
	Hostname string                `json:"hostname"`
	IP       net.IP                `json:"ip"`
	Port     int                   `json:"port"`
	Site     string                `json:"site"`
	Nics     map[string]*HVNic     `json:"nics" db:"-"`
	Storages map[string]*HVStorage `json:"storages" db:"-"`
	VMs      map[string]*VM        `json:"vms" db:"-"`
	Created  time.Time             `json:"created"`
	Updated  time.Time             `json:"updated"`
	Remarks  string                `json:"remarks"`
	Status   string                `json:"status" db:"-"`
	Version  string                `json:"version" db:"-"`
}

type HVNic struct {
	Mutex   sync.Mutex `db:"-"`
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
	Mutex     sync.Mutex `db:"-"`
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

	HVs, collectErr := pgx.CollectRows(rows, pgx.RowToStructByName[HV])

	cloud.mutex.Lock()
	defer cloud.mutex.Unlock()
	cloud.HVs = make(map[string]*HV)
	for i := range HVs {
		hvid := HVs[i].ID.String()
		cloud.HVs[hvid] = &HVs[i]
	}

	if collectErr != nil {
		return fmt.Errorf("Error collecting hv: %w", collectErr)
	}

	return
}
