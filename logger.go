package logx

type Logger interface {
	Debug(msg string, args ...Arg)
	Info(msg string, args ...Arg)
	Warn(msg string, args ...Arg)
	Error(msg string, args ...Arg)
	Fatal(msg string, args ...Arg)
	Panic(msg string, args ...Arg)
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
