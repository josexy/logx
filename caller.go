package logx

import (
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

type callerField struct {
	enable    bool
	skipDepth int
	funcName  bool
	fileName  bool
	lineNum   bool
	color     bool
}

func (c *callerField) toArg() arg {
	fileName, funcName, lineNum := c.Value()
	sm := make([]Pair, 0, 3)
	if fileName.has() {
		sm = append(sm, Pair{Key: "file", Value: fileName.get()})
	}
	if funcName.has() {
		sm = append(sm, Pair{Key: "func", Value: funcName.get()})
	}
	if lineNum.has() {
		sm = append(sm, Pair{Key: "line", Value: lineNum.get()})
	}
	return arg{key: "caller", typ: sortedMapArg, inner: innerArg{SM: sm}}
}

func (c *callerField) Value() (fileName, funcName optval[string], lineNum optval[int]) {
	if !c.enable {
		return
	}
	pc, file, line, ok := runtime.Caller(c.skipDepth)
	if !ok {
		return
	}
	if c.fileName {
		fileName.set(filepath.Base(file))
	}
	if c.funcName {
		name := runtime.FuncForPC(pc).Name()
		parts := strings.Split(name, "/")
		funcName.set(parts[len(parts)-1])
	}
	if c.lineNum {
		lineNum.set(line)
	}
	return
}

func (c *callerField) String() string {
	fileName, funcName, lineNum := c.Value()
	var w strings.Builder
	if fileName.has() {
		w.WriteString(fileName.get())
	}
	if funcName.has() {
		if fileName.has() {
			w.WriteByte(':')
		}
		w.WriteString(funcName.get())
	}
	if lineNum.has() {
		w.WriteByte('#')
		w.WriteString(strconv.Itoa(lineNum.get()))
	}
	out := w.String()
	if c.color {
		out = Yellow(out)
	}
	return out
}
