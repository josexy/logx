package logx

import (
	"fmt"
	"strconv"
	"time"
)

type timeField struct {
	enable bool
	fn     func(t time.Time) any
	now    time.Time
	color  bool
}

func (t *timeField) format() Field {
	var field Field
	switch val := t.value(); val.(type) {
	case string:
		field = String("time", val.(string))
	case int64:
		field = Int64("time", val.(int64))
	default:
		if t, ok := val.(time.Time); ok {
			field = Time("time", t)
		} else {
			field = String("time", fmt.Sprintf("%v", val))
		}
	}
	return field
}

func (t *timeField) value() (out any) {
	if !t.enable {
		return
	}
	if t.fn == nil {
		return
	}
	out = t.fn(t.now)
	return
}

func (t *timeField) String() string {
	field := t.format()
	var out string
	switch field.Type {
	case StringType:
		out = field.StringValue
	case Int64Type:
		out = strconv.FormatInt(field.IntValue, 10)
	}
	if t.color && len(out) > 0 {
		out = Blue(out)
	}
	return out
}
