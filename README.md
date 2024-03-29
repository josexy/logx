## simple and colorful golang logger

```shell
go get -u github.com/josexy/logx
```

## usage

```go
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
			logx.Pair{"key3", logx.M{"list": []any{10, "20", true, 30.10, nil}}}),
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
	)

	loggerSimple.Error("this is an error message")
}

```

output

![](./screenshots/example.jpg)
