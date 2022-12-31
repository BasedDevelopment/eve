package libvirt

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/digitalocean/go-libvirt"
	"github.com/digitalocean/go-libvirt/socket/dialers"
	"github.com/google/uuid"
)

type Libvirt struct {
	conn *libvirt.Libvirt
}

func InitHV(ip net.IP, port int) *Libvirt {
	conn := libvirt.NewWithDialer(dialers.NewRemote(
		ip.String(),
		dialers.UsePort(strconv.Itoa(port)),
		dialers.WithRemoteTimeout(time.Second*2),
	))

	return &Libvirt{conn}
}

func (l Libvirt) IsConnected() bool {
	// Check if the connection is alive
	ok := true
	select {
	case _, ok = <-l.conn.Disconnected():
	default:
	}
	return ok
}

func (l Libvirt) Connect() (error, string) {
	if err := l.conn.Connect(); err != nil {
		err = fmt.Errorf("Failed to communicate with libvirt: %v", err)
		return err, ""
	}

	v, err := l.conn.Version()
	if err != nil {
		err = fmt.Errorf("Failed to get libvirt version: %v", err)
		return err, ""
	}
	return nil, v
}

func (l Libvirt) GetVMs() (vms []uuid.UUID, err error) {
	// Fetches list of all defined domains
	// Won't be used to populate the HV's VM list, instead to check for inconsistencies
	doms, _, err := l.conn.ConnectListAllDomains(1, libvirt.ConnectListDomainsPersistent)
	if err != nil {
		return
	}
	for _, dom := range doms {
		vmuuidbytes := fmt.Sprintf("%x", dom.UUID)
		vmuuid, _ := uuid.Parse(vmuuidbytes)
		vms = append(vms, vmuuid)
	}
	return
}

//func (l Libvirt) GetVM(vmid string) (err error) {
//	vmid := uuid.Parse(vmid)
//	dom := l.conn.DomainLookupByUUIDString(vmid)
//}
