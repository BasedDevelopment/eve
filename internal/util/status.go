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
	StatusUnknown Status = iota
	StatusRunning
	StatusBlocked
	StatusPaused
	StatusShutdown
	StatusShutoff
	StatusCrashed
	StatusPMSuspended
)

func (s Status) String() string {
	switch s {
	case StatusUnknown:
		return "unknown"
	case StatusRunning:
		return "running"
	case StatusBlocked:
		return "blocked"
	case StatusPaused:
		return "paused"
	case StatusShutdown:
		return "shutdown"
	case StatusShutoff:
		return "shutoff"
	case StatusCrashed:
		return "crashed"
	case StatusPMSuspended:
		return "pmsuspended"
	default:
		return "unknown"
	}
}
