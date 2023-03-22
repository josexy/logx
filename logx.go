package logx

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/fatih/color"
	"github.com/josexy/logx/internal"
)

type LoggerX struct {
	isDiscard int32
	lockLevel internal.LevelType
	mu        sync.Mutex
	out       io.Writer
	buf       *bytes.Buffer
	cfg       *internal.Config
}

// NewNop nothing to do
func NewNop() *LoggerX {
	return &LoggerX{
		isDiscard: 1,
		out:       io.Discard,
	}
}

func NewDevelopment(opts ...ConfigOption) *LoggerX {
	return newLogger(internal.LevelDebug, opts...)
}

func NewProduction(opts ...ConfigOption) *LoggerX {
	return newLogger(internal.LevelInfo, opts...)
}

func newLogger(lvl internal.LevelType, opts ...ConfigOption) *LoggerX {
	defaultCfg := new(internal.Config)
	for _, opt := range opts {
		opt.applyTo(defaultCfg)
	}
	defaultCfg.Encoder.Init()
	return &LoggerX{
		lockLevel: lvl,
		buf:       bytes.NewBuffer(make([]byte, 128)),
		cfg:       defaultCfg,
		out:       color.Output,
	}
}

func (l *LoggerX) SetOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.out = w
	isDiscard := int32(0)
	if w == io.Discard {
		isDiscard = 1
	}
	atomic.StoreInt32(&l.isDiscard, isDiscard)
}

func (l *LoggerX) Debug(msg string, args ...Arg) {
	if internal.LevelDebug < l.lockLevel {
		return
	}
	if atomic.LoadInt32(&l.isDiscard) != 0 {
		return
	}
	l.output(internal.LevelDebug, msg, args...)
}

func (l *LoggerX) Info(msg string, args ...Arg) {
	if internal.LevelInfo < l.lockLevel {
		return
	}
	if atomic.LoadInt32(&l.isDiscard) != 0 {
		return
	}
	l.output(internal.LevelInfo, msg, args...)
}

func (l *LoggerX) Warn(msg string, args ...Arg) {
	if internal.LevelWarn < l.lockLevel {
		return
	}
	if atomic.LoadInt32(&l.isDiscard) != 0 {
		return
	}
	l.output(internal.LevelWarn, msg, args...)
}

func (l *LoggerX) Error(msg string, args ...Arg) {
	if internal.LevelError < l.lockLevel {
		return
	}
	if atomic.LoadInt32(&l.isDiscard) != 0 {
		return
	}
	l.output(internal.LevelError, msg, args...)
}

func (l *LoggerX) Panic(msg string, args ...Arg) {
	if internal.LevelPanic < l.lockLevel {
		return
	}
	if atomic.LoadInt32(&l.isDiscard) != 0 {
		return
	}
	l.output(internal.LevelPanic, msg, args...)
	panic(msg)
}

func (l *LoggerX) Fatal(msg string, args ...Arg) {
	if atomic.LoadInt32(&l.isDiscard) != 0 {
		return
	}
	if internal.LevelFatal < l.lockLevel {
		return
	}
	l.output(internal.LevelFatal, msg, args...)
	os.Exit(1)
}

func (l *LoggerX) Debugf(format string, args ...any) {
	if internal.LevelDebug < l.lockLevel {
		return
	}
	if atomic.LoadInt32(&l.isDiscard) != 0 {
		return
	}
	l.output(internal.LevelDebug, fmt.Sprintf(format, args...))
}

func (l *LoggerX) Infof(format string, args ...any) {
	if internal.LevelInfo < l.lockLevel {
		return
	}
	if atomic.LoadInt32(&l.isDiscard) != 0 {
		return
	}
	l.output(internal.LevelInfo, fmt.Sprintf(format, args...))
}

func (l *LoggerX) Warnf(format string, args ...any) {
	if internal.LevelWarn < l.lockLevel {
		return
	}
	if atomic.LoadInt32(&l.isDiscard) != 0 {
		return
	}
	l.output(internal.LevelWarn, fmt.Sprintf(format, args...))
}

func (l *LoggerX) Errorf(format string, args ...any) {
	if internal.LevelError < l.lockLevel {
		return
	}
	if atomic.LoadInt32(&l.isDiscard) != 0 {
		return
	}
	l.output(internal.LevelError, fmt.Sprintf(format, args...))
}

func (l *LoggerX) Panicf(format string, args ...any) {
	if internal.LevelPanic < l.lockLevel {
		return
	}
	if atomic.LoadInt32(&l.isDiscard) != 0 {
		return
	}
	value := fmt.Sprintf(format, args...)
	l.output(internal.LevelPanic, value)
	panic(value)
}

func (l *LoggerX) Fatalf(format string, args ...any) {
	if atomic.LoadInt32(&l.isDiscard) != 0 {
		return
	}
	if internal.LevelFatal < l.lockLevel {
		return
	}
	l.output(internal.LevelFatal, fmt.Sprintf(format, args...))
	os.Exit(1)
}

func (l *LoggerX) ErrorBy(err error) {
	if err == nil {
		return
	}
	if internal.LevelError < l.lockLevel {
		return
	}
	if atomic.LoadInt32(&l.isDiscard) != 0 {
		return
	}
	l.output(internal.LevelError, err.Error())
}

func (l *LoggerX) PanicBy(err error) {
	if err == nil {
		return
	}
	if internal.LevelPanic < l.lockLevel {
		return
	}
	if atomic.LoadInt32(&l.isDiscard) != 0 {
		return
	}
	l.output(internal.LevelPanic, err.Error())
	panic(err)
}

func (l *LoggerX) FatalBy(err error) {
	if err == nil {
		return
	}
	if atomic.LoadInt32(&l.isDiscard) != 0 {
		return
	}
	if internal.LevelFatal < l.lockLevel {
		return
	}
	l.output(internal.LevelFatal, err.Error())
	os.Exit(1)
}

func (l *LoggerX) output(lvl internal.LevelType, msg string, args ...Arg) {
	if l.cfg.Encoder == nil {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	if l.cfg.TimeItem.Enable {
		l.cfg.TimeItem.Now = time.Now()
	}
	if l.cfg.LevelItem.Enable {
		l.cfg.LevelItem.Typ = lvl
	}

	// reset the buffer
	l.buf.Reset()

	_ = l.cfg.Encoder.Encode(l.buf, msg, args...)

	if l.buf.Len() > 0 && l.buf.Bytes()[l.buf.Len()-1] != '\n' {
		l.buf.WriteByte('\n')
	}
	_, _ = l.out.Write(l.buf.Bytes())
}
