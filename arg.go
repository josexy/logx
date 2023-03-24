package logx

import (
	"time"

	"github.com/josexy/logx/internal"
)

type Arg = internal.Arg

func String(key string, value string) Arg          { return internal.String(key, value) }
func Bool(key string, value bool) Arg              { return internal.Bool(key, value) }
func Int8(key string, value int8) Arg              { return internal.Int8(key, value) }
func Int16(key string, value int16) Arg            { return internal.Int16(key, value) }
func Int32(key string, value int32) Arg            { return internal.Int32(key, value) }
func Int64(key string, value int64) Arg            { return internal.Int64(key, value) }
func Int(key string, value int) Arg                { return internal.Int(key, value) }
func UInt8(key string, value uint8) Arg            { return internal.UInt8(key, value) }
func UInt16(key string, value uint16) Arg          { return internal.UInt16(key, value) }
func UInt32(key string, value uint32) Arg          { return internal.UInt32(key, value) }
func UInt64(key string, value uint64) Arg          { return internal.UInt64(key, value) }
func UInt(key string, value uint) Arg              { return internal.UInt(key, value) }
func Float32(key string, value float32) Arg        { return internal.Float32(key, value) }
func Float64(key string, value float64) Arg        { return internal.Float64(key, value) }
func Time(key string, value time.Time) Arg         { return internal.Time(key, value) }
func Duration(key string, value time.Duration) Arg { return internal.Duration(key, value) }
func Error(key string, value error) Arg            { return internal.Error(key, value) }
