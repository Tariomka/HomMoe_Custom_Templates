package content

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

type Button struct {
	Text     string
	Click    *widget.Clickable
	Disabled bool
}

func NewButton(text string, click *widget.Clickable) *Button {
	return &Button{
		Text:  text,
		Click: click,
	}
}

// .
// .
// .
// .
// .
// .
// .
func (this *Button) Layout(gtx layout.Context, theme *material.Theme) layout.Dimensions {
	if this.Disabled {
		gtx = gtx.Disabled()
	}
	return material.Clickable(gtx, this.Click, func(gtx layout.Context) layout.Dimensions {
		macro := op.Record(gtx.Ops)
		dims := layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				label := material.Body1(theme, this.Text)
				label.Color = themes.ColorGoldBright
				label.TextSize = unit.Sp(14)
				label.Font = font.Font{Weight: font.SemiBold}
				if this.Disabled {
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
		if this.Disabled {
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
