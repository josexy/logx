package logx

type Logger interface {
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
	PanicWith(err error)
	ErrorWith(err error)
	FatalWith(err error)
	WithFields(fields ...Field) Logger
}
