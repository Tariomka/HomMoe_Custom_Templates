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
	label  string
	button *widget.Clickable

	isDisabledCallback func() bool
	clickHandler       func()
}

func NewButton(label string, isDisabledCallback func() bool, clickHandler func()) *Button {
	return &Button{
		label:              label,
		isDisabledCallback: isDisabledCallback,
		clickHandler:       clickHandler,
	}
}

func (this *Button) GetWidget(theme *material.Theme) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		isDisabled := this.isDisabledCallback()
		if isDisabled {
			gtx = gtx.Disabled()
		}
		return material.Clickable(gtx, this.button, func(gtx layout.Context) layout.Dimensions {
			macro := op.Record(gtx.Ops)
			dims := layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					label := material.Body1(theme, this.label)
					label.Color = themes.ColorGoldBright
					label.TextSize = unit.Sp(14)
					label.Font = font.Font{Weight: font.SemiBold}
					if isDisabled {
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
			if isDisabled {
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

func (this *Button) HandleClicks(gtx layout.Context) {
	if this.button.Clicked(gtx) {
		this.clickHandler()
	}
}
