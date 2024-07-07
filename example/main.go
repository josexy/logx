package main

import (
	"errors"
	"io"
	"runtime"
	"time"

	"github.com/fatih/color"
	"github.com/josexy/logx"
)

func main() {
	logCtx := logx.NewLogContext().
		WithColorfulset(true, logx.TextColorAttri{}).
		WithLevel(true, true).
		WithCaller(true, true, true, true).
		WithWriter(logx.AddSync(color.Output)).
		// WithWriter(logx.AddSync(nil)).
		WithEncoder(logx.Console).
		WithTime(true, func(t time.Time) any { return t.Format(time.DateTime) })

	loggerSimple := logCtx.BuildConsoleLogger(logx.LevelTrace)
	loggerSimple.Trace("this is a trace message", logx.String("key", "value"), logx.Int("key", 2222))
	loggerSimple.Debug("this is a debug message")
	loggerSimple.Info("this is an info message")
	loggerSimple.Warn("this is a warning message")
	loggerSimple.Error("this is an error message")

	logCtx = logCtx.Copy().
		WithEncoder(logx.Json).
		WithEscapeQuote(true)

	// file, err := os.Create("test.log")
	// if err != nil {
	// 	panic(err)
	// }
	// defer file.Close()
	// loggerJson := logCtx.WithEncoder(logx.Console).BuildFileLogger(logx.LevelInfo, io.MultiWriter(file, os.Stdout))

	loggerJson := logCtx.
		WithEncoder(logx.Console).
		WithFields(logx.String("os", runtime.GOOS), logx.String("arch", runtime.GOARCH)).
		BuildConsoleLogger(logx.LevelTrace)
	loggerJson.Trace("this is a trace message")
	loggerJson.Debug("this is a debug message")
	loggerJson.Info("this is an info message")
	loggerJson.Warn("this is a warning message")
	loggerJson.Error("this is an error message")
	loggerJson.Info("this is a info message",
		logx.String("string", "string"),
		logx.Bool("bool", false),
		logx.Bool("bool2", true),
		logx.Int8("int8", 10),
		logx.Int16("int16", -20),
		logx.Int32("int32", -30),
		logx.Int64("int64", 40),
		logx.Int("int", 50),
		logx.UInt8("uint8", 60),
		logx.UInt16("uint16", 70),
		logx.UInt32("uint32", 80),
		logx.UInt64("uint64", 90),
		logx.UInt("uint", 100),
		logx.Float32("float32", 1234.45),
		logx.Float64("float64", 1234.4567),
		logx.Time("ts", time.Now().Add(time.Duration(time.Hour))),
		logx.Duration("duration", time.Duration(time.Hour+30*time.Minute+40*time.Second)),
		logx.Error("err", errors.New("error message")),
		logx.Error("err2", nil),
		logx.String("a", `"message"`),
		logx.String("b", `"message`),
		logx.String("c", `message"`),
		logx.String("d", `'message'`),
		logx.String("e", `'message`),
		logx.String("f", `message'`),
		logx.String("g", `message."test".message`),
		logx.String("h", `message.'test'.message`),
		logx.String("i", `"hello" '\n' "\n"`),
		logx.Any("any", "hello"),
		logx.Any("any2", time.Now()),
		logx.Any("any3", &struct{ k, v string }{k: "key", v: "value"}),
		logx.Any("any4", loggerJson),
		logx.Any("any5", "\"xxx\""),
		logx.Any("any6", `"""`),
		logx.ArrayT("slice1",
			logx.Object("z", logx.Bool("b", false)),
			logx.Object("a", logx.String("b", "b")), logx.Array("b", true, 20), logx.ArrayT("c", "c", "c")),
		logx.ArrayT("slice2", 110, 20, 300, 1000),
		logx.Array("slice3", true, nil, false, 112233, 1122.33, "hello world", time.Now(), nil, io.EOF),
		logx.Array("slice4", []string{"hello", "world", "golang"}),
		logx.Array("slice5", "hello", []int{10, 20, 30, 40}, []bool{false, true, false}),
		logx.Array("slice6", []struct {
			string
			int
		}{{"hello", 10}, {"world", 20}}),
		logx.Object("object1", logx.Bool("bool", false), logx.String("string", "string"), logx.Int("integer", 100)),
		logx.Object("object2", logx.Time("time", time.Now()), logx.UInt32("uint32", 88999), logx.Error("err", io.EOF),
			logx.Error("err2", nil), logx.Duration("duration", time.Millisecond*200), logx.Float64("float64", 11.2222223)),
		logx.Object("object3", logx.Array("arr1", "str", 123, false, time.Now().Add(time.Hour)),
			logx.Object("obj", logx.Object("obj2", logx.Array("arr", "xx", 12000), logx.Object("obj3", logx.Int("int", 2222))))),
	)
	loggerJson.Info("info",
		logx.Object("obj"),
		logx.Array("arr"), logx.ArrayT("arr2", io.EOF, nil, io.ErrShortBuffer))

	loggerJson.Trace("trace", logx.Array("arr",
		logx.ArrayT("arr1", logx.ArrayT("arr2", logx.ArrayT("arr3", logx.Array("arr5", 666)))),
		logx.Array("arr1", 100, 200, logx.ArrayT("x", 10, 20, 30), logx.ArrayT("y", "ff", "gg"), logx.Object("xx", logx.String("xx", "ttt"))),
		logx.Array("arr2",
			logx.ArrayT("arr3", logx.Array("arr4", false, 1.11), logx.Array("arr5", 20, "hello", true)),
			200, "hello",
		)))
	loggerJson.Trace("trace", logx.Object("obj", logx.Object("obj2", logx.Object("obj3"))))
	loggerJson.Infof("hello %s", "world")

	loggerSimple.Error("this is an error message")
	loggerSimple.Tracef("hello %s", "world")
}
