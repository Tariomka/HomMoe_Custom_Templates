package widgets

import (
	"image"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/constants"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/themes"
)

// NewButtonWidget returns a Widget that renders a button with the given text
func NewButtonWidget(theme *material.Theme, label string, button *widget.Clickable, disabled bool) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		if disabled {
			gtx = gtx.Disabled()
		}
		return material.Clickable(gtx, button, func(gtx layout.Context) layout.Dimensions {
			macro := op.Record(gtx.Ops)
			dims := layout.Inset{Top: unit.Dp(5), Bottom: unit.Dp(5), Left: unit.Dp(10), Right: unit.Dp(10)}.
				Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					label := material.Caption(theme, label)
					label.Color = themes.ColorText
					if disabled {
						label.Color = themes.ColorTextDim
					}
					return layout.Center.Layout(gtx, label.Layout)
				})
			call := macro.Stop()
			radius := gtx.Dp(constants.DefaultRoundness)
			rect := image.Rectangle{Max: dims.Size}
			paint.FillShape(gtx.Ops, themes.ColorButton,
				clip.UniformRRect(rect, radius).Op(gtx.Ops))
			border := themes.ColorBorder
			if !disabled && button.Hovered() {
				border = themes.ColorHover
			}
			paint.FillShape(gtx.Ops, border, clip.Stroke{
				Path:  clip.UniformRRect(rect, radius).Path(gtx.Ops),
				Width: float32(gtx.Dp(1)),
			}.Op())
			call.Add(gtx.Ops)
			return dims
		})
	}
}

// NewToggleButtonWidget returns a segment-style toggle button matching the
// SegmentButtonGroup visuals (tier selector, game-mode selector): a compact
// 12sp label that turns gold while active.
func NewToggleButtonWidget(theme *material.Theme, label string, button *widget.Clickable, active bool) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		return material.Clickable(gtx, button, func(gtx layout.Context) layout.Dimensions {
			bgColor := themes.ColorInput
			fgColor := themes.ColorTextDim
			border := themes.ColorBorder
			if button.Hovered() && !active {
				border = themes.ColorHover
			}
			if active {
				bgColor = themes.ColorPrimaryButton
				fgColor = themes.ColorAccentBright
				border = themes.ColorAccent
			}
			macro := op.Record(gtx.Ops)
			dims := layout.UniformInset(constants.DefaultPaddingSmall).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				label := material.Caption(theme, label)
				label.Color = fgColor
				return label.Layout(gtx)
			})
			call := macro.Stop()
			radius := gtx.Dp(constants.DefaultRoundness)
			rect := image.Rectangle{Max: dims.Size}
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

func NewGoldButtonWidget(theme *material.Theme, label string, button *widget.Clickable, disabled bool) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		if disabled {
			gtx = gtx.Disabled()
		}
		return material.Clickable(gtx, button, func(gtx layout.Context) layout.Dimensions {
			macro := op.Record(gtx.Ops)
			dims := layout.UniformInset(constants.DefaultPaddingLarge).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				label := material.Body2(theme, label)
				label.Font = font.Font{Weight: font.SemiBold}
				label.Color = themes.ColorAccentBright
				if disabled {
					label.Color = themes.ColorTextDim
				}
				return layout.Center.Layout(gtx, label.Layout)
			})
			call := macro.Stop()
			radius := gtx.Dp(constants.DefaultRoundness)
			rect := image.Rectangle{Max: dims.Size}
			bgColor := themes.ColorPrimaryButton
			border := themes.ColorAccent
			if disabled {
				bgColor = themes.ColorButtonDisabled
				border = themes.ColorBorderDisabled
			} else if button.Hovered() {
				border = themes.ColorAccentBright
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
