package internal

import (
	"bytes"
)

type SimpleEncoder struct {
	*Config
}

func (enc *SimpleEncoder) Init() {
	if enc.CallerItem.Enable {
		enc.CallerItem.Skip = 4
	}
}

func (enc *SimpleEncoder) Encode(w *bytes.Buffer, msg string, _ ...Arg) error {
	if enc.LevelItem.Enable {
		w.WriteByte('[')
		w.WriteString(enc.LevelItem.String())
		w.WriteByte(']')
	}
	if enc.TimeItem.Enable {
		w.WriteByte('[')
		w.WriteString(enc.TimeItem.String())
		w.WriteByte(']')
	}
	if enc.CallerItem.Enable {
		w.WriteByte('[')
		w.WriteString(enc.CallerItem.String())
		w.WriteByte(']')
	}
	if enc.LevelItem.Enable || enc.TimeItem.Enable || enc.CallerItem.Enable {
		w.WriteByte(' ')
	}
	_, err := w.WriteString(msg)
	return err
}
