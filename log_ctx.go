package logx

import (
	"bytes"
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

func (lc *LogContext) WithLevel(enable, lower bool) *LogContext {
	lc.levelF.enable = enable
	lc.levelF.lower = lower
	return lc
}

func (lc *LogContext) WithTime(enable bool, format func(time.Time) any) *LogContext {
	if format == nil {
		format = func(t time.Time) any {
			return t.Format(time.RFC3339)
		}
	}
	lc.timeF.enable = enable
	lc.timeF.fn = format
	return lc
}

func (lc *LogContext) WithCaller(enable, fileName, funcName, lineNumber bool) *LogContext {
	lc.callerF.enable = enable
	lc.callerF.fileName = fileName
	lc.callerF.funcName = funcName
	lc.callerF.lineNum = lineNumber
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
	return &LoggerX{
		logCtx:   lc,
		logLevel: level,
		pool: sync.Pool{
			New: func() any { return bytes.NewBuffer(make([]byte, 0, 1024)) },
		},
	}
}

func (lc *LogContext) BuildFileLogger(level LevelType, writer io.Writer) Logger {
	// For file logger, need to disable color attributes
	lc = lc.WithColorfulset(false, TextColorAttri{}).WithWriter(AddSync(writer))
	if lc.enc != nil {
		lc.enc.Init()
	}
	return &LoggerX{
		logCtx:   lc,
		logLevel: level,
		pool: sync.Pool{
			New: func() any { return bytes.NewBuffer(make([]byte, 0, 1024)) },
		},
	}
}
