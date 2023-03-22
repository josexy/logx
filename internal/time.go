package internal

import (
	"time"
)

type TimeItem struct {
	Enable bool
	Format func(t time.Time) string
	Now    time.Time
	Color  bool
}

func (t *TimeItem) String() (out string) {
	if !t.Enable {
		return
	}
	if t.Format == nil {
		return
	}
	out = t.Format(t.Now)
	if t.Color {
		out = Blue(out)
	}
	return
}
