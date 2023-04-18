package logx

import (
	"time"

	"github.com/josexy/logx/internal"
)

type ConfigOption interface{ applyTo(*internal.Config) }
type configOptionFn func(*internal.Config)

func (fn configOptionFn) applyTo(c *internal.Config) { fn(c) }

func WithColor(enable bool) ConfigOption {
	return configOptionFn(func(c *internal.Config) {
		c.LevelItem.Color = enable
		c.TimeItem.Color = enable
		c.CallerItem.Color = enable
	})
}

func WithLevel(enable, lower bool) ConfigOption {
	return configOptionFn(func(c *internal.Config) {
		c.LevelItem.Enable = enable
		c.LevelItem.Lower = lower
	})
}

func WithTime(enable bool, format func(time.Time) string) ConfigOption {
	return configOptionFn(func(c *internal.Config) {
		if format == nil {
			format = func(t time.Time) string {
				return t.Format(time.DateTime)
			}
		}
		c.TimeItem.Enable = enable
		c.TimeItem.Format = format
	})
}

func WithCaller(enable, fileName, funcName, lineNumber bool) ConfigOption {
	return configOptionFn(func(c *internal.Config) {
		c.CallerItem.Enable = enable
		c.CallerItem.File = fileName
		c.CallerItem.Func = funcName
		c.CallerItem.Line = lineNumber
	})
}

func WithEscapeQuote(enable bool) ConfigOption {
	return configOptionFn(func(c *internal.Config) {
		c.EscapeQuote = enable
	})
}

func WithJsonEncoder() ConfigOption {
	return configOptionFn(func(c *internal.Config) {
		c.Encoder = &internal.JsonEncoder{Config: c}
	})
}

func WithSimpleEncoder() ConfigOption {
	return configOptionFn(func(c *internal.Config) {
		c.Encoder = &internal.SimpleEncoder{Config: c}
	})
}
