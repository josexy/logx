package logx

import "time"

type TimeOption struct {
	// time key, default: time
	TimeKey string
	// convert time to int64 timestamp with UnixNano, default: false
	Timestamp bool
	// time layout, default: time.DateTime
	Layout string
}

type timeField struct {
	option      TimeOption
	enable      bool
	color       bool
	stringColor ColorAttr
	numberColor ColorAttr
}

func (t *timeField) AppendField(enc *JsonEncoder, ti time.Time) {
	enc.writeFieldKey(t.option.TimeKey)
	enc.writeSplitColon()
	t.AppendTime(enc, ti)
}

func (t *timeField) AppendTime(enc *JsonEncoder, ti time.Time) {
	if !t.option.Timestamp {
		enc.writeQuote()
	}
	t.AppendTimePrimitive(enc.buf, ti)
	if !t.option.Timestamp {
		enc.writeQuote()
	}
}

func (t *timeField) AppendTimePrimitive(buf *Buffer, ti time.Time) {
	if t.color {
		if t.option.Timestamp {
			appendColorWithFunc(buf, t.numberColor, func(buf *Buffer) { buf.AppendInt(ti.UnixNano()) })
		} else {
			appendColorWithFunc(buf, t.stringColor, func(buf *Buffer) { buf.AppendTime(ti, t.option.Layout) })
		}
		return
	}
	if t.option.Timestamp {
		buf.AppendInt(ti.UnixNano())
		return
	}
	buf.AppendTime(ti, t.option.Layout)
}
