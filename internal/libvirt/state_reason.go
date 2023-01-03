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
