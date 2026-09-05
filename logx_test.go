package logx

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"math"
	"net/netip"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestStringQuote(t *testing.T) {
	buf := NewBuffer(make([]byte, 0, 8))
	buf.AppendString("test-")
	appendQuoteString(buf, "hello world")
	buf.AppendInt(100)
	if buf.String() != `test-hello world100` {
		t.Errorf("expect %s, got %s", "test-hello world", buf.String())
	}

	buf.Reset()
	buf.AppendString("test2-")
	appendQuoteString(buf, `"hello+"hahaha"+world"`)
	buf.AppendInt(100)
	if buf.String() != `test2-\"hello+\"hahaha\"+world\"100` {
		t.Errorf("expect %s, got %s", `test2-\"hello+\"hahaha\"+world\"100`, buf.String())
	}
}

func TestConsoleEncoderOutput(t *testing.T) {
	ent := entry{
		level:   LevelInfo,
		time:    time.Date(2026, 6, 18, 9, 18, 35, 0, time.UTC),
		message: "console message",
	}

	tests := []struct {
		name   string
		ctx    *LogContext
		fields []Field
		want   string
	}{
		{
			name: "prefix and fields",
			ctx: NewLogContext().
				WithFields(String("app", "logx")).
				WithTimeKey(true, TimeOption{}).
				WithLevelKey(true, LevelOption{}),
			fields: []Field{String("event", "start"), Bool("ok", true)},
			want:   "2026-06-18 09:18:35\tINFO\tconsole message\t{\"app\":\"logx\",\"event\":\"start\",\"ok\":true}",
		},
		{
			name: "message only",
			ctx:  NewLogContext(),
			want: "console message",
		},
		{
			name: "prefix only",
			ctx:  NewLogContext().WithFields(Int("count", 2)),
			want: "console message\t{\"count\":2}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := encodeForTest(t, tt.ctx, Console, ent, tt.fields); got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestConsoleEncoderColorSettings(t *testing.T) {
	oldNoColor := NoColor
	NoColor = false
	defer func() { NoColor = oldNoColor }()

	got := encodeForTest(t, NewLogContext().
		WithColorfulset(true, TextColorAttri{}).
		WithLevelKey(true, LevelOption{}).
		WithTimeKey(true, TimeOption{}), Console, entry{
		level:   LevelWarn,
		time:    time.Date(2026, 6, 18, 9, 18, 35, 0, time.UTC),
		message: "colored",
	}, []Field{Bool("ok", true)})

	for _, want := range []string{"\x1b[32m", "\x1b[33m", "\x1b[34m", "\x1b[0m"} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected colored console output to contain %q, got %q", want, got)
		}
	}
}

func TestJsonEncoderFieldTypes(t *testing.T) {
	loc := time.FixedZone("TST", 2*60*60)
	entTime := time.Date(2026, 6, 18, 9, 18, 35, 0, loc)
	fieldTime := time.Date(2026, 6, 18, 10, 20, 30, 0, loc)
	got := decodeJSONLog(t, encodeForTest(t, NewLogContext().
		WithFields(String("service", "logx")).
		WithLevelKey(true, LevelOption{}).
		WithTimeKey(true, TimeOption{Layout: time.RFC3339}).
		WithEscapeQuote(true).
		WithReflectValue(true), Json, entry{
		level:   LevelInfo,
		time:    entTime,
		message: "json message",
	}, []Field{
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
		Float32("float32", 1.5),
		Float64("float64", 2.25),
		Time("ts", fieldTime),
		Duration("duration", 90*time.Second),
		Error("err", errors.New("error message")),
		Error("err2", nil),
		String("quoted", `"message"`),
		Any("any", "hello"),
		Any("anyTime", fieldTime),
		Array("slice", true, nil, false, 112233, 1122.33, "hello world", fieldTime, nil, io.EOF),
		Any("ips", []netip.Addr{netip.MustParseAddr("1.1.1.1"), netip.MustParseAddr("2.1.1.1")}),
		Object("object", Bool("bool", false), String("string", "string"), Int("integer", 100)),
		Object("nested", Array("arr", "str", 123, false, fieldTime),
			Object("obj", Object("child", Array("arr", "xx", 12000), Object("leaf", Int("int", 2222))))),
	}))

	if got["level"] != "INFO" || got["time"] != entTime.Format(time.RFC3339) || got["msg"] != "json message" {
		t.Fatalf("unexpected prompt fields: %#v", got)
	}
	if got["service"] != "logx" || got["string"] != "string" || got["quoted"] != `"message"` {
		t.Fatalf("unexpected string fields: %#v", got)
	}
	if got["bool"] != false || got["bool2"] != true {
		t.Fatalf("unexpected bool fields: %#v", got)
	}
	for key, want := range map[string]float64{
		"int8": 10, "int16": -20, "int32": -30, "int64": 40, "int": 50,
		"uint8": 60, "uint16": 70, "uint32": 80, "uint64": 90, "uint": 100,
		"float32": 1.5, "float64": 2.25,
	} {
		if got[key] != want {
			t.Fatalf("%s = %v, want %v", key, got[key], want)
		}
	}
	if got["ts"] != fieldTime.Format(time.RFC3339) || got["anyTime"] != fieldTime.Format(time.RFC3339) {
		t.Fatalf("unexpected time fields: %#v", got)
	}
	if got["duration"] != "1m30s" || got["err"] != "error message" || got["err2"] != nil {
		t.Fatalf("unexpected duration/error fields: %#v", got)
	}
	slice := got["slice"].([]any)
	if slice[0] != true || slice[1] != nil || slice[2] != false || slice[5] != "hello world" || slice[6] != fieldTime.Format(time.RFC3339) || slice[8] != "EOF" {
		t.Fatalf("unexpected slice: %#v", slice)
	}
	ips := got["ips"].([]any)
	if ips[0] != "1.1.1.1" || ips[1] != "2.1.1.1" {
		t.Fatalf("unexpected ips: %#v", ips)
	}
	object := got["object"].(map[string]any)
	if object["bool"] != false || object["string"] != "string" || object["integer"] != float64(100) {
		t.Fatalf("unexpected object: %#v", object)
	}
	nested := got["nested"].(map[string]any)
	child := nested["obj"].(map[string]any)["child"].(map[string]any)
	if child["leaf"].(map[string]any)["int"] != float64(2222) {
		t.Fatalf("unexpected nested object: %#v", nested)
	}
}

func TestJsonEncoderAlwaysProducesValidJSON(t *testing.T) {
	oldNoColor := NoColor
	NoColor = false
	defer func() { NoColor = oldNoColor }()

	invalidUTF8 := string([]byte{'a', 0xff, 'b'})
	ctx := NewLogContext().
		WithColorfulset(true, TextColorAttri{}).
		WithTimeKey(true, TimeOption{Layout: `2006"01`})
	ent := entry{
		level:   LevelInfo,
		time:    time.Date(2026, time.July, 12, 1, 2, 3, 0, time.UTC),
		message: "quoted \"message\"\\with\nnewline\a",
	}

	encoded := encodeForTest(t, ctx, Json, ent, []Field{
		String("quoted\"key", "value\vwith\tcontrols"),
		String("invalidUTF8", invalidUTF8),
		Float64("nan", math.NaN()),
		Float64("positiveInfinity", math.Inf(1)),
		Float64("negativeInfinity", math.Inf(-1)),
		Object("object", String("valid", "field"), Field{}),
	})

	if !json.Valid([]byte(encoded)) {
		t.Fatalf("Json encoder produced invalid JSON: %q", encoded)
	}
	if strings.Contains(encoded, "\x1b[") {
		t.Fatalf("Json encoder emitted ANSI color codes: %q", encoded)
	}

	got := decodeJSONLog(t, encoded)
	if got["msg"] != ent.message {
		t.Fatalf("msg = %q, want %q", got["msg"], ent.message)
	}
	if got["time"] != `2026"07` {
		t.Fatalf("time = %q, want %q", got["time"], `2026"07`)
	}
	if got[`quoted"key`] != "value\vwith\tcontrols" {
		t.Fatalf("quoted field = %q", got[`quoted"key`])
	}
	if got["invalidUTF8"] != "a\ufffdb" {
		t.Fatalf("invalidUTF8 = %q, want %q", got["invalidUTF8"], "a\ufffdb")
	}
	if got["nan"] != "NaN" || got["positiveInfinity"] != "+Inf" || got["negativeInfinity"] != "-Inf" {
		t.Fatalf("unexpected non-finite float encoding: %#v", got)
	}
	if object := got["object"].(map[string]any); object["valid"] != "field" || len(object) != 1 {
		t.Fatalf("unexpected object: %#v", object)
	}
}

func TestTextEncoderExampleOutput(t *testing.T) {
	loc := time.FixedZone("CST", 8*60*60)
	ent := entry{
		level:   LevelInfo,
		time:    time.Date(2026, 6, 18, 9, 18, 35, 764000000, loc),
		message: "AssetServer Info:",
	}
	logCtx := NewLogContext().
		WithTimeKey(true, TimeOption{Layout: "2006-01-02T15:04:05.000Z07:00"}).
		WithLevelKey(true, LevelOption{})

	got := encodeForTest(t, logCtx, Text, ent, []Field{
		Bool("middleware", true),
		Bool("handler", true),
		String("devServerURL", "http://localhost:9245"),
	})
	want := `time=2026-06-18T09:18:35.764+08:00 level=INFO msg="AssetServer Info:" middleware=true handler=true devServerURL=http://localhost:9245`
	if got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
}

func TestTextEncoderDefaultTimeLayoutIsQuoted(t *testing.T) {
	ent := entry{
		level:   LevelInfo,
		time:    time.Date(2026, 6, 18, 9, 18, 35, 0, time.UTC),
		message: "msg",
	}
	logCtx := NewLogContext().
		WithTimeKey(true, TimeOption{})

	want := `time="2026-06-18 09:18:35" msg=msg`
	if got := encodeForTest(t, logCtx, Text, ent, nil); got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
}

func TestTimeFieldUsesDefaultLayoutWithoutPromptTime(t *testing.T) {
	value := time.Date(2026, time.July, 12, 14, 30, 45, 0, time.UTC)
	ent := entry{level: LevelInfo, time: time.Unix(0, 0), message: "message"}

	jsonLog := decodeJSONLog(t, encodeForTest(t, NewLogContext(), Json, ent, []Field{Time("eventTime", value)}))
	if got, want := jsonLog["eventTime"], value.Format(time.DateTime); got != want {
		t.Fatalf("JSON eventTime = %q, want %q", got, want)
	}

	if got, want := encodeTextForTest(t, NewLogContext(), ent, []Field{Time("eventTime", value)}), `msg=message eventTime="2026-07-12 14:30:45"`; got != want {
		t.Fatalf("Text eventTime output = %q, want %q", got, want)
	}
}

func TestTextEncoderQuotingRules(t *testing.T) {
	ent := entry{level: LevelInfo, time: time.Unix(0, 0), message: "ok"}
	fields := []Field{
		String("empty", ""),
		String("space", "a b"),
		String("equal", "a=b"),
		String("quote", `"ab"`),
		String("tab", "a\tb"),
		String("url", "http://localhost:9245"),
		String("unicode", "µåπ"),
	}

	want := `msg=ok empty="" space="a b" equal="a=b" quote="\"ab\"" tab="a\tb" url=http://localhost:9245 unicode=µåπ`
	if got := encodeForTest(t, NewLogContext(), Text, ent, fields); got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
}

func TestTextEncoderObjectFlatteningAndFieldOrder(t *testing.T) {
	ent := entry{level: LevelInfo, time: time.Unix(0, 0), message: "message"}
	logCtx := NewLogContext().
		WithFields(String("prefix", "first"))
	got := encodeForTest(t, logCtx, Text, ent, []Field{
		Object("obj", String("a", "b"), Object("child", Bool("ok", true))),
		String("event", "second"),
	})

	want := `msg=message prefix=first obj.a=b obj.child.ok=true event=second`
	if got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
}

func TestTextEncoderKeepsPlainTextWithoutColor(t *testing.T) {
	ent := entry{level: LevelInfo, time: time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC), message: "message"}
	logCtx := NewLogContext().
		WithLevelKey(true, LevelOption{}).
		WithTimeKey(true, TimeOption{})

	got := encodeForTest(t, logCtx, Text, ent, nil)
	if strings.Contains(got, "\x1b[") {
		t.Fatalf("text encoder should not emit ANSI colors when color is disabled, got %q", got)
	}
	want := `time="1970-01-01 00:00:00" level=INFO msg=message`
	if got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
}

func TestTextEncoderColorSettings(t *testing.T) {
	oldNoColor := NoColor
	NoColor = false
	defer func() { NoColor = oldNoColor }()

	ent := entry{level: LevelInfo, time: time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC), message: "message"}
	logCtx := NewLogContext().
		WithColorfulset(true, TextColorAttri{}).
		WithLevelKey(true, LevelOption{}).
		WithTimeKey(true, TimeOption{})
	got := encodeForTest(t, logCtx, Text, ent, []Field{
		Bool("ok", true),
		Int("count", 10),
		Float64("cost", 1.25),
	})

	for _, want := range []string{
		"\x1b[34m",
		"\x1b[32m",
		"\x1b[33m",
		"\x1b[31m",
		"\x1b[36m",
		"\x1b[0m",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected colored text output to contain %q, got %q", want, got)
		}
	}
}

func TestTextEncoderCaller(t *testing.T) {
	buffer := bytes.NewBuffer(nil)
	logger := NewLogContext().
		WithCallerKey(true, CallerOption{}).
		WithWriter(AddSync(buffer)).
		WithEncoder(Text).
		Build()

	_, _, line, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	logger.Info("caller")

	want := fmt.Sprintf("caller.file=logx/logx_test.go:%d msg=caller\n", line+4)
	if got := buffer.String(); got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
}

func TestCallerFieldsForJSONAndConsole(t *testing.T) {
	buffer := bytes.NewBuffer(nil)
	jsonLogger := NewLogContext().
		WithCallerKey(true, CallerOption{Formatter: ShortFileFunc}).
		WithWriter(AddSync(buffer)).
		WithEncoder(Json).
		Build()
	jsonLogger.Info("caller")

	got := decodeJSONLog(t, strings.TrimSpace(buffer.String()))
	caller := got["caller"].(map[string]any)
	if !strings.Contains(caller["file"].(string), "logx/logx_test.go:") {
		t.Fatalf("unexpected JSON caller file: %#v", caller)
	}
	if !strings.Contains(caller["func"].(string), "logx.TestCallerFieldsForJSONAndConsole") {
		t.Fatalf("unexpected JSON caller func: %#v", caller)
	}

	buffer.Reset()
	consoleLogger := NewLogContext().
		WithCallerKey(true, CallerOption{}).
		WithWriter(AddSync(buffer)).
		WithEncoder(Console).
		Build()
	consoleLogger.Info("caller")

	line := buffer.String()
	if !strings.Contains(line, "logx/logx_test.go:") || !strings.HasSuffix(line, "\tcaller\n") {
		t.Fatalf("unexpected console caller output: %q", line)
	}
}

func TestLogContextCopyAndInvalidEncoder(t *testing.T) {
	for _, tt := range []EncoderType{Console, Json, Text} {
		t.Run(strconv.Itoa(int(tt)), func(t *testing.T) {
			copied := NewLogContext().WithEncoder(tt).Copy()
			switch tt {
			case Console:
				if _, ok := copied.enc.(*ConsoleEncoder); !ok {
					t.Fatalf("expected copied encoder to be *ConsoleEncoder, got %T", copied.enc)
				}
			case Json:
				if _, ok := copied.enc.(*JsonEncoder); !ok {
					t.Fatalf("expected copied encoder to be *JsonEncoder, got %T", copied.enc)
				}
			case Text:
				if _, ok := copied.enc.(*TextEncoder); !ok {
					t.Fatalf("expected copied encoder to be *TextEncoder, got %T", copied.enc)
				}
			}
		})
	}

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for unsupported encoder")
		} else if !strings.Contains(r.(string), "not support other log encoder") {
			t.Fatalf("unexpected panic: %v", r)
		}
	}()
	NewLogContext().WithEncoder(EncoderType(0))
}

func TestLogContextCopyHasIndependentLevel(t *testing.T) {
	level := NewAtomicLevel(LevelInfo)
	original := NewLogContext().WithAtomicLevel(level)
	copied := original.Copy()

	level.SetLevel(LevelDebug)
	if got := copied.AtomicLevel().Level(); got != LevelInfo {
		t.Fatalf("copied level changed with original: got %v, want %v", got, LevelInfo)
	}

	copied.WithLevel(LevelError)
	if got := original.AtomicLevel().Level(); got != LevelDebug {
		t.Fatalf("original level changed with copy: got %v, want %v", got, LevelDebug)
	}
}

type recordingWriteSyncer struct {
	bytes.Buffer
	syncs int
}

func (w *recordingWriteSyncer) Sync() error {
	w.syncs++
	return nil
}

type errorWriteSyncer struct {
	writeErr   error
	syncErr    error
	shortWrite bool
	syncs      int
}

func (w *errorWriteSyncer) Write(p []byte) (int, error) {
	if w.writeErr != nil {
		return 0, w.writeErr
	}
	if w.shortWrite && len(p) > 0 {
		return len(p) - 1, nil
	}
	return len(p), nil
}

func (w *errorWriteSyncer) Sync() error {
	w.syncs++
	return w.syncErr
}

type countingStringer struct {
	calls *int
}

func (s countingStringer) String() string {
	(*s.calls)++
	return "formatted"
}

func encodeForTest(t *testing.T, ctx *LogContext, encoderType EncoderType, ent entry, fields []Field) string {
	t.Helper()
	ctx.WithMsgKey(ctx.msgKey)
	ctx.WithEncoder(encoderType)
	ctx.enc.Init()

	buf, err := ctx.enc.Encode(ent, fields)
	if err != nil {
		t.Fatal(err)
	}
	defer bufPool.Put(buf)
	return string(buf.Bytes())
}

func encodeTextForTest(t *testing.T, ctx *LogContext, ent entry, fields []Field) string {
	t.Helper()
	return encodeForTest(t, ctx, Text, ent, fields)
}

func decodeJSONLog(t *testing.T, data string) map[string]any {
	t.Helper()
	var got map[string]any
	if err := json.Unmarshal([]byte(data), &got); err != nil {
		t.Fatalf("invalid JSON %q: %v", data, err)
	}
	return got
}

func TestBufferIOHelpers(t *testing.T) {
	if got := NewBuffer(nil).String(); got != "" {
		t.Fatalf("empty Buffer.String() = %q, want empty", got)
	}

	buf := NewBuffer(make([]byte, 0, 4))
	buf.AppendBytes([]byte("go"))

	if n, err := buf.Write([]byte("lang")); err != nil || n != 4 {
		t.Fatalf("Write() = (%d, %v), want (4, nil)", n, err)
	}
	if err := buf.WriteByte('!'); err != nil {
		t.Fatalf("WriteByte() error = %v", err)
	}
	if n, err := buf.WriteString("\n"); err != nil || n != 1 {
		t.Fatalf("WriteString() = (%d, %v), want (1, nil)", n, err)
	}
	if buf.Cap() < buf.Len() {
		t.Fatalf("Cap() = %d, Len() = %d", buf.Cap(), buf.Len())
	}

	clone := buf.Clone()
	buf.TrimNewline()
	if got, want := buf.String(), "golang!"; got != want {
		t.Fatalf("after TrimNewline() got %q, want %q", got, want)
	}
	if got, want := clone.String(), "golang!\n"; got != want {
		t.Fatalf("clone changed with original buffer: got %q, want %q", got, want)
	}
}

func TestAddSyncAndLock(t *testing.T) {
	if got := AddSync(nil); got != nil {
		t.Fatalf("AddSync(nil) = %T, want nil", got)
	}
	if got := Lock(nil); got != nil {
		t.Fatalf("Lock(nil) = %T, want nil", got)
	}

	raw := bytes.NewBuffer(nil)
	wrapped := AddSync(raw)
	if n, err := wrapped.Write([]byte("wrapped")); err != nil || n != len("wrapped") {
		t.Fatalf("wrapped Write() = (%d, %v)", n, err)
	}
	if err := wrapped.Sync(); err != nil {
		t.Fatalf("wrapped Sync() error = %v", err)
	}
	if got := raw.String(); got != "wrapped" {
		t.Fatalf("wrapped writer got %q", got)
	}

	rec := &recordingWriteSyncer{}
	if got := AddSync(rec); got != rec {
		t.Fatalf("AddSync(existing WriteSyncer) should return the same value")
	}

	locked := Lock(rec)
	if n, err := locked.Write([]byte("lock")); err != nil || n != len("lock") {
		t.Fatalf("locked Write() = (%d, %v)", n, err)
	}
	if err := locked.Sync(); err != nil {
		t.Fatalf("locked Sync() error = %v", err)
	}
	if got := rec.String(); got != "lock" {
		t.Fatalf("locked writer got %q", got)
	}
	if rec.syncs != 1 {
		t.Fatalf("Sync() count = %d, want 1", rec.syncs)
	}
	if got := Lock(locked); got != locked {
		t.Fatalf("Lock() should not wrap an already locked writer")
	}
}

func TestLoggerLevelFilteringAndFormattedMethods(t *testing.T) {
	buffer := bytes.NewBuffer(nil)
	logger := NewLogContext().
		WithLevel(LevelWarn).
		WithWriter(AddSync(buffer)).
		WithEncoder(Text).
		Build()

	logger.Tracef("trace %d", 1)
	logger.Debug("debug")
	logger.Info("info")
	logger.Warnf("warn %d", 2)
	logger.ErrorWith(nil)

	want := "msg=\"warn 2\"\nmsg=<nil>\n"
	if got := buffer.String(); got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestLoggerEnabledSkipsDisabledFormatting(t *testing.T) {
	buffer := bytes.NewBuffer(nil)
	logger := NewLogContext().
		WithLevel(LevelWarn).
		WithWriter(AddSync(buffer)).
		WithEncoder(Text).
		Build()

	if logger.Enabled(LevelInfo) {
		t.Fatal("Info should be disabled at Warn level")
	}
	if !logger.Enabled(LevelError) {
		t.Fatal("Error should be enabled at Warn level")
	}

	calls := 0
	value := countingStringer{calls: &calls}
	logger.Infof("disabled %s", value)
	if calls != 0 {
		t.Fatalf("disabled Infof formatted its arguments %d time(s)", calls)
	}
	logger.Errorf("enabled %s", value)
	if calls != 1 {
		t.Fatalf("enabled Errorf formatted its arguments %d time(s), want 1", calls)
	}
	if got, want := buffer.String(), "msg=\"enabled formatted\"\n"; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestLoggerSyncAndHighSeverityFlush(t *testing.T) {
	writer := &recordingWriteSyncer{}
	logger := NewLogContext().
		WithLevel(LevelTrace).
		WithWriter(writer).
		WithEncoder(Text).
		Build()

	if err := logger.Sync(); err != nil {
		t.Fatalf("Sync() error = %v", err)
	}
	if writer.syncs != 1 {
		t.Fatalf("Sync() count = %d, want 1", writer.syncs)
	}

	func() {
		defer func() { _ = recover() }()
		logger.Panic("flush before panic")
	}()
	if writer.syncs != 2 {
		t.Fatalf("Panic() Sync count = %d, want 2", writer.syncs)
	}
}

func TestLoggerReportsOutputErrors(t *testing.T) {
	writeErr := errors.New("write failed")
	syncErr := errors.New("sync failed")
	tests := []struct {
		name      string
		writer    *errorWriteSyncer
		log       func(Logger)
		wantError error
		wantSyncs int
	}{
		{
			name:      "encode error",
			writer:    &errorWriteSyncer{},
			log:       func(logger Logger) { logger.Info("message", Field{}) },
			wantError: errInvalidFieldType,
		},
		{
			name:      "write error",
			writer:    &errorWriteSyncer{writeErr: writeErr},
			log:       func(logger Logger) { logger.Info("message") },
			wantError: writeErr,
		},
		{
			name:      "short write",
			writer:    &errorWriteSyncer{shortWrite: true},
			log:       func(logger Logger) { logger.Info("message") },
			wantError: io.ErrShortWrite,
		},
		{
			name:   "high severity sync error",
			writer: &errorWriteSyncer{syncErr: syncErr},
			log: func(logger Logger) {
				defer func() { _ = recover() }()
				logger.Panic("message")
			},
			wantError: syncErr,
			wantSyncs: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got error
			logger := NewLogContext().
				WithWriter(tt.writer).
				WithErrorHandler(func(err error) { got = err }).
				WithEncoder(Text).
				Build()

			tt.log(logger)
			if !errors.Is(got, tt.wantError) {
				t.Fatalf("reported error = %v, want %v", got, tt.wantError)
			}
			if tt.writer.syncs != tt.wantSyncs {
				t.Fatalf("Sync() count = %d, want %d", tt.writer.syncs, tt.wantSyncs)
			}
		})
	}
}

func TestPanicAttemptsSyncAfterOutputFailure(t *testing.T) {
	writeErr := errors.New("write failed")
	tests := []struct {
		name      string
		writer    *errorWriteSyncer
		fields    []Field
		wantError error
	}{
		{
			name:      "encode error",
			writer:    &errorWriteSyncer{},
			fields:    []Field{{}},
			wantError: errInvalidFieldType,
		},
		{
			name:      "write error",
			writer:    &errorWriteSyncer{writeErr: writeErr},
			wantError: writeErr,
		},
		{
			name:      "short write",
			writer:    &errorWriteSyncer{shortWrite: true},
			wantError: io.ErrShortWrite,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got error
			logger := NewLogContext().
				WithWriter(tt.writer).
				WithErrorHandler(func(err error) { got = err }).
				WithEncoder(Text).
				Build()

			func() {
				defer func() { _ = recover() }()
				logger.Panic("message", tt.fields...)
			}()

			if !errors.Is(got, tt.wantError) {
				t.Fatalf("reported error = %v, want %v", got, tt.wantError)
			}
			if tt.writer.syncs != 1 {
				t.Fatalf("Sync() count = %d, want 1", tt.writer.syncs)
			}
		})
	}
}

func TestPanicJoinsWriteAndSyncErrors(t *testing.T) {
	writeErr := errors.New("write failed")
	syncErr := errors.New("sync failed")
	writer := &errorWriteSyncer{writeErr: writeErr, syncErr: syncErr}
	var got error
	logger := NewLogContext().
		WithWriter(writer).
		WithErrorHandler(func(err error) { got = err }).
		WithEncoder(Text).
		Build()

	func() {
		defer func() { _ = recover() }()
		logger.Panic("message")
	}()

	if !errors.Is(got, writeErr) || !errors.Is(got, syncErr) {
		t.Fatalf("reported error = %v, want joined write and sync errors", got)
	}
	if writer.syncs != 1 {
		t.Fatalf("Sync() count = %d, want 1", writer.syncs)
	}
}

func TestLoggerSyncReturnsWriterError(t *testing.T) {
	syncErr := errors.New("sync failed")
	logger := NewLogContext().
		WithWriter(&errorWriteSyncer{syncErr: syncErr}).
		WithEncoder(Text).
		Build()

	if err := logger.Sync(); !errors.Is(err, syncErr) {
		t.Fatalf("Sync() error = %v, want %v", err, syncErr)
	}
}

func TestAtomicLevelUpdatesExistingLogger(t *testing.T) {
	buffer := bytes.NewBuffer(nil)
	level := NewAtomicLevel(LevelWarn)
	logger := NewLogContext().
		WithAtomicLevel(level).
		WithWriter(AddSync(buffer)).
		WithEncoder(Text).
		Build()

	logger.Info("hidden info")
	logger.Warn("visible warn")

	level.SetLevel(LevelTrace)
	logger.Debug("visible debug")

	level.SetLevel(LevelError)
	logger.Warn("hidden warn")
	logger.Error("visible error")

	want := strings.Join([]string{
		"msg=\"visible warn\"",
		"msg=\"visible debug\"",
		"msg=\"visible error\"",
		"",
	}, "\n")
	if got := buffer.String(); got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestLoggerMethodsWriteExpectedMessages(t *testing.T) {
	buffer := bytes.NewBuffer(nil)
	logger := NewLogContext().
		WithLevel(LevelTrace).
		WithWriter(AddSync(buffer)).
		WithEncoder(Text).
		Build()

	logger.Trace("trace")
	logger.Debugf("debug %d", 1)
	logger.Infof("info %s", "ok")
	logger.Warn("warn")
	logger.Error("error")
	logger.Errorf("errorf %s", "ok")
	logger.ErrorWith(errors.New("error with"))

	func() {
		defer func() {
			if got := recover(); got != "panic" {
				t.Fatalf("Panic() recovered %v, want panic", got)
			}
		}()
		logger.Panic("panic")
	}()
	func() {
		defer func() {
			err, ok := recover().(error)
			if !ok || err.Error() != "panic with" {
				t.Fatalf("PanicWith() recovered unexpected value")
			}
		}()
		logger.PanicWith(errors.New("panic with"))
	}()

	want := strings.Join([]string{
		"msg=trace",
		"msg=\"debug 1\"",
		"msg=\"info ok\"",
		"msg=warn",
		"msg=error",
		"msg=\"errorf ok\"",
		"msg=\"error with\"",
		"msg=panic",
		"msg=\"panic with\"",
		"",
	}, "\n")
	if got := buffer.String(); got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestLoggerWithClonesPreFields(t *testing.T) {
	buffer := bytes.NewBuffer(nil)
	parent := NewLogContext().
		WithFields(String("app", "root")).
		WithWriter(AddSync(buffer)).
		WithEncoder(Text).
		Build()
	child := parent.With(String("request", "abc"))

	parent.Info("parent")
	child.Info("child", Int("status", 200))
	parent.Info("parent2")

	want := strings.Join([]string{
		"msg=parent app=root",
		"msg=child app=root request=abc status=200",
		"msg=parent2 app=root",
		"",
	}, "\n")
	if got := buffer.String(); got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestBuildFreezesConfigurationAndCopiesNewFields(t *testing.T) {
	firstBuffer := bytes.NewBuffer(nil)
	secondBuffer := bytes.NewBuffer(nil)
	fields := []Field{String("scope", "initial")}
	ctx := NewLogContext().
		WithNewFields(fields...).
		WithWriter(AddSync(firstBuffer)).
		WithEncoder(Text)

	fields[0] = String("scope", "mutated through caller slice")
	logger := ctx.Build()
	ctx.WithMsgKey("changed").
		WithNewFields(String("scope", "changed after build")).
		WithWriter(AddSync(secondBuffer)).
		WithEncoder(Json)

	logger.Info("message")
	if got, want := firstBuffer.String(), "msg=message scope=initial\n"; got != want {
		t.Fatalf("built logger changed with its source context: got %q, want %q", got, want)
	}
	if got := secondBuffer.String(); got != "" {
		t.Fatalf("built logger wrote to a later configured writer: %q", got)
	}
}

func TestLoggerPanicMethods(t *testing.T) {
	buffer := bytes.NewBuffer(nil)
	logger := NewLogContext().WithWriter(AddSync(buffer)).WithEncoder(Text).Build()

	logger.PanicWith(nil)
	if got := buffer.String(); got != "" {
		t.Fatalf("PanicWith(nil) wrote %q, want empty output", got)
	}

	func() {
		defer func() {
			if got := recover(); got != "panic 7" {
				t.Fatalf("Panicf() recovered %v, want %q", got, "panic 7")
			}
		}()
		logger.Panicf("panic %d", 7)
	}()

	if got, want := buffer.String(), "msg=\"panic 7\"\n"; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestLoggerDiscardAndNilWriterSkipOutput(t *testing.T) {
	nilWriterLogger := NewLogContext().WithEncoder(Text).Build()
	if nilWriterLogger.Enabled(LevelInfo) {
		t.Fatal("logger with nil writer should be disabled")
	}
	nilWriterLogger.Info("nil writer")

	for _, tt := range []struct {
		name   string
		writer WriteSyncer
	}{
		{name: "AddSync", writer: AddSync(io.Discard)},
		{name: "Lock", writer: Lock(AddSync(io.Discard))},
	} {
		t.Run(tt.name, func(t *testing.T) {
			discardLogger := NewLogContext().
				WithWriter(tt.writer).
				WithEncoder(Text).
				Build()
			if discardLogger.Enabled(LevelInfo) {
				t.Fatal("logger writing to io.Discard should be disabled")
			}
			calls := 0
			discardLogger.Infof("discard %s", countingStringer{calls: &calls})
			if calls != 0 {
				t.Fatalf("discard logger formatted arguments %d time(s)", calls)
			}
		})
	}
}

func TestTextEncoderOptionsAndFieldTypes(t *testing.T) {
	loc := time.FixedZone("TST", 2*60*60)
	entTime := time.Date(2026, 6, 18, 12, 0, 0, 123, loc)
	fieldTime := time.Date(2026, 6, 18, 13, 14, 15, 0, loc)
	ctx := NewLogContext().
		WithNewFields(String("prefix", "new")).
		WithMsgKey("message").
		WithLevelKey(true, LevelOption{LevelKey: "severity", LowerKey: true}).
		WithTimeKey(true, TimeOption{TimeKey: "ts", Timestamp: true, Layout: time.RFC3339})

	got := encodeTextForTest(t, ctx, entry{level: LevelWarn, time: entTime, message: "hello"}, []Field{
		Int8("i8", -8),
		Int16("i16", -16),
		Int32("i32", -32),
		Int64("i64", -64),
		Int("i", -1),
		UInt8("u8", 8),
		UInt16("u16", 16),
		UInt32("u32", 32),
		UInt64("u64", 64),
		UInt("u", 1),
		Float32("f32", 1.5),
		Float64("f64", 2.25),
		Time("when", fieldTime),
		Duration("dur", 90*time.Second),
		Error("err", io.EOF),
		Error("nilerr", nil),
		Field{Key: "nil", Type: NilType},
		Any("anyString", "value"),
		Any("anyBool", true),
		Any("anyInt", int64(-3)),
		Any("anyUint", uint(3)),
		Any("anyFloat", 1.25),
		Any("anyTime", fieldTime),
		Any("anyDuration", time.Second),
		Any("anyError", io.ErrClosedPipe),
		Any("anyNil", nil),
		Any("singleField", String("nested", "value")),
		Any("fieldList", []Field{String("nested", "value"), Int("count", 2)}),
		Array("array", "x", 2),
	})

	want := "ts=" + strconv.FormatInt(entTime.UnixNano(), 10) +
		" severity=warn message=hello prefix=new" +
		" i8=-8 i16=-16 i32=-32 i64=-64 i=-1" +
		" u8=8 u16=16 u32=32 u64=64 u=1" +
		" f32=1.5 f64=2.25 when=" + strconv.FormatInt(fieldTime.UnixNano(), 10) + " dur=1m30s" +
		" err=EOF nilerr=<nil> nil=<nil> anyString=value anyBool=true" +
		" anyInt=-3 anyUint=3 anyFloat=1.25 anyTime=" + strconv.FormatInt(fieldTime.UnixNano(), 10) +
		" anyDuration=1s anyError=\"io: read/write on closed pipe\" anyNil=<nil>" +
		" singleField=\"nested=value\" fieldList=\"nested=value count=2\" array=\"[x 2]\""
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestTextEncoderInvalidFieldsReturnError(t *testing.T) {
	ent := entry{level: LevelInfo, time: time.Unix(0, 0), message: "msg"}

	for _, encType := range []EncoderType{Console, Json, Text} {
		t.Run(strconv.Itoa(int(encType)), func(t *testing.T) {
			ctx := NewLogContext().WithEncoder(encType)
			ctx.enc.Init()

			buf, err := ctx.enc.Encode(ent, []Field{{Key: "bad"}})
			if !errors.Is(err, errInvalidFieldType) {
				t.Fatalf("Encode() error = %v, want %v", err, errInvalidFieldType)
			}
			if buf != nil {
				t.Fatalf("Encode() buffer = %v, want nil", buf)
			}
		})
	}
}

func TestEncodersInvalidPrefixFieldsReturnError(t *testing.T) {
	ent := entry{level: LevelInfo, time: time.Unix(0, 0), message: "msg"}

	for _, encType := range []EncoderType{Console, Json, Text} {
		t.Run(strconv.Itoa(int(encType)), func(t *testing.T) {
			ctx := NewLogContext().WithFields(Field{}).WithEncoder(encType)
			ctx.enc.Init()

			buf, err := ctx.enc.Encode(ent, nil)
			if !errors.Is(err, errInvalidFieldType) {
				t.Fatalf("Encode() error = %v, want %v", err, errInvalidFieldType)
			}
			if buf != nil {
				t.Fatalf("Encode() buffer = %v, want nil", buf)
			}
		})
	}
}

func TestJsonEncoderAnyMapsAndReflectSlices(t *testing.T) {
	type customInt int

	ctx := NewLogContext().
		WithReflectValue(true)

	got := decodeJSONLog(t, encodeForTest(t, ctx, Json, entry{level: LevelInfo, time: time.Unix(0, 0), message: "json"}, []Field{
		Any("anyMap", map[string]any{
			"str": "value",
			"num": 2,
			"arr": []any{"x", true, nil},
		}),
		Any("strMap", map[string]string{"a": "b"}),
		Any("multiMap", map[string][]string{"h": {"x", "y"}}),
		Any("durations", []time.Duration{time.Second, 2 * time.Second}),
		Any("errors", []error{io.EOF, nil}),
		Any("strings", []string{"a", "b"}),
		Any("bools", []bool{true, false}),
		ArrayT("typedInts", 1, 2, 3),
		Any("singleField", String("nested", "value")),
		Any("customSlice", []customInt{1, 2}),
	}))

	if got["msg"] != "json" {
		t.Fatalf("msg = %v, want json", got["msg"])
	}
	anyMap := got["anyMap"].(map[string]any)
	if anyMap["str"] != "value" || anyMap["num"] != float64(2) {
		t.Fatalf("unexpected anyMap: %#v", anyMap)
	}
	arr := anyMap["arr"].([]any)
	if arr[0] != "x" || arr[1] != true || arr[2] != nil {
		t.Fatalf("unexpected anyMap.arr: %#v", arr)
	}
	if got["strMap"].(map[string]any)["a"] != "b" {
		t.Fatalf("unexpected strMap: %#v", got["strMap"])
	}
	multi := got["multiMap"].(map[string]any)["h"].([]any)
	if multi[0] != "x" || multi[1] != "y" {
		t.Fatalf("unexpected multiMap: %#v", got["multiMap"])
	}
	durations := got["durations"].([]any)
	if durations[0] != "1s" || durations[1] != "2s" {
		t.Fatalf("unexpected durations: %#v", durations)
	}
	errs := got["errors"].([]any)
	if errs[0] != "EOF" || errs[1] != nil {
		t.Fatalf("unexpected errors: %#v", errs)
	}
	stringsValue := got["strings"].([]any)
	if stringsValue[0] != "a" || stringsValue[1] != "b" {
		t.Fatalf("unexpected strings: %#v", stringsValue)
	}
	bools := got["bools"].([]any)
	if bools[0] != true || bools[1] != false {
		t.Fatalf("unexpected bools: %#v", bools)
	}
	typedInts := got["typedInts"].([]any)
	if typedInts[0] != float64(1) || typedInts[1] != float64(2) || typedInts[2] != float64(3) {
		t.Fatalf("unexpected typedInts: %#v", typedInts)
	}
	singleField := got["singleField"].(map[string]any)
	if singleField["nested"] != "value" {
		t.Fatalf("unexpected singleField: %#v", singleField)
	}
	custom := got["customSlice"].([]any)
	if custom[0] != "1" || custom[1] != "2" {
		t.Fatalf("unexpected customSlice: %#v", custom)
	}
}

func TestFieldConstructorsAndTextQuotingEdges(t *testing.T) {
	outOfRange := time.Date(1600, 1, 1, 0, 0, 0, 0, time.UTC)
	field := Time("old", outOfRange)
	if field.Type != TimeFullType || field.AnyValue != outOfRange {
		t.Fatalf("Time() for out-of-range value = %#v", field)
	}

	cases := []struct {
		name  string
		value string
		want  bool
	}{
		{name: "plain", value: "abc-123:/", want: false},
		{name: "empty", value: "", want: true},
		{name: "space", value: "a b", want: true},
		{name: "equal", value: "a=b", want: true},
		{name: "quote", value: `"`, want: true},
		{name: "delete", value: string([]byte{0x7f}), want: true},
		{name: "invalid utf8", value: string([]byte{0xff}), want: true},
		{name: "unicode space", value: "a\u00a0b", want: true},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			if got := textNeedsQuoting(tt.value); got != tt.want {
				t.Fatalf("textNeedsQuoting(%q) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func FuzzJsonEncoderStrings(f *testing.F) {
	f.Add("message", "key", "value")
	f.Add("quoted \"message\"", "quoted\"key", "line one\nline two")
	f.Add(string([]byte{0xff}), string([]byte{0xfe}), "\a\v\u2028\u2029")

	f.Fuzz(func(t *testing.T, msg, key, value string) {
		ctx := NewLogContext().WithEncoder(Json)
		ctx.WithMsgKey(ctx.msgKey)
		ctx.enc.Init()

		buf, err := ctx.enc.Encode(entry{
			level:   LevelInfo,
			time:    time.Unix(0, 0),
			message: msg,
		}, []Field{String(key, value)})
		if err != nil {
			t.Fatal(err)
		}
		defer bufPool.Put(buf)
		if !json.Valid(buf.Bytes()) {
			t.Fatalf("invalid JSON for msg=%q key=%q value=%q: %q", msg, key, value, buf.Bytes())
		}
	})
}

func TestAtomicLevelDynamicUpdate(t *testing.T) {
	var buffer bytes.Buffer
	atomicLevel := NewAtomicLevel(LevelInfo)
	logger := NewLogContext().
		WithAtomicLevel(atomicLevel).
		WithWriter(AddSync(&buffer)).
		WithEncoder(Console).
		Build()

	logger.Debug("debug hidden")
	if got := buffer.String(); got != "" {
		t.Fatalf("expected debug log to be skipped, got %q", got)
	}

	atomicLevel.SetLevel(LevelDebug)
	logger.Debug("debug visible")
	if got, want := buffer.String(), "debug visible\n"; got != want {
		t.Fatalf("unexpected log output: got %q, want %q", got, want)
	}
}

func TestAtomicLevelSetWarnSkipsInfo(t *testing.T) {
	var buffer bytes.Buffer
	atomicLevel := NewAtomicLevel(LevelTrace)
	logger := NewLogContext().
		WithAtomicLevel(atomicLevel).
		WithWriter(AddSync(&buffer)).
		WithEncoder(Console).
		Build()

	atomicLevel.SetLevel(LevelWarn)
	logger.Info("info hidden")
	if got := buffer.String(); got != "" {
		t.Fatalf("expected info log to be skipped, got %q", got)
	}

	logger.Warn("warn visible")
	if got, want := buffer.String(), "warn visible\n"; got != want {
		t.Fatalf("unexpected log output: got %q, want %q", got, want)
	}
}

func TestAtomicLevelUpdatePropagatesToChildLogger(t *testing.T) {
	var buffer bytes.Buffer
	atomicLevel := NewAtomicLevel(LevelInfo)
	logger := NewLogContext().
		WithAtomicLevel(atomicLevel).
		WithWriter(AddSync(&buffer)).
		WithEncoder(Console).
		Build()
	child := logger.With(String("scope", "child"))

	atomicLevel.SetLevel(LevelDebug)
	logger.Debug("parent visible")
	child.Debug("child visible")

	if got, want := buffer.String(), "parent visible\nchild visible\t{\"scope\":\"child\"}\n"; got != want {
		t.Fatalf("dynamic level did not propagate to child logger, got %q, want %q", got, want)
	}
}

func TestAtomicLevelConcurrentSetLevelAndLog(t *testing.T) {
	atomicLevel := NewAtomicLevel(LevelInfo)
	logger := NewLogContext().
		WithAtomicLevel(atomicLevel).
		WithWriter(AddSync(nullWriter{})).
		WithEncoder(Console).
		Build()

	levels := []LevelType{
		LevelTrace,
		LevelDebug,
		LevelInfo,
		LevelWarn,
		LevelError,
		LevelFatal,
		LevelPanic,
	}

	var wg sync.WaitGroup
	start := make(chan struct{})
	wg.Add(2)

	go func() {
		defer wg.Done()
		<-start
		for i := 0; i < 1000; i++ {
			logger.Debug("debug")
			logger.Info("info")
			logger.Warn("warn")
		}
	}()

	go func() {
		defer wg.Done()
		<-start
		for i := 0; i < 1000; i++ {
			atomicLevel.SetLevel(levels[i%len(levels)])
			_ = atomicLevel.Level()
		}
	}()

	close(start)
	wg.Wait()
}

type nullWriter struct{}

func (w nullWriter) Write(b []byte) (n int, err error) { return len(b), nil }

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
	logger := NewLogContext().WithWriter(AddSync(nullWriter{})).WithEncoder(Console).Build()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("this is a message")
	}
}

func BenchmarkJsonLogger(b *testing.B) {
	// disable level/time/caller attributes
	logger := NewLogContext().WithWriter(AddSync(nullWriter{})).WithEncoder(Json).Build()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("this is a message")
	}
}

func BenchmarkConsoleLoggerWithEscapeQuote(b *testing.B) {
	// disable level/time/caller attributes
	logger := NewLogContext().WithWriter(AddSync(nullWriter{})).WithEscapeQuote(true).WithEncoder(Console).Build()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info(`"this is a message"`)
	}
}

func BenchmarkJsonLoggerWithEscapeQuote(b *testing.B) {
	// disable level/time/caller attributes
	logger := NewLogContext().WithWriter(AddSync(nullWriter{})).WithEscapeQuote(true).WithEncoder(Json).Build()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info(`"this is a message"`)
	}
}

func BenchmarkConsoleLoggerWithSimpleField(b *testing.B) {
	// disable level/time/caller attributes
	logger := NewLogContext().WithWriter(AddSync(nullWriter{})).WithEncoder(Console).Build()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("this is a message", Int("key", 100))
	}
}

func BenchmarkJsonLoggerWithSimpleField(b *testing.B) {
	// disable level/time/caller attributes
	logger := NewLogContext().WithWriter(AddSync(nullWriter{})).WithEncoder(Json).Build()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("this is a message", Int("key", 100))
	}
}

func BenchmarkConsoleLoggerWithField(b *testing.B) {
	// disable level/time/caller attributes
	logger := NewLogContext().WithWriter(AddSync(nullWriter{})).WithEncoder(Console).Build()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("this is a message", String("key", "value"), Int("int", 10000))
	}
}

func BenchmarkJsonLoggerWithField(b *testing.B) {
	// disable level/time/caller attributes
	logger := NewLogContext().WithWriter(AddSync(nullWriter{})).WithEncoder(Json).Build()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("this is a message", String("key", "value"), Int("int", 10000))
	}
}

func BenchmarkConsoleLoggerWithPrefixField(b *testing.B) {
	logger := NewLogContext().
		WithLevelKey(true, LevelOption{}).
		WithTimeKey(true, TimeOption{}).
		WithWriter(AddSync(nullWriter{})).WithEncoder(Console).Build()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("this is a message", String("key", "value"), Int("int", 10000))
	}
}

func BenchmarkJsonLoggerWithPrefixField(b *testing.B) {
	logger := NewLogContext().
		WithLevelKey(true, LevelOption{}).
		WithTimeKey(true, TimeOption{}).
		WithWriter(AddSync(nullWriter{})).WithEncoder(Json).Build()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("this is a message", String("key", "value"), Int("int", 10000))
	}
}

func BenchmarkColorConsoleLoggerWithPrefixField(b *testing.B) {
	logger := NewLogContext().
		WithColorfulset(true, TextColorAttri{}).
		WithLevelKey(true, LevelOption{}).
		WithTimeKey(true, TimeOption{}).
		WithWriter(AddSync(nullWriter{})).WithEncoder(Console).Build()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("this is a message", String("key", "value"), Int("int", 10000))
	}
}

func BenchmarkColorJsonLoggerWithPrefixField(b *testing.B) {
	logger := NewLogContext().
		WithColorfulset(true, TextColorAttri{}).
		WithLevelKey(true, LevelOption{}).
		WithTimeKey(true, TimeOption{}).
		WithWriter(AddSync(nullWriter{})).WithEncoder(Json).Build()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("this is a message", String("key", "value"), Int("int", 10000))
	}
}

func BenchmarkSlogJsonLoggerWithPrefixField(b *testing.B) {
	logger := slog.New(slog.NewJSONHandler(nullWriter{}, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.LogAttrs(context.Background(), slog.LevelInfo, "this is a message", slog.String("key", "value"), slog.Int("int", 10000))
	}
}

func BenchmarkJsonLoggerWithReflectValueField(b *testing.B) {
	type Int int
	// disable level/time/caller attributes
	logger := NewLogContext().WithReflectValue(true).WithWriter(AddSync(nullWriter{})).WithEncoder(Json).Build()
	arr := []Int{10, 20, 30}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("this is a message", ArrayT("key", arr...))
	}
}

func BenchmarkConsoleLoggerWithEscapeQuoteWithField(b *testing.B) {
	// disable level/time/caller attributes
	logger := NewLogContext().WithWriter(AddSync(nullWriter{})).WithEscapeQuote(true).WithEncoder(Json).Build()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info(`"this is a message"`, String(`"key"`, `"value"`))
	}
}

func BenchmarkJsonLoggerWithEscapeQuoteWithField(b *testing.B) {
	// disable level/time/caller attributes
	logger := NewLogContext().WithWriter(AddSync(nullWriter{})).WithEscapeQuote(true).WithEncoder(Json).Build()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info(`"this is a message"`, String(`"key"`, `"value"`))
	}
}
