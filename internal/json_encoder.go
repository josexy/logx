package internal

import (
	"bytes"
	"fmt"
	"strconv"
	"time"
	"unsafe"

	"github.com/fatih/color"
)

type colorWrapper struct {
	keyColor    *color.Color
	stringColor *color.Color
	boolColor   *color.Color
	floatColor  *color.Color
	numberColor *color.Color
}

type JsonEncoder struct {
	colorWrapper
	*Config
	buf      *bytes.Buffer
	consumer *argsConsumer
}

func (enc *JsonEncoder) Init() {
	if enc.LevelItem.Color {
		enc.colorWrapper = colorWrapper{
			keyColor:    colorMap[BlueAttr],
			stringColor: colorMap[GreenAttr],
			boolColor:   colorMap[YellowAttr],
			floatColor:  colorMap[CyanAttr],
			numberColor: colorMap[RedAttr],
		}
	}
	if enc.CallerItem.Enable {
		enc.CallerItem.Skip = 5
	}
}

func (enc *JsonEncoder) Encode(w *bytes.Buffer, msg string, args ...Arg) error {
	enc.buf = w
	enc.consumer = &argsConsumer{args: args}
	enc.buf.Reset()

	enc.beginObject()
	enc.writePrefix()
	enc.writeMsg(msg)
	enc.writeSplitComma(enc.consumer)
	for enc.consumer.hasNext() {
		enc.writeKeyValue2(enc.consumer)
	}
	enc.endObject()
	return nil
}

func (enc *JsonEncoder) beginObject() { enc.buf.WriteByte('{') }

func (enc *JsonEncoder) endObject() { enc.buf.WriteByte('}') }

func (enc *JsonEncoder) beginArray() { enc.buf.WriteByte('[') }

func (enc *JsonEncoder) endArray() { enc.buf.WriteByte(']') }

func (enc *JsonEncoder) writePrefix() {
	enc.consumer.index = -1
	if enc.LevelItem.Enable {
		enc.writeKeyValue1("level", enc.LevelItem.String())
	}
	if enc.TimeItem.Enable {
		enc.writeKeyValue1("time", enc.TimeItem.String())
	}
	if enc.CallerItem.Enable {
		enc.writeKeyValue1("caller", enc.CallerItem.String())
	}
	enc.consumer.index = 0
}

func (enc *JsonEncoder) writeMsg(msg string) {
	enc.writeKey("msg", true)
	enc.buf.WriteByte(':')
	enc.writeString(msg, true)
}

func (enc *JsonEncoder) writeSplitComma(consumer *argsConsumer) {
	if consumer.index != len(consumer.args) {
		enc.buf.WriteByte(',')
	}
}

func (enc *JsonEncoder) writeKeyValue1(key, value string) {
	// "key": "value"
	enc.writeKey(key, true)
	enc.buf.WriteByte(':')
	enc.writeString(value, false)
	enc.writeSplitComma(enc.consumer)
}

func (enc *JsonEncoder) wrapKey(key string) string {
	if enc.keyColor != nil {
		return enc.keyColor.Sprint(key)
	}
	return key
}

func (enc *JsonEncoder) writeValue(arg *Arg) {
	switch arg.typ {
	case StringArg:
		enc.writeString(arg.inner.string, true)
	case BoolArg:
		enc.writeBool(arg.inner.bool)
	case Int8Arg:
		enc.writeInt8(arg.inner.int8)
	case Int16Arg:
		enc.writeInt16(arg.inner.int16)
	case Int32Arg:
		enc.writeInt32(arg.inner.int32)
	case Int64Arg:
		enc.writeInt64(arg.inner.int64)
	case IntArg:
		enc.writeInt(arg.inner.int)
	case UInt8Arg:
		enc.writeUint8(arg.inner.uint8)
	case UInt16Arg:
		enc.writeUint16(arg.inner.uint16)
	case UInt32Arg:
		enc.writeUint32(arg.inner.uint32)
	case UInt64Arg:
		enc.writeUint64(arg.inner.uint64)
	case UIntArg:
		enc.writeUint(arg.inner.uint)
	case Float32Arg:
		enc.writeFloat32(arg.inner.float32)
	case Float64Arg:
		enc.writeFloat64(arg.inner.float64)
	case TimeArg:
		enc.writeTime(arg.inner.Time)
	case DurationArg:
		enc.writeDuration(arg.inner.Duration)
	case ErrorArg:
		enc.writeError(arg.inner.error)
	case MapArg:
		enc.writeMap(arg.inner.M)
	case SliceArg:
		enc.writeSlice(arg.inner.L)
	case NilArg:
		enc.writeNull()
	}
}

func (enc *JsonEncoder) writeKeyValue2(consumer *argsConsumer) {
	arg, err := consumer.getNext()
	if err != nil {
		return
	}
	// "key": VALUE
	enc.writeKey(arg.key, true)
	enc.buf.WriteByte(':')
	enc.writeValue(&arg)
	enc.writeSplitComma(consumer)
}

func (enc *JsonEncoder) wrapString(value string) string {
	if enc.stringColor != nil {
		return enc.stringColor.Sprint(value)
	}
	return value
}

func (enc *JsonEncoder) wrapBool(value string) string {
	if enc.boolColor != nil {
		return enc.boolColor.Sprint(value)
	}
	return value
}

func (enc *JsonEncoder) wrapFloat(value string) string {
	if enc.floatColor != nil {
		return enc.floatColor.Sprint(value)
	}
	return value
}

func (enc *JsonEncoder) wrapNumber(value string) string {
	if enc.numberColor != nil {
		return enc.numberColor.Sprint(value)
	}
	return value
}

func (enc *JsonEncoder) writeKey(key string, isWrap bool) {
	enc.buf.WriteByte('"')
	if enc.Config.EscapeQuote && isWrap {
		buf := make([]byte, 0, 3*len(key)/2)
		data := strconv.AppendQuote(buf, key)
		key = bytesToString(data[1 : len(data)-1])
	}
	if isWrap {
		enc.buf.WriteString(enc.wrapKey(key))
	} else {
		enc.buf.WriteString(key)
	}
	enc.buf.WriteByte('"')
}

func (enc *JsonEncoder) writeNull() {
	enc.buf.WriteString(enc.wrapString("null"))
}

func bytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func (enc *JsonEncoder) writeString(value string, isWrap bool) {
	enc.buf.WriteByte('"')
	if enc.Config.EscapeQuote && isWrap {
		buf := make([]byte, 0, 3*len(value)/2)
		data := strconv.AppendQuote(buf, value)
		value = bytesToString(data[1 : len(data)-1])
	}
	if isWrap {
		enc.buf.WriteString(enc.wrapString(value))
	} else {
		enc.buf.WriteString(value)
	}
	enc.buf.WriteByte('"')
}

func (enc *JsonEncoder) writeBool(value bool) {
	if value {
		enc.buf.WriteString(enc.wrapBool("true"))
	} else {
		enc.buf.WriteString(enc.wrapBool("false"))
	}
}

func (enc *JsonEncoder) writeInt8(value int8) {
	enc.buf.WriteString(enc.wrapNumber(enc.wrapNumber(strconv.FormatInt(int64(value), 10))))
}

func (enc *JsonEncoder) writeInt16(value int16) {
	enc.buf.WriteString(enc.wrapNumber(strconv.FormatInt(int64(value), 10)))
}

func (enc *JsonEncoder) writeInt32(value int32) {
	enc.buf.WriteString(enc.wrapNumber(strconv.FormatInt(int64(value), 10)))
}

func (enc *JsonEncoder) writeInt64(value int64) {
	enc.buf.WriteString(enc.wrapNumber(strconv.FormatInt(value, 10)))
}

func (enc *JsonEncoder) writeInt(value int) {
	enc.buf.WriteString(enc.wrapNumber(strconv.FormatInt(int64(value), 10)))
}

func (enc *JsonEncoder) writeUint8(value uint8) {
	enc.buf.WriteString(enc.wrapNumber(strconv.FormatUint(uint64(value), 10)))
}

func (enc *JsonEncoder) writeUint16(value uint16) {
	enc.buf.WriteString(enc.wrapNumber(strconv.FormatUint(uint64(value), 10)))
}

func (enc *JsonEncoder) writeUint32(value uint32) {
	enc.buf.WriteString(enc.wrapNumber(strconv.FormatUint(uint64(value), 10)))
}

func (enc *JsonEncoder) writeUint64(value uint64) {
	enc.buf.WriteString(enc.wrapNumber(strconv.FormatUint(uint64(value), 10)))
}

func (enc *JsonEncoder) writeUint(value uint) {
	enc.buf.WriteString(enc.wrapNumber(strconv.FormatUint(uint64(value), 10)))
}

func (enc *JsonEncoder) writeFloat32(value float32) {
	enc.buf.WriteString(enc.wrapFloat(strconv.FormatFloat(float64(value), 'f', 3, 32)))
}

func (enc *JsonEncoder) writeFloat64(value float64) {
	enc.buf.WriteString(enc.wrapFloat(strconv.FormatFloat(value, 'f', 3, 64)))
}

func (enc *JsonEncoder) writeTime(value time.Time) {
	enc.writeString(value.Format(time.DateTime), true)
}

func (enc *JsonEncoder) writeDuration(value time.Duration) {
	enc.writeString(value.String(), true)
}

func (enc *JsonEncoder) writeError(value error) {
	if value == nil {
		enc.writeNull()
	} else {
		enc.writeString(value.Error(), true)
	}
}

func (enc *JsonEncoder) writeMap(value M) {
	if value == nil {
		enc.writeNull()
	} else {
		args := make([]Arg, 0, len(value))
		for k, v := range value {
			args = append(args, convert(k, v))
		}
		enc.beginObject()
		consumer := &argsConsumer{args: args}
		for consumer.hasNext() {
			enc.writeKeyValue2(consumer)
		}
		enc.endObject()
	}
}

func (enc *JsonEncoder) writeSlice(value []any) {
	if value == nil {
		enc.writeNull()
	} else {
		args := make([]Arg, 0, len(value))
		for _, v := range value {
			args = append(args, convert("", v))
		}
		enc.beginArray()
		consumer := &argsConsumer{args: args}
		for consumer.hasNext() {
			arg, _ := consumer.getNext()
			enc.writeValue(&arg)
			enc.writeSplitComma(consumer)
		}
		enc.endArray()
	}
}

func convert(k string, v any) Arg {
	switch ele := v.(type) {
	case string:
		return Arg{key: k, typ: StringArg, inner: innerArg{string: ele}}
	case bool:
		return Arg{key: k, typ: BoolArg, inner: innerArg{bool: ele}}
	case int8:
		return Arg{key: k, typ: Int8Arg, inner: innerArg{int8: ele}}
	case int16:
		return Arg{key: k, typ: Int16Arg, inner: innerArg{int16: ele}}
	case int32:
		return Arg{key: k, typ: Int32Arg, inner: innerArg{int32: ele}}
	case int64:
		return Arg{key: k, typ: Int64Arg, inner: innerArg{int64: ele}}
	case int:
		return Arg{key: k, typ: IntArg, inner: innerArg{int: ele}}
	case uint8:
		return Arg{key: k, typ: UInt8Arg, inner: innerArg{uint8: ele}}
	case uint16:
		return Arg{key: k, typ: UInt16Arg, inner: innerArg{uint16: ele}}
	case uint32:
		return Arg{key: k, typ: UInt32Arg, inner: innerArg{uint32: ele}}
	case uint64:
		return Arg{key: k, typ: UInt64Arg, inner: innerArg{uint64: ele}}
	case uint:
		return Arg{key: k, typ: UIntArg, inner: innerArg{uint: ele}}
	case float32:
		return Arg{key: k, typ: Float32Arg, inner: innerArg{float32: ele}}
	case float64:
		return Arg{key: k, typ: Float64Arg, inner: innerArg{float64: ele}}
	case time.Time:
		return Arg{key: k, typ: TimeArg, inner: innerArg{Time: ele}}
	case time.Duration:
		return Arg{key: k, typ: DurationArg, inner: innerArg{Duration: ele}}
	case error:
		return Arg{key: k, typ: ErrorArg, inner: innerArg{error: ele}}
	case M:
		return Arg{key: k, typ: MapArg, inner: innerArg{M: ele}}
	case []any:
		return Arg{key: k, typ: SliceArg, inner: innerArg{L: ele}}
	case nil:
		return Arg{key: k, typ: NilArg}
	default:
		return Arg{key: k, typ: StringArg, inner: innerArg{string: fmt.Sprintf("%v", v)}}
	}
}
