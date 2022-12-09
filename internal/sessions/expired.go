package sessions

import "time"

func isExpired(expiry time.Time) bool {
	return expiry.After(time.Now())
}
