package widgets

import (
	"image"
	"image/color"

	"gioui.org/font"
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
func NewButtonWidget(theme *material.Theme, label string, button *widget.Clickable, disabled bool) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		if disabled {
			gtx = gtx.Disabled()
		}
		return material.Clickable(gtx, button, func(gtx layout.Context) layout.Dimensions {
			macro := op.Record(gtx.Ops)
			dims := layout.Inset{Top: unit.Dp(5), Bottom: unit.Dp(5), Left: unit.Dp(10), Right: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					label := material.Body2(theme, label)
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

func NewGoldButtonWidget(theme *material.Theme, label string, button *widget.Clickable, disabled bool) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		if disabled {
			gtx = gtx.Disabled()
		}
		return material.Clickable(gtx, button, func(gtx layout.Context) layout.Dimensions {
			macro := op.Record(gtx.Ops)
			dims := layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					label := material.Body1(theme, label)
					label.Color = themes.ColorGoldBright
					label.TextSize = unit.Sp(14)
					label.Font = font.Font{Weight: font.SemiBold}
					if disabled {
						label.Color = themes.ColorTextDim
					}
					return label.Layout(gtx)
				})
			})
			call := macro.Stop()
			radius := gtx.Dp(3)
			rect := image.Rectangle{Max: dims.Size}
			bgColor := themes.ColorGenerate
			border := themes.ColorGold
			if disabled {
				bgColor = color.NRGBA{R: 0x3A, G: 0x30, B: 0x20, A: 0xFF}
				border = color.NRGBA{R: 0x4A, G: 0x40, B: 0x30, A: 0xFF}
			}
			paint.FillShape(gtx.Ops, bgColor, clip.UniformRRect(rect, radius).Op(gtx.Ops))
			paint.FillShape(gtx.Ops, border, clip.Stroke{
				Path:  clip.UniformRRect(rect, radius).Path(gtx.Ops),
				Width: float32(gtx.Dp(1)),
			}.Op())
			call.Add(gtx.Ops)
			return dims
		})
	}
}
