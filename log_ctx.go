package logx

import (
	"io"
	"sync"
	"time"
)

type LogContext struct {
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
	return &LogContext{}
}

func (lc *LogContext) Copy() *LogContext {
	newLogCtx := new(LogContext)
	*newLogCtx = *lc
	if len(lc.preFields) > 0 {
		newLogCtx.preFields = make([]Field, 0, len(lc.preFields))
		newLogCtx.preFields = append(newLogCtx.preFields, lc.preFields...)
	}
	switch lc.enc.(type) {
	case *JsonEncoder:
		newLogCtx = newLogCtx.WithEncoder(Json)
	case *ConsoleEncoder:
		newLogCtx = newLogCtx.WithEncoder(Console)
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

func (lc *LogContext) WithLevel(enable bool, option LevelOption) *LogContext {
	lc.levelF.enable = enable
	if enable {
		if len(option.LevelKey) == 0 {
			option.LevelKey = "level"
		}
		lc.levelF.option = option
	}
	return lc
}

func (lc *LogContext) WithTime(enable bool, option TimeOption) *LogContext {
	lc.timeF.enable = enable
	if enable {
		if len(option.TimeKey) == 0 {
			option.TimeKey = "time"
		}
		if option.Formatter == nil {
			option.Formatter = func(t time.Time) any { return t.Format(time.RFC3339) }
		}
		lc.timeF.option = option
	}
	return lc
}

func (lc *LogContext) WithCaller(enable bool, option CallerOption) *LogContext {
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
	default:
		panic("not support other log encoder")
	}
	return lc
}

func (lc *LogContext) BuildConsoleLogger(level LevelType) Logger {
	if lc.enc != nil {
		lc.enc.Init()
	}
	lc.WithMsgKey(lc.msgKey)
	return &LoggerX{
		logCtx:   lc,
		logLevel: level,
		pool: sync.Pool{
			New: func() any { return NewBuffer(make([]byte, 0, 1024)) },
		},
	}
}

func (lc *LogContext) BuildFileLogger(level LevelType, writer io.Writer) Logger {
	// For file logger, need to disable color attributes
	lc = lc.WithColorfulset(false, TextColorAttri{}).WithMsgKey(lc.msgKey).WithWriter(AddSync(writer))
	if lc.enc != nil {
		lc.enc.Init()
	}
	return &LoggerX{
		logCtx:   lc,
		logLevel: level,
		pool: sync.Pool{
			New: func() any { return NewBuffer(make([]byte, 0, 1024)) },
		},
	}
}
