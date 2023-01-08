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
	"encoding/xml"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/BasedDevelopment/eve/internal/util"
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

func (l Libvirt) Connect() error {
	if err := l.conn.Connect(); err != nil {
		return fmt.Errorf("failed to communicate with libvirt: %v", err)
	}
	return nil
}

func (l Libvirt) Close() error {
	return l.conn.Disconnect()
}

func (l Libvirt) GetHVQemuVersion() (string, error) {
	v, err := l.conn.ConnectGetVersion()
	if err != nil {
		return "", fmt.Errorf("failed to get libvirt version: %v", err)
	}
	return util.VerFromDec(int(v)), nil
}

func (l Libvirt) GetHVLibvirtVersion() (string, error) {
	v, err := l.conn.ConnectGetLibVersion()
	if err != nil {
		return "", fmt.Errorf("failed to get libvirt version: %v", err)
	}
	return util.VerFromDec(int(v)), nil
}

func (l Libvirt) GetHVSpecs() (specs HVSpecs, err error) {
	hvXml, err := l.conn.ConnectGetSysinfo(0)
	if err != nil {
		return
	}

	hvXmlBytes := []byte(hvXml)
	err = xml.Unmarshal([]byte(hvXmlBytes), &specs)
	if err != nil {
		return
	}
	return
}

func (l Libvirt) GetHVStats() (arch string, memoryTotal uint64, memoryFree uint64, cpus int32, mhz int32, nodes int32, sockets int32, cores int32, threads int32, err error) {
	// Go returns memory in KiB
	model, memoryKiB, cpus, mhz, nodes, sockets, cores, threads, err := l.conn.NodeGetInfo()
	if err != nil {
		return
	}

	// Libvirt returns the arch as [32]int8, so we need to convert it to a string
	archl := []string{}
	for _, i := range model {
		if i == 0 {
			continue
		}
		archl = append(archl, strconv.Itoa(int(i)))
	}
	arch = strings.Join(archl, "")

	// Convirt memory to bytes
	memoryTotal = memoryKiB * 1024

	memoryFree, err = l.conn.NodeGetFreeMemory()
	if err != nil {
		return
	}

	return
}

func (l Libvirt) GetHVBrs() (hvnics []HVNicSpecs, err error) {
	// List active interfaces
	nics, _, err := l.conn.ConnectListAllInterfaces(0, 2)
	if err != nil {
		return
	}

	for i := range nics {
		// Get the XML for each interface and marshall it into a list of HVNicSpecs
		nicxml, err := l.conn.InterfaceGetXMLDesc(nics[i], 0)

		nicXmlBytes := []byte(nicxml)

		var nic HVNicSpecs
		err = xml.Unmarshal([]byte(nicXmlBytes), &nic)
		if err != nil {
			return hvnics, err
		}
		hvnics = append(hvnics, nic)
	}
	return
}
