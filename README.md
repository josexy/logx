# logx

A simple, colorful and flexible logging library for Go.

## Features

- 🎨 Colorful console output
- 📊 Multiple log levels support (Trace, Debug, Info, Warn, Error)
- 🔍 Customizable caller information
- ⚙️ Flexible configuration options
- 🎯 Structured logging with key-value pairs
- ⏰ Customizable timestamp format

## Installation

```shell
go get github.com/josexy/logx
```

## Usage

### Basic Example

```go
func main() {
	// Create a simple console logger with default settings
	logger := logx.NewLogContext().WithLevel(logx.LevelTrace).WithEncoder(logx.Console).WithWriter(os.Stdout).Build()
	logger.Info("Hello logx!")
}
```

### Advanced Configuration

```go
func main() {
	logCtx := logx.NewLogContext().
		WithLevel(logx.LevelTrace).
		WithColorfulset(true, logx.TextColorAttri{}).                                                           // Enable colored output
		WithLevelKey(true, logx.LevelOption{}).                                                                 // Show log level
		WithCallerKey(true, logx.CallerOption{}).                                                               // Show caller information
		WithWriter(logx.AddSync(logx.Output)).                                                                 // Set output writer
		WithEncoder(logx.Console).                                                                              // Use console encoder
		WithTimeKey(true, logx.TimeOption{Formatter: func(t time.Time) any { return t.Format(time.DateTime) }}) // Customize time format

	logger := logCtx.Build()

	// Different log levels with structured fields
	logger.Trace("this is a trace message", logx.String("key", "value"), logx.Int("key", 2222))
	logger.Debug("this is a debug message")
	logger.Info("this is an info message")
	logger.Warn("this is a warning message")
	logger.Error("this is an error message")
	logger.With(logx.String("os", runtime.GOOS)).Debug("this is a debug message")

	logger2 := logCtx.Copy().WithFields(logx.String("arch", runtime.GOARCH)).WithEncoder(logx.Json).Build()
	logger2.Debug("this is a debug message")
}
```

output:

```
2025-03-21 16:12:21     TRACE   example/main.go:213     this is a trace message {"key":"value","key":2222}
2025-03-21 16:12:21     DEBUG   example/main.go:214     this is a debug message
2025-03-21 16:12:21     INFO    example/main.go:215     this is an info message
2025-03-21 16:12:21     WARN    example/main.go:216     this is a warning message
2025-03-21 16:12:21     ERROR   example/main.go:217     this is an error message
2025-03-21 16:12:21     DEBUG   example/main.go:218     this is a debug message {"os":"linux"}
{"level":"DEBUG","time":"2025-03-21 16:12:21","caller":{"file":"example/main.go:221"},"arch":"amd64","msg":"this is a debug message"}
```

## License

[MIT](LICENSE)