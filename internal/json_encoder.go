package internal

import (
	"bytes"
	"strconv"
	"time"

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
	buf       *bytes.Buffer
	numArgs   int
	totalArgs int
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

func (enc *JsonEncoder) calcPrefixCount() (count int) {
	if enc.LevelItem.Enable {
		count++
	}
	if enc.TimeItem.Enable {
		count++
	}
	if enc.CallerItem.Enable {
		count++
	}
	return
}

func (enc *JsonEncoder) Encode(w *bytes.Buffer, msg string, args ...Arg) error {
	enc.buf = w
	enc.numArgs = 0
	enc.totalArgs = enc.calcPrefixCount() + 1 + len(args)

	enc.begin()
	enc.writePrefix()
	enc.writeMsg(msg)
	for i := 0; i < len(args); i++ {
		enc.writeKeyValue2(&args[i])
	}
	enc.end()
	return nil
}

func (enc *JsonEncoder) begin() {
	enc.buf.Reset()
	enc.buf.WriteByte('{')
}

func (enc *JsonEncoder) end() {
	enc.buf.WriteByte('}')
}

func (enc *JsonEncoder) writePrefix() {
	if enc.LevelItem.Enable {
		enc.writeKeyValue1("level", enc.LevelItem.String())
	}
	if enc.TimeItem.Enable {
		enc.writeKeyValue1("time", enc.TimeItem.String())
	}
	if enc.CallerItem.Enable {
		enc.writeKeyValue1("caller", enc.CallerItem.String())
	}
}

func (enc *JsonEncoder) writeMsg(msg string) {
	enc.writeKey("msg", true)
	enc.buf.WriteByte(':')
	enc.writeString(msg, true)
	enc.writeSplitComma()
}

func (enc *JsonEncoder) writeSplitComma() {
	if enc.numArgs+1 == enc.totalArgs {
		return
	}
	enc.buf.WriteByte(',')
	enc.numArgs++
}

func (enc *JsonEncoder) writeKeyValue1(key, value string) {
	// "key": "value"
	enc.writeKey(key, true)
	enc.buf.WriteByte(':')
	enc.writeString(value, false)
	enc.writeSplitComma()
}

func (enc *JsonEncoder) writeKeyValue2(arg *Arg) {
	// "key": VALUE
	enc.writeKey(arg.key, true)
	enc.buf.WriteByte(':')
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
	}
	enc.writeSplitComma()
}

func (enc *JsonEncoder) wrapKey(key string) string {
	if enc.keyColor != nil {
		return enc.keyColor.Sprint(key)
	}
	return key
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
	if isWrap {
		enc.buf.WriteString(enc.wrapKey(key))
	} else {
		enc.buf.WriteString(key)
	}
	enc.buf.WriteByte('"')
}

func (enc *JsonEncoder) writeString(value string, isWrap bool) {
	enc.buf.WriteByte('"')
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
		enc.writeString("<nil>", true)
	} else {
		enc.writeString(value.Error(), true)
	}
}
