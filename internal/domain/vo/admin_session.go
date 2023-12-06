package vo

import "time"

type AdminSession struct {
	Session string
	Expiry  time.Time
}
