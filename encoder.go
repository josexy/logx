package logx

import "time"

type EncoderType byte

const (
	Console EncoderType = 1 << iota
	Json
)

type entry struct {
	level   LevelType
	time    time.Time
	message string
}

type encoder interface {
	Init()
	Encode(entry, []Field) (*Buffer, error)
}
