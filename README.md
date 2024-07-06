## simple and colorful golang logger

```shell
go get github.com/josexy/logx
```

## usage

```go
func main() {
	logCtx := logx.NewLogContext().
		WithColorfulset(true, logx.TextColorAttri{}).
		WithLevel(true, true).
		WithCaller(true, true, true, true).
		WithWriter(logx.AddSync(color.Output)).
		WithEncoder(logx.Console).
		WithTime(true, func(t time.Time) any { return t.Format(time.DateTime) })

	loggerSimple := logCtx.BuildConsoleLogger(logx.LevelTrace)
	loggerSimple.Trace("this is a trace message", logx.String("key", "value"), logx.Int("key", 2222))
	loggerSimple.Debug("this is a debug message")
	loggerSimple.Info("this is an info message")
	loggerSimple.Warn("this is a warning message")
	loggerSimple.Error("this is an error message")
}
```
