package sub

import (
	"github.com/fatih/color"
	"github.com/josexy/logx"
)

func TestLogger() {
	logCtx := logx.NewLogContext().
		WithFields(logx.String("module", "sub")).
		WithMsgKey("message").
		WithColorfulset(false, logx.TextColorAttri{}).
		WithLevel(true, logx.LevelOption{LowerKey: true}).
		WithCaller(true, logx.CallerOption{Formatter: logx.FullFile}).
		WithWriter(logx.AddSync(color.Output)).
		// WithWriter(logx.AddSync(nil)).
		WithEncoder(logx.Console).
		WithTime(true, logx.TimeOption{})

	logger := logCtx.BuildConsoleLogger(logx.LevelTrace)
	logger.Trace("hello world")

	loggerJson := logCtx.
		WithCaller(true, logx.CallerOption{
			FuncKey:   "function",
			Formatter: logx.FullFileFunc,
		}).
		WithEncoder(logx.Json).BuildConsoleLogger(logx.LevelInfo)
	loggerJson.Info("hello world")
}
