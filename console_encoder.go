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
		enc.timeF.AppendTimePrimitive(buf, ent.time)
		buf.AppendByte(ConsoleEncoderSplitCharacter)
	}
	if enc.levelF.enable {
		enc.levelF.AppendPrimitive(buf, ent.level)
		buf.AppendByte(ConsoleEncoderSplitCharacter)
	}
	if enc.callerF.enable {
		enc.callerF.AppendPrimitive(buf)
		buf.AppendByte(ConsoleEncoderSplitCharacter)
	}

	if enc.escapeQuote {
		buf.bs = appendEscapedControl(buf.bs, ent.message)
	} else {
		buf.AppendString(ent.message)
	}

	n1 := len(fields)
	n2 := len(enc.preFields)
	if n1 == 0 && n2 == 0 {
		return buf, nil
	}
	buf.AppendByte(ConsoleEncoderSplitCharacter)

	jsonEnc.writeBeginObject()
	wrotePrefix, prefixErr := jsonEnc.writePrefixFields()
	if prefixErr != nil {
		putBuffer(buf)
		return nil, prefixErr
	}
	if n1 == 0 {
		jsonEnc.writeEndObject()
		ret = buf
		return
	}
	if wrotePrefix {
		jsonEnc.writeSplitComma()
	}
	for i := 0; i < n1; i++ {
		if err = jsonEnc.writeField(&fields[i]); err != nil {
			putBuffer(buf)
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
