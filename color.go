package logx

import (
	"os"

	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
)

type ColorAttr uint8

const (
	format   = "\x1b["
	unformat = "\x1b[0m"
	bold     = 1
)

const (
	BlackAttr ColorAttr = iota + 30
	RedAttr
	GreenAttr
	YellowAttr
	BlueAttr
	MagentaAttr
	CyanAttr
	WhiteAttr
)

const (
	HiBlackAttr ColorAttr = iota + 90
	HiRedAttr
	HiGreenAttr
	HiYellowAttr
	HiBlueAttr
	HiMagentaAttr
	HiCyanAttr
	HiWhiteAttr
)

var (
	NoColor = os.Getenv("NO_COLOR") != "" || os.Getenv("TERM") == "dumb" ||
		(!isatty.IsTerminal(os.Stdout.Fd()) && !isatty.IsCygwinTerminal(os.Stdout.Fd()))

	Output = colorable.NewColorableStdout()
)

type TextColorAttri struct {
	KeyColor     ColorAttr
	StringColor  ColorAttr
	BooleanColor ColorAttr
	FloatColor   ColorAttr
	NumberColor  ColorAttr
}

func (c ColorAttr) appendTo(buf *Buffer) {
	if c <= 0 {
		return
	}
	buf.AppendInt(int64(c))
	if c >= HiBlackAttr && c <= HiWhiteAttr {
		buf.AppendByte(';')
		buf.AppendInt(int64(bold))
	}
}

type colorfulset struct {
	enable bool
	attr   TextColorAttri
}

func (c *colorfulset) init() {
	if !c.enable {
		return
	}
	if c.attr.KeyColor == 0 || c.attr.KeyColor > HiWhiteAttr {
		c.attr.KeyColor = BlueAttr
	}
	if c.attr.StringColor == 0 || c.attr.StringColor > HiWhiteAttr {
		c.attr.StringColor = GreenAttr
	}
	if c.attr.BooleanColor == 0 || c.attr.BooleanColor > HiWhiteAttr {
		c.attr.BooleanColor = YellowAttr
	}
	if c.attr.FloatColor == 0 || c.attr.FloatColor > HiWhiteAttr {
		c.attr.FloatColor = CyanAttr
	}
	if c.attr.NumberColor == 0 || c.attr.NumberColor > HiWhiteAttr {
		c.attr.NumberColor = RedAttr
	}
}

func appendColor(buf *Buffer, color ColorAttr, s string) {
	buf.AppendString(format)
	color.appendTo(buf)
	buf.AppendByte('m')
	buf.AppendString(s)
	buf.AppendString(unformat)
}

func appendColorWithFunc(buf *Buffer, color ColorAttr, append func(*Buffer)) {
	buf.AppendString(format)
	color.appendTo(buf)
	buf.AppendByte('m')
	append(buf)
	buf.AppendString(unformat)
}
