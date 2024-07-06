package logx

import (
	"errors"
	"math"
	"time"
)

const (
	NoneType FieldType = iota
	StringType
	BoolType
	Int8Type
	Int16Type
	Int32Type
	Int64Type
	IntType
	Uint8Type
	Uint16Type
	Uint32Type
	Uint64Type
	UintType
	Float32Type
	Float64Type
	TimeFullType
	TimeType
	DurationType
	ErrorType
	ObjectType
	ArrayType
	NilType
	AnyType
)

var (
	_minTimeInt64 = time.Unix(0, math.MinInt64)
	_maxTimeInt64 = time.Unix(0, math.MaxInt64)
)

var errInvalidFieldType = errors.New("invalid field type")

type (
	FieldType byte

	Field struct {
		Key         string
		Type        FieldType
		IntValue    int64
		StringValue string
		AnyValue    any
		NoWrap      bool
	}
)

func String(key string, value string) Field {
	return Field{Key: key, Type: StringType, StringValue: value}
}

func Bool(key string, value bool) Field {
	var intvalue int64
	if value {
		intvalue = 1
	}
	return Field{Key: key, Type: BoolType, IntValue: intvalue}
}

func Int8(key string, value int8) Field {
	return Field{Key: key, Type: Int8Type, IntValue: int64(value)}
}

func Int16(key string, value int16) Field {
	return Field{Key: key, Type: Int16Type, IntValue: int64(value)}
}

func Int32(key string, value int32) Field {
	return Field{Key: key, Type: Int32Type, IntValue: int64(value)}
}

func Int64(key string, value int64) Field {
	return Field{Key: key, Type: Int64Type, IntValue: value}
}

func Int(key string, value int) Field {
	return Field{Key: key, Type: IntType, IntValue: int64(value)}
}

func UInt8(key string, value uint8) Field {
	return Field{Key: key, Type: Uint8Type, IntValue: int64(value)}
}

func UInt16(key string, value uint16) Field {
	return Field{Key: key, Type: Uint16Type, IntValue: int64(value)}
}

func UInt32(key string, value uint32) Field {
	return Field{Key: key, Type: Uint32Type, IntValue: int64(value)}
}

func UInt64(key string, value uint64) Field {
	return Field{Key: key, Type: Uint64Type, IntValue: int64(value)}
}

func UInt(key string, value uint) Field {
	return Field{Key: key, Type: UintType, IntValue: int64(value)}
}

func Float32(key string, value float32) Field {
	return Field{Key: key, Type: Float32Type, IntValue: int64(math.Float32bits(value))}
}

func Float64(key string, value float64) Field {
	return Field{Key: key, Type: Float64Type, IntValue: int64(math.Float64bits(value))}
}

func Time(key string, value time.Time) Field {
	if value.Before(_minTimeInt64) || value.After(_maxTimeInt64) {
		return Field{Key: key, Type: TimeFullType, AnyValue: value}
	}
	return Field{Key: key, Type: TimeType, IntValue: value.UnixNano(), AnyValue: value.Location()}
}

func Duration(key string, value time.Duration) Field {
	return Field{Key: key, Type: DurationType, IntValue: int64(value)}
}

func Error(key string, value error) Field {
	return Field{Key: key, Type: ErrorType, AnyValue: value}
}

func Object(key string, value ...Field) Field {
	return Field{Key: key, Type: ObjectType, AnyValue: value}
}

func Array(key string, value ...any) Field {
	return Field{Key: key, Type: ArrayType, AnyValue: value}
}

func ArrayT[T any](key string, value ...T) Field {
	return Field{Key: key, Type: ArrayType, AnyValue: value}
}

func Any(key string, value any) Field {
	return Field{Key: key, Type: AnyType, AnyValue: value}
}
