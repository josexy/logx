package logx

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

type LoggerX struct {
	mu        sync.Mutex
	lockLevel LevelType
	buf       *bytes.Buffer
	logCtx    *logContext
}

func (l *LoggerX) print(level LevelType, msg string, args ...arg) {
	// discard the log
	if l.logCtx.writer == nil || l.logCtx.writer == io.Discard {
		return
	}
	if l.skipLevelLog(level) {
		return
	}
	l.output(level, msg, args...)
}

func (l *LoggerX) Trace(msg string, args ...arg) { l.print(LevelTrace, msg, args...) }

func (l *LoggerX) Debug(msg string, args ...arg) { l.print(LevelDebug, msg, args...) }

func (l *LoggerX) Info(msg string, args ...arg) { l.print(LevelInfo, msg, args...) }

func (l *LoggerX) Warn(msg string, args ...arg) { l.print(LevelWarn, msg, args...) }

func (l *LoggerX) Error(msg string, args ...arg) { l.print(LevelError, msg, args...) }

func (l *LoggerX) Fatal(msg string, args ...arg) {
	l.print(LevelFatal, msg, args...)
	os.Exit(1)
}

func (l *LoggerX) Panic(msg string, args ...arg) {
	l.print(LevelPanic, msg, args...)
	panic(msg)
}

func (l *LoggerX) Tracef(format string, args ...any) {
	l.print(LevelTrace, fmt.Sprintf(format, args...))
}

func (l *LoggerX) Debugf(format string, args ...any) {
	l.print(LevelDebug, fmt.Sprintf(format, args...))
}

func (l *LoggerX) Infof(format string, args ...any) { l.print(LevelInfo, fmt.Sprintf(format, args...)) }

func (l *LoggerX) Warnf(format string, args ...any) { l.print(LevelWarn, fmt.Sprintf(format, args...)) }

func (l *LoggerX) Errorf(format string, args ...any) {
	l.print(LevelError, fmt.Sprintf(format, args...))
}

func (l *LoggerX) Fatalf(format string, args ...any) {
	l.print(LevelFatal, fmt.Sprintf(format, args...))
	os.Exit(1)
}

func (l *LoggerX) Panicf(format string, args ...any) {
	value := fmt.Sprintf(format, args...)
	l.print(LevelPanic, value)
	panic(value)
}

func (l *LoggerX) ErrorBy(err error) {
	value := "<nil>"
	if err != nil {
		value = err.Error()
	}
	l.print(LevelError, value)
}

func (l *LoggerX) PanicBy(err error) {
	if err == nil {
		return
	}
	l.print(LevelPanic, err.Error())
	panic(err)
}

func (l *LoggerX) FatalBy(err error) {
	if err == nil {
		return
	}
	l.print(LevelFatal, err.Error())
	os.Exit(1)
}

func (l *LoggerX) skipLevelLog(expect LevelType) bool {
	return l.lockLevel > expect
}

func (l *LoggerX) output(level LevelType, msg string, args ...arg) {
	if l.logCtx.encoder == nil {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	if l.logCtx.timeField.enable {
		l.logCtx.timeField.now = time.Now()
	}
	if l.logCtx.levelField.enable {
		l.logCtx.levelField.typ = level
	}

	// reset the buffer
	l.buf.Reset()

	_ = l.logCtx.encoder.Encode(l.buf, msg, args...)

	if l.buf.Len() > 0 && l.buf.Bytes()[l.buf.Len()-1] != '\n' {
		l.buf.WriteByte('\n')
	}
	_, _ = l.logCtx.writer.Write(l.buf.Bytes())

	// for file writer, need to flush the buffer data to file
	if flusher, ok := l.logCtx.writer.(*bufio.Writer); ok {
		flusher.Flush()
	}
}
