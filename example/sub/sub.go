package sub

import (
	"encoding/json"
	"runtime"
	"time"

	"github.com/fatih/color"
	"github.com/josexy/logx"
)

func loggerInfo(lc *logx.LogContext, msg string) {
	logger := lc.WithCallerKey(true, logx.CallerOption{
		Formatter:  logx.ShortFileFunc,
		CallerSkip: 1,
	}).WithEncoder(logx.Json).Build()

	logger.Info(msg)
	logger.Debug(msg)

	data, _ := json.Marshal(logx.CallerOption{CallerKey: "caller"})
	var r1 map[string]any
	json.Unmarshal(data, &r1)
	mapRes := map[string]any{
		"id":           10001,
		"name":         "guest",
		"created_time": time.Now(),
		"info": map[string]any{
			"k1": "127.0.0.1",
			"k2": logx.Array("array", "10", 20, false, []string{"a", "b"},
				[]map[string]any{{"k3": 10, "ok": true}, r1, {"k3": 20}, {"k3": 30}}),
			"k3": []int{10, 20, 30},
			"k4": r1,
			"k5": []map[string]any{r1, r1}, // use map[string]any type formatter
			"k6": map[any]any{"xx": true},  // fallback to final any type formatter
			"kvs": map[string]string{
				"version": "v1",
				"env":     "dev",
			},
		},
		"kvs": map[string]string{
			"version": "v1",
			"env":     "dev",
		},
	}
	logger.Warn(msg, logx.Any("any", mapRes))
	logger.Warn(msg, logx.Array("list", mapRes, []map[string]any{mapRes, mapRes}))
}

func TestLogger() {
	logCtx := logx.NewLogContext().
		WithFields(logx.String("module", "sub")).
		WithMsgKey("message").
		WithColorfulset(true, logx.TextColorAttri{}).
		WithLevel(logx.LevelTrace).
		WithLevelKey(true, logx.LevelOption{LowerKey: true}).
		WithCallerKey(true, logx.CallerOption{Formatter: logx.FullFile}).
		WithWriter(logx.Lock(logx.AddSync(color.Output))).
		// WithWriter(logx.AddSync(nil)).
		WithEncoder(logx.Console).
		WithTimeKey(true, logx.TimeOption{})

	logger := logCtx.Build().With(logx.Int("id", 1000))
	logger.Trace("hello world", logx.String("key", "value"))
	logger.With(logx.String("os", runtime.GOOS)).Trace("hello world", logx.String("key", "value"))
	logger.With(logx.String("arch", runtime.GOARCH)).Trace("hello world")

	loggerJson := logCtx.
		WithLevel(logx.LevelInfo).
		WithCallerKey(true, logx.CallerOption{
			FuncKey:   "function",
			Formatter: logx.FullFileFunc,
		}).
		WithEncoder(logx.Json).Build()
	loggerJson.Info("hello world")

	loggerJson2 := logCtx.Copy().WithFields(logx.String("namespace", "default")).Build()
	loggerJson2.Info("hello world", logx.String("key", "value"))

	loggerJson3 := logCtx.Copy().WithNewFields(logx.String("svc", "default")).Build()
	loggerJson3.Info("hello world", logx.String("key", "value"))

	loggerJson2.Info("hello world")

	loggerInfo(logCtx.Copy(), "hello world")
}
