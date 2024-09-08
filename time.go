package logx

import (
	"fmt"
	"strconv"
	"time"
)

type TimeOption struct {
	// time key, default: time
	TimeKey string
	// time formatter, return value: string, int64, time.Time, any
	Formatter func(t time.Time) any
}

type timeField struct {
	enable bool
	color  bool
	now    time.Time
	option TimeOption
}

func (t *timeField) formatJson(enc *JsonEncoder) {
	enc.writeFieldKey(t.option.TimeKey)
	enc.buf.WriteByte(':')

	switch val := t.value(); val.(type) {
	case string:
		enc.writeFieldString(val.(string))
	case int64:
		enc.writeFieldInt64(val.(int64))
	default:
		if ts, ok := val.(time.Time); ok {
			enc.writeFieldTime(ts)
		} else {
			enc.writeFieldString(fmt.Sprintf("%v", val))
		}
	}
}

func (t *timeField) value() (out any) {
	if !t.enable {
		return
	}
	out = t.option.Formatter(t.now)
	return
}

func (t *timeField) String() string {
	var out string
	switch val := t.value(); val.(type) {
	case string:
		out = val.(string)
	case int64:
		out = strconv.FormatInt(val.(int64), 10)
	default:
		if ts, ok := val.(time.Time); ok {
			out = ts.Format(time.DateTime)
		} else {
			out = fmt.Sprintf("%v", val)
		}
	}
	if t.color && len(out) > 0 {
		out = Blue(out)
	}
	return out
}
