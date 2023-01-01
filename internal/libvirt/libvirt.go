package libvirt

import (
	"encoding/hex"
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

// Initializes a Libvirt object for later connections
func Init(ip net.IP, port int) *Libvirt {
	conn := libvirt.NewWithDialer(dialers.NewRemote(
		ip.String(),
		dialers.UsePort(strconv.Itoa(port)),
		dialers.WithRemoteTimeout(time.Second*2),
	))

	return &Libvirt{conn}
}

func (l Libvirt) IsConnected() bool {
	return l.conn.IsConnected()
}

func (l Libvirt) Connect() (string, error) {
	if err := l.conn.Connect(); err != nil {
		return "", fmt.Errorf("failed to communicate with libvirt: %v", err)
	}

	v, err := l.conn.ConnectGetVersion()
	if err != nil {
		return "", fmt.Errorf("failed to get libvirt version: %v", err)
	}

	return strconv.FormatInt(int64(v), 10), nil
}

func (l Libvirt) GetVMs() (vms []uuid.UUID, err error) {
	// Fetches list of all defined domains
	// Won't be used to populate the HV's VM list, instead to check for inconsistencies
	doms, _, err := l.conn.ConnectListAllDomains(1, libvirt.ConnectListDomainsPersistent)
	if err != nil {
		return
	}
	for _, dom := range doms {
		vmUuidStr := hex.EncodeToString(dom.UUID[:])
		vmUuid, err := uuid.Parse(vmUuidStr)
		if err != nil {
			return vms, err
		}
		vms = append(vms, vmUuid)
	}
	return
}
