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
	errorHandler ErrorHandler
}

func NewLogContext() *LogContext {
	return &LogContext{
		level:        NewAtomicLevel(LevelTrace),
		timeF:        timeField{option: TimeOption{Layout: time.DateTime}},
		errorHandler: defaultErrorHandler,
	}
}

// Copy snapshots the logger configuration with an independent dynamic level.
func (lc *LogContext) Copy() *LogContext {
	return lc.copy(false)
}

func (lc *LogContext) copy(shareLevel bool) *LogContext {
	newLogCtx := new(LogContext)
	*newLogCtx = *lc
	if shareLevel {
		newLogCtx.level = lc.AtomicLevel()
	} else {
		newLogCtx.level = NewAtomicLevel(lc.AtomicLevel().Level())
	}
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

func (lc *LogContext) copySharedLevel() *LogContext {
	return lc.copy(true)
}

func (lc *LogContext) WithFields(fields ...Field) *LogContext {
	lc.preFields = append(lc.preFields, fields...)
	return lc
}

// WithNewFields replaces inherited fields and copies the supplied slice.
func (lc *LogContext) WithNewFields(fields ...Field) *LogContext {
	lc.preFields = append([]Field(nil), fields...)
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

// WithEscapeQuote controls string escaping in the JSON fragments embedded in
// Console output. The Json encoder always performs standards-compliant escaping.
func (lc *LogContext) WithEscapeQuote(enable bool) *LogContext {
	lc.escapeQuote = enable
	return lc
}

func (lc *LogContext) WithReflectValue(enable bool) *LogContext {
	lc.reflectValue = enable
	return lc
}

// WithErrorHandler configures how internal encoding, write, and automatic
// high-severity sync errors are reported. Passing nil restores the default
// handler, which writes the error to standard error.
func (lc *LogContext) WithErrorHandler(handler ErrorHandler) *LogContext {
	if handler == nil {
		handler = defaultErrorHandler
	}
	lc.errorHandler = handler
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
		lc.enc = &JsonEncoder{LogContext: lc, jsonOutput: true}
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

// Build returns a logger with a snapshot of the current configuration. The
// configured AtomicLevel remains shared for runtime updates.
func (lc *LogContext) Build() Logger {
	built := lc.copySharedLevel()
	if built.level == nil {
		built.level = NewAtomicLevel(LevelTrace)
	}
	if built.timeF.option.Layout == "" {
		built.timeF.option.Layout = time.DateTime
	}
	built.WithMsgKey(built.msgKey)
	if built.errorHandler == nil {
		built.errorHandler = defaultErrorHandler
	}
	if built.enc != nil {
		built.enc.Init()
	}
	return &LoggerX{logCtx: built}
}
