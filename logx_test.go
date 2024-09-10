package logx

import (
	"errors"
	"io"
	"log"
	"log/slog"
	"net/netip"
	"runtime"
	"testing"
	"time"

	"github.com/fatih/color"
)

func TestConsoleLogger(t *testing.T) {
	logCtx := NewLogContext().
		WithColorfulset(true, TextColorAttri{}).
		WithLevel(true, LevelOption{}).
		WithCaller(true, CallerOption{}).
		WithTime(true, TimeOption{}).
		WithWriter(AddSync(color.Output)).
		WithFields(String("arch", runtime.GOARCH), Bool("bool", true)).
		WithEncoder(Console)

	logger := logCtx.BuildConsoleLogger(LevelTrace)
	logger.Debug("")
	logger.Debug("this is a debug message")
	logger.Info("this is an info message")
	logger.Warn("this is a warning message")
	logger.Error("this is an error message")
	logger.Debugf("hello %s", "golang")
	logger.Infof("time :%v", time.Now())
	logger.Errorf("this is an error message: %v", io.EOF)
	logger.ErrorWith(errors.New("error"))

	logger.Debug("debug", String("os", runtime.GOOS), Time("ts", time.Now()))
	func() {
		logger.Info("info", String("os", runtime.GOOS), Time("ts", time.Now()))
	}()
	func() {
		defer func() {
			if err := recover(); err != nil {
				t.Logf("panic: %v", err)
			}
		}()
		logger.Panic("panic", Int("code", 255))
	}()

	newLogCtx := logCtx.Copy().WithTime(false, TimeOption{}).WithCaller(false, CallerOption{}).WithEncoder(Json)
	newLogger := newLogCtx.BuildConsoleLogger(LevelInfo)
	newLogger.Trace("trace", Time("ts", time.Now()))
	newLogger.Debugf("debug")
	newLogger.Info("info", Time("ts", time.Now()))
	newLogger.ErrorWith(io.EOF)

	newLogCtx.WithLevel(true, LevelOption{LevelKey: "ts", LowerKey: true}).WithFields()

	newLogger2 := newLogCtx.Copy().WithTime(true, TimeOption{}).
		WithLevel(true, LevelOption{}).WithEncoder(Console).BuildConsoleLogger(LevelInfo)
	newLogger.Info("info", Time("ts", time.Now()))
	newLogger2.Info("info", Time("ts", time.Now()))
}

func TestJsonLogger(t *testing.T) {
	logger := NewLogContext().
		WithColorfulset(true, TextColorAttri{}).
		WithFields(String("os", runtime.GOOS), String("arch", runtime.GOARCH)).
		WithLevel(true, LevelOption{}).
		WithCaller(true, CallerOption{Formatter: FullFileFunc}).
		WithWriter(AddSync(color.Output)).
		WithTime(true, TimeOption{Formatter: func(t time.Time) any { return t.Format(time.DateTime) }}).
		WithEncoder(Json).
		WithReflectValue(true).
		BuildConsoleLogger(LevelTrace)

	logger.Trace("this is a trace message")
	logger.Debug("this is a debug message")
	logger.Info("this is an info message")
	logger.Warn("this is a warning message")
	logger.Error("this is an error message")
	logger.Info("this is a info message",
		String("string", "string"),
		Bool("bool", false),
		Bool("bool2", true),
		Int8("int8", 10),
		Int16("int16", -20),
		Int32("int32", -30),
		Int64("int64", 40),
		Int("int", 50),
		UInt8("uint8", 60),
		UInt16("uint16", 70),
		UInt32("uint32", 80),
		UInt64("uint64", 90),
		UInt("uint", 100),
		Float32("float32", 1234.45),
		Float64("float64", 1234.4567),
		Time("ts", time.Now().Add(time.Duration(time.Hour))),
		Duration("duration", time.Duration(time.Hour+30*time.Minute+40*time.Second)),
		Error("err", errors.New("error message")),
		Error("err2", nil),
		String("a", `"message"`),
		String("b", `"message`),
		String("c", `message"`),
		String("d", `'message'`),
		String("e", `'message`),
		String("f", `message'`),
		String("g", `message."test".message`),
		String("h", `message.'test'.message`),
		String("i", `"hello" '\n' "\n"`),
		Any("any", "hello"),
		Any("any2", time.Now()),
		Any("any3", &struct{ k, v string }{k: "key", v: "value"}),
		Any("any4", logger),
		Any("any5", "\"xxx\""),
		Any("any6", `"""`),
		ArrayT("slice1",
			Object("z", Bool("b", false)),
			Object("a", String("b", "b")), Array("b", true, 20), ArrayT("c", "c", "c")),
		ArrayT("slice2", 110, 20, 300, 1000),
		Array("slice3", true, nil, false, 112233, 1122.33, "hello world", time.Now(), nil, io.EOF),
		Array("slice4", []string{"hello", "world", "golang"}),
		Array("slice5", "hello", []int{10, 20, 30, 40}, []bool{false, true, false}),
		Array("slice6", []struct {
			string
			int
		}{{"hello", 10}, {"world", 20}}),
		Object("object1", Bool("bool", false), String("string", "string"), Int("integer", 100)),
		Object("object2", Time("time", time.Now()), UInt32("uint32", 88999), Error("err", io.EOF),
			Error("err2", nil), Duration("duration", time.Millisecond*200), Float64("float64", 11.2222223)),
		Object("object3", Array("arr1", "str", 123, false, time.Now().Add(time.Hour)),
			Object("obj", Object("obj2", Array("arr", "xx", 12000), Object("obj3", Int("int", 2222))))),
	)
	logger.Info("info",
		Object("obj"),
		Array("arr"), ArrayT("arr2", io.EOF, nil, io.ErrShortBuffer))

	logger.Trace("trace", Array("arr",
		ArrayT("arr1", ArrayT("arr2", ArrayT("arr3", Array("arr5", 666)))),
		Array("arr1", 100, 200, ArrayT("x", 10, 20, 30), ArrayT("y", "ff", "gg"), Object("xx", String("xx", "ttt"))),
		Array("arr2",
			ArrayT("arr3", Array("arr4", false, 1.11), Array("arr5", 20, "hello", true)),
			200, "hello",
		)))
	logger.Info("info", ArrayT("ips",
		[]netip.Addr{netip.MustParseAddr("1.1.1.1"), netip.MustParseAddr("2.1.1.1")},
		[]netip.Addr{netip.MustParseAddr("2.2.2.2")},
	))
	logger.Info("info", Array("ips",
		netip.MustParseAddr("8.8.8.8"),
		[]netip.Addr{netip.MustParseAddr("1.1.1.1"), netip.MustParseAddr("2.1.1.1")},
		[]netip.Addr{netip.MustParseAddr("2.2.2.2")},
	))
	logger.Trace("trace", Object("obj", Object("obj2", Object("obj3"))))
	logger.Infof("hello %s", "world")
}

type nullWriter struct{}

func (w nullWriter) Write(b []byte) (n int, err error) { return }

func BenchmarkStdPrintLogger(b *testing.B) {
	logger := log.New(nullWriter{}, "", 0)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Println("this is a message")
	}
}

func BenchmarkStdWriterLogger(b *testing.B) {
	logger := log.New(nullWriter{}, "", 0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Writer().Write([]byte("this is a message\n"))
	}
}

func BenchmarkSlogTextLogger(b *testing.B) {
	logger := slog.New(slog.NewTextHandler(nullWriter{}, nil))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("this is a message")
	}
}

func BenchmarkSlogJsonLogger(b *testing.B) {
	logger := slog.New(slog.NewJSONHandler(nullWriter{}, nil))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("this is a message")
	}
}

func BenchmarkConsoleLogger(b *testing.B) {
	// disable level/time/caller attributes
	logger := NewLogContext().WithWriter(AddSync(nullWriter{})).WithEncoder(Console).BuildConsoleLogger(LevelTrace)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("this is a message")
	}
}

func BenchmarkJsonLogger(b *testing.B) {
	// disable level/time/caller attributes
	logger := NewLogContext().WithWriter(AddSync(nullWriter{})).WithEncoder(Json).BuildConsoleLogger(LevelTrace)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("this is a message")
	}
}

func BenchmarkConsoleLoggerWithEscapeQuote(b *testing.B) {
	// disable level/time/caller attributes
	logger := NewLogContext().WithWriter(AddSync(nullWriter{})).WithEscapeQuote(true).WithEncoder(Console).BuildConsoleLogger(LevelTrace)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info(`"this is a message"`)
	}
}

func BenchmarkJsonLoggerWithEscapeQuote(b *testing.B) {
	// disable level/time/caller attributes
	logger := NewLogContext().WithWriter(AddSync(nullWriter{})).WithEscapeQuote(true).WithEncoder(Json).BuildConsoleLogger(LevelTrace)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info(`"this is a message"`)
	}
}

func BenchmarkConsoleLoggerWithSimpleField(b *testing.B) {
	// disable level/time/caller attributes
	logger := NewLogContext().WithWriter(AddSync(nullWriter{})).WithEncoder(Console).BuildConsoleLogger(LevelTrace)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("this is a message", Int("key", 100))
	}
}

func BenchmarkJsonLoggerWithSimpleField(b *testing.B) {
	// disable level/time/caller attributes
	logger := NewLogContext().WithWriter(AddSync(nullWriter{})).WithEncoder(Json).BuildConsoleLogger(LevelTrace)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("this is a message", Int("key", 100))
	}
}

func BenchmarkConsoleLoggerWithField(b *testing.B) {
	// disable level/time/caller attributes
	logger := NewLogContext().WithWriter(AddSync(nullWriter{})).WithEncoder(Console).BuildConsoleLogger(LevelTrace)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("this is a message", String("key", "value"), Int("int", 10000))
	}
}

func BenchmarkJsonLoggerWithField(b *testing.B) {
	// disable level/time/caller attributes
	logger := NewLogContext().WithWriter(AddSync(nullWriter{})).WithEncoder(Json).BuildConsoleLogger(LevelTrace)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("this is a message", String("key", "value"), Int("int", 10000))
	}
}

func BenchmarkJsonLoggerWithReflectValueField(b *testing.B) {
	type Int int
	// disable level/time/caller attributes
	logger := NewLogContext().WithReflectValue(true).WithWriter(AddSync(nullWriter{})).WithEncoder(Json).BuildConsoleLogger(LevelTrace)
	arr := []Int{10, 20, 30}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("this is a message", ArrayT("key", arr...))
	}
}

func BenchmarkConsoleLoggerWithEscapeQuoteWithField(b *testing.B) {
	// disable level/time/caller attributes
	logger := NewLogContext().WithWriter(AddSync(nullWriter{})).WithEscapeQuote(true).WithEncoder(Json).BuildConsoleLogger(LevelTrace)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info(`"this is a message"`, String(`"key"`, `"value"`))
	}
}

func BenchmarkJsonLoggerWithEscapeQuoteWithField(b *testing.B) {
	// disable level/time/caller attributes
	logger := NewLogContext().WithWriter(AddSync(nullWriter{})).WithEscapeQuote(true).WithEncoder(Json).BuildConsoleLogger(LevelTrace)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info(`"this is a message"`, String(`"key"`, `"value"`))
	}
}
