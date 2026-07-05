package widgets

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/themes"
)

// NewLabeledCheckboxRowWidget returns a Widget that renders a checkbox with a label as a clickable row.
func NewLabeledCheckboxRowWidget(theme *material.Theme, boolValue *widget.Bool, label string) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{Top: unit.Dp(3), Bottom: unit.Dp(3)}.
			Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				checkbox := material.CheckBox(theme, boolValue, label)
				checkbox.Color = themes.ColorText
				checkbox.IconColor = themes.ColorAccent
				checkbox.TextSize = unit.Sp(12)
				return checkbox.Layout(gtx)
			})
	}
}
