package internal

import "bytes"

type Encoder interface {
	Init()
	Encode(*bytes.Buffer, string, ...Arg) error
}
