package libvirt

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/digitalocean/go-libvirt"
	"github.com/ericzty/eve/internal/db"
)

func Init(HV *db.HV) (err error) {
	c, err := net.DialTimeout("tcp", HV.IP.String()+":"+strconv.Itoa(HV.Port), 3*time.Second)
	if err != nil {
		return fmt.Errorf("Failed to dial to libvirt tcp socket: %v", err)
	}
	l := libvirt.New(c)
	if err := l.Connect(); err != nil {
		return fmt.Errorf("Failed to communicate with libvirt: %v", err)
	}
	v, err := l.Version()
	if err != nil {
		return fmt.Errorf("Failed to get libvirt version: %v", err)
	}
	defer l.Disconnect()
	defer c.Close()
	HV.Version = v
	HV.Status = "online"
	return
}
