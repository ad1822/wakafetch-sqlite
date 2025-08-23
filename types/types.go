package types

import "time"

type DailyActivity struct {
	Count int64
	Date  time.Time
}
