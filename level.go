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
	option LevelOption
}

func (lvl *levelField) formatJson(enc *JsonEncoder, ent *entry) {
	enc.writeFieldKey(lvl.option.LevelKey)
	enc.writeSplitColon()
	enc.writeFieldStringPrimitive(lvl.String(ent))
}

func (lvl *levelField) String(ent *entry) (out string) {
	if !lvl.enable {
		return
	}
	if lvl.option.LowerKey {
		out = levelTypeLowerMap[ent.level]
	} else {
		out = levelTypeUpperMap[ent.level]
	}
	if lvl.color && len(out) > 0 {
		out = levelTypeColorMap[ent.level](out)
	}
	return
}
