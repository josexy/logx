package logx

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

type WriteSyncer interface {
	io.Writer
	Sync() error
}

type writerWrapper struct{ io.Writer }

func (w writerWrapper) Sync() error { return nil }

func AddSync(w io.Writer) WriteSyncer {
	if w == nil {
		return nil
	}
	switch w := w.(type) {
	case WriteSyncer:
		return w
	default:
		return writerWrapper{w}
	}
}

type lockedWriteSyncer struct {
	sync.Mutex
	ws WriteSyncer
}

func (s *lockedWriteSyncer) Write(bs []byte) (int, error) {
	s.Lock()
	n, err := s.ws.Write(bs)
	s.Unlock()
	return n, err
}

func (s *lockedWriteSyncer) Sync() error {
	s.Lock()
	err := s.ws.Sync()
	s.Unlock()
	return err
}

// Lock wraps a WriteSyncer in a mutex to make it safe for concurrent use. In
// particular, *os.Files must be locked before use.
// See zap log
func Lock(ws WriteSyncer) WriteSyncer {
	if _, ok := ws.(*lockedWriteSyncer); ok {
		// no need to layer on another lock
		return ws
	}
	return &lockedWriteSyncer{ws: ws}
}

type LoggerX struct {
	logCtx *LogContext
}

func (l *LoggerX) print(level LevelType, msg string, fields []Field) {
	// discard the log
	if l.logCtx.writer == nil || l.logCtx.writer == io.Discard {
		return
	}
	if l.skipLevelLog(level) {
		return
	}
	l.output(level, msg, fields)
}

func (l *LoggerX) Trace(msg string, fields ...Field) { l.print(LevelTrace, msg, fields) }

func (l *LoggerX) Debug(msg string, fields ...Field) { l.print(LevelDebug, msg, fields) }

func (l *LoggerX) Info(msg string, fields ...Field) { l.print(LevelInfo, msg, fields) }

func (l *LoggerX) Warn(msg string, fields ...Field) { l.print(LevelWarn, msg, fields) }

func (l *LoggerX) Error(msg string, fields ...Field) { l.print(LevelError, msg, fields) }

func (l *LoggerX) Fatal(msg string, fields ...Field) {
	l.print(LevelFatal, msg, fields)
	os.Exit(1)
}

func (l *LoggerX) Panic(msg string, fields ...Field) {
	l.print(LevelPanic, msg, fields)
	panic(msg)
}

func (l *LoggerX) Tracef(format string, args ...any) {
	l.print(LevelTrace, fmt.Sprintf(format, args...), nil)
}

func (l *LoggerX) Debugf(format string, args ...any) {
	l.print(LevelDebug, fmt.Sprintf(format, args...), nil)
}

func (l *LoggerX) Infof(format string, args ...any) {
	l.print(LevelInfo, fmt.Sprintf(format, args...), nil)
}

func (l *LoggerX) Warnf(format string, args ...any) {
	l.print(LevelWarn, fmt.Sprintf(format, args...), nil)
}

func (l *LoggerX) Errorf(format string, args ...any) {
	l.print(LevelError, fmt.Sprintf(format, args...), nil)
}

func (l *LoggerX) Fatalf(format string, args ...any) {
	l.print(LevelFatal, fmt.Sprintf(format, args...), nil)
	os.Exit(1)
}

func (l *LoggerX) Panicf(format string, args ...any) {
	value := fmt.Sprintf(format, args...)
	l.print(LevelPanic, value, nil)
	panic(value)
}

func (l *LoggerX) ErrorWith(err error) {
	value := "<nil>"
	if err != nil {
		value = err.Error()
	}
	l.print(LevelError, value, nil)
}

func (l *LoggerX) PanicWith(err error) {
	if err == nil {
		return
	}
	l.print(LevelPanic, err.Error(), nil)
	panic(err)
}

func (l *LoggerX) FatalWith(err error) {
	if err == nil {
		return
	}
	l.print(LevelFatal, err.Error(), nil)
	os.Exit(1)
}

func (l *LoggerX) skipLevelLog(expect LevelType) bool {
	return l.logCtx.levelT > expect
}

func (l *LoggerX) clone() *LoggerX {
	clone := &LoggerX{logCtx: l.logCtx.Copy()}
	if clone.logCtx.enc != nil {
		clone.logCtx.enc.Init()
	}
	return clone
}

func (l *LoggerX) With(fields ...Field) Logger {
	nl := l.clone()
	nl.logCtx = nl.logCtx.WithFields(fields...)
	return nl
}

func (l *LoggerX) output(level LevelType, msg string, fields []Field) {
	if l.logCtx.enc == nil {
		return
	}

	ent := entry{
		level:   level,
		message: msg,
		time:    time.Now(),
	}

	var buf *Buffer
	var err error
	if buf, err = l.logCtx.enc.Encode(ent, fields); err != nil {
		return
	}
	if buf.Len() > 0 && buf.Bytes()[buf.Len()-1] != '\n' {
		buf.AppendByte('\n')
	}
	l.logCtx.writer.Write(buf.Bytes())
	bufPool.Put(buf)

	if l.logCtx.levelT > LevelError {
		_ = l.logCtx.writer.Sync()
	}
}
