# logx

A small structured logging library for Go with console, JSON, and text encoders.

## Features

- Console, JSON, and text output formats
- Standards-compliant JSON string and non-finite float encoding
- Trace, debug, info, warn, error, fatal, and panic levels
- Structured fields with primitive values, arrays, objects, errors, durations, and time values
- Optional time, level, and caller fields
- Runtime log level updates shared by parent and child loggers with `AtomicLevel`
- Custom field keys, timestamp layouts, and caller formatting
- Inherited logger fields with `With`, `WithFields`, and `WithNewFields`
- Optional ANSI color output for terminal logs
- Explicit `Sync`, automatic fatal/panic flushing, and configurable internal error handling

## Installation

```shell
go get github.com/josexy/logx
```

## Quick Start

```go
package main

import (
	"os"
	"time"

	"github.com/josexy/logx"
)

func main() {
	logger := logx.NewLogContext().
		WithLevel(logx.LevelTrace).
		WithLevelKey(true, logx.LevelOption{}).
		WithTimeKey(true, logx.TimeOption{Layout: time.DateTime}).
		WithWriter(logx.Lock(logx.AddSync(os.Stdout))).
		WithEncoder(logx.Console).
		Build()

	logger.Info("server started",
		logx.Int("port", 8080),
		logx.String("url", "http://localhost:8080"),
	)
}
```

Example console output:

```text
2026-06-18 14:11:48	INFO	server started	{"port":8080,"url":"http://localhost:8080"}
```

## Console Encoder

Use the console encoder for readable local development logs. It prints time, level, caller, message, and structured fields.

```go
logger := logx.NewLogContext().
	WithLevel(logx.LevelTrace).
	WithColorfulset(true, logx.TextColorAttri{
		NumberColor: logx.CyanAttr,
	}).
	WithLevelKey(true, logx.LevelOption{}).
	WithCallerKey(true, logx.CallerOption{Formatter: logx.ShortFile}).
	WithTimeKey(true, logx.TimeOption{}).
	WithWriter(logx.Lock(logx.AddSync(logx.Output))).
	WithEncoder(logx.Console).
	WithEscapeQuote(true).
	Build()

logger.Trace("trace message", logx.String("quoted", `"value"`), logx.Int("attempt", 1))
logger.Error("error message", logx.Error("err", io.EOF))
```

## Dynamic Level

Use `AtomicLevel` when you need to update the minimum log level at runtime without rebuilding existing loggers.

```go
level := logx.NewAtomicLevel(logx.LevelWarn)

logger := logx.NewLogContext().
	WithAtomicLevel(level).
	WithLevelKey(true, logx.LevelOption{}).
	WithTimeKey(true, logx.TimeOption{}).
	WithWriter(logx.Lock(logx.AddSync(logx.Output))).
	WithEncoder(logx.Console).
	Build()

logger.Info("not printed")
logger.Warn("printed before level update")

level.SetLevel(logx.LevelDebug)
logger.Debug("printed after level update")
logger.Info("also printed after level update")

child := logger.With(logx.String("component", "worker"))
child.Debug("child loggers observe the same dynamic level")
```

## JSON Encoder

Use the JSON encoder for structured logs consumed by log collectors.

```go
logger := logx.NewLogContext().
	WithLevel(logx.LevelTrace).
	WithLevelKey(true, logx.LevelOption{}).
	WithTimeKey(true, logx.TimeOption{Layout: time.RFC3339}).
	WithCallerKey(true, logx.CallerOption{Formatter: logx.ShortFileFunc}).
	WithWriter(logx.Lock(logx.AddSync(logx.Output))).
	WithFields(
		logx.String("service", "logx-example"),
		logx.String("env", "dev"),
	).
	WithEncoder(logx.Json).
	WithReflectValue(true).
	Build()

logger.Info("request handled",
	logx.Int("status", 200),
	logx.Duration("latency", 23*time.Millisecond),
	logx.Object("user",
		logx.Int("id", 10001),
		logx.String("name", "guest"),
	),
)
```

Example JSON output:

```json
{"level":"INFO","time":"2026-06-18T14:11:48+08:00","service":"logx-example","env":"dev","msg":"request handled","status":200,"latency":"23ms","user":{"id":10001,"name":"guest"}}
```

## Text Encoder

Use the text encoder for key-value logs similar to Go's `slog.TextHandler`.

```go
logger := logx.NewLogContext().
	WithLevel(logx.LevelTrace).
	WithTimeKey(true, logx.TimeOption{Layout: time.RFC3339Nano}).
	WithLevelKey(true, logx.LevelOption{}).
	WithCallerKey(true, logx.CallerOption{Formatter: logx.ShortFileFunc}).
	WithWriter(logx.Lock(logx.AddSync(logx.Output))).
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
```

Example text output:

```text
time=2026-06-18T14:11:48.7469115+08:00 level=INFO caller.file=example/main.go:97 caller.func=main.runTextExample msg="http server started" service=logx-example middleware=true handler=true port=8080 url=http://localhost:8080 runtime.os=windows runtime.arch=amd64
```

## Reliability and Shutdown

`Json` always escapes keys, messages, and string fields as valid JSON. ANSI
colors are automatically disabled for JSON output, even if color is enabled on
the context.

Call `Sync` during a normal shutdown when the writer buffers data. Fatal and
panic entries are synced automatically before their terminal action runs.

```go
logger := logx.NewLogContext().
	WithWriter(logx.Lock(logx.AddSync(os.Stderr))).
	WithErrorHandler(func(err error) {
		fmt.Fprintf(os.Stderr, "logger failure: %v\n", err)
	}).
	WithEncoder(logx.Json).
	Build()

defer func() {
	if err := logger.Sync(); err != nil {
		fmt.Fprintf(os.Stderr, "logger sync failure: %v\n", err)
	}
}()
```

Without a custom error handler, encoding, writing, short-write, and automatic
sync errors are reported to standard error. `Build` snapshots the logger
configuration; use the shared `AtomicLevel` for runtime level changes.

Compatibility note: the exported `Logger` interface now includes `Enabled` and
`Sync`. Custom logger implementations, adapters, and mocks must implement both
methods when upgrading.

## License

[MIT](LICENSE)
