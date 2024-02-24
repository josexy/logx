package logx

type Logger interface {
	Trace(msg string, args ...arg)
	Debug(msg string, args ...arg)
	Info(msg string, args ...arg)
	Warn(msg string, args ...arg)
	Error(msg string, args ...arg)
	Fatal(msg string, args ...arg)
	Panic(msg string, args ...arg)
	Tracef(format string, args ...any)
	Debugf(format string, args ...any)
	Infof(format string, args ...any)
	Warnf(format string, args ...any)
	Errorf(format string, args ...any)
	Panicf(format string, args ...any)
	Fatalf(format string, args ...any)
	PanicBy(err error)
	ErrorBy(err error)
	FatalBy(err error)
}
