package themes

import (
	"image/color"

	"gioui.org/font/gofont"
	"gioui.org/text"
	"gioui.org/widget/material"
)

var (
	ColorBackground     = color.NRGBA{R: 0x1A, G: 0x14, B: 0x10, A: 0xFF}
	ColorPanel          = color.NRGBA{R: 0x22, G: 0x1E, B: 0x15, A: 0xFF}
	ColorInput          = color.NRGBA{R: 0x2C, G: 0x26, B: 0x19, A: 0xFF}
	ColorBorder         = color.NRGBA{R: 0x5A, G: 0x4A, B: 0x28, A: 0xFF}
	ColorGold           = color.NRGBA{R: 0xC9, G: 0xA8, B: 0x4C, A: 0xFF}
	ColorGoldDim        = color.NRGBA{R: 0x8A, G: 0x6E, B: 0x30, A: 0xFF}
	ColorGoldBright     = color.NRGBA{R: 0xE0, G: 0xC0, B: 0x60, A: 0xFF}
	ColorText           = color.NRGBA{R: 0xE8, G: 0xD5, B: 0xA3, A: 0xFF}
	ColorTextDim        = color.NRGBA{R: 0x9A, G: 0x8A, B: 0x6A, A: 0xFF}
	ColorError          = color.NRGBA{R: 0xFF, G: 0x70, B: 0x70, A: 0xFF}
	ColorWarnBackground = color.NRGBA{R: 0x2E, G: 0x24, B: 0x10, A: 0xFF}
	ColorWarnBorder     = color.NRGBA{R: 0x6A, G: 0x50, B: 0x20, A: 0xFF}
	ColorWarnText       = color.NRGBA{R: 0xD4, G: 0xA8, B: 0x43, A: 0xFF}
	ColorGenerate       = color.NRGBA{R: 0x7A, G: 0x5A, B: 0x18, A: 0xFF}
)

// NewTheme builds a material.Theme
func NewTheme() *material.Theme {
	theme := material.NewTheme()
	theme.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	theme.Palette = material.Palette{
		Bg:         ColorBackground,
		Fg:         ColorText,
		ContrastBg: ColorGenerate,
		ContrastFg: ColorText,
	}
	return theme
}
