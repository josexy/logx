package logx

import (
	"runtime"
	"strconv"
	"strings"
)

type CallerFormatter uint8

const (
	// package/file:line
	ShortFile CallerFormatter = iota
	// /full/path/to/package/file:line
	FullFile
	// package/file:line package.func
	ShortFileFunc
	// /full/path/to/package/file:line package.func
	FullFileFunc
)

type CallerOption struct {
	// caller key, default: "caller"
	CallerKey string
	// file key of caller, default: "file"
	FileKey string
	// function key of caller, default: "func"
	FuncKey string
	// caller formatter, default: ShortFileCaller
	Formatter CallerFormatter
	// caller skips increases the number of callers skipped by caller annotation.
	// when building wrappers around the Logger, supplying this Option prevents logx from always
	// reporting the wrapper code as the caller. default: 0
	CallerSkip int
}

type callerField struct {
	enable    bool
	skipDepth int
	color     bool
	option    CallerOption
}

func (c *callerField) formatJson(enc *JsonEncoder) {
	fileName, funcName := c.value()
	fields := make([]Field, 0, 2)
	if len(fileName) > 0 {
		fields = append(fields, String(c.option.FileKey, fileName))
	}
	if len(funcName) > 0 {
		fields = append(fields, String(c.option.FuncKey, funcName))
	}
	enc.writeFieldKey(c.option.CallerKey)
	enc.writeSplitColon()
	enc.writeFieldObject(fields)
}

func (c *callerField) value() (fileName, funcName string) {
	if !c.enable {
		return
	}
	pc, file, line, ok := runtime.Caller(c.skipDepth)
	if !ok {
		return
	}
	if c.option.Formatter == ShortFile || c.option.Formatter == ShortFileFunc {
		if idx := strings.LastIndexByte(file, '/'); idx != -1 {
			if idx = strings.LastIndexByte(file[:idx], '/'); idx != -1 {
				file = file[idx+1:]
			}
		}
	}
	fileName = file + ":" + strconv.FormatInt(int64(line), 10)

	if c.option.Formatter == ShortFileFunc || c.option.Formatter == FullFileFunc {
		funcName = runtime.FuncForPC(pc).Name()
		if idx := strings.LastIndexByte(funcName, '/'); idx != -1 {
			funcName = funcName[idx+1:]
		}
	}
	return
}

func (c *callerField) String() string {
	fileName, _ := c.value()
	if c.color {
		fileName = Yellow(fileName)
	}
	return fileName
}
