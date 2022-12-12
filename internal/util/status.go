package util

type Status int64

const (
	STATUS_OFFLINE Status = iota
	STATUS_ONLINE
	STATUS_SUSPENDED
)

func (s Status) String() string {
	switch s {
	case STATUS_OFFLINE:
		return "offline"
	case STATUS_ONLINE:
		return "online"
	case STATUS_SUSPENDED:
		return "suspended"
	}

	return "unknown"
}
