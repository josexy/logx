package internal

type LevelType uint8

const (
	LevelDebug LevelType = 1 << iota
	LevelInfo
	LevelWarn
	LevelError
	LevelPanic
	LevelFatal
)

type LevelItem struct {
	Enable bool
	Lower  bool
	Typ    LevelType
	Color  bool
}

func (lvl *LevelItem) String() (out string) {
	if !lvl.Enable {
		return
	}
	switch lvl.Typ {
	case LevelDebug:
		if lvl.Lower {
			out = "debug"
		} else {
			out = "DEBUG"
		}
		if lvl.Color {
			out = HiCyan(out)
		}
	case LevelInfo:
		if lvl.Lower {
			out = "info"
		} else {
			out = "INFO"
		}
		if lvl.Color {
			out = Green(out)
		}
	case LevelWarn:
		if lvl.Lower {
			out = "warn"
		} else {
			out = "WARN"
		}
		if lvl.Color {
			out = Yellow(out)
		}
	case LevelError:
		if lvl.Lower {
			out = "error"
		} else {
			out = "ERROR"
		}
		if lvl.Color {
			out = Red(out)
		}
	case LevelPanic:
		if lvl.Lower {
			out = "panic"
		} else {
			out = "PANIC"
		}
		if lvl.Color {
			out = HiYellow(out)
		}
	case LevelFatal:
		if lvl.Lower {
			out = "fatal"
		} else {
			out = "FATAL"
		}
		if lvl.Color {
			out = HiRed(out)
		}
	default:
		out = "unknown"
	}
	return
}
