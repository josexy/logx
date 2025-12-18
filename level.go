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
	levelTypeColorMap = map[LevelType]ColorAttr{
		LevelTrace: MagentaAttr,
		LevelDebug: HiCyanAttr,
		LevelInfo:  GreenAttr,
		LevelWarn:  YellowAttr,
		LevelError: RedAttr,
		LevelFatal: HiRedAttr,
		LevelPanic: HiYellowAttr,
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

func (lvl *levelField) appendWithJson(enc *JsonEncoder, ent *entry) {
	enc.writeFieldKey(lvl.option.LevelKey)
	enc.writeSplitColon()
	enc.writeQuote()
	lvl.append(enc.buf, ent)
	enc.writeQuote()
}

func (lvl *levelField) append(buf *Buffer, ent *entry) {
	var level string
	if lvl.option.LowerKey {
		level = levelTypeLowerMap[ent.level]
	} else {
		level = levelTypeUpperMap[ent.level]
	}
	if lvl.color {
		appendColor(buf, levelTypeColorMap[ent.level], level)
		return
	}
	buf.AppendString(level)
}
