package controllers

import (
	"net"
	"time"

	"github.com/google/uuid"
)

type VM struct {
	ID       uuid.UUID
	Hostname string
	User     *Profile `db:"-"`
	CPU      int
	Memory   int
	Nics     []VMNic     `db:"-"`
	Storages []VMStorage `db:"-"`
	Created  time.Time
	Updated  time.Time
	Remarks  string
	State    string `db:"-"`
}

type VMNic struct {
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
	ID      uuid.UUID
	Size    int
	Created time.Time
	Updated time.Time
	Remarks string
}
