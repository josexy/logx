package logx

import "unicode/utf8"

const lowerhex = "0123456789abcdef"

// Modified from strconv/quote.go
func appendQuotedWith(buf []byte, s string) []byte {
	if cap(buf)-len(buf) < len(s) {
		nBuf := make([]byte, len(buf), len(buf)+len(s))
		copy(nBuf, buf)
		buf = nBuf
	}
	start := 0
	for i := 0; i < len(s); {
		if b := s[i]; b < utf8.RuneSelf {
			if b >= 0x20 && b != '\\' && b != '"' {
				i++
				continue
			}

			buf = append(buf, s[start:i]...)
			buf = appendEscapedASCII(buf, b)
			i++
			start = i
			continue
		}

		r, size := utf8.DecodeRuneInString(s[i:])
		if r == utf8.RuneError && size == 1 {
			buf = append(buf, s[start:i]...)
			buf = append(buf, `\ufffd`...)
			i++
			start = i
			continue
		}
		if r == '\u2028' || r == '\u2029' {
			buf = append(buf, s[start:i]...)
			buf = append(buf, `\u202`...)
			buf = append(buf, lowerhex[r&0x0f])
			i += size
			start = i
			continue
		}
		i += size
	}
	return append(buf, s[start:]...)
}

// appendEscapedControl appends s with ASCII control characters and the
// U+2028/U+2029 line separators escaped JSON-style. Everything else, including
// quotes and backslashes, is appended as-is. The Console encoder uses it so a
// message cannot inject line breaks or column separators.
func appendEscapedControl(buf []byte, s string) []byte {
	start := 0
	for i := 0; i < len(s); {
		b := s[i]
		if b < utf8.RuneSelf {
			if b >= 0x20 {
				i++
				continue
			}
			buf = append(buf, s[start:i]...)
			buf = appendEscapedASCII(buf, b)
			i++
			start = i
			continue
		}
		r, size := utf8.DecodeRuneInString(s[i:])
		if r == '\u2028' || r == '\u2029' {
			buf = append(buf, s[start:i]...)
			buf = append(buf, `\u202`...)
			buf = append(buf, lowerhex[r&0x0f])
			i += size
			start = i
			continue
		}
		i += size
	}
	return append(buf, s[start:]...)
}

// appendEscapedASCII appends the JSON escape sequence for b, which must be a
// control character, a backslash or a double quote.
func appendEscapedASCII(buf []byte, b byte) []byte {
	switch b {
	case '\\', '"':
		return append(buf, '\\', b)
	case '\b':
		return append(buf, `\b`...)
	case '\f':
		return append(buf, `\f`...)
	case '\n':
		return append(buf, `\n`...)
	case '\r':
		return append(buf, `\r`...)
	case '\t':
		return append(buf, `\t`...)
	default:
		buf = append(buf, `\u00`...)
		return append(buf, lowerhex[b>>4], lowerhex[b&0x0f])
	}
}

func jsonBytesNeedEscaping(value []byte) bool {
	for i := 0; i < len(value); {
		b := value[i]
		if b < utf8.RuneSelf {
			if b < 0x20 || b == '\\' || b == '"' {
				return true
			}
			i++
			continue
		}
		r, size := utf8.DecodeRune(value[i:])
		if (r == utf8.RuneError && size == 1) || r == '\u2028' || r == '\u2029' {
			return true
		}
		i += size
	}
	return false
}
