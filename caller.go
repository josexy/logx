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

func (c *callerField) format() Field {
	fileName, funcName, lineNum := c.value()
	sm := make([]Field, 0, 3)
	if fileName.has() {
		sm = append(sm, String("file", fileName.get()))
	}
	if funcName.has() {
		sm = append(sm, String("func", funcName.get()))
	}
	if lineNum.has() {
		sm = append(sm, Int("line", lineNum.get()))
	}
	return Object("caller", sm...)
}

func (c *callerField) value() (fileName, funcName optval[string], lineNum optval[int]) {
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
	fileName, funcName, lineNum := c.value()
	var w strings.Builder
	var packageName string
	if funcName.has() {
		parts := strings.Split(funcName.value, ".")
		if len(parts) > 1 {
			packageName = parts[0]
		}
	}
	if fileName.has() {
		if packageName != "" {
			w.WriteString(packageName)
			w.WriteByte('/')
		}
		w.WriteString(fileName.get())
	}
	if lineNum.has() {
		w.WriteByte(':')
		w.WriteString(strconv.Itoa(lineNum.get()))
	}
	out := w.String()
	if c.color {
		out = Yellow(out)
	}
	return out
}
