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
	"encoding/hex"
	"encoding/xml"
	"fmt"

	"github.com/digitalocean/go-libvirt"
	"github.com/google/uuid"
)

type Dom struct {
	Dom libvirt.Domain
}

func (l Libvirt) GetVMs() (doms []libvirt.Domain, err error) {
	// Fetches list of all defined domains
	// Won't be used to populate the HV's VM list, instead to check for inconsistencies
	doms, _, err = l.conn.ConnectListAllDomains(1, libvirt.ConnectListDomainsPersistent)
	if err != nil {
		return
	}
	for _, dom := range doms {
		l.GetVMSpecs(dom)
		l.GetVMState(dom)
	}
	return
}

func (l Libvirt) GetVMFromUUID(vmId uuid.UUID) (dom Dom, err error) {
	// Fetches a domain from a UUID
	vmIdHex, _ := hex.DecodeString(vmId.String())
	var libvirtUUID libvirt.UUID
	copy(libvirtUUID[:], vmIdHex[:])
	domain, err := l.conn.DomainLookupByUUID(libvirtUUID)
	dom = Dom{domain}
	return
}

func (l Libvirt) GetVMSpecs(dom libvirt.Domain) (err error) {
	domXml, err := l.conn.DomainGetXMLDesc(dom, 0)
	if err != nil {
		return
	}
	fmt.Println(domXml)
	var specs domSpecs
	domXmlBytes := []byte(domXml)
	err = xml.Unmarshal([]byte(domXmlBytes), &specs)
	if err != nil {
		return
	}
	fmt.Println(specs)
	return
}

func (l Libvirt) GetVMState(dom libvirt.Domain) (state int32, reason int32, err error) {
	state, reason, err = l.conn.DomainGetState(dom, 0)
	fmt.Println("state")
	fmt.Println(state, reason)
	return
}
