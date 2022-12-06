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

func InitHVs(HV *controllers.HV) (err error) {
	l := libvirt.NewWithDialer(dialers.NewRemote(
		HV.IP.String(),
		dialers.UsePort(strconv.Itoa(HV.Port)),
		dialers.WithRemoteTimeout(time.Second*2),
	))

	if err := l.Connect(); err != nil {
		HV.Status = util.STATUS_OFFLINE
		return fmt.Errorf("Failed to communicate with libvirt: %v", err)
	}

	v, err := l.Version()

	if err != nil {
		HV.Status = util.STATUS_OFFLINE
		return fmt.Errorf("Failed to get libvirt version: %v", err)
	}

	defer l.Disconnect()

	HV.Status = util.STATUS_ONLINE
	HV.Version = v

	return
}
