package main

import (
	"errors"
	"fmt"
	"io"
	"net/netip"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/josexy/logx"
)

func main() {
	runConsoleExample()
	runDynamicLevelExample()
	runJSONExample()
	runTextExample()
	runDerivedLoggerExamples()
	runConcurrentFileExample()
}

func runConsoleExample() {
	section("console encoder")

	logger := logx.NewLogContext().
		WithLevel(logx.LevelTrace).
		WithColorfulset(true, logx.TextColorAttri{
			NumberColor: logx.CyanAttr,
		}).
		WithLevelKey(true, logx.LevelOption{}).
		WithCallerKey(true, logx.CallerOption{Formatter: logx.ShortFile}).
		WithTimeKey(true, logx.TimeOption{}).
		WithWriter(stdout()).
		WithEncoder(logx.Console).
		WithEscapeQuote(true).
		Build()

	logger.Trace("trace message", logx.String("quoted", `"value"`), logx.Int("attempt", 1))
	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warn("warn message")
	logger.Error("error message", logx.Error("err", errors.New("example error")))

	requestLogger := logger.With(
		logx.String("request_id", "req-001"),
		logx.String("component", "console-example"),
	)
	requestLogger.Info("request handled",
		logx.Int("status", 200),
		logx.Duration("latency", 23*time.Millisecond),
	)
}

func runDynamicLevelExample() {
	section("dynamic level")

	level := logx.NewAtomicLevel(logx.LevelWarn)
	logger := logx.NewLogContext().
		WithAtomicLevel(level).
		WithLevelKey(true, logx.LevelOption{}).
		WithTimeKey(true, logx.TimeOption{}).
		WithWriter(stdout()).
		WithEncoder(logx.Console).
		Build()

	logger.Info("not printed before level update")
	logger.Warn("warn before level update")

	level.SetLevel(logx.LevelDebug)
	logger.Debug("debug after level update")
	logger.Info("info after level update")
}

func runJSONExample() {
	section("json encoder")

	now := time.Now()
	logger := logx.NewLogContext().
		WithLevel(logx.LevelTrace).
		WithLevelKey(true, logx.LevelOption{}).
		WithTimeKey(true, logx.TimeOption{Layout: time.RFC3339}).
		WithCallerKey(true, logx.CallerOption{Formatter: logx.ShortFileFunc}).
		WithWriter(stdout()).
		WithFields(
			logx.String("service", "logx-example"),
			logx.String("os", runtime.GOOS),
			logx.String("arch", runtime.GOARCH),
		).
		WithEncoder(logx.Json).
		WithEscapeQuote(true).
		WithReflectValue(true).
		Build()

	logger.Info("primitive fields", primitiveFields(now)...)
	logger.Info("nested fields", nestedFields(now)...)
	logger.Info("network fields", networkFields()...)
	logger.Infof("formatted %s", "message")
	logger.ErrorWith(io.EOF)
}

func runTextExample() {
	section("text encoder")

	logger := logx.NewLogContext().
		WithLevel(logx.LevelTrace).
		WithColorfulset(true, logx.TextColorAttri{}).
		WithTimeKey(true, logx.TimeOption{Layout: time.RFC3339Nano}).
		WithCallerKey(true, logx.CallerOption{Formatter: logx.ShortFileFunc}).
		WithLevelKey(true, logx.LevelOption{}).
		WithWriter(stdout()).
		WithFields(logx.String("service", "logx-example")).
		WithEncoder(logx.Text).
		Build()

	logger.Info("http server started",
		logx.Bool("middleware", true),
		logx.Bool("handler", true),
		logx.Int("port", 8080),
		logx.String("url", "http://localhost:8080"),
		logx.Object("runtime",
			logx.String("os", runtime.GOOS),
			logx.String("arch", runtime.GOARCH),
		),
	)

	logger.Warn("special text values",
		logx.String("empty", ""),
		logx.String("space", "hello world"),
		logx.String("quote", `"hello"`),
		logx.Duration("timeout", time.Second),
		logx.Error("err", io.EOF),
	)
}

func runConcurrentFileExample() {
	section("concurrent file writer")

	file, err := os.CreateTemp("", "logx-example-*.log")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	logger := logx.NewLogContext().
		WithLevel(logx.LevelInfo).
		WithLevelKey(true, logx.LevelOption{}).
		WithTimeKey(true, logx.TimeOption{}).
		WithCallerKey(true, logx.CallerOption{}).
		WithWriter(logx.Lock(logx.AddSync(file))).
		WithEncoder(logx.Console).
		Build()

	var wg sync.WaitGroup
	for workerID := 1; workerID <= 4; workerID++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for jobID := 1; jobID <= 3; jobID++ {
				logger.Info("job processed",
					logx.Int("worker_id", workerID),
					logx.Int("job_id", jobID),
				)
			}
		}(workerID)
	}
	wg.Wait()

	fmt.Printf("wrote concurrent log example to %s\n", file.Name())
}

func runDerivedLoggerExamples() {
	section("derived logger examples")

	ctx := logx.NewLogContext().
		WithLevel(logx.LevelTrace).
		WithColorfulset(true, logx.TextColorAttri{}).
		WithLevelKey(true, logx.LevelOption{LowerKey: true}).
		WithTimeKey(true, logx.TimeOption{}).
		WithCallerKey(true, logx.CallerOption{Formatter: logx.ShortFile}).
		WithWriter(stdout()).
		WithFields(logx.String("module", "main")).
		WithEncoder(logx.Console)

	runDerivedConsoleExample(ctx.Copy())
	runDerivedJSONExample(ctx.Copy())
	runWrappedLoggerExample(ctx.Copy(), "message from wrapper")
}

func runDerivedConsoleExample(ctx *logx.LogContext) {
	logger := ctx.Build().With(logx.Int("id", 1000))

	logger.Trace("console logger with inherited fields", logx.String("key", "value"))
	logger.With(logx.String("os", runtime.GOOS)).Debug("debug from child logger")
	logger.With(logx.String("arch", runtime.GOARCH)).Info("info from child logger")
}

func runDerivedJSONExample(ctx *logx.LogContext) {
	logger := ctx.
		WithLevel(logx.LevelInfo).
		WithMsgKey("message").
		WithCallerKey(true, logx.CallerOption{
			FuncKey:   "function",
			Formatter: logx.ShortFileFunc,
		}).
		WithEncoder(logx.Json).
		Build()

	logger.Info("json logger with custom message key", logx.String("key", "value"))

	loggerWithReplacedFields := ctx.Copy().
		WithNewFields(logx.String("namespace", "default")).
		WithEncoder(logx.Json).
		Build()
	loggerWithReplacedFields.Info("json logger with replaced fields", logx.String("key", "value"))
}

func runWrappedLoggerExample(ctx *logx.LogContext, msg string) {
	logger := ctx.
		WithCallerKey(true, logx.CallerOption{
			Formatter:  logx.ShortFileFunc,
			CallerSkip: 1,
		}).
		WithEncoder(logx.Json).
		Build()

	logger.Info(msg)
	logger.Warn(msg, logx.Any("payload", samplePayload()))
	logger.Warn(msg, logx.Array("payloads", samplePayload(), []map[string]any{
		samplePayload(),
		samplePayload(),
	}))
}

func primitiveFields(now time.Time) []logx.Field {
	return []logx.Field{
		logx.String("string", "value"),
		logx.Bool("bool", true),
		logx.Int("int", 50),
		logx.UInt("uint", 100),
		logx.Float64("float64", 1234.4567),
		logx.Time("time", now),
		logx.Duration("duration", time.Hour+30*time.Minute+40*time.Second),
		logx.Error("error", errors.New("example error")),
		logx.Error("nil_error", nil),
	}
}

func nestedFields(now time.Time) []logx.Field {
	return []logx.Field{
		logx.Array("events",
			"created",
			123,
			false,
			now.Add(time.Hour),
			io.EOF,
		),
		logx.ArrayT("numbers", 10, 20, 30),
		logx.Object("user",
			logx.Int("id", 10001),
			logx.String("name", "guest"),
			logx.Object("profile",
				logx.Bool("active", true),
				logx.ArrayT("roles", "reader", "writer"),
			),
		),
		logx.Any("metadata", map[string]any{
			"env":     "dev",
			"version": "v1",
			"features": []string{
				"console",
				"json",
				"text",
			},
		}),
	}
}

func networkFields() []logx.Field {
	addresses := []netip.Addr{
		netip.MustParseAddr("1.1.1.1"),
		netip.MustParseAddr("8.8.8.8"),
	}

	return []logx.Field{
		logx.ArrayT("addresses", addresses...),
		logx.Any("address_list", addresses),
		logx.Array("mixed",
			netip.MustParseAddr("127.0.0.1"),
			addresses,
			[]string{"http", "grpc"},
		),
	}
}

func samplePayload() map[string]any {
	return map[string]any{
		"id":           10001,
		"name":         "guest",
		"created_time": time.Date(2026, 6, 18, 9, 18, 35, 0, time.Local),
		"info": map[string]any{
			"host": "127.0.0.1",
			"ports": []int{
				8080,
				9090,
			},
			"features": []map[string]any{
				{"name": "console", "enabled": true},
				{"name": "json", "enabled": true},
				{"name": "text", "enabled": true},
			},
			"labels": map[string]string{
				"version": "v1",
				"env":     "dev",
			},
			"headers": map[string][]string{
				"accept": {"application/json"},
				"trace":  {"trace-001", "trace-002"},
			},
		},
	}
}

func stdout() logx.WriteSyncer {
	return logx.Lock(logx.AddSync(logx.Output))
}

func section(title string) {
	fmt.Printf("\n--- %s ---\n", title)
}
