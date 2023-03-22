package internal

import (
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

type CallerItem struct {
	Skip   int
	Enable bool
	Func   bool
	File   bool
	Line   bool
	Color  bool
}

func (c *CallerItem) String() string {
	if !c.Enable {
		return ""
	}
	var w strings.Builder
	pc, file, line, ok := runtime.Caller(c.Skip)
	if !ok {
		return ""
	}
	if c.File {
		w.WriteString(filepath.Base(file))
	}
	if c.Func {
		if c.File {
			w.WriteByte(':')
		}
		ls := strings.Split(runtime.FuncForPC(pc).Name(), ".")
		w.WriteString(ls[len(ls)-1])
	}
	if c.Line {
		w.WriteByte('#')
		w.WriteString(strconv.Itoa(line))
	}
	out := w.String()
	if c.Color {
		out = Yellow(out)
	}
	return out
}
