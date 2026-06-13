package widgets

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/themes"
)

// NewDimmedLabelWidget returns a Widget that renders a small dimmed description line
func NewDimmedLabelWidget(theme *material.Theme, description string) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		label := material.Body2(theme, description)
		label.Color = themes.ColorTextDim
		label.TextSize = unit.Sp(12)
		return label.Layout(gtx)
	}
}
