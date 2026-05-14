package widgets

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/Tariomka/hommoe_custom_templates/internal/gui/components/themes"
)

// NewButtonWidget returns a Widget that renders a button with the given text
func NewButtonWidget(theme *material.Theme, text string, button *widget.Clickable, disabled bool) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		if disabled {
			gtx = gtx.Disabled()
		}
		return material.Clickable(gtx, button, func(gtx layout.Context) layout.Dimensions {
			macro := op.Record(gtx.Ops)
			dims := layout.Inset{Top: unit.Dp(5), Bottom: unit.Dp(5), Left: unit.Dp(10), Right: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					label := material.Body2(theme, text)
					label.Color = themes.ColorText
					label.TextSize = unit.Sp(12)
					if disabled {
						label.Color = themes.ColorTextDim
					}
					return label.Layout(gtx)
				})
			})
			call := macro.Stop()
			radius := gtx.Dp(3)
			rect := image.Rectangle{Max: dims.Size}
			paint.FillShape(gtx.Ops, color.NRGBA{R: 0x2A, G: 0x2A, B: 0x2A, A: 0xFF},
				clip.UniformRRect(rect, radius).Op(gtx.Ops))
			paint.FillShape(gtx.Ops, themes.ColorBorder, clip.Stroke{
				Path:  clip.UniformRRect(rect, radius).Path(gtx.Ops),
				Width: float32(gtx.Dp(1)),
			}.Op())
			call.Add(gtx.Ops)
			return dims
		})
	}
}
