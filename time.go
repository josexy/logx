package logx

import (
	"fmt"
	"strconv"
	"time"
)

type timeField struct {
	enable bool
	format func(t time.Time) any
	now    time.Time
	color  bool
}

func (t *timeField) toArg() arg {
	arg := arg{key: "time"}
	switch val := t.Value(); val.(type) {
	case string:
		arg.typ = stringArg
		arg.inner = innerArg{string: val.(string)}
	case int64:
		arg.typ = int64Arg
		arg.inner = innerArg{int64: val.(int64)}
	default:
		arg.typ = stringArg
		arg.inner = innerArg{string: fmt.Sprintf("%v", val)}
	}
	return arg
}

func (t *timeField) Value() (out any) {
	if !t.enable {
		return
	}
	if t.format == nil {
		return
	}
	out = t.format(t.now)
	return
}

func (t *timeField) String() string {
	arg := t.toArg()
	var out string
	switch arg.typ {
	case stringArg:
		out = arg.inner.string
	case int64Arg:
		out = strconv.FormatInt(arg.inner.int64, 10)
	}
	if t.color && len(out) > 0 {
		out = Blue(out)
	}
	return out
}
