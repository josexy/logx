package logx

type Logger interface {
	Trace(msg string, fields ...Field)
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)
	Panic(msg string, fields ...Field)
	Tracef(format string, fields ...any)
	Debugf(format string, fields ...any)
	Infof(format string, fields ...any)
	Warnf(format string, fields ...any)
	Errorf(format string, fields ...any)
	Panicf(format string, fields ...any)
	Fatalf(format string, fields ...any)
	PanicWith(err error)
	ErrorWith(err error)
	FatalWith(err error)
	WithFields(fields ...Field) Logger
}
