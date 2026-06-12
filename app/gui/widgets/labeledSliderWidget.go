package widgets

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/themes"
)

// NewLabeledSliderWidget returns a Widget that renders a slider with a label on the right
func NewLabeledSliderWidget(theme *material.Theme, slider *widget.Float, value string) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				slider := material.Slider(theme, slider)
				slider.Color = themes.ColorAccent
				return slider.Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Width: unit.Dp(6)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Min.X = gtx.Dp(32)
				label := material.Body1(theme, value)
				label.Color = themes.ColorAccent
				label.TextSize = unit.Sp(13)
				label.Alignment = 1
				return label.Layout(gtx)
			}),
		)
	}
}
