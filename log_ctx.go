package logx

import (
	"bytes"
	"io"
	"sync"
	"time"
)

type logContext struct {
	levelField
	timeField
	callerField
	encoder
	colorfulset
	writer      WriteSyncer
	preFields   []Field
	escapeQuote bool
}

func NewLogContext() *logContext {
	return &logContext{}
}

func (lc *logContext) Copy() *logContext {
	newLogCtx := new(logContext)
	*newLogCtx = *lc
	newLogCtx.preFields = make([]Field, 0, len(lc.preFields))
	newLogCtx.preFields = append(newLogCtx.preFields, lc.preFields...)
	switch lc.encoder.(type) {
	case *JsonEncoder:
		newLogCtx = newLogCtx.WithEncoder(Json)
	case *ConsoleEncoder:
		newLogCtx = newLogCtx.WithEncoder(Console)
	}
	return newLogCtx
}

func (lc *logContext) WithFields(fields ...Field) *logContext {
	lc.preFields = fields
	return lc
}

func (lc *logContext) WithColorfulset(enable bool, attr TextColorAttri) *logContext {
	lc.levelField.color = enable
	lc.timeField.color = enable
	lc.callerField.color = enable
	lc.colorfulset.enable = enable
	lc.colorfulset.TextColorAttri = attr
	return lc
}

func (lc *logContext) WithLevel(enable, lower bool) *logContext {
	lc.levelField.enable = enable
	lc.levelField.lower = lower
	return lc
}

func (lc *logContext) WithTime(enable bool, format func(time.Time) any) *logContext {
	if format == nil {
		format = func(t time.Time) any {
			return t.Format(time.RFC3339)
		}
	}
	lc.timeField.enable = enable
	lc.timeField.fn = format
	return lc
}

func (lc *logContext) WithCaller(enable, fileName, funcName, lineNumber bool) *logContext {
	lc.callerField.enable = enable
	lc.callerField.fileName = fileName
	lc.callerField.funcName = funcName
	lc.callerField.lineNum = lineNumber
	return lc
}

func (lc *logContext) WithEscapeQuote(enable bool) *logContext {
	lc.escapeQuote = enable
	return lc
}

func (lc *logContext) WithWriter(writer WriteSyncer) *logContext {
	lc.writer = writer
	return lc
}

func (lc *logContext) WithEncoder(encoder EncoderType) *logContext {
	switch encoder {
	case Console:
		lc.encoder = &ConsoleEncoder{logContext: lc}
	case Json:
		lc.encoder = &JsonEncoder{logContext: lc}
	default:
		panic("not support other log encoder")
	}
	return lc
}

func (lc *logContext) BuildConsoleLogger(level LevelType) Logger {
	if lc.encoder != nil {
		lc.encoder.Init()
	}
	return &LoggerX{
		logCtx:   lc,
		logLevel: level,
		pool: sync.Pool{
			New: func() any { return bytes.NewBuffer(make([]byte, 0, 1024)) },
		},
	}
}

func (lc *logContext) BuildFileLogger(level LevelType, writer io.Writer) Logger {
	// For file logger, need to disable color attributes
	lc = lc.WithColorfulset(false, TextColorAttri{}).WithWriter(AddSync(writer))
	if lc.encoder != nil {
		lc.encoder.Init()
	}
	return &LoggerX{
		logCtx:   lc,
		logLevel: level,
		pool: sync.Pool{
			New: func() any { return bytes.NewBuffer(make([]byte, 0, 1024)) },
		},
	}
}
