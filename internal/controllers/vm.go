package controllers

import (
	"net"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgtype"
)

type VM struct {
	ID       uuid.UUID
	Hostname string
	User     *Profile `db:"-"`
	CPU      int
	Memory   int
	Nics     map[string]VMNic     `db:"-"`
	Storages map[string]VMStorage `db:"-"`
	Created  time.Time
	Updated  time.Time
	Remarks  pgtype.Text
	State    string `db:"-"`
}

type VMNic struct {
	ID      uuid.UUID
	name    string
	MAC     string
	IP      []net.IP `db:"-"`
	Created time.Time
	Updated time.Time
	Remarks pgtype.Text
	State   string `db:"-"`
}

type VMStorage struct {
	ID      uuid.UUID
	Size    int
	Created time.Time
	Updated time.Time
	Remarks pgtype.Text
}
