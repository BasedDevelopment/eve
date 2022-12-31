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

func InitHV(ip net.IP, port int) (l *libvirt.Libvirt) {
	l = libvirt.NewWithDialer(dialers.NewRemote(
		ip.String(),
		dialers.UsePort(strconv.Itoa(port)),
		dialers.WithRemoteTimeout(time.Second*2),
	))

	return
}

func IsConnected(l *libvirt.Libvirt) bool {
	// Check if the connection is alive
	ok := true
	select {
	case _, ok = <-l.Disconnected():
	default:
	}
	return ok
}

func Connect(l *libvirt.Libvirt) (error, string) {
	if err := l.Connect(); err != nil {
		err = fmt.Errorf("Failed to communicate with libvirt: %v", err)
		return err, ""
	}

	v, err := l.Version()
	if err != nil {
		err = fmt.Errorf("Failed to get libvirt version: %v", err)
		return err, ""
	}
	return nil, v
}

func GetVMs(l *libvirt.Libvirt) (vms []uuid.UUID, err error) {
	// Fetcheds list of all defined domains
	doms, _, err := l.ConnectListAllDomains(1, libvirt.ConnectListDomainsPersistent)
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
