package logx

type LevelType uint8

const (
	LevelTrace LevelType = iota
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
	LevelPanic
)

var (
	levelTypeLowerMap = map[LevelType]string{
		LevelTrace: "trace",
		LevelDebug: "debug",
		LevelInfo:  "info",
		LevelWarn:  "warn",
		LevelError: "error",
		LevelFatal: "fatal",
		LevelPanic: "panic",
	}
	levelTypeUpperMap = map[LevelType]string{
		LevelTrace: "TRACE",
		LevelDebug: "DEBUG",
		LevelInfo:  "INFO",
		LevelWarn:  "WARN",
		LevelError: "ERROR",
		LevelFatal: "FATAL",
		LevelPanic: "PANIC",
	}
	levelTypeColorMap = map[LevelType]func(string) string{
		LevelTrace: Magenta,
		LevelDebug: HiCyan,
		LevelInfo:  Green,
		LevelWarn:  Yellow,
		LevelError: Red,
		LevelFatal: HiRed,
		LevelPanic: HiYellow,
	}
)

type LevelOption struct {
	// level key, default: "level"
	LevelKey string
	// lower level key
	LowerKey bool
}

type levelField struct {
	enable bool
	color  bool
	typ    LevelType
	option LevelOption
}

func (lvl *levelField) formatJson(enc *JsonEncoder) {
	enc.writeFieldKey(lvl.option.LevelKey)
	enc.writeSplitColon()
	enc.writeFieldStringPrimitive(lvl.String())
}

func (lvl *levelField) String() (out string) {
	if !lvl.enable {
		return
	}
	if lvl.option.LowerKey {
		out = levelTypeLowerMap[lvl.typ]
	} else {
		out = levelTypeUpperMap[lvl.typ]
	}
	if lvl.color && len(out) > 0 {
		out = levelTypeColorMap[lvl.typ](out)
	}
	return
}
