package logx

import (
	"bytes"
	"fmt"
	"strconv"
	"time"
	"unsafe"

	"github.com/fatih/color"
)

type colorAttri struct {
	keyColor    *color.Color
	stringColor *color.Color
	boolColor   *color.Color
	floatColor  *color.Color
	numberColor *color.Color
}

type JsonEncoder struct {
	*logContext
	*colorAttri
	buf      *bytes.Buffer
	consumer *argsConsumer
}

func (enc *JsonEncoder) Init() {
	enc.consumer = new(argsConsumer)
	if enc.callerField.enable {
		enc.callerField.skipDepth = 7
	}
	enc.colorAttri = nil
	if enc.levelField.color && enc.timeField.color && enc.callerField.color {
		enc.colorAttri = &colorAttri{
			keyColor:    colorMap[BlueAttr],
			stringColor: colorMap[GreenAttr],
			boolColor:   colorMap[YellowAttr],
			floatColor:  colorMap[CyanAttr],
			numberColor: colorMap[RedAttr],
		}
	}
}

func (enc *JsonEncoder) Encode(w *bytes.Buffer, msg string, args ...arg) error {
	enc.buf = w
	enc.buf.Reset()
	enc.consumer.reset()

	enc.beginObject()
	enc.basePrompt()
	enc.withPrefix()
	enc.consumer.put(arg{key: "msg", typ: stringArg, inner: innerArg{string: msg}})
	enc.consumer.put(args...)
	for enc.consumer.hasNext() {
		if err := enc.keyValue(enc.consumer); err != nil {
			return err
		}
	}
	enc.endObject()
	return nil
}

func (enc *JsonEncoder) beginObject() { enc.buf.WriteByte('{') }

func (enc *JsonEncoder) endObject() { enc.buf.WriteByte('}') }

func (enc *JsonEncoder) beginArray() { enc.buf.WriteByte('[') }

func (enc *JsonEncoder) endArray() { enc.buf.WriteByte(']') }

func (enc *JsonEncoder) basePrompt() {
	if enc.levelField.enable {
		enc.consumer.put(enc.levelField.toArg())
	}
	if enc.timeField.enable {
		enc.consumer.put(enc.timeField.toArg())
	}
	if enc.callerField.enable {
		enc.consumer.put(enc.callerField.toArg())
	}
}

func (enc *JsonEncoder) withPrefix() {
	if len(enc.preKeyValues) == 0 {
		return
	}
	for i := 0; i < len(enc.preKeyValues); i++ {
		enc.consumer.put(convert(enc.preKeyValues[i].Key, enc.preKeyValues[i].Value))
	}
}

func (enc *JsonEncoder) splitComma(consumer *argsConsumer) {
	if consumer.index != len(consumer.args) {
		enc.buf.WriteByte(',')
	}
}

func (enc *JsonEncoder) wrapKey(key string) string {
	if enc.colorAttri != nil && enc.colorAttri.keyColor != nil {
		return enc.colorAttri.keyColor.Sprint(key)
	}
	return key
}

func (enc *JsonEncoder) value(arg *arg) {
	switch arg.typ {
	case stringArg:
		enc.string(arg.inner.string, !arg.nowrap)
	case boolArg:
		enc.bool(arg.inner.bool)
	case int8Arg:
		enc.int8(arg.inner.int8)
	case int16Arg:
		enc.int16(arg.inner.int16)
	case int32Arg:
		enc.int32(arg.inner.int32)
	case int64Arg:
		enc.int64(arg.inner.int64)
	case intArg:
		enc.int(arg.inner.int)
	case uint8Arg:
		enc.uint8(arg.inner.uint8)
	case uint16Arg:
		enc.uint16(arg.inner.uint16)
	case uint32Arg:
		enc.uint32(arg.inner.uint32)
	case uint64Arg:
		enc.uint64(arg.inner.uint64)
	case uintArg:
		enc.uint(arg.inner.uint)
	case float32Arg:
		enc.float32(arg.inner.float32)
	case float64Arg:
		enc.float64(arg.inner.float64)
	case timeArg:
		enc.time(arg.inner.Time)
	case durationArg:
		enc.duration(arg.inner.Duration)
	case errorArg:
		enc.error(arg.inner.error)
	case mapArg:
		enc.dict(arg.inner.M)
	case sortedMapArg:
		enc.sortedDict(arg.inner.SM)
	case sliceArg:
		enc.slice(arg.inner.L)
	case nilArg:
		enc.null()
	case anyArg:
		enc.any(arg.inner.any)
	}
}

func (enc *JsonEncoder) keyValue(consumer *argsConsumer) error {
	arg, err := consumer.getNext()
	if err != nil {
		return err
	}
	if arg.typ == noneArg {
		return errInvalidType
	}
	// "key": VALUE
	enc.key(arg.key, true)
	enc.buf.WriteByte(':')
	enc.value(arg)
	enc.splitComma(consumer)
	return nil
}

func (enc *JsonEncoder) wrapString(value string) string {
	if enc.colorAttri != nil && enc.colorAttri.stringColor != nil {
		return enc.colorAttri.stringColor.Sprint(value)
	}
	return value
}

func (enc *JsonEncoder) wrapBool(value string) string {
	if enc.colorAttri != nil && enc.colorAttri.boolColor != nil {
		return enc.colorAttri.boolColor.Sprint(value)
	}
	return value
}

func (enc *JsonEncoder) wrapFloat(value string) string {
	if enc.colorAttri != nil && enc.colorAttri.floatColor != nil {
		return enc.colorAttri.floatColor.Sprint(value)
	}
	return value
}

func (enc *JsonEncoder) wrapNumber(value string) string {
	if enc.colorAttri != nil && enc.colorAttri.numberColor != nil {
		return enc.colorAttri.numberColor.Sprint(value)
	}
	return value
}

func (enc *JsonEncoder) key(key string, isWrap bool) {
	enc.buf.WriteByte('"')
	if enc.escapeQuote && isWrap {
		key = quoteString(key)
	}
	if isWrap {
		enc.buf.WriteString(enc.wrapKey(key))
	} else {
		enc.buf.WriteString(key)
	}
	enc.buf.WriteByte('"')
}

func (enc *JsonEncoder) null() {
	enc.buf.WriteString(enc.wrapString("null"))
}

func (enc *JsonEncoder) any(value any) {
	enc.string(fmt.Sprintf("%+v", value), true)
}

func (enc *JsonEncoder) string(value string, wrap bool) {
	enc.buf.WriteByte('"')
	if enc.escapeQuote && wrap {
		value = quoteString(value)
	}
	if wrap {
		enc.buf.WriteString(enc.wrapString(value))
	} else {
		enc.buf.WriteString(value)
	}
	enc.buf.WriteByte('"')
}

func (enc *JsonEncoder) bool(value bool) {
	if value {
		enc.buf.WriteString(enc.wrapBool("true"))
	} else {
		enc.buf.WriteString(enc.wrapBool("false"))
	}
}

func (enc *JsonEncoder) int8(value int8) {
	enc.buf.WriteString(enc.wrapNumber(enc.wrapNumber(strconv.FormatInt(int64(value), 10))))
}

func (enc *JsonEncoder) int16(value int16) {
	enc.buf.WriteString(enc.wrapNumber(strconv.FormatInt(int64(value), 10)))
}

func (enc *JsonEncoder) int32(value int32) {
	enc.buf.WriteString(enc.wrapNumber(strconv.FormatInt(int64(value), 10)))
}

func (enc *JsonEncoder) int64(value int64) {
	enc.buf.WriteString(enc.wrapNumber(strconv.FormatInt(value, 10)))
}

func (enc *JsonEncoder) int(value int) {
	enc.buf.WriteString(enc.wrapNumber(strconv.FormatInt(int64(value), 10)))
}

func (enc *JsonEncoder) uint8(value uint8) {
	enc.buf.WriteString(enc.wrapNumber(strconv.FormatUint(uint64(value), 10)))
}

func (enc *JsonEncoder) uint16(value uint16) {
	enc.buf.WriteString(enc.wrapNumber(strconv.FormatUint(uint64(value), 10)))
}

func (enc *JsonEncoder) uint32(value uint32) {
	enc.buf.WriteString(enc.wrapNumber(strconv.FormatUint(uint64(value), 10)))
}

func (enc *JsonEncoder) uint64(value uint64) {
	enc.buf.WriteString(enc.wrapNumber(strconv.FormatUint(uint64(value), 10)))
}

func (enc *JsonEncoder) uint(value uint) {
	enc.buf.WriteString(enc.wrapNumber(strconv.FormatUint(uint64(value), 10)))
}

func (enc *JsonEncoder) float32(value float32) {
	enc.buf.WriteString(enc.wrapFloat(strconv.FormatFloat(float64(value), 'f', 3, 32)))
}

func (enc *JsonEncoder) float64(value float64) {
	enc.buf.WriteString(enc.wrapFloat(strconv.FormatFloat(value, 'f', 3, 64)))
}

func (enc *JsonEncoder) time(value time.Time) {
	enc.string(value.Format(time.DateTime), true)
}

func (enc *JsonEncoder) duration(value time.Duration) {
	enc.string(value.String(), true)
}

func (enc *JsonEncoder) error(value error) {
	if value == nil {
		enc.null()
	} else {
		enc.string(value.Error(), true)
	}
}

func (enc *JsonEncoder) dict(value M) {
	if value == nil {
		enc.null()
	} else {
		args := make([]arg, 0, len(value))
		for k, v := range value {
			args = append(args, convert(k, v))
		}
		enc.beginObject()
		consumerForJson := &argsConsumer{args: args}
		for consumerForJson.hasNext() {
			enc.keyValue(consumerForJson)
		}
		enc.endObject()
	}
}

func (enc *JsonEncoder) sortedDict(value []Pair) {
	if value == nil {
		enc.null()
	} else {
		args := make([]arg, 0, len(value))
		for _, p := range value {
			args = append(args, convert(p.Key, p.Value))
		}
		enc.beginObject()
		consumerForJson := &argsConsumer{args: args}
		for consumerForJson.hasNext() {
			enc.keyValue(consumerForJson)
		}
		enc.endObject()
	}
}

func (enc *JsonEncoder) slice(value []any) {
	if value == nil {
		enc.null()
	} else {
		args := make([]arg, 0, len(value))
		for _, v := range value {
			args = append(args, convert("", v))
		}
		enc.beginArray()
		consumer := &argsConsumer{args: args}
		for consumer.hasNext() {
			arg, _ := consumer.getNext()
			enc.value(arg)
			enc.splitComma(consumer)
		}
		enc.endArray()
	}
}

func convert(k string, v any) arg {
	switch ele := v.(type) {
	case string:
		return arg{key: k, typ: stringArg, inner: innerArg{string: ele}}
	case bool:
		return arg{key: k, typ: boolArg, inner: innerArg{bool: ele}}
	case int8:
		return arg{key: k, typ: int8Arg, inner: innerArg{int8: ele}}
	case int16:
		return arg{key: k, typ: int16Arg, inner: innerArg{int16: ele}}
	case int32:
		return arg{key: k, typ: int32Arg, inner: innerArg{int32: ele}}
	case int64:
		return arg{key: k, typ: int64Arg, inner: innerArg{int64: ele}}
	case int:
		return arg{key: k, typ: intArg, inner: innerArg{int: ele}}
	case uint8:
		return arg{key: k, typ: uint8Arg, inner: innerArg{uint8: ele}}
	case uint16:
		return arg{key: k, typ: uint16Arg, inner: innerArg{uint16: ele}}
	case uint32:
		return arg{key: k, typ: uint32Arg, inner: innerArg{uint32: ele}}
	case uint64:
		return arg{key: k, typ: uint64Arg, inner: innerArg{uint64: ele}}
	case uint:
		return arg{key: k, typ: uintArg, inner: innerArg{uint: ele}}
	case float32:
		return arg{key: k, typ: float32Arg, inner: innerArg{float32: ele}}
	case float64:
		return arg{key: k, typ: float64Arg, inner: innerArg{float64: ele}}
	case time.Time:
		return arg{key: k, typ: timeArg, inner: innerArg{Time: ele}}
	case time.Duration:
		return arg{key: k, typ: durationArg, inner: innerArg{Duration: ele}}
	case error:
		return arg{key: k, typ: errorArg, inner: innerArg{error: ele}}
	case M:
		return arg{key: k, typ: mapArg, inner: innerArg{M: ele}}
	case []Pair:
		return arg{key: k, typ: sortedMapArg, inner: innerArg{SM: ele}}
	case []any:
		return arg{key: k, typ: sliceArg, inner: innerArg{L: ele}}
	case nil:
		return arg{key: k, typ: nilArg}
	default:
		return arg{key: k, typ: anyArg, inner: innerArg{any: v}}
	}
}

func quoteString(value string) string {
	buf := make([]byte, 0, 3*len(value)/2)
	data := strconv.AppendQuote(buf, value)
	return bytesToString(data[1 : len(data)-1])
}

func bytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
