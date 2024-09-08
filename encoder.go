package logx

type EncoderType byte

const (
	Console EncoderType = 1 << iota
	Json
)

type encoder interface {
	Init()
	Encode(*Buffer, string, ...Field) error
}

type sliceFields struct {
	idx    int
	fields []Field
}

func (l *sliceFields) size() int { return len(l.fields) }

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
