package widgets

import (
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/themes"
)

// NewSectionHeaderWidget returns a Widget that renders a section title
func NewSectionHeaderWidget(theme *material.Theme, title string) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		label := material.Body1(theme, "◆  "+title)
		label.Color = themes.ColorAccent
		label.Font = font.Font{Weight: font.SemiBold}
		label.TextSize = unit.Sp(13)
		return label.Layout(gtx)
	}
}
