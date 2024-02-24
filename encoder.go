package logx

import "bytes"

type EncoderType byte

const (
	// Simple Encode does't support key-value args, only message field can be set
	Simple EncoderType = 1 << iota
	Json
)

type encoder interface {
	Init()
	Encode(*bytes.Buffer, string, ...arg) error
}
