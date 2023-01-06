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

	"github.com/BasedDevelopment/eve/internal/util"
	"github.com/digitalocean/go-libvirt"
	"github.com/google/uuid"
)

// Domain is a Virtual Machine in libvirt
type Dom struct {
	Dom libvirt.Domain
}

// Fetches list of all defined domains
// Won't be used to populate the HV's VM list, instead to check for inconsistencies
func (l Libvirt) GetVMs() (vms map[uuid.UUID]Dom, err error) {
	doms, _, err := l.conn.ConnectListAllDomains(1, libvirt.ConnectListDomainsPersistent)
	if err != nil {
		return
	}
	vms = make(map[uuid.UUID]Dom)
	for _, dom := range doms {
		vmuuidstr := hex.EncodeToString(dom.UUID[:])
		vmuuid := uuid.MustParse(vmuuidstr)
		vms[vmuuid] = Dom{dom}
	}
	return
}

// Fetch list of all undefined domains
func (l Libvirt) GetUndefinedVMs() (vms map[uuid.UUID]Dom, err error) {
	doms, _, err := l.conn.ConnectListAllDomains(1, libvirt.ConnectListDomainsTransient)
	if err != nil {
		return
	}
	vms = make(map[uuid.UUID]Dom)
	for _, dom := range doms {
		vmuuidstr := hex.EncodeToString(dom.UUID[:])
		vmuuid := uuid.MustParse(vmuuidstr)
		vms[vmuuid] = Dom{dom}
	}
	return
}

// Fetches a domain from a UUID
func (l Libvirt) GetVMFromUUID(vmId uuid.UUID) (dom Dom, err error) {
	vmIdHex, _ := hex.DecodeString(vmId.String())
	var libvirtUUID libvirt.UUID
	copy(libvirtUUID[:], vmIdHex[:])
	domain, err := l.conn.DomainLookupByUUID(libvirtUUID)
	dom = Dom{domain}
	return
}

// Fetch VM specs from libvirt, will be used to check consistency
func (l Libvirt) GetVMSpecs(dom Dom) (specs DomSpecs, err error) {
	domXml, err := l.conn.DomainGetXMLDesc(dom.Dom, 0)
	if err != nil {
		return
	}

	domXmlBytes := []byte(domXml)
	err = xml.Unmarshal([]byte(domXmlBytes), &specs)
	if err != nil {
		return
	}
	return
}

// Get the state of a domain(vm)
func (l Libvirt) GetVMState(dom Dom) (stateInt util.Status, stateStr string, reasonStr string, err error) {
	//https://libvirt.org/html/libvirt-libvirt-domain.html#virDomainState
	//https://libvirt.org/html/libvirt-libvirt-domain.html#virDomainRunningReason
	lState, lReason, err := l.conn.DomainGetState(dom.Dom, 0)
	stateInt = util.Status(lState)
	stateStr, reasonStr = getStateReason(lState, lReason)
	return
}

func (l Libvirt) VMStart(dom Dom) (err error) {
	return l.conn.DomainCreate(dom.Dom)
}

func (l Libvirt) VMReboot(dom Dom) error {
	return l.conn.DomainReboot(dom.Dom, libvirt.DomainRebootDefault)
}

func (l Libvirt) VMPowerOff(dom Dom) error {
	return l.conn.DomainShutdownFlags(dom.Dom, libvirt.DomainShutdownDefault)
}

func (l Libvirt) VMStop(dom Dom) error {
	return l.conn.DomainDestroy(dom.Dom)
}

func (l Libvirt) VMReset(dom Dom) error {
	return l.conn.DomainReset(dom.Dom, 0)
}
