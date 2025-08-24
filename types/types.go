package types

import (
	"time"

	"github.com/fatih/color"
)

type DailyActivity struct {
	Count int64
	Date  time.Time
}

type NamedColor struct {
	Name string
	Attr color.Attribute
}
