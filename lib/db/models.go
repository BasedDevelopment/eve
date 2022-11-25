package db

import (
	"net"
	"time"

	"github.com/google/uuid"
)

type Profile struct {
	ID        uuid.UUID
	Name      string
	Email     string
	Password  string
	Disabled  bool
	LastLogin time.Time `db:"last_login"`
	Created   time.Time
	Updated   time.Time
	Remarks   string
}

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
