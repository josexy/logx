package logx

import (
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"sync"
	"time"
)

type WriteSyncer interface {
	io.Writer
	Sync() error
}

type writerWrapper struct{ io.Writer }

func (w writerWrapper) Sync() error { return nil }

type discardWriteSyncer struct{}

func (discardWriteSyncer) Write(p []byte) (int, error) { return len(p), nil }

func (discardWriteSyncer) Sync() error { return nil }

var discardWriterType = reflect.TypeOf(io.Discard)

func isDiscardWriteSyncer(ws WriteSyncer) bool {
	switch w := ws.(type) {
	case discardWriteSyncer:
		return true
	case writerWrapper:
		return reflect.TypeOf(w.Writer) == discardWriterType
	case *lockedWriteSyncer:
		return isDiscardWriteSyncer(w.ws)
	default:
		return false
	}
}

func AddSync(w io.Writer) WriteSyncer {
	if w == nil {
		return nil
	}
	if reflect.TypeOf(w) == discardWriterType {
		return discardWriteSyncer{}
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
	if ws == nil {
		return nil
	}
	if isDiscardWriteSyncer(ws) {
		return ws
	}
	if _, ok := ws.(*lockedWriteSyncer); ok {
		// no need to layer on another lock
		return ws
	}
	return &lockedWriteSyncer{ws: ws}
}

type LoggerX struct {
	logCtx *LogContext
}

func defaultErrorHandler(err error) {
	_, _ = fmt.Fprintf(os.Stderr, "logx: %v\n", err)
}

func (l *LoggerX) reportError(err error) {
	if err == nil {
		return
	}
	handler := l.logCtx.errorHandler
	if handler == nil {
		handler = defaultErrorHandler
	}
	handler(err)
}

// Enabled reports whether a log entry at level would be written.
func (l *LoggerX) Enabled(level LevelType) bool {
	if l == nil || l.logCtx == nil || level > LevelPanic {
		return false
	}
	if l.logCtx.writer == nil || isDiscardWriteSyncer(l.logCtx.writer) || l.logCtx.enc == nil {
		return false
	}
	return l.logCtx.AtomicLevel().Level() <= level
}

func (l *LoggerX) print(level LevelType, msg string, fields []Field) {
	if !l.Enabled(level) {
		return
	}
	l.reportError(l.output(level, msg, fields))
}

func (l *LoggerX) finishTerminal(logErr error) {
	var syncErr error
	if err := l.Sync(); err != nil {
		syncErr = fmt.Errorf("sync log entry: %w", err)
	}
	l.reportError(errors.Join(logErr, syncErr))
}

func (l *LoggerX) terminal(level LevelType, msg string, fields []Field) {
	var logErr error
	if l.Enabled(level) {
		logErr = l.output(level, msg, fields)
	}
	l.finishTerminal(logErr)
}

func (l *LoggerX) terminalf(level LevelType, format string, args ...any) {
	var logErr error
	if l.Enabled(level) {
		logErr = l.output(level, fmt.Sprintf(format, args...), nil)
	}
	l.finishTerminal(logErr)
}

func (l *LoggerX) Trace(msg string, fields ...Field) { l.print(LevelTrace, msg, fields) }

func (l *LoggerX) Debug(msg string, fields ...Field) { l.print(LevelDebug, msg, fields) }

func (l *LoggerX) Info(msg string, fields ...Field) { l.print(LevelInfo, msg, fields) }

func (l *LoggerX) Warn(msg string, fields ...Field) { l.print(LevelWarn, msg, fields) }

func (l *LoggerX) Error(msg string, fields ...Field) { l.print(LevelError, msg, fields) }

func (l *LoggerX) Fatal(msg string, fields ...Field) {
	l.terminal(LevelFatal, msg, fields)
	os.Exit(1)
}

func (l *LoggerX) Panic(msg string, fields ...Field) {
	l.terminal(LevelPanic, msg, fields)
	panic(msg)
}

func (l *LoggerX) Tracef(format string, args ...any) {
	l.printf(LevelTrace, format, args...)
}

func (l *LoggerX) Debugf(format string, args ...any) {
	l.printf(LevelDebug, format, args...)
}

func (l *LoggerX) Infof(format string, args ...any) {
	l.printf(LevelInfo, format, args...)
}

func (l *LoggerX) Warnf(format string, args ...any) {
	l.printf(LevelWarn, format, args...)
}

func (l *LoggerX) Errorf(format string, args ...any) {
	l.printf(LevelError, format, args...)
}

func (l *LoggerX) printf(level LevelType, format string, args ...any) {
	if !l.Enabled(level) {
		return
	}
	if err := l.output(level, fmt.Sprintf(format, args...), nil); err != nil {
		l.reportError(err)
	}
}

func (l *LoggerX) Fatalf(format string, args ...any) {
	l.terminalf(LevelFatal, format, args...)
	os.Exit(1)
}

func (l *LoggerX) Panicf(format string, args ...any) {
	value := fmt.Sprintf(format, args...)
	l.terminal(LevelPanic, value, nil)
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
	l.terminal(LevelPanic, err.Error(), nil)
	panic(err)
}

func (l *LoggerX) FatalWith(err error) {
	if err == nil {
		return
	}
	l.terminal(LevelFatal, err.Error(), nil)
	os.Exit(1)
}

func (l *LoggerX) clone() *LoggerX {
	clone := &LoggerX{logCtx: l.logCtx.copySharedLevel()}
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

// Sync flushes buffered log entries in the configured writer.
func (l *LoggerX) Sync() error {
	if l == nil || l.logCtx == nil || l.logCtx.writer == nil {
		return nil
	}
	return l.logCtx.writer.Sync()
}

func (l *LoggerX) output(level LevelType, msg string, fields []Field) error {
	if l.logCtx.enc == nil {
		return nil
	}

	ent := entry{
		level:   level,
		message: msg,
		time:    time.Now(),
	}

	var buf *Buffer
	var err error
	if buf, err = l.logCtx.enc.Encode(ent, fields); err != nil {
		return fmt.Errorf("encode log entry: %w", err)
	}
	if buf == nil {
		return fmt.Errorf("encode log entry: encoder returned a nil buffer")
	}
	defer bufPool.Put(buf)
	if buf.Len() > 0 && buf.Bytes()[buf.Len()-1] != '\n' {
		buf.AppendByte('\n')
	}
	n, err := l.logCtx.writer.Write(buf.Bytes())
	if err != nil {
		return fmt.Errorf("write log entry: %w", err)
	}
	if n != buf.Len() {
		return fmt.Errorf("write log entry: %w", io.ErrShortWrite)
	}

	return nil
}
