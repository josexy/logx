package logx

import (
	"bytes"
	"fmt"
	"math"
	"strconv"
	"time"
	"unsafe"
)

type JsonEncoder struct {
	*logContext
	buf          *bytes.Buffer
	fieldsRanger sliceFields
}

func (enc *JsonEncoder) Init() {
	enc.fieldsRanger = sliceFields{fields: make([]Field, 0, 64)}

	if enc.callerField.enable {
		enc.callerField.skipDepth = 7
	}
	enc.colorfulset.init()
}

func (enc *JsonEncoder) Encode(buf *bytes.Buffer, msg string, fields ...Field) error {
	enc.buf = buf
	enc.fieldsRanger.reset()

	enc.writeBeginObject()
	enc.addPromptFields()
	enc.addPrefixFields()
	enc.fieldsRanger.put(String("msg", msg))
	enc.fieldsRanger.put(fields...)

	enc.fieldsRanger.writeRangeFields(enc.writeField, enc.writeSplitComma)

	enc.writeEndObject()
	return nil
}

func (enc *JsonEncoder) writeBeginObject() { enc.buf.WriteByte('{') }

func (enc *JsonEncoder) writeEndObject() { enc.buf.WriteByte('}') }

func (enc *JsonEncoder) writeBeginArray() { enc.buf.WriteByte('[') }

func (enc *JsonEncoder) writeEndArray() { enc.buf.WriteByte(']') }

func (enc *JsonEncoder) addPromptFields() {
	if enc.levelField.enable {
		enc.fieldsRanger.put(enc.levelField.format())
	}
	if enc.timeField.enable {
		enc.fieldsRanger.put(enc.timeField.format())
	}
	if enc.callerField.enable && (enc.callerField.fileName || enc.callerField.funcName || enc.callerField.lineNum) {
		enc.fieldsRanger.put(enc.callerField.format())
	}
}

func (enc *JsonEncoder) addPrefixFields() {
	if len(enc.preFields) == 0 {
		return
	}
	enc.fieldsRanger.put(enc.preFields...)
}

func (enc *JsonEncoder) writeSplitComma() {
	enc.buf.WriteByte(',')
}

func (enc *JsonEncoder) wrapKey(key string) string {
	if enc.colorfulset.enable && enc.colorAttri.keyColor != nil {
		return enc.colorAttri.keyColor.Sprint(key)
	}
	return key
}

func (enc *JsonEncoder) writeFieldValue(field Field) {
	switch field.Type {
	case StringType:
		enc.writeFieldString(field.StringValue, !field.NoWrap)
	case BoolType:
		enc.writeFieldBool(field.IntValue == 1)
	case Int8Type:
		enc.writeFieldInt8(int8(field.IntValue))
	case Int16Type:
		enc.writeFieldInt16(int16(field.IntValue))
	case Int32Type:
		enc.writeFieldInt32(int32(field.IntValue))
	case Int64Type:
		enc.writeFieldInt64(field.IntValue)
	case IntType:
		enc.writeFieldInt(int(field.IntValue))
	case Uint8Type:
		enc.writeFieldUint8(uint8(field.IntValue))
	case Uint16Type:
		enc.writeFieldUint16(uint16(field.IntValue))
	case Uint32Type:
		enc.writeFieldUint32(uint32(field.IntValue))
	case Uint64Type:
		enc.writeFieldUint64(uint64(field.IntValue))
	case UintType:
		enc.writeFieldUint(uint(field.IntValue))
	case Float32Type:
		enc.writeFieldFloat32(math.Float32frombits(uint32(field.IntValue)))
	case Float64Type:
		enc.writeFieldFloat64(math.Float64frombits(uint64(field.IntValue)))
	case TimeType:
		if field.AnyValue != nil {
			enc.writeFieldTime(time.Unix(0, field.IntValue).In(field.AnyValue.(*time.Location)))
		} else {
			enc.writeFieldTime(time.Unix(0, field.IntValue))
		}
	case TimeFullType:
		enc.writeFieldTime(field.AnyValue.(time.Time))
	case DurationType:
		enc.writeFieldDuration(time.Duration(field.IntValue))
	case ErrorType:
		if field.AnyValue == nil {
			enc.writeFieldNil()
			return
		}
		enc.writeFieldError(field.AnyValue.(error))
	case ObjectType:
		enc.writeFieldObject(field.AnyValue.([]Field))
	case ArrayType:
		enc.writeFieldArray(field.AnyValue)
	case NilType:
		enc.writeFieldNil()
	case AnyType:
		enc.writeFieldAny(field.AnyValue)
	}
}

func (enc *JsonEncoder) writeField(field Field) error {
	if field.Type == NoneType {
		return errInvalidFieldType
	}
	enc.writeFieldKey(field.Key, true)
	enc.buf.WriteByte(':')
	enc.writeFieldValue(field)
	return nil
}

func (enc *JsonEncoder) wrapString(value string) string {
	if enc.colorfulset.enable && enc.colorAttri.stringColor != nil {
		return enc.colorAttri.stringColor.Sprint(value)
	}
	return value
}

func (enc *JsonEncoder) wrapBool(value string) string {
	if enc.colorfulset.enable && enc.colorAttri.boolColor != nil {
		return enc.colorAttri.boolColor.Sprint(value)
	}
	return value
}

func (enc *JsonEncoder) wrapFloat(value string) string {
	if enc.colorfulset.enable && enc.colorAttri.floatColor != nil {
		return enc.colorAttri.floatColor.Sprint(value)
	}
	return value
}

func (enc *JsonEncoder) wrapNumber(value string) string {
	if enc.colorfulset.enable && enc.colorAttri.numberColor != nil {
		return enc.colorAttri.numberColor.Sprint(value)
	}
	return value
}

func (enc *JsonEncoder) writeFieldKey(key string, isWrap bool) {
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

func (enc *JsonEncoder) writeFieldString(value string, wrap bool) {
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

func (enc *JsonEncoder) writeFieldBool(value bool) {
	if value {
		enc.buf.WriteString(enc.wrapBool("true"))
	} else {
		enc.buf.WriteString(enc.wrapBool("false"))
	}
}

func (enc *JsonEncoder) writeFieldInt8(value int8) {
	enc.buf.WriteString(enc.wrapNumber(enc.wrapNumber(strconv.FormatInt(int64(value), 10))))
}

func (enc *JsonEncoder) writeFieldInt16(value int16) {
	enc.buf.WriteString(enc.wrapNumber(strconv.FormatInt(int64(value), 10)))
}

func (enc *JsonEncoder) writeFieldInt32(value int32) {
	enc.buf.WriteString(enc.wrapNumber(strconv.FormatInt(int64(value), 10)))
}

func (enc *JsonEncoder) writeFieldInt64(value int64) {
	enc.buf.WriteString(enc.wrapNumber(strconv.FormatInt(value, 10)))
}

func (enc *JsonEncoder) writeFieldInt(value int) {
	enc.buf.WriteString(enc.wrapNumber(strconv.FormatInt(int64(value), 10)))
}

func (enc *JsonEncoder) writeFieldUint8(value uint8) {
	enc.buf.WriteString(enc.wrapNumber(strconv.FormatUint(uint64(value), 10)))
}

func (enc *JsonEncoder) writeFieldUint16(value uint16) {
	enc.buf.WriteString(enc.wrapNumber(strconv.FormatUint(uint64(value), 10)))
}

func (enc *JsonEncoder) writeFieldUint32(value uint32) {
	enc.buf.WriteString(enc.wrapNumber(strconv.FormatUint(uint64(value), 10)))
}

func (enc *JsonEncoder) writeFieldUint64(value uint64) {
	enc.buf.WriteString(enc.wrapNumber(strconv.FormatUint(uint64(value), 10)))
}

func (enc *JsonEncoder) writeFieldUint(value uint) {
	enc.buf.WriteString(enc.wrapNumber(strconv.FormatUint(uint64(value), 10)))
}

func (enc *JsonEncoder) writeFieldFloat32(value float32) {
	enc.buf.WriteString(enc.wrapFloat(strconv.FormatFloat(float64(value), 'f', 3, 32)))
}

func (enc *JsonEncoder) writeFieldFloat64(value float64) {
	enc.buf.WriteString(enc.wrapFloat(strconv.FormatFloat(value, 'f', 3, 64)))
}

func (enc *JsonEncoder) writeFieldTime(value time.Time) {
	enc.writeFieldString(value.Format(time.DateTime), true)
}

func (enc *JsonEncoder) writeFieldDuration(value time.Duration) {
	enc.writeFieldString(value.String(), true)
}

func (enc *JsonEncoder) writeFieldError(value error) {
	if value == nil {
		enc.writeFieldNil()
	} else {
		enc.writeFieldString(value.Error(), true)
	}
}

func (enc *JsonEncoder) writeFieldObject(value []Field) {
	enc.writeBeginObject()
	consumer := sliceFields{fields: value}
	consumer.writeRangeFields(enc.writeField, enc.writeSplitComma)
	enc.writeEndObject()
}

func (enc *JsonEncoder) writeFieldSingleObject(value Field) {
	enc.writeBeginObject()
	enc.writeField(value)
	enc.writeEndObject()
}

func (enc *JsonEncoder) writeFieldArray(value any) {
	enc.writeBeginArray()
	enc.writeFieldAny(value)
	enc.writeEndArray()
}

func (enc *JsonEncoder) writeFieldNil() {
	enc.buf.WriteString(enc.wrapString("null"))
}

func writeFieldArrayListFor[T any](value []T, wf func(T), lf func()) {
	for i := 0; i < len(value); i++ {
		wf(value[i])
		if i+1 != len(value) {
			lf()
		}
	}
}

func (enc *JsonEncoder) writeFieldAny(value any) {
	switch v := value.(type) {
	case []string:
		writeFieldArrayListFor(v, func(s string) { enc.writeFieldString(s, true) }, enc.writeSplitComma)
	case []bool:
		writeFieldArrayListFor(v, func(s bool) { enc.writeFieldBool(s) }, enc.writeSplitComma)
	case []int8:
		writeFieldArrayListFor(v, func(s int8) { enc.writeFieldInt8(s) }, enc.writeSplitComma)
	case []int16:
		writeFieldArrayListFor(v, func(s int16) { enc.writeFieldInt16(s) }, enc.writeSplitComma)
	case []int32:
		writeFieldArrayListFor(v, func(s int32) { enc.writeFieldInt32(s) }, enc.writeSplitComma)
	case []int64:
		writeFieldArrayListFor(v, func(s int64) { enc.writeFieldInt64(s) }, enc.writeSplitComma)
	case []int:
		writeFieldArrayListFor(v, func(s int) { enc.writeFieldInt(s) }, enc.writeSplitComma)
	case []uint8:
		writeFieldArrayListFor(v, func(s uint8) { enc.writeFieldUint8(s) }, enc.writeSplitComma)
	case []uint16:
		writeFieldArrayListFor(v, func(s uint16) { enc.writeFieldUint16(s) }, enc.writeSplitComma)
	case []uint32:
		writeFieldArrayListFor(v, func(s uint32) { enc.writeFieldUint32(s) }, enc.writeSplitComma)
	case []uint64:
		writeFieldArrayListFor(v, func(s uint64) { enc.writeFieldUint64(s) }, enc.writeSplitComma)
	case []uint:
		writeFieldArrayListFor(v, func(s uint) { enc.writeFieldUint(s) }, enc.writeSplitComma)
	case []float32:
		writeFieldArrayListFor(v, func(s float32) { enc.writeFieldFloat32(s) }, enc.writeSplitComma)
	case []float64:
		writeFieldArrayListFor(v, func(s float64) { enc.writeFieldFloat64(s) }, enc.writeSplitComma)
	case []time.Time:
		writeFieldArrayListFor(v, func(s time.Time) { enc.writeFieldTime(s) }, enc.writeSplitComma)
	case []time.Duration:
		writeFieldArrayListFor(v, func(s time.Duration) { enc.writeFieldDuration(s) }, enc.writeSplitComma)
	case []error:
		writeFieldArrayListFor(v, func(s error) { enc.writeFieldError(s) }, enc.writeSplitComma)
	case []Field:
		enc.writeFieldObject(v)
	case []any:
		writeFieldArrayListFor(v, func(s any) { enc.writeFieldAny(s) }, enc.writeSplitComma)
	case string:
		enc.writeFieldString(v, true)
	case bool:
		enc.writeFieldBool(v)
	case int8:
		enc.writeFieldInt8(v)
	case int16:
		enc.writeFieldInt16(v)
	case int32:
		enc.writeFieldInt32(v)
	case int64:
		enc.writeFieldInt64(v)
	case int:
		enc.writeFieldInt(v)
	case uint8:
		enc.writeFieldUint8(v)
	case uint16:
		enc.writeFieldUint16(v)
	case uint32:
		enc.writeFieldUint32(v)
	case uint64:
		enc.writeFieldUint64(v)
	case uint:
		enc.writeFieldUint(v)
	case float32:
		enc.writeFieldFloat32(v)
	case float64:
		enc.writeFieldFloat64(v)
	case time.Time:
		enc.writeFieldTime(v)
	case time.Duration:
		enc.writeFieldDuration(v)
	case error:
		enc.writeFieldError(v)
	case Field:
		enc.writeFieldSingleObject(v)
	case nil:
		enc.writeFieldNil()
	default:
		enc.writeFieldString(fmt.Sprintf("%+v", value), true)
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
