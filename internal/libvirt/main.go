package libvirt

import (
	"fmt"
	"strconv"
	"time"

	"github.com/digitalocean/go-libvirt"
	"github.com/digitalocean/go-libvirt/socket/dialers"
	"github.com/ericzty/eve/internal/controllers"
	"github.com/ericzty/eve/internal/util"
)

func InitHV(HV *controllers.HV) error {
	l := libvirt.NewWithDialer(dialers.NewRemote(
		HV.IP.String(),
		dialers.UsePort(strconv.Itoa(HV.Port)),
		dialers.WithRemoteTimeout(time.Second*2),
	))

	HV.Libvirt = l
	return EnsureConn(HV)
}

func EnsureConn(HV *controllers.HV) error {
	// Ensures the connection is alive
	_, ok := <-HV.Libvirt.Disconnected()
	if !ok {
		// Connection lost
		return connect(HV)
	}
	return nil
}

func connect(HV *controllers.HV) error {
	l := HV.Libvirt
	if err := l.Connect(); err != nil {
		HV.Status = util.STATUS_OFFLINE
		return fmt.Errorf("Failed to communicate with libvirt: %v", err)
	}

	if v, err := l.Version(); err != nil {
		HV.Status = util.STATUS_OFFLINE
		return fmt.Errorf("Failed to get libvirt version: %v", err)
	} else {
		HV.Version = v
		HV.Status = util.STATUS_ONLINE
		return nil
	}
}
