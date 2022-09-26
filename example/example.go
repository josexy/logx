package main

import (
	"github.com/josexy/logx"
)

var str = "hello golang"

func logxx() {
	logx.Info("%s", str)
	logx.Debug("%s", str)
	logx.Error("%s", str)
}

func main() {
	logx.SetFlags(logx.FlagPrefix | logx.FlagDatetime | logx.FlagLineNumber | logx.FlagFunction)
	// logx.SetOutput(io.Discard)
	logxx()
	logx.DisableColor = true
	logxx()
	logx.DisableColor = false
	logxx()
	logx.DisableDebug = true
	logxx()
	logx.DisableDebug = false
	logxx()
}
