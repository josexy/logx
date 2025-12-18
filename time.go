package logx

type TimeOption struct {
	// time key, default: time
	TimeKey string
	// convert time to int64 timestamp seconds, default: false
	Timestamp bool
	// time layout, default: time.DateTime
	Layout string
}

type timeField struct {
	enable bool
	color  bool
	option TimeOption
}

func (t *timeField) appendWithJson(enc *JsonEncoder, ent *entry) {
	enc.writeFieldKey(t.option.TimeKey)
	enc.writeSplitColon()

	if !t.option.Timestamp {
		enc.writeQuote()
	}
	t.value(enc.buf, ent)
	if !t.option.Timestamp {
		enc.writeQuote()
	}
}

func (t *timeField) value(buf *Buffer, ent *entry) {
	if t.color {
		if t.option.Timestamp {
			appendColorWithFunc(buf, RedAttr, func() { buf.AppendInt(ent.time.Unix()) })
		} else {
			appendColorWithFunc(buf, BlueAttr, func() { buf.AppendTime(ent.time, t.option.Layout) })
		}
		return
	}
	if t.option.Timestamp {
		buf.AppendInt(ent.time.Unix())
		return
	}
	buf.AppendTime(ent.time, t.option.Layout)
}

func (t *timeField) append(buf *Buffer, ent *entry) {
	t.value(buf, ent)
}
