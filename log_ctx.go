package logx

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"time"
)

type logContext struct {
	levelField
	timeField
	callerField
	encoder
	writer      io.Writer
	escapeQuote bool
}

func NewLogContext() *logContext {
	return &logContext{}
}

func (lc *logContext) Copy() *logContext {
	newLogCtx := new(logContext)
	*newLogCtx = *lc
	return newLogCtx
}

func (lc *logContext) WithColor(enable bool) *logContext {
	lc.levelField.color = enable
	lc.timeField.color = enable
	lc.callerField.color = enable
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
			return t.Format(time.DateTime)
		}
	}
	lc.timeField.enable = enable
	lc.timeField.format = format
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

func (lc *logContext) WithWriter(writer io.Writer) *logContext {
	lc.writer = writer
	return lc
}

func (lc *logContext) WithEncoder(encoder EncoderType) *logContext {
	switch encoder {
	case Simple:
		lc.encoder = &SimpleEncoder{logContext: lc}
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
		lockLevel: level,
		buf:       bytes.NewBuffer(make([]byte, 128)),
		logCtx:    lc,
	}
}

func (lc *logContext) BuildFileLogger(level LevelType, fsWriter *os.File) Logger {
	// For file logger, need to disable color attributes
	lc = lc.WithColor(false).WithWriter(bufio.NewWriter(fsWriter))

	if lc.encoder != nil {
		lc.encoder.Init()
	}
	return &LoggerX{
		lockLevel: level,
		buf:       bytes.NewBuffer(make([]byte, 128)),
		logCtx:    lc,
	}
}
