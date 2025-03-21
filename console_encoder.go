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
		enc.callerF.skipDepth += enc.callerF.option.CallerSkip
	}
}

func (enc *ConsoleEncoder) Encode(ent entry, fields []Field) (ret *Buffer, err error) {
	jsonEnc := enc.jsonEncoder.clone()
	defer putJsonEncoder(jsonEnc)

	buf := jsonEnc.buf

	if enc.timeF.enable {
		buf.WriteString(enc.timeF.String(&ent))
		buf.WriteByte(ConsoleEncoderSplitCharacter)
	}
	if enc.levelF.enable {
		buf.WriteString(enc.levelF.String(&ent))
		buf.WriteByte(ConsoleEncoderSplitCharacter)
	}
	if enc.callerF.enable {
		buf.WriteString(enc.callerF.String())
		buf.WriteByte(ConsoleEncoderSplitCharacter)
	}

	buf.WriteString(ent.message)

	n1 := len(fields)
	n2 := len(enc.preFields)
	if n1 == 0 && n2 == 0 {
		return buf, nil
	}
	buf.WriteByte(ConsoleEncoderSplitCharacter)

	jsonEnc.writeBeginObject()
	jsonEnc.writePrefixFields()
	if n1 == 0 {
		jsonEnc.writeEndObject()
		ret = buf
		return
	}
	if n2 > 0 {
		jsonEnc.writeSplitComma()
	}
	for i := 0; i < n1; i++ {
		if err = jsonEnc.writeField(&fields[i]); err != nil {
			bufPool.Put(buf)
			return
		}
		if i+1 != n1 {
			jsonEnc.writeSplitComma()
		}
	}
	jsonEnc.writeEndObject()
	ret = buf
	return
}
