package logx

var ConsoleEncoderSplitCharacter = byte('\t')

type ConsoleEncoder struct {
	*LogContext
	jsonEncoder *JsonEncoder
}

func (enc *ConsoleEncoder) Init() {
	enc.jsonEncoder = &JsonEncoder{
		LogContext: enc.LogContext,
	}
	enc.jsonEncoder.Init()
	if enc.callerF.enable {
		enc.callerF.skipDepth = 6
	}
}

func (enc *ConsoleEncoder) Encode(buf *Buffer, msg string, fields ...Field) error {
	if enc.timeF.enable {
		buf.WriteString(enc.timeF.String())
		buf.WriteByte(ConsoleEncoderSplitCharacter)
	}
	if enc.levelF.enable {
		buf.WriteString(enc.levelF.String())
		buf.WriteByte(ConsoleEncoderSplitCharacter)
	}
	if enc.callerF.enable {
		buf.WriteString(enc.callerF.String())
		buf.WriteByte(ConsoleEncoderSplitCharacter)
	}

	buf.WriteString(msg)

	enc.jsonEncoder.buf = buf
	enc.jsonEncoder.fieldsRanger.reset()
	enc.jsonEncoder.addPrefixFields()
	enc.jsonEncoder.fieldsRanger.put(fields...)
	if enc.jsonEncoder.fieldsRanger.size() == 0 {
		return nil
	}

	buf.WriteByte(ConsoleEncoderSplitCharacter)
	enc.jsonEncoder.writeBeginObject()
	err := enc.jsonEncoder.fieldsRanger.writeRangeFields(enc.jsonEncoder.writeField, enc.jsonEncoder.writeSplitComma)
	if err != nil {
		return err
	}
	enc.jsonEncoder.writeEndObject()
	return nil
}
