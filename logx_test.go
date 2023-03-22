package logx

import (
	"errors"
	"log"
	"testing"
	"time"
)

func TestSimpleLogger(t *testing.T) {
	logger := NewDevelopment(
		WithColor(true),
		WithLevel(true, false),
		WithCaller(true, true, true, true),
		WithSimpleEncoder(),
		WithTime(true, func(t time.Time) string { return t.Format(time.DateTime) }),
	)
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
	logger := NewDevelopment(
		WithColor(true),
		WithLevel(true, true),
		WithCaller(true, true, true, true),
		WithJsonEncoder(),
		WithTime(true, func(t time.Time) string { return t.Format(time.TimeOnly) }),
	)
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
	)
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

func BenchmarkSimpleLogger(b *testing.B) {
	// disable level/time/caller attributes
	logger := NewDevelopment(WithSimpleEncoder())
	logger.SetOutput(nullWriter{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("this is a message")
	}
}

func BenchmarkJsonLogger(b *testing.B) {
	// disable level/time/caller attributes
	logger := NewDevelopment(WithJsonEncoder())
	logger.SetOutput(nullWriter{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("this is a message")
	}
}
