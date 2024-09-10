package logx

type EncoderType byte

const (
	Console EncoderType = 1 << iota
	Json
)

type encoder interface {
	Init()
	Encode(*Buffer, string, []Field) error
}
