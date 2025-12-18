package logx

import (
	"fmt"
	"math"
	"reflect"
	"sync"
	"time"
)

var jsonPool = sync.Pool{New: func() any { return &JsonEncoder{} }}

type JsonEncoder struct {
	*LogContext
	buf *Buffer
}

func (enc *JsonEncoder) Init() {
	if enc.callerF.enable {
		enc.callerF.skipDepth = 7
		enc.callerF.skipDepth += enc.callerF.option.CallerSkip
	}
	enc.colors.init()
}

func (enc *JsonEncoder) clone() *JsonEncoder {
	clone := jsonPool.Get().(*JsonEncoder)
	clone.LogContext = enc.LogContext
	clone.buf = bufPool.Get().(*Buffer)
	clone.buf.Reset()
	return clone
}

func putJsonEncoder(enc *JsonEncoder) {
	enc.LogContext = nil
	enc.buf = nil
	jsonPool.Put(enc)
}

func (enc *JsonEncoder) Encode(ent entry, fields []Field) (ret *Buffer, err error) {
	nenc := enc.clone()
	defer putJsonEncoder(nenc)

	nenc.writeBeginObject()
	nenc.writePromptFields(&ent)
	if nenc.writePrefixFields() {
		nenc.writeSplitComma()
	}
	nenc.writeMsg(ent.message)

	n := len(fields)
	if n == 0 {
		nenc.writeEndObject()
		ret = nenc.buf
		return
	}
	nenc.writeSplitComma()

	for i := 0; i < n; i++ {
		if err = nenc.writeField(&fields[i]); err != nil {
			bufPool.Put(nenc.buf)
			return
		}
		if i+1 != n {
			nenc.writeSplitComma()
		}
	}
	nenc.writeEndObject()
	ret = nenc.buf
	return
}

func (enc *JsonEncoder) writeBeginObject() { enc.buf.AppendByte('{') }

func (enc *JsonEncoder) writeEndObject() { enc.buf.AppendByte('}') }

func (enc *JsonEncoder) writeBeginArray() { enc.buf.AppendByte('[') }

func (enc *JsonEncoder) writeEndArray() { enc.buf.AppendByte(']') }

func (enc *JsonEncoder) writePromptFields(ent *entry) {
	if enc.levelF.enable {
		enc.levelF.appendWithJson(enc, ent)
		enc.writeSplitComma()
	}
	if enc.timeF.enable {
		enc.timeF.appendWithJson(enc, ent)
		enc.writeSplitComma()
	}
	if enc.callerF.enable {
		enc.callerF.appendWithJson(enc)
		enc.writeSplitComma()
	}
}

func (enc *JsonEncoder) writeMsg(msg string) {
	enc.writeFieldKey(enc.msgKey)
	enc.writeSplitColon()
	enc.writeFieldString(msg)
}

func (enc *JsonEncoder) writePrefixFields() bool {
	n := len(enc.preFields)
	if n == 0 {
		return false
	}
	for i := 0; i < n; i++ {
		enc.writeField(&enc.preFields[i])
		if i+1 != n {
			enc.writeSplitComma()
		}
	}
	return true
}

func (enc *JsonEncoder) writeQuote() {
	enc.buf.AppendByte('"')
}

func (enc *JsonEncoder) writeSplitComma() {
	enc.buf.AppendByte(',')
}

func (enc *JsonEncoder) writeSplitColon() {
	enc.buf.AppendByte(':')
}

func (enc *JsonEncoder) writeFieldValue(field *Field) {
	switch field.Type {
	case StringType:
		enc.writeFieldString(field.StringValue)
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

func (enc *JsonEncoder) writeField(field *Field) error {
	if field.Type == NoneType {
		return errInvalidFieldType
	}
	enc.writeFieldKey(field.Key)
	enc.writeSplitColon()
	enc.writeFieldValue(field)
	return nil
}

func (enc *JsonEncoder) colorEnabled() bool { return enc.colors.enable }

func (enc *JsonEncoder) writeFieldKey(key string) {
	enc.writeQuote()
	if enc.escapeQuote {
		writeFieldWrapper(enc, enc.colors.attr.KeyColor, func() { appendQuoteString(enc.buf, key) })
	} else {
		writeFieldWrapper(enc, enc.colors.attr.KeyColor, func() { enc.buf.AppendString(key) })
	}
	enc.writeQuote()
}

func (enc *JsonEncoder) writeFieldString(value string) {
	enc.writeQuote()
	if enc.escapeQuote {
		writeFieldWrapper(enc, enc.colors.attr.StringColor, func() { appendQuoteString(enc.buf, value) })
	} else {
		writeFieldWrapper(enc, enc.colors.attr.StringColor, func() { enc.buf.AppendString(value) })
	}
	enc.writeQuote()
}

func (enc *JsonEncoder) writeFieldBool(value bool) {
	writeFieldWrapper(enc, enc.colors.attr.BooleanColor, func() { enc.buf.AppendBool(value) })
}

func writeFieldWrapper(enc *JsonEncoder, color ColorAttr, appendFn func()) {
	if enc.colorEnabled() {
		// \x1b[30mAAAAAAAAA\x1b[0m
		appendColorWithFunc(enc.buf, color, appendFn)
		return
	}
	appendFn()
}

func (enc *JsonEncoder) writeFieldInt8(value int8) {
	writeFieldWrapper(enc, enc.colors.attr.NumberColor, func() { enc.buf.AppendInt(int64(value)) })
}

func (enc *JsonEncoder) writeFieldInt16(value int16) {
	writeFieldWrapper(enc, enc.colors.attr.NumberColor, func() { enc.buf.AppendInt(int64(value)) })
}

func (enc *JsonEncoder) writeFieldInt32(value int32) {
	writeFieldWrapper(enc, enc.colors.attr.NumberColor, func() { enc.buf.AppendInt(int64(value)) })
}

func (enc *JsonEncoder) writeFieldInt64(value int64) {
	writeFieldWrapper(enc, enc.colors.attr.NumberColor, func() { enc.buf.AppendInt(value) })
}

func (enc *JsonEncoder) writeFieldInt(value int) {
	writeFieldWrapper(enc, enc.colors.attr.NumberColor, func() { enc.buf.AppendInt(int64(value)) })
}

func (enc *JsonEncoder) writeFieldUint8(value uint8) {
	writeFieldWrapper(enc, enc.colors.attr.NumberColor, func() { enc.buf.AppendUint(uint64(value)) })
}

func (enc *JsonEncoder) writeFieldUint16(value uint16) {
	writeFieldWrapper(enc, enc.colors.attr.NumberColor, func() { enc.buf.AppendUint(uint64(value)) })
}

func (enc *JsonEncoder) writeFieldUint32(value uint32) {
	writeFieldWrapper(enc, enc.colors.attr.NumberColor, func() { enc.buf.AppendUint(uint64(value)) })
}

func (enc *JsonEncoder) writeFieldUint64(value uint64) {
	writeFieldWrapper(enc, enc.colors.attr.NumberColor, func() { enc.buf.AppendUint(value) })
}

func (enc *JsonEncoder) writeFieldUint(value uint) {
	writeFieldWrapper(enc, enc.colors.attr.NumberColor, func() { enc.buf.AppendUint(uint64(value)) })
}

func (enc *JsonEncoder) writeFieldFloat32(value float32) {
	writeFieldWrapper(enc, enc.colors.attr.NumberColor, func() { enc.buf.AppendFloat(float64(value), 32) })
}

func (enc *JsonEncoder) writeFieldFloat64(value float64) {
	writeFieldWrapper(enc, enc.colors.attr.NumberColor, func() { enc.buf.AppendFloat(value, 64) })
}

func (enc *JsonEncoder) writeFieldTime(value time.Time) {
	enc.writeFieldString(value.Format(time.DateTime))
}

func (enc *JsonEncoder) writeFieldDuration(value time.Duration) {
	enc.writeFieldString(value.String())
}

func (enc *JsonEncoder) writeFieldError(value error) {
	if value == nil {
		enc.writeFieldNil()
	} else {
		enc.writeFieldString(value.Error())
	}
}

func (enc *JsonEncoder) writeMapObjectForAnyValue(value map[string]any) {
	enc.writeBeginObject()
	// the key-values are unsorted!!!
	i, n := 0, len(value)
	for k, v := range value {
		enc.writeFieldKey(k)
		enc.writeSplitColon()
		enc.writeFieldAny(v)
		if i+1 != n {
			enc.writeSplitComma()
		}
		i++
	}
	enc.writeEndObject()
}

func (enc *JsonEncoder) writeMapObjectForStringValue(value map[string]string) {
	enc.writeBeginObject()
	// the key-values are unsorted!!!
	i, n := 0, len(value)
	for k, v := range value {
		enc.writeFieldKey(k)
		enc.writeSplitColon()
		enc.writeFieldString(v)
		if i+1 != n {
			enc.writeSplitComma()
		}
		i++
	}
	enc.writeEndObject()
}

func (enc *JsonEncoder) writeMapObjectForMultipleValue(value map[string][]string) {
	enc.writeBeginObject()
	// the key-values are unsorted!!!
	i, n := 0, len(value)
	for k, v := range value {
		enc.writeFieldKey(k)
		enc.writeSplitColon()
		enc.writeFieldArray(v)
		if i+1 != n {
			enc.writeSplitComma()
		}
		i++
	}
	enc.writeEndObject()
}

func (enc *JsonEncoder) writeFieldObject(value []Field) {
	enc.writeBeginObject()
	n := len(value)
	for i := 0; i < n; i++ {
		if err := enc.writeField(&value[i]); err != nil {
			continue
		}
		if i+1 != n {
			enc.writeSplitComma()
		}
	}
	enc.writeEndObject()
}

func (enc *JsonEncoder) writeFieldSingleObject(value Field) {
	enc.writeBeginObject()
	enc.writeField(&value)
	enc.writeEndObject()
}

func (enc *JsonEncoder) writeFieldArray(value any) {
	enc.writeFieldAny(value)
}

func (enc *JsonEncoder) writeFieldNil() {
	writeFieldWrapper(enc, enc.colors.attr.StringColor, func() { enc.buf.AppendString("null") })
}

func writeFieldArrayListFor[T any](value []T, enc *JsonEncoder, wf func(T), lf func()) {
	enc.writeBeginArray()
	for i := 0; i < len(value); i++ {
		wf(value[i])
		if i+1 != len(value) {
			lf()
		}
	}
	enc.writeEndArray()
}

func writeFieldArrayListForReflectValue(value reflect.Value, enc *JsonEncoder, wf func(any), lf func()) {
	enc.writeBeginArray()
	n := value.Len()
	for i := 0; i < n; i++ {
		idxv := value.Index(i)
		if !idxv.CanInterface() {
			continue
		}
		wf(idxv.Interface())
		if i+1 != n {
			lf()
		}
	}
	enc.writeEndArray()
}

func (enc *JsonEncoder) writeFieldAny(value any) {
	switch v := value.(type) {
	case []string:
		writeFieldArrayListFor(v, enc, func(s string) { enc.writeFieldString(s) }, enc.writeSplitComma)
	case []bool:
		writeFieldArrayListFor(v, enc, func(s bool) { enc.writeFieldBool(s) }, enc.writeSplitComma)
	case []int8:
		writeFieldArrayListFor(v, enc, func(s int8) { enc.writeFieldInt8(s) }, enc.writeSplitComma)
	case []int16:
		writeFieldArrayListFor(v, enc, func(s int16) { enc.writeFieldInt16(s) }, enc.writeSplitComma)
	case []int32:
		writeFieldArrayListFor(v, enc, func(s int32) { enc.writeFieldInt32(s) }, enc.writeSplitComma)
	case []int64:
		writeFieldArrayListFor(v, enc, func(s int64) { enc.writeFieldInt64(s) }, enc.writeSplitComma)
	case []int:
		writeFieldArrayListFor(v, enc, func(s int) { enc.writeFieldInt(s) }, enc.writeSplitComma)
	case []uint8:
		writeFieldArrayListFor(v, enc, func(s uint8) { enc.writeFieldUint8(s) }, enc.writeSplitComma)
	case []uint16:
		writeFieldArrayListFor(v, enc, func(s uint16) { enc.writeFieldUint16(s) }, enc.writeSplitComma)
	case []uint32:
		writeFieldArrayListFor(v, enc, func(s uint32) { enc.writeFieldUint32(s) }, enc.writeSplitComma)
	case []uint64:
		writeFieldArrayListFor(v, enc, func(s uint64) { enc.writeFieldUint64(s) }, enc.writeSplitComma)
	case []uint:
		writeFieldArrayListFor(v, enc, func(s uint) { enc.writeFieldUint(s) }, enc.writeSplitComma)
	case []float32:
		writeFieldArrayListFor(v, enc, func(s float32) { enc.writeFieldFloat32(s) }, enc.writeSplitComma)
	case []float64:
		writeFieldArrayListFor(v, enc, func(s float64) { enc.writeFieldFloat64(s) }, enc.writeSplitComma)
	case []time.Time:
		writeFieldArrayListFor(v, enc, func(s time.Time) { enc.writeFieldTime(s) }, enc.writeSplitComma)
	case []time.Duration:
		writeFieldArrayListFor(v, enc, func(s time.Duration) { enc.writeFieldDuration(s) }, enc.writeSplitComma)
	case []error:
		writeFieldArrayListFor(v, enc, func(s error) { enc.writeFieldError(s) }, enc.writeSplitComma)
	case []map[string]any:
		writeFieldArrayListFor(v, enc, func(s map[string]any) { enc.writeMapObjectForAnyValue(s) }, enc.writeSplitComma)
	case []map[string]string:
		writeFieldArrayListFor(v, enc, func(s map[string]string) { enc.writeMapObjectForStringValue(s) }, enc.writeSplitComma)
	case []map[string][]string:
		writeFieldArrayListFor(v, enc, func(s map[string][]string) { enc.writeMapObjectForMultipleValue(s) }, enc.writeSplitComma)
	case []Field:
		enc.writeFieldObject(v)
	case []any:
		writeFieldArrayListFor(v, enc, func(s any) { enc.writeFieldAny(s) }, enc.writeSplitComma)
	case string:
		enc.writeFieldString(v)
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
	case map[string]any:
		enc.writeMapObjectForAnyValue(v)
	case map[string]string:
		enc.writeMapObjectForStringValue(v)
	case map[string][]string:
		enc.writeMapObjectForMultipleValue(v)
	case error:
		enc.writeFieldError(v)
	case Field:
		enc.writeFieldSingleObject(v)
	case nil:
		enc.writeFieldNil()
	default:
		if enc.reflectValue && reflect.TypeOf(value).Kind() == reflect.Slice {
			valueOf := reflect.ValueOf(value)
			writeFieldArrayListForReflectValue(valueOf, enc, enc.writeFieldAny, enc.writeSplitComma)
		} else {
			enc.writeFieldString(fmt.Sprintf("%v", value))
		}
	}
}

func appendQuoteString(buf *Buffer, value string) {
	if len(value) == 0 {
		return
	}
	buf.TryGrow(3 * len(value) / 2)
	buf.AppendQuote(value)
}
