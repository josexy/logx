package logx

import "sync/atomic"

type LevelType uint8

const (
	LevelTrace LevelType = iota
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
	LevelPanic
)

var (
	levelTypeLowerMap = map[LevelType]string{
		LevelTrace: "trace",
		LevelDebug: "debug",
		LevelInfo:  "info",
		LevelWarn:  "warn",
		LevelError: "error",
		LevelFatal: "fatal",
		LevelPanic: "panic",
	}
	levelTypeUpperMap = map[LevelType]string{
		LevelTrace: "TRACE",
		LevelDebug: "DEBUG",
		LevelInfo:  "INFO",
		LevelWarn:  "WARN",
		LevelError: "ERROR",
		LevelFatal: "FATAL",
		LevelPanic: "PANIC",
	}
	levelTypeColorMap = map[LevelType]ColorAttr{
		LevelTrace: MagentaAttr,
		LevelDebug: HiCyanAttr,
		LevelInfo:  GreenAttr,
		LevelWarn:  YellowAttr,
		LevelError: RedAttr,
		LevelFatal: HiRedAttr,
		LevelPanic: HiYellowAttr,
	}
)

type LevelOption struct {
	// level key, default: "level"
	LevelKey string
	// lower level key
	LowerKey bool
}

type AtomicLevel struct {
	level atomic.Uint32
}

func NewAtomicLevel(level LevelType) *AtomicLevel {
	atomicLevel := &AtomicLevel{}
	atomicLevel.SetLevel(level)
	return atomicLevel
}

func (lvl *AtomicLevel) Level() LevelType {
	if lvl == nil {
		return LevelTrace
	}
	return LevelType(lvl.level.Load())
}

func (lvl *AtomicLevel) SetLevel(level LevelType) {
	if lvl == nil {
		return
	}
	lvl.level.Store(uint32(level))
}

type levelField struct {
	option LevelOption
	enable bool
	color  bool
}

func (lvl *levelField) AppendField(enc *JsonEncoder, level LevelType) {
	enc.writeFieldKey(lvl.option.LevelKey)
	enc.writeSplitColon()
	enc.writeQuote()
	lvl.AppendPrimitive(enc.buf, level)
	enc.writeQuote()
}

func (lvl *levelField) AppendPrimitive(buf *Buffer, level LevelType) {
	var levelStr string
	if lvl.option.LowerKey {
		levelStr = levelTypeLowerMap[level]
	} else {
		levelStr = levelTypeUpperMap[level]
	}
	if lvl.color {
		appendColor(buf, levelTypeColorMap[level], levelStr)
		return
	}
	buf.AppendString(levelStr)
}
