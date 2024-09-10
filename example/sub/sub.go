package sub

import (
	"runtime"

	"github.com/fatih/color"
	"github.com/josexy/logx"
)

func TestLogger() {
	logCtx := logx.NewLogContext().
		WithFields(logx.String("module", "sub")).
		WithMsgKey("message").
		WithColorfulset(true, logx.TextColorAttri{}).
		WithLevel(true, logx.LevelOption{LowerKey: true}).
		WithCaller(true, logx.CallerOption{Formatter: logx.FullFile}).
		WithWriter(logx.AddSync(color.Output)).
		// WithWriter(logx.AddSync(nil)).
		WithEncoder(logx.Console).
		WithTime(true, logx.TimeOption{})

	logger := logCtx.BuildConsoleLogger(logx.LevelTrace).With(logx.Int("id", 1000))
	logger.Trace("hello world")
	logger.With(logx.String("os", runtime.GOOS)).Trace("hello world")
	logger.With(logx.String("arch", runtime.GOARCH)).Trace("hello world")

	loggerJson := logCtx.
		WithCaller(true, logx.CallerOption{
			FuncKey:   "function",
			Formatter: logx.FullFileFunc,
		}).
		WithEncoder(logx.Json).BuildConsoleLogger(logx.LevelInfo)
	loggerJson.Info("hello world")

	loggerJson2 := logCtx.Copy().WithFields(logx.String("namespace", "default")).BuildConsoleLogger(logx.LevelTrace)
	loggerJson2.Info("hello world")

	loggerJson3 := logCtx.Copy().WithNewFields(logx.String("svc", "default")).BuildConsoleLogger(logx.LevelTrace)
	loggerJson3.Info("hello world")

	loggerJson2.Info("hello world")
}
