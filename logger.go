package logx

// ErrorHandler receives errors produced while encoding, writing, or syncing a
// log entry. Implementations must be safe for concurrent use.
type ErrorHandler func(error)

type Logger interface {
	Enabled(level LevelType) bool
	Trace(msg string, fields ...Field)
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)
	Panic(msg string, fields ...Field)
	Tracef(format string, args ...any)
	Debugf(format string, args ...any)
	Infof(format string, args ...any)
	Warnf(format string, args ...any)
	Errorf(format string, args ...any)
	Panicf(format string, args ...any)
	Fatalf(format string, args ...any)
	// PanicWith logs err's message at panic level, then panics with err.
	// It does nothing when err is nil.
	PanicWith(err error)
	// ErrorWith logs err's message at error level. A nil err is logged as
	// "<nil>".
	ErrorWith(err error)
	// FatalWith logs err's message at fatal level, then calls os.Exit(1).
	// It does nothing, and does not exit, when err is nil.
	FatalWith(err error)
	With(fields ...Field) Logger
	Sync() error
}
