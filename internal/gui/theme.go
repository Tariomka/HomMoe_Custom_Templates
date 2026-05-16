package gui

import (
	"image/color"

	"gioui.org/font/gofont"
	"gioui.org/text"
	"gioui.org/widget/material"
)

// Medieval palette — ported from Themes/MedievalTheme.xaml.
var (
	colBackground = color.NRGBA{R: 0x1A, G: 0x14, B: 0x10, A: 0xFF}
	colPanel      = color.NRGBA{R: 0x22, G: 0x1E, B: 0x15, A: 0xFF}
	colInput      = color.NRGBA{R: 0x2C, G: 0x26, B: 0x19, A: 0xFF}
	colBorder     = color.NRGBA{R: 0x5A, G: 0x4A, B: 0x28, A: 0xFF}
	colGold       = color.NRGBA{R: 0xC9, G: 0xA8, B: 0x4C, A: 0xFF}
	colGoldDim    = color.NRGBA{R: 0x8A, G: 0x6E, B: 0x30, A: 0xFF}
	colGoldBright = color.NRGBA{R: 0xE0, G: 0xC0, B: 0x60, A: 0xFF}
	colText       = color.NRGBA{R: 0xE8, G: 0xD5, B: 0xA3, A: 0xFF}
	colTextDim    = color.NRGBA{R: 0x9A, G: 0x8A, B: 0x6A, A: 0xFF}
	colError      = color.NRGBA{R: 0xFF, G: 0x70, B: 0x70, A: 0xFF}
	colGenerate   = color.NRGBA{R: 0x7A, G: 0x5A, B: 0x18, A: 0xFF}
)

// NewTheme builds a material.Theme tuned to the medieval palette.
func NewTheme() *material.Theme {
	theme := material.NewTheme()
	theme.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	theme.Palette = material.Palette{
		Bg:         colBackground,
		Fg:         colText,
		ContrastBg: colGenerate,
		ContrastFg: colText,
	}
	return theme
}
