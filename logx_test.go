package logx

import (
	"errors"
	"log"
	"log/slog"
	"testing"
	"time"

	"github.com/fatih/color"
)

func TestSimpleLogger(t *testing.T) {
	logger := NewLogContext().WithColor(true).
		WithLevel(true, false).
		WithCaller(true, true, true, true).
		WithWriter(color.Output).
		WithTime(true, func(t time.Time) any { return t.Format(time.DateTime) }).
		WithEncoder(Simple).BuildConsoleLogger(LevelTrace)

	logger.Debug("")
	logger.Debug("this is a debug message")
	logger.Info("this is an info message")
	logger.Warn("this is a warning message")
	logger.Error("this is an error message")
	logger.Debugf("hello %s", "golang")
	logger.Infof("time :%v", time.Now())
	logger.ErrorBy(errors.New("error"))
	// logger.PanicBy(errors.New("panic"))
	// logger.FatalBy(errors.New("fatal"))
	// logger.Panic("this is a panic message")
	// logger.Fatal("this is a fatal message")
}

func TestJsonLogger(t *testing.T) {
	logger := NewLogContext().WithColor(true).
		WithLevel(true, false).
		WithCaller(true, true, true, true).
		WithWriter(color.Output).
		WithTime(true, func(t time.Time) any { return t.Format(time.DateTime) }).
		WithEncoder(Json).BuildConsoleLogger(LevelTrace)

	logger.Debug("this is a debug message",
		String("string", "string"),
		Bool("bool", false),
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
		SortedMap("sortedmap1"),
		SortedMap("sortedmap2",
			Pair{Key: "key1", Value: "value1"},
			Pair{"key2", 100},
			Pair{"key3", M{"list": []any{10, "20", true, 30.10, nil}}}),
		Map("map", M{
			"name":     "tony",
			"age":      34,
			"time":     time.Now(),
			"float32":  123.123,
			"duration": time.Duration(time.Hour + 30*time.Minute + 40*time.Second),
			"err":      nil,
			"slice2":   []any{100, "hello", time.Now(), false, 10.2223, nil},
			"map2": M{
				"name2": "mike",
				"age2":  20,
			},
		}),
		Slice("slice", []any{100, "hello", time.Now(), false, 10.2223, nil}),
		Map("map2", nil),
		Slice("slice2", nil),
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
	)

	logger.Info(`"hello" \n "\n"`)
	logger.Info("")
	logger.Info("this is an info message")
	logger.Warn("this is a warning message")
	logger.Error("this is an error message")
	logger.Debugf("hello %s", "golang")
	logger.Infof("time :%v", time.Now())
	logger.ErrorBy(errors.New("error"))
	// logger.Panic("this is a panic message")
	// logger.Fatal("this is a fatal message")
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

func BenchmarkSimpleLogger(b *testing.B) {
	// disable level/time/caller attributes
	logger := NewLogContext().WithWriter(nullWriter{}).WithEncoder(Simple).BuildConsoleLogger(LevelTrace)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("this is a message")
	}
}

func BenchmarkJsonLogger(b *testing.B) {
	// disable level/time/caller attributes
	logger := NewLogContext().WithWriter(nullWriter{}).WithEncoder(Json).BuildConsoleLogger(LevelTrace)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("this is a message")
	}
}

func BenchmarkJsonLoggerWithEscapeQuote(b *testing.B) {
	// disable level/time/caller attributes
	logger := NewLogContext().WithWriter(nullWriter{}).WithEscapeQuote(true).WithEncoder(Json).BuildConsoleLogger(LevelTrace)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info(`"this is a message"`)
	}
}
