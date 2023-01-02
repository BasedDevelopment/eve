/*
 * eve - management toolkit for libvirt servers
 * Copyright (C) 2022-2023  BNS Services LLC

 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.

 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.

 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package libvirt

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/digitalocean/go-libvirt"
	"github.com/digitalocean/go-libvirt/socket/dialers"
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
