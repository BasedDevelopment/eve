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

func getStateReason(state int32, reason int32) (stateStr string, reasonStr string) {
	switch state {
	// https://pkg.go.dev/github.com/digitalocean/go-libvirt#DomainState
	case 0:
		stateStr = "NoState"
		switch reason {
		case 0:
			reasonStr = "Unknown"
		}
	case 1:
		stateStr = "Running"
		switch reason {
		// https://pkg.go.dev/github.com/digitalocean/go-libvirt#DomainRunningReason
		case 0:
			reasonStr = "Unknown"
		case 1:
			reasonStr = "Booted"
		case 2:
			reasonStr = "Migrated"
		case 3:
			reasonStr = "Restored"
		case 4:
			reasonStr = "FromSnapshot"
		case 5:
			reasonStr = "Unpaused"
		case 6:
			reasonStr = "MigrationCanceled"
		case 7:
			reasonStr = "SaveCanceled"
		case 8:
			reasonStr = "Wakeup"
		case 9:
			reasonStr = "Crashed"
		case 10:
			reasonStr = "Postcopy"
		}
	case 2:
		stateStr = "Blocked"
		switch reason {
		// https://pkg.go.dev/github.com/digitalocean/go-libvirt#DomainBlockedReason
		case 0:
			reasonStr = "Unknown"
		}
	case 3:
		stateStr = "Paused"
		// https://pkg.go.dev/github.com/digitalocean/go-libvirt#DomainPausedReason
		switch reason {
		case 0:
			reasonStr = "Unknown"
		case 1:
			reasonStr = "User"
		case 2:
			reasonStr = "Migration"
		case 3:
			reasonStr = "Save"
		case 4:
			reasonStr = "Dump"
		case 5:
			reasonStr = "IOError"
		case 6:
			reasonStr = "Watchdog"
		case 7:
			reasonStr = "FromSnapshot"
		case 8:
			reasonStr = "ShuttingDown"
		case 9:
			reasonStr = "Snapshot"
		case 10:
			reasonStr = "Crashed"
		case 11:
			reasonStr = "StartingUp"
		case 12:
			reasonStr = "Postcopy"
		case 13:
			reasonStr = "PostcopyFailed"
		}
	case 4:
		stateStr = "Shutdown"
		// https://pkg.go.dev/github.com/digitalocean/go-libvirt#DomainShutdownReason
		switch reason {
		case 0:
			reasonStr = "Unknown"
		case 1:
			reasonStr = "User"
		}
	case 5:
		stateStr = "Shutoff"
		// https://pkg.go.dev/github.com/digitalocean/go-libvirt#DomainShutoffReason
		switch reason {
		case 0:
			reasonStr = "Unknown"
		case 1:
			reasonStr = "Shutdown"
		case 2:
			reasonStr = "Destroyed"
		case 3:
			reasonStr = "Crashed"
		case 4:
			reasonStr = "Migrated"
		case 5:
			reasonStr = "Saved"
		case 6:
			reasonStr = "Failed"
		case 7:
			reasonStr = "FromSnapshot"
		case 8:
			reasonStr = "Daemon"
		}
	case 6:
		stateStr = "Crashed"
		// https://pkg.go.dev/github.com/digitalocean/go-libvirt#DomainCrashedReason
		switch reason {
		case 0:
			reasonStr = "Unknown"
		case 1:
			reasonStr = "Panicked"
		}
	case 7:
		stateStr = "PMSuspended"
		// https://pkg.go.dev/github.com/digitalocean/go-libvirt#DomainPMSuspendedReason
		switch reason {
		case 0:
			reasonStr = "Unknown"
		}
	default:
		stateStr = "Unknown"
	}
	return stateStr, reasonStr
}
