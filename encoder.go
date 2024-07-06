package logx

import "bytes"

type EncoderType byte

const (
	// Simple Encode does't support key-value args, only message field can be set
	Console EncoderType = 1 << iota
	Json
)

type encoder interface {
	Init()
	Encode(*bytes.Buffer, string, ...Field) error
}

type sliceFields struct {
	idx    int
	fields []Field
}

func (l *sliceFields) reset() {
	l.idx = 0
	l.fields = l.fields[:0]
}

func (l *sliceFields) put(fields ...Field) {
	if len(fields) > 0 {
		l.fields = append(l.fields, fields...)
	}
}

func (l *sliceFields) writeRangeFields(fn func(f Field) error, lastfn func()) error {
	n := len(l.fields)
	for i := 0; i < n; i++ {
		if err := fn(l.fields[i]); err != nil {
			return err
		}
		if l.idx+1 != n {
			lastfn()
		}
		l.idx++
	}
	return nil
}
