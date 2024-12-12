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
	logger := logx.NewLogContext().WithWriter(os.Stdout).BuildConsoleLogger(logx.LevelTrace)
	logger.Info("Hello logx!")
}
```

### Advanced Configuration

```go
func main() {
	logCtx = logx.NewLogContext().
		WithColorfulset(true, logx.TextColorAttri{}). // Enable colored output
		WithLevel(true, logx.LevelOption{}).          // Show log level
		WithCaller(true, logx.CallerOption{}).        // Show caller information
		WithWriter(logx.AddSync(color.Output)).       // Set output writer
		WithEncoder(logx.Console).                    // Use console encoder
		WithTime(true, logx.TimeOption{Formatter: func(t time.Time) any { return t.Format(time.DateTime) }}) 	// Customize time format

	logger := logCtx.BuildConsoleLogger(logx.LevelTrace)

	// Different log levels with structured fields
	logger.Trace("this is a trace message", logx.String("key", "value"), logx.Int("key", 2222))
	logger.Debug("this is a debug message")
	logger.Info("this is an info message")
	logger.Warn("this is a warning message")
	logger.Error("this is an error message")
}
```

## License

[MIT](LICENSE)