package logx

import (
	"errors"
	"time"
)

type (
	argType byte

	M    map[string]any
	Pair struct {
		Key   string
		Value any
	}
)

const (
	noneArg argType = iota
	stringArg
	boolArg
	int8Arg
	int16Arg
	int32Arg
	int64Arg
	intArg
	uint8Arg
	uint16Arg
	uint32Arg
	uint64Arg
	uintArg
	float32Arg
	float64Arg
	timeArg
	durationArg
	errorArg
	mapArg
	sortedMapArg
	sliceArg
	nilArg
)

type innerArg struct {
	string
	bool
	int8
	int16
	int32
	int64
	int
	uint8
	uint16
	uint32
	uint64
	uint
	float32
	float64
	time.Time
	time.Duration
	error
	M
	SM []Pair
	L  []any
}

var errEoc = errors.New("end of consumer")

type arg struct {
	key    string
	typ    argType
	inner  innerArg
	nowrap bool
}

type argsConsumer struct {
	index int
	args  []arg
}

func (l *argsConsumer) reset() {
	l.index = 0
	l.args = l.args[:0]
}

func (l *argsConsumer) put(arg ...arg) {
	if len(arg) > 0 {
		l.args = append(l.args, arg...)
	}
}

func (l *argsConsumer) hasNext() bool { return l.index < len(l.args) }

func (l *argsConsumer) getNext() (arg *arg, err error) {
	if l.index >= len(l.args) {
		err = errEoc
		return
	}
	arg = &l.args[l.index]
	l.index++
	return
}

func String(key string, value string) arg {
	return arg{
		key:   key,
		typ:   stringArg,
		inner: innerArg{string: value},
	}
}

func Bool(key string, value bool) arg {
	return arg{
		key:   key,
		typ:   boolArg,
		inner: innerArg{bool: value},
	}
}

func Int8(key string, value int8) arg {
	return arg{
		key:   key,
		typ:   int8Arg,
		inner: innerArg{int8: value},
	}
}

func Int16(key string, value int16) arg {
	return arg{
		key:   key,
		typ:   int16Arg,
		inner: innerArg{int16: value},
	}
}

func Int32(key string, value int32) arg {
	return arg{
		key:   key,
		typ:   int32Arg,
		inner: innerArg{int32: value},
	}
}

func Int64(key string, value int64) arg {
	return arg{
		key:   key,
		typ:   int64Arg,
		inner: innerArg{int64: value},
	}
}

func Int(key string, value int) arg {
	return arg{
		key:   key,
		typ:   intArg,
		inner: innerArg{int: value},
	}
}

func UInt8(key string, value uint8) arg {
	return arg{
		key:   key,
		typ:   uint8Arg,
		inner: innerArg{uint8: value},
	}
}

func UInt16(key string, value uint16) arg {
	return arg{
		key:   key,
		typ:   uint16Arg,
		inner: innerArg{uint16: value},
	}
}

func UInt32(key string, value uint32) arg {
	return arg{
		key:   key,
		typ:   uint32Arg,
		inner: innerArg{uint32: value},
	}
}

func UInt64(key string, value uint64) arg {
	return arg{
		key:   key,
		typ:   uint64Arg,
		inner: innerArg{uint64: value},
	}
}

func UInt(key string, value uint) arg {
	return arg{
		key:   key,
		typ:   uintArg,
		inner: innerArg{uint: value},
	}
}

func Float32(key string, value float32) arg {
	return arg{
		key:   key,
		typ:   float32Arg,
		inner: innerArg{float32: value},
	}
}

func Float64(key string, value float64) arg {
	return arg{
		key:   key,
		typ:   float64Arg,
		inner: innerArg{float64: value},
	}
}

func Time(key string, value time.Time) arg {
	return arg{
		key:   key,
		typ:   timeArg,
		inner: innerArg{Time: value},
	}
}

func Duration(key string, value time.Duration) arg {
	return arg{
		key:   key,
		typ:   durationArg,
		inner: innerArg{Duration: value},
	}
}

func Error(key string, value error) arg {
	return arg{
		key:   key,
		typ:   errorArg,
		inner: innerArg{error: value},
	}
}

func Map(key string, value M) arg {
	return arg{
		key:   key,
		typ:   mapArg,
		inner: innerArg{M: value},
	}
}

func SortedMap(key string, value ...Pair) arg {
	return arg{
		key:   key,
		typ:   sortedMapArg,
		inner: innerArg{SM: value},
	}
}

func Slice(key string, value []any) arg {
	return arg{
		key:   key,
		typ:   sliceArg,
		inner: innerArg{L: value},
	}
}
