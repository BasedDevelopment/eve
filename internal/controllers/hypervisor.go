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
	Mutex    sync.Mutex `db:"-"`
	ID       uuid.UUID
	Hostname string
	IP       net.IP
	Port     int
	Site     string
	Nics     []HVNic     `db:"-"`
	Storages []HVStorage `db:"-"`
	VMs      []VM        `db:"-"`
	Created  time.Time
	Updated  time.Time
	Remarks  string
	Status   string `db:"-"`
	Version  string `db:"-"`
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

func GetHVs(cloud *Cloud) (err error) {
	// Reading HVs
	rows, queryErr := db.Pool.Query(context.Background(), "SELECT * FROM hv")

	if queryErr != nil {
		return fmt.Errorf("Error reading hv: %w", queryErr)
	}

	HVs, collectErr := pgx.CollectRows(rows, pgx.RowToStructByName[HV])

	cloud.Mutex.Lock()
	cloud.HVs = make(map[string]*HV)
	for i := range HVs {
		hvid := HVs[i].ID.String()
		cloud.HVs[hvid] = &HVs[i]
	}
	cloud.Mutex.Unlock()

	if collectErr != nil {
		return fmt.Errorf("Error collecting hv: %w", collectErr)
	}

	return
}
