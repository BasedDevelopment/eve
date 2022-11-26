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
	/*
		inputRows := [][]interface{}{
			{"85833bb8-2f0a-4b1e-981f-f9cb3597904c", "dorm0.sit.eric.si", "10.10.9.4", 16509, "sit"},
			{"70ada0c5-b641-4633-88bc-2d58c8a387a5", "broke.sit.eric.si", "10.10.222.1", 16509, "sit"},
		}

		copyCount, err := pool.CopyFrom(context.Background(), pgx.Identifier{"hv"}, []string{"id", "hostname", "ip", "port", "site"}, pgx.CopyFromRows(inputRows))
		if err != nil {
			fmt.Println("Error copying from rows: ", err)
		}
		if int(copyCount) != len(inputRows) {
			fmt.Println("Len mismatch", copyCount, len(inputRows))
		} else {
			fmt.Println("Copied", copyCount, "rows")
		}
	*/

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
