package logx

import "github.com/fatih/color"

const (
	NoneAttr = iota
	GreenAttr
	YellowAttr
	BlueAttr
	RedAttr
	CyanAttr
	MagentaAttr
	WhiteAttr
	HiGreenAttr
	HiYellowAttr
	HiBlueAttr
	HiRedAttr
	HiCyanAttr
	HiMagentaAttr
	HiWhiteAttr
)

type colorAttri struct {
	keyColor    *color.Color
	stringColor *color.Color
	boolColor   *color.Color
	floatColor  *color.Color
	numberColor *color.Color
}

type colorfulset struct {
	enable bool
	TextColorAttri
	colorAttri
}

func (c *colorfulset) init() {
	if !c.enable {
		return
	}
	if c.KeyColor == 0 {
		c.KeyColor = BlueAttr
	}
	if c.StringColor == 0 {
		c.StringColor = GreenAttr
	}
	if c.BooleanColor == 0 {
		c.BooleanColor = YellowAttr
	}
	if c.FloatColor == 0 {
		c.FloatColor = CyanAttr
	}
	if c.NumberColor == 0 {
		c.NumberColor = RedAttr
	}
	c.keyColor = colorMap[c.KeyColor]
	c.stringColor = colorMap[c.StringColor]
	c.boolColor = colorMap[c.BooleanColor]
	c.floatColor = colorMap[c.FloatColor]
	c.numberColor = colorMap[c.NumberColor]
}

type TextColorAttri struct {
	KeyColor     uint8
	StringColor  uint8
	BooleanColor uint8
	FloatColor   uint8
	NumberColor  uint8
}

var (
	colorMap = map[uint8]*color.Color{
		GreenAttr:     color.New(color.FgGreen),
		YellowAttr:    color.New(color.FgYellow),
		BlueAttr:      color.New(color.FgBlue),
		RedAttr:       color.New(color.FgRed),
		CyanAttr:      color.New(color.FgCyan),
		MagentaAttr:   color.New(color.FgMagenta),
		WhiteAttr:     color.New(color.FgWhite),
		HiGreenAttr:   color.New(color.FgHiGreen, color.Bold),
		HiYellowAttr:  color.New(color.FgHiYellow, color.Bold),
		HiBlueAttr:    color.New(color.FgHiBlue, color.Bold),
		HiRedAttr:     color.New(color.FgHiRed, color.Bold),
		HiCyanAttr:    color.New(color.FgHiCyan, color.Bold),
		HiMagentaAttr: color.New(color.FgHiMagenta, color.Bold),
		HiWhiteAttr:   color.New(color.FgHiWhite, color.Bold),
	}
)

func Green(msg string) string {
	return colorMap[GreenAttr].Sprint(msg)
}

func Yellow(msg string) string {
	return colorMap[YellowAttr].Sprint(msg)
}

func Blue(msg string) string {
	return colorMap[BlueAttr].Sprint(msg)
}

func Red(msg string) string {
	return colorMap[RedAttr].Sprint(msg)
}

func Cyan(msg string) string {
	return colorMap[CyanAttr].Sprint(msg)
}

func Magenta(msg string) string {
	return colorMap[MagentaAttr].Sprint(msg)
}

func White(msg string) string {
	return colorMap[WhiteAttr].Sprint(msg)
}

func HiGreen(msg string) string {
	return colorMap[HiGreenAttr].Sprint(msg)
}

func HiYellow(msg string) string {
	return colorMap[HiYellowAttr].Sprint(msg)
}

func HiBlue(msg string) string {
	return colorMap[HiBlueAttr].Sprint(msg)
}

func HiRed(msg string) string {
	return colorMap[HiRedAttr].Sprint(msg)
}

func HiCyan(msg string) string {
	return colorMap[HiCyanAttr].Sprint(msg)
}

func HiMagenta(msg string) string {
	return colorMap[HiMagentaAttr].Sprint(msg)
}

func HiWhite(msg string) string {
	return colorMap[HiWhiteAttr].Sprint(msg)
}
