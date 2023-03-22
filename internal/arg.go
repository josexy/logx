package internal

import "time"

type ArgType uint8

const (
	StringArg ArgType = iota
	BoolArg
	Int8Arg
	Int16Arg
	Int32Arg
	Int64Arg
	IntArg
	UInt8Arg
	UInt16Arg
	UInt32Arg
	UInt64Arg
	UIntArg
	Float32Arg
	Float64Arg
	TimeArg
	DurationArg
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
}

type Arg struct {
	key   string
	typ   ArgType
	inner innerArg
}

func String(key string, value string) Arg {
	return Arg{
		key:   key,
		typ:   StringArg,
		inner: innerArg{string: value},
	}
}

func Bool(key string, value bool) Arg {
	return Arg{
		key:   key,
		typ:   BoolArg,
		inner: innerArg{bool: value},
	}
}

func Int8(key string, value int8) Arg {
	return Arg{
		key:   key,
		typ:   Int8Arg,
		inner: innerArg{int8: value},
	}
}

func Int16(key string, value int16) Arg {
	return Arg{
		key:   key,
		typ:   Int16Arg,
		inner: innerArg{int16: value},
	}
}

func Int32(key string, value int32) Arg {
	return Arg{
		key:   key,
		typ:   Int32Arg,
		inner: innerArg{int32: value},
	}
}

func Int64(key string, value int64) Arg {
	return Arg{
		key:   key,
		typ:   Int64Arg,
		inner: innerArg{int64: value},
	}
}

func Int(key string, value int) Arg {
	return Arg{
		key:   key,
		typ:   IntArg,
		inner: innerArg{int: value},
	}
}

func UInt8(key string, value uint8) Arg {
	return Arg{
		key:   key,
		typ:   UInt8Arg,
		inner: innerArg{uint8: value},
	}
}

func UInt16(key string, value uint16) Arg {
	return Arg{
		key:   key,
		typ:   UInt16Arg,
		inner: innerArg{uint16: value},
	}
}

func UInt32(key string, value uint32) Arg {
	return Arg{
		key:   key,
		typ:   UInt32Arg,
		inner: innerArg{uint32: value},
	}
}

func UInt64(key string, value uint64) Arg {
	return Arg{
		key:   key,
		typ:   UInt64Arg,
		inner: innerArg{uint64: value},
	}
}

func UInt(key string, value uint) Arg {
	return Arg{
		key:   key,
		typ:   UIntArg,
		inner: innerArg{uint: value},
	}
}

func Float32(key string, value float32) Arg {
	return Arg{
		key:   key,
		typ:   Float32Arg,
		inner: innerArg{float32: value},
	}
}

func Float64(key string, value float64) Arg {
	return Arg{
		key:   key,
		typ:   Float64Arg,
		inner: innerArg{float64: value},
	}
}

func Time(key string, value time.Time) Arg {
	return Arg{
		key:   key,
		typ:   TimeArg,
		inner: innerArg{Time: value},
	}
}

func Duration(key string, value time.Duration) Arg {
	return Arg{
		key:   key,
		typ:   DurationArg,
		inner: innerArg{Duration: value},
	}
}
