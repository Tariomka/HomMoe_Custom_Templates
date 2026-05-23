package widgets

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/Tariomka/hommoe_custom_templates/internal/gui/components/themes"
)

// NewLabeledRowWidget returns a Widget that renders a row with specified
// label on the left and control widget on the right.
//
// Note: labelWidthPixels is density-independent pixel value
func NewLabeledRowWidget(theme *material.Theme, label string, labelWidthPixels int, control layout.Widget) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Min.X = gtx.Dp(unit.Dp(labelWidthPixels))
				label := material.Body1(theme, label)
				label.Color = themes.ColorText
				label.TextSize = unit.Sp(13)
				return label.Layout(gtx)
			}),
			layout.Flexed(1, control))
	}
}
