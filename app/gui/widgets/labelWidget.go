package widgets

import (
	"image/color"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget/material"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/themes"
)

// NewLabelWidget returns a Widget that renders a small text line
func NewLabelWidget(theme *material.Theme, text string, color color.NRGBA) layout.Widget {
	label := material.Caption(theme, text)
	label.Color = color
	return label.Layout
}

// NewStyledLabelWidget returns a Widget that renders a small text line with a specific style
func NewStyledLabelWidget(theme *material.Theme, text string, color color.NRGBA, style font.Font) layout.Widget {
	label := material.Caption(theme, text)
	label.Color = color
	label.Font = style
	return label.Layout
}

// NewAlignedLabelWidget returns a Widget that renders a small text line with a specific alignment
func NewAlignedLabelWidget(theme *material.Theme, text string, color color.NRGBA, alignment text.Alignment) layout.Widget {
	label := material.Caption(theme, text)
	label.Color = color
	label.Alignment = alignment
	return label.Layout
}

// NewDimmedLabelWidget returns a Widget that renders a small dimmed text line
func NewDimmedLabelWidget(theme *material.Theme, text string) layout.Widget {
	label := material.Caption(theme, text)
	label.Color = themes.ColorTextDim
	return label.Layout
}
