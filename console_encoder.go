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

func (enc *ConsoleEncoder) Encode(buf *Buffer, msg string, fields []Field) error {
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

	n1 := len(fields)
	n2 := len(enc.preFields)
	if n1 == 0 && n2 == 0 {
		return nil
	}
	buf.WriteByte(ConsoleEncoderSplitCharacter)

	enc.jsonEncoder.buf = buf
	enc.jsonEncoder.writeBeginObject()
	enc.jsonEncoder.writePrefixFields()
	if n1 == 0 {
		enc.jsonEncoder.writeEndObject()
		return nil
	}
	if n2 > 0 {
		enc.jsonEncoder.writeSplitComma()
	}
	for i := 0; i < n1; i++ {
		if err := enc.jsonEncoder.writeField(&fields[i]); err != nil {
			return err
		}
		if i+1 != n1 {
			enc.jsonEncoder.writeSplitComma()
		}
	}
	enc.jsonEncoder.writeEndObject()
	return nil
}
