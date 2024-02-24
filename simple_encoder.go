package logx

import (
	"bytes"
)

type SimpleEncoder struct {
	*logContext
}

func (enc *SimpleEncoder) Init() {
	if enc.callerField.enable {
		enc.callerField.skipDepth = 6
	}
}

func (enc *SimpleEncoder) Encode(w *bytes.Buffer, msg string, _ ...arg) error {
	if enc.levelField.enable {
		w.WriteByte('[')
		w.WriteString(enc.levelField.String())
		w.WriteByte(']')
	}
	if enc.timeField.enable {
		w.WriteByte('[')
		w.WriteString(enc.timeField.String())
		w.WriteByte(']')
	}
	if enc.callerField.enable {
		w.WriteByte('[')
		w.WriteString(enc.callerField.String())
		w.WriteByte(']')
	}
	if enc.levelField.enable || enc.timeField.enable || enc.callerField.enable {
		w.WriteByte(' ')
	}
	_, err := w.WriteString(msg)
	return err
}
