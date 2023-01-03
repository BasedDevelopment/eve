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

package util

type Status int8

// Status constants for HVs, VMs, etc..
// https://libvirt.org/html/libvirt-libvirt-domain.html#virDomainState
const (
	STATUS_UNKNOWN Status = iota
	STATUS_RUNNING
	STATUS_BLOCKED
	STATUS_PAUSED
	STATUS_SHUTDOWN
	STATUS_SHUTOFF
	STATUS_CRASHED
	STATUS_PMSUSPENDED
	STATUS_LAST
)

func (s Status) String() string {
	switch s {
	case STATUS_UNKNOWN:
		return "unknown"
	case STATUS_RUNNING:
		return "running"
	case STATUS_BLOCKED:
		return "blocked"
	case STATUS_PAUSED:
		return "paused"
	case STATUS_SHUTDOWN:
		return "shutdown"
	case STATUS_SHUTOFF:
		return "shutoff"
	case STATUS_CRASHED:
		return "crashed"
	case STATUS_PMSUSPENDED:
		return "pmsuspended"
	case STATUS_LAST:
		return "last"
	default:
		return "unknown"
	}
}
