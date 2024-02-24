package logx

type LevelType uint8

const (
	LevelTrace LevelType = 1 << iota
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
	LevelPanic
)

type levelField struct {
	enable bool
	lower  bool
	color  bool
	typ    LevelType
}

func (lvl *levelField) toArg() arg {
	return arg{
		key:    "level",
		typ:    stringArg,
		inner:  innerArg{string: lvl.String()},
		nowrap: true,
	}
}

func (lvl *levelField) String() (out string) {
	if !lvl.enable {
		return
	}
	switch lvl.typ {
	case LevelTrace:
		if lvl.lower {
			out = "trace"
		} else {
			out = "TRACE"
		}
		if lvl.color {
			out = Magenta(out)
		}
	case LevelDebug:
		if lvl.lower {
			out = "debug"
		} else {
			out = "DEBUG"
		}
		if lvl.color {
			out = HiCyan(out)
		}
	case LevelInfo:
		if lvl.lower {
			out = "info"
		} else {
			out = "INFO"
		}
		if lvl.color {
			out = Green(out)
		}
	case LevelWarn:
		if lvl.lower {
			out = "warn"
		} else {
			out = "WARN"
		}
		if lvl.color {
			out = Yellow(out)
		}
	case LevelError:
		if lvl.lower {
			out = "error"
		} else {
			out = "ERROR"
		}
		if lvl.color {
			out = Red(out)
		}
	case LevelFatal:
		if lvl.lower {
			out = "fatal"
		} else {
			out = "FATAL"
		}
		if lvl.color {
			out = HiRed(out)
		}
	case LevelPanic:
		if lvl.lower {
			out = "panic"
		} else {
			out = "PANIC"
		}
		if lvl.color {
			out = HiYellow(out)
		}
	default:
		out = "unknown"
	}
	return
}
