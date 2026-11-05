package logx

import (
	"errors"
	"fmt"
	"math"
	"reflect"
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
		StringValue string
		AnyValue    any
		IntValue    int64
		Type        FieldType
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

// valid reports whether f can be encoded: its type is known and AnyValue holds
// what that type requires. Fields made by the constructors in this file are
// always valid; hand-built Fields may not be, and encoding them must not panic.
func (f *Field) valid() bool {
	switch f.Type {
	case NoneType:
		return false
	case TimeType:
		if f.AnyValue == nil {
			return true
		}
		loc, ok := f.AnyValue.(*time.Location)
		return ok && loc != nil
	case TimeFullType:
		_, ok := f.AnyValue.(time.Time)
		return ok
	case ErrorType:
		if f.AnyValue == nil {
			return true
		}
		_, ok := f.AnyValue.(error)
		return ok
	case ObjectType:
		_, ok := f.AnyValue.([]Field)
		return ok
	default:
		return f.Type <= AnyType
	}
}

// errorMessage calls err.Error() without letting it crash the caller. It
// reports isNil for a nil error and for a typed-nil pointer whose Error method
// panics (for example a nil *MyErr returned as error). Any other panic raised
// by Error is rendered into the message, the way fmt reports such panics.
func errorMessage(err error) (msg string, isNil bool) {
	if err == nil {
		return "", true
	}
	defer func() {
		if r := recover(); r != nil {
			if v := reflect.ValueOf(err); v.Kind() == reflect.Pointer && v.IsNil() {
				msg, isNil = "", true
				return
			}
			msg, isNil = fmt.Sprintf("%%!v(PANIC=Error method: %v)", r), false
		}
	}()
	return err.Error(), false
}

// errorText is errorMessage with nil errors rendered as "<nil>".
func errorText(err error) string {
	msg, isNil := errorMessage(err)
	if isNil {
		return "<nil>"
	}
	return msg
}
