package main

import (
	"errors"
	"time"

	"github.com/fatih/color"
	"github.com/josexy/logx"
)

func main() {
	logCtx := logx.NewLogContext().
		WithColor(true).
		WithLevel(true, true).
		WithCaller(true, true, true, true).
		WithWriter(color.Output).
		WithEncoder(logx.Simple).
		WithTime(true, func(t time.Time) any { return t.Format(time.DateTime) })

	loggerSimple := logCtx.BuildConsoleLogger(logx.LevelTrace)
	loggerSimple.Trace("this is a trace message")
	loggerSimple.Debug("this is a debug message")
	loggerSimple.Info("this is an info message")
	loggerSimple.Warn("this is a warning message")
	loggerSimple.Error("this is an error message")

	logCtx = logCtx.Copy().
		WithPrefix(logx.Pair{"scope", "test"}, logx.Pair{"version", "v1"}, logx.Pair{"list", []any{100, "a", "b"}}).
		WithEncoder(logx.Json).
		WithEscapeQuote(true).
		WithCaller(true, true, true, true).
		WithTime(true, func(t time.Time) any { return t.Unix() })

	// file, err := os.Create("test.log")
	// if err != nil {
	// 	panic(err)
	// }
	// defer file.Close()
	// loggerJson := logCtx.BuildFileLogger(logx.LevelInfo, file)
	loggerJson := logCtx.BuildConsoleLogger(logx.LevelInfo)
	loggerJson.Trace("this is a trace message")
	loggerJson.Debug("this is a debug message")
	loggerJson.Info("this is an info message")
	loggerJson.Warn("this is a warning message")
	loggerJson.Error("this is an error message")
	loggerJson.Info("this is a info message",
		logx.String("string", "string"),
		logx.Bool("bool", false),
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
		logx.SortedMap("sortedmap1"),
		logx.SortedMap("sortedmap2",
			logx.Pair{Key: "key1", Value: "value1"},
			logx.Pair{"key2", 100},
			logx.Pair{"key3", map[string]int{"a": 10, "b": 20}},
			logx.Pair{"key4", logx.M{"list": []any{10, "20", true, 30.10, nil}}}),
		logx.Map("map", logx.M{
			"name":     "tony",
			"age":      34,
			"time":     time.Now(),
			"float32":  123.123,
			"duration": time.Duration(time.Hour + 30*time.Minute + 40*time.Second),
			"err":      nil,
			"slice2":   []any{100, "hello", time.Now(), false, 10.2223, nil},
			"map2": logx.M{
				"name2": "mike",
				"age2":  20,
			},
		}),
		logx.Slice("slice", []any{100, "hello", time.Now(), false, 10.2223, nil}),
		logx.Map("map2", nil),
		logx.Slice("slice2", nil),
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
		logx.Slice2("slice3", true, false, 112233, 1122.33, "hello world", time.Now(), nil),
		logx.Slice3("slice4", []string{"hello", "world", "golang"}),
		logx.Slice3("slice5", []int{10, 20, 30, 40}),
		logx.Slice3("slice6", []struct {
			string
			int
		}{{"hello", 10}, {"world", 20}}),
	)
	loggerJson.Infof("hello %s", "world")

	loggerSimple.Error("this is an error message")
	loggerSimple.Tracef("hello %s", "world")

	newLogCtx := logCtx.Copy().WithTime(false, nil).WithPrefix(logx.Pair{"ip", "11.22.33.44"})
	newLogCtx.BuildConsoleLogger(logx.LevelTrace).Debug("hello world")
	newLogCtx2 := newLogCtx.Copy().WithPrefix(logx.Pair{"host", "localhost"})
	newLogCtx2.BuildConsoleLogger(logx.LevelTrace).Debug("hello world")
}
