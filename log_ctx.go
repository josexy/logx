package logx

import (
	"sync"
	"time"
)

var bufPool = sync.Pool{
	New: func() any { return NewBuffer(make([]byte, 0, 1024)) },
}

type LogContext struct {
	level        *AtomicLevel
	levelF       levelField
	timeF        timeField
	callerF      callerField
	enc          encoder
	colors       colorfulset
	writer       WriteSyncer
	preFields    []Field
	msgKey       string
	escapeQuote  bool
	reflectValue bool
}

func NewLogContext() *LogContext {
	return &LogContext{level: NewAtomicLevel(LevelTrace)}
}

func (lc *LogContext) Copy() *LogContext {
	newLogCtx := new(LogContext)
	*newLogCtx = *lc
	newLogCtx.level = NewAtomicLevel(lc.AtomicLevel().Level())
	if len(lc.preFields) > 0 {
		newLogCtx.preFields = make([]Field, 0, len(lc.preFields))
		newLogCtx.preFields = append(newLogCtx.preFields, lc.preFields...)
	}
	switch lc.enc.(type) {
	case *JsonEncoder:
		newLogCtx = newLogCtx.WithEncoder(Json)
	case *ConsoleEncoder:
		newLogCtx = newLogCtx.WithEncoder(Console)
	case *TextEncoder:
		newLogCtx = newLogCtx.WithEncoder(Text)
	}
	return newLogCtx
}

func (lc *LogContext) WithFields(fields ...Field) *LogContext {
	lc.preFields = append(lc.preFields, fields...)
	return lc
}

func (lc *LogContext) WithNewFields(fields ...Field) *LogContext {
	lc.preFields = fields
	return lc
}

func (lc *LogContext) WithColorfulset(enable bool, attr TextColorAttri) *LogContext {
	if NoColor {
		enable = false
	}
	lc.levelF.color = enable
	lc.timeF.color = enable
	lc.callerF.color = enable
	lc.colors.enable = enable
	lc.colors.attr = attr
	return lc
}

func (lc *LogContext) WithMsgKey(key string) *LogContext {
	if len(key) == 0 {
		key = "msg"
	}
	if lc.msgKey != key {
		lc.msgKey = key
	}
	return lc
}

func (lc *LogContext) WithLevelKey(enable bool, option LevelOption) *LogContext {
	lc.levelF.enable = enable
	if enable {
		if len(option.LevelKey) == 0 {
			option.LevelKey = "level"
		}
		lc.levelF.option = option
	}
	return lc
}

func (lc *LogContext) WithTimeKey(enable bool, option TimeOption) *LogContext {
	lc.timeF.enable = enable
	if enable {
		if len(option.TimeKey) == 0 {
			option.TimeKey = "time"
		}
		if len(option.Layout) == 0 {
			option.Layout = time.DateTime
		}
		lc.timeF.option = option
	}
	return lc
}

func (lc *LogContext) WithCallerKey(enable bool, option CallerOption) *LogContext {
	lc.callerF.enable = enable
	if enable {
		if len(option.CallerKey) == 0 {
			option.CallerKey = "caller"
		}
		if len(option.FileKey) == 0 {
			option.FileKey = "file"
		}
		if len(option.FuncKey) == 0 {
			option.FuncKey = "func"
		}
		if option.Formatter > FullFileFunc {
			option.Formatter = FullFileFunc
		}
		lc.callerF.option = option
	}
	return lc
}

func (lc *LogContext) WithEscapeQuote(enable bool) *LogContext {
	lc.escapeQuote = enable
	return lc
}

func (lc *LogContext) WithReflectValue(enable bool) *LogContext {
	lc.reflectValue = enable
	return lc
}

func (lc *LogContext) WithWriter(writer WriteSyncer) *LogContext {
	lc.writer = writer
	return lc
}

func (lc *LogContext) WithEncoder(encoder EncoderType) *LogContext {
	switch encoder {
	case Console:
		lc.enc = &ConsoleEncoder{LogContext: lc}
	case Json:
		lc.enc = &JsonEncoder{LogContext: lc}
	case Text:
		lc.enc = &TextEncoder{LogContext: lc}
	default:
		panic("not support other log encoder")
	}
	return lc
}

func (lc *LogContext) WithLevel(level LevelType) *LogContext {
	lc.AtomicLevel().SetLevel(level)
	return lc
}

func (lc *LogContext) WithAtomicLevel(level *AtomicLevel) *LogContext {
	if level == nil {
		level = NewAtomicLevel(LevelTrace)
	}
	lc.level = level
	return lc
}

func (lc *LogContext) AtomicLevel() *AtomicLevel {
	if lc.level == nil {
		lc.level = NewAtomicLevel(LevelTrace)
	}
	return lc.level
}

func (lc *LogContext) Build() Logger {
	lc.WithMsgKey(lc.msgKey)
	if lc.enc != nil {
		lc.enc.Init()
	}
	return &LoggerX{logCtx: lc}
}
