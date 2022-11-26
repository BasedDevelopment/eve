package controllers

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/ericzty/eve/internal/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type HV struct {
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
	ID        uuid.UUID
	Size      int
	Used      int `db:"-"`
	Available int `db:"-"`
	Created   time.Time
	Updated   time.Time
	Remarks   string
}

func GetHVs() (HVs []HV, err error) {
	// Reading HVs
	rows, queryErr := db.Pool.Query(context.Background(), "SELECT * FROM hv")

	if queryErr != nil {
		return nil, fmt.Errorf("Error reading hv: %w", queryErr)
	}

	HVs, collectErr := pgx.CollectRows(rows, pgx.RowToStructByName[HV])

	if collectErr != nil {
		return nil, fmt.Errorf("Error collecting hv: %w", collectErr)
	}

	return
}
