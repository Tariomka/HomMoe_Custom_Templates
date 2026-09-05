package themes

import (
	"gioui.org/font/gofont"
	"gioui.org/text"
	"gioui.org/widget/material"
)

// NewTheme builds the application's material.Theme from the Crimson Night palette.
func NewTheme() *material.Theme {
	theme := material.NewTheme()
	theme.Shaper = text.NewShaper(text.NoSystemFonts(), text.WithCollection(gofont.Collection()))
	theme.Palette = material.Palette{
		Bg:         ColorsBase.Background,
		Fg:         ColorsBase.Text,
		ContrastBg: ColorsBase.PrimaryButton,
		ContrastFg: ColorsBase.Text,
	}
	return theme
}
