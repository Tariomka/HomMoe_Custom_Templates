package widgets

import (
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/themes"
)

// NewLabeledSliderWidget returns a Widget that renders a slider with a label on the right.
func NewLabeledSliderWidget(theme *material.Theme, slider *widget.Float, value string) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				slider := material.Slider(theme, slider)
				slider.Color = themes.ColorAccent
				return slider.Layout(gtx)
			}),
			layout.Rigid(NewHorizontalSpacerWidget(6)),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Min.X = gtx.Dp(32)
				return NewLabelBuilder(theme).WithSizeDefault().
					WithText(value).WithColor(themes.ColorAccent).WithAlignment(text.End).Build(gtx)
			}),
		)
	}
}
