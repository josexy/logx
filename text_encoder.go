package logx

import (
	"fmt"
	"math"
	"strconv"
	"sync"
	"time"
	"unicode"
	"unicode/utf8"
)

var textPool = sync.Pool{New: func() any { return &TextEncoder{} }}

type TextEncoder struct {
	*LogContext
	buf          *Buffer
	disableColor bool
}

func (enc *TextEncoder) Init() {
	enc.LogContext.WithMsgKey(enc.msgKey)
	enc.colors.init()
	if enc.callerF.enable {
		enc.callerF.skipDepth = 6
		enc.callerF.skipDepth += enc.callerF.option.CallerSkip
	}
}

func (enc *TextEncoder) clone() *TextEncoder {
	clone := textPool.Get().(*TextEncoder)
	clone.LogContext = enc.LogContext
	clone.buf = bufPool.Get().(*Buffer)
	clone.buf.Reset()
	return clone
}

func putTextEncoder(enc *TextEncoder) {
	enc.LogContext = nil
	enc.buf = nil
	enc.disableColor = false
	textPool.Put(enc)
}

func (enc *TextEncoder) Encode(ent entry, fields []Field) (ret *Buffer, err error) {
	nenc := enc.clone()
	defer putTextEncoder(nenc)

	if nenc.timeF.enable {
		nenc.writeFieldKey(nenc.timeF.option.TimeKey)
		nenc.writeTimePrimitive(ent.time)
	}
	if nenc.levelF.enable {
		nenc.writeFieldKey(nenc.levelF.option.LevelKey)
		nenc.writeLevelPrimitive(ent.level)
	}
	if nenc.callerF.enable {
		nenc.writeCaller()
	}

	nenc.writeFieldKey(nenc.msgKey)
	nenc.writeTextString(ent.message)
	if err = nenc.writeFields("", nenc.preFields); err != nil {
		bufPool.Put(nenc.buf)
		return
	}
	if err = nenc.writeFields("", fields); err != nil {
		bufPool.Put(nenc.buf)
		return
	}
	ret = nenc.buf
	return
}

func (enc *TextEncoder) writeFields(prefix string, fields []Field) error {
	for i := range fields {
		if err := enc.writeField(prefix, &fields[i]); err != nil {
			return err
		}
	}
	return nil
}

func (enc *TextEncoder) writeField(prefix string, field *Field) error {
	if field.Type == NoneType {
		return errInvalidFieldType
	}
	key := field.Key
	if prefix != "" {
		key = prefix + "." + key
	}
	if field.Type == ObjectType {
		return enc.writeFields(key, field.AnyValue.([]Field))
	}
	enc.writeFieldKey(key)
	enc.writeFieldValue(field)
	return nil
}

func (enc *TextEncoder) writeCaller() {
	fileName, funcName := enc.callerF.value()
	if fileName != "" {
		enc.writeFieldKey(enc.callerF.option.CallerKey + "." + enc.callerF.option.FileKey)
		enc.writeTextStringWithColor(fileName, enc.colors.attr.StringColor)
	}
	if funcName != "" {
		enc.writeFieldKey(enc.callerF.option.CallerKey + "." + enc.callerF.option.FuncKey)
		enc.writeTextStringWithColor(funcName, enc.colors.attr.StringColor)
	}
}

func (enc *TextEncoder) writeFieldKey(key string) {
	if enc.buf.Len() > 0 {
		enc.buf.AppendByte(' ')
	}
	enc.writeTextStringWithColor(key, enc.colors.attr.KeyColor)
	enc.buf.AppendByte('=')
}

func (enc *TextEncoder) writeFieldValue(field *Field) {
	switch field.Type {
	case StringType:
		enc.writeTextStringWithColor(field.StringValue, enc.colors.attr.StringColor)
	case BoolType:
		enc.writeBool(field.IntValue == 1)
	case Int8Type:
		enc.writeInt(int64(int8(field.IntValue)))
	case Int16Type:
		enc.writeInt(int64(int16(field.IntValue)))
	case Int32Type:
		enc.writeInt(int64(int32(field.IntValue)))
	case Int64Type:
		enc.writeInt(field.IntValue)
	case IntType:
		enc.writeInt(int64(int(field.IntValue)))
	case Uint8Type:
		enc.writeUint(uint64(uint8(field.IntValue)))
	case Uint16Type:
		enc.writeUint(uint64(uint16(field.IntValue)))
	case Uint32Type:
		enc.writeUint(uint64(uint32(field.IntValue)))
	case Uint64Type:
		enc.writeUint(uint64(field.IntValue))
	case UintType:
		enc.writeUint(uint64(uint(field.IntValue)))
	case Float32Type:
		enc.writeFloat(float64(math.Float32frombits(uint32(field.IntValue))), 32)
	case Float64Type:
		enc.writeFloat(math.Float64frombits(uint64(field.IntValue)), 64)
	case TimeType:
		if field.AnyValue != nil {
			enc.writeTimePrimitive(time.Unix(0, field.IntValue).In(field.AnyValue.(*time.Location)))
		} else {
			enc.writeTimePrimitive(time.Unix(0, field.IntValue))
		}
	case TimeFullType:
		enc.writeTimePrimitive(field.AnyValue.(time.Time))
	case DurationType:
		enc.writeTextStringWithColor(time.Duration(field.IntValue).String(), enc.colors.attr.StringColor)
	case ErrorType:
		if field.AnyValue == nil {
			enc.writeTextStringWithColor("<nil>", enc.colors.attr.StringColor)
			return
		}
		enc.writeTextStringWithColor(field.AnyValue.(error).Error(), enc.colors.attr.StringColor)
	case ArrayType:
		enc.writeAnyValue(field.AnyValue)
	case NilType:
		enc.writeTextStringWithColor("<nil>", enc.colors.attr.StringColor)
	case AnyType:
		enc.writeAnyValue(field.AnyValue)
	}
}

func (enc *TextEncoder) writeTimePrimitive(value time.Time) {
	if enc.timeF.option.Timestamp {
		enc.writeInt(value.UnixNano())
		return
	}
	buf := bufPool.Get().(*Buffer)
	buf.Reset()
	buf.AppendTime(value, enc.timeF.option.Layout)
	valueString := string(buf.Bytes())
	bufPool.Put(buf)
	enc.writeTextStringWithColor(valueString, enc.colors.attr.StringColor)
}

func (enc *TextEncoder) writeLevelPrimitive(level LevelType) {
	appendLevel := func(buf *Buffer) {
		if enc.levelF.option.LowerKey {
			buf.AppendString(levelTypeLowerMap[level])
			return
		}
		buf.AppendString(levelTypeUpperMap[level])
	}
	if enc.colorEnabled() {
		appendColorWithFunc(enc.buf, levelTypeColorMap[level], appendLevel)
		return
	}
	appendLevel(enc.buf)
}

func (enc *TextEncoder) writeBool(value bool) {
	if enc.colorEnabled() {
		appendColorWithFunc(enc.buf, enc.colors.attr.BooleanColor, func(buf *Buffer) { buf.AppendBool(value) })
		return
	}
	enc.buf.AppendBool(value)
}

func (enc *TextEncoder) writeInt(value int64) {
	if enc.colorEnabled() {
		appendColorWithFunc(enc.buf, enc.colors.attr.NumberColor, func(buf *Buffer) { buf.AppendInt(value) })
		return
	}
	enc.buf.AppendInt(value)
}

func (enc *TextEncoder) writeUint(value uint64) {
	if enc.colorEnabled() {
		appendColorWithFunc(enc.buf, enc.colors.attr.NumberColor, func(buf *Buffer) { buf.AppendUint(value) })
		return
	}
	enc.buf.AppendUint(value)
}

func (enc *TextEncoder) writeFloat(value float64, bitSize int) {
	if enc.colorEnabled() {
		appendColorWithFunc(enc.buf, enc.colors.attr.FloatColor, func(buf *Buffer) { buf.AppendFloat(value, bitSize) })
		return
	}
	enc.buf.AppendFloat(value, bitSize)
}

func (enc *TextEncoder) writeTextStringWithColor(value string, color ColorAttr) {
	if enc.colorEnabled() {
		appendColorWithFunc(enc.buf, color, func(buf *Buffer) { writeTextString(buf, value) })
		return
	}
	writeTextString(enc.buf, value)
}

func (enc *TextEncoder) colorEnabled() bool {
	return enc.colors.enable && !enc.disableColor
}

func (enc *TextEncoder) writeTextString(value string) {
	enc.writeTextStringWithColor(value, enc.colors.attr.StringColor)
}

func writeTextString(buf *Buffer, value string) {
	if textNeedsQuoting(value) {
		buf.bs = strconv.AppendQuote(buf.bs, value)
		return
	}
	buf.AppendString(value)
}

func textNeedsQuoting(value string) bool {
	if len(value) == 0 {
		return true
	}
	for i := 0; i < len(value); {
		b := value[i]
		if b < utf8.RuneSelf {
			if b == ' ' || b == '=' || b == '"' || b < ' ' || b == 0x7f {
				return true
			}
			i++
			continue
		}
		r, size := utf8.DecodeRuneInString(value[i:])
		if r == utf8.RuneError || unicode.IsSpace(r) || !unicode.IsPrint(r) {
			return true
		}
		i += size
	}
	return false
}

func (enc *TextEncoder) writeAnyValue(value any) {
	switch v := value.(type) {
	case []Field:
		enc.writeTextString(enc.textFieldsString(v))
	case Field:
		enc.writeTextString(enc.textFieldsString([]Field{v}))
	case error:
		if v == nil {
			enc.writeTextString("<nil>")
		} else {
			enc.writeTextString(v.Error())
		}
	case time.Time:
		enc.writeTimePrimitive(v)
	case time.Duration:
		enc.writeTextString(v.String())
	case string:
		enc.writeTextString(v)
	case bool:
		enc.writeBool(v)
	case int8:
		enc.writeInt(int64(v))
	case int16:
		enc.writeInt(int64(v))
	case int32:
		enc.writeInt(int64(v))
	case int64:
		enc.writeInt(v)
	case int:
		enc.writeInt(int64(v))
	case uint8:
		enc.writeUint(uint64(v))
	case uint16:
		enc.writeUint(uint64(v))
	case uint32:
		enc.writeUint(uint64(v))
	case uint64:
		enc.writeUint(v)
	case uint:
		enc.writeUint(uint64(v))
	case float32:
		enc.writeFloat(float64(v), 32)
	case float64:
		enc.writeFloat(v, 64)
	case nil:
		enc.writeTextString("<nil>")
	default:
		enc.writeTextString(fmt.Sprintf("%v", value))
	}
}

func (enc *TextEncoder) textFieldsString(fields []Field) string {
	nenc := enc.clone()
	nenc.disableColor = true
	defer func() {
		bufPool.Put(nenc.buf)
		putTextEncoder(nenc)
	}()
	for i := range fields {
		if fields[i].Type == NoneType {
			continue
		}
		_ = nenc.writeField("", &fields[i])
	}
	if nenc.buf.Len() == 0 {
		return ""
	}
	return string(nenc.buf.Bytes())
}
