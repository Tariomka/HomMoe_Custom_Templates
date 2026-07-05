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
	textColor := themes.ColorText
	backgroundColor := themes.ColorButton
	borderColor := themes.ColorBorder
	return func(gtx layout.Context) layout.Dimensions {
		if disabled {
			gtx = gtx.Disabled()
			textColor = themes.ColorTextDim
		}
		if !disabled && button.Hovered() {
			borderColor = themes.ColorHover
		}

		return material.Clickable(gtx, button, func(gtx layout.Context) layout.Dimensions {
			macro := op.Record(gtx.Ops)
			buttonDimensions := newButtonInset().Layout(gtx, NewLabelWidget(theme, label, textColor))
			call := macro.Stop()
			radius := gtx.Dp(constants.DefaultRoundness)
			rect := image.Rectangle{Max: buttonDimensions.Size}
			paint.FillShape(gtx.Ops, backgroundColor, clip.UniformRRect(rect, radius).Op(gtx.Ops))
			paint.FillShape(gtx.Ops, borderColor, clip.Stroke{
				Path:  clip.UniformRRect(rect, radius).Path(gtx.Ops),
				Width: float32(gtx.Dp(1)),
			}.Op())
			call.Add(gtx.Ops)
			return buttonDimensions
		})
	}
}

// NewToggleButtonWidget returns a Widget that renders a toggle button with the given text.
// If the button is active, it looks like Bright button widget, otherwise it looks like a standard button widget.
func NewToggleButtonWidget(theme *material.Theme, label string, button *widget.Clickable, active bool) layout.Widget {
	textColor := themes.ColorText
	backgroundColor := themes.ColorButton
	borderColor := themes.ColorBorder
	return func(gtx layout.Context) layout.Dimensions {
		if button.Hovered() {
			borderColor = themes.ColorHover
		}
		if active {
			textColor = themes.ColorAccentBright
			backgroundColor = themes.ColorPrimaryButton
			borderColor = themes.ColorAccent
		}

		return material.Clickable(gtx, button, func(gtx layout.Context) layout.Dimensions {
			macro := op.Record(gtx.Ops)
			dims := newButtonInset().Layout(gtx, NewLabelWidget(theme, label, textColor))
			call := macro.Stop()
			radius := gtx.Dp(constants.DefaultRoundness)
			rect := image.Rectangle{Max: dims.Size}
			paint.FillShape(gtx.Ops, backgroundColor, clip.UniformRRect(rect, radius).Op(gtx.Ops))
			paint.FillShape(gtx.Ops, borderColor, clip.Stroke{
				Path:  clip.UniformRRect(rect, radius).Path(gtx.Ops),
				Width: float32(gtx.Dp(1)),
			}.Op())
			call.Add(gtx.Ops)
			return dims
		})
	}
}

// NewSegmentButtonWidget returns a Widget that renders a segment button with the given text.
// If the button is active, it looks like Bright button widget, otherwise it looks like a disabled button widget.
func NewSegmentButtonWidget(theme *material.Theme, label string, button *widget.Clickable, active bool) layout.Widget {
	textColor := themes.ColorTextDim
	backgroundColor := themes.ColorInput
	borderColor := themes.ColorBorder
	return func(gtx layout.Context) layout.Dimensions {
		if button.Hovered() && !active {
			borderColor = themes.ColorHover
		}
		if active {
			backgroundColor = themes.ColorPrimaryButton
			textColor = themes.ColorAccentBright
			borderColor = themes.ColorAccent
		}

		return material.Clickable(gtx, button, func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(unit.Dp(2)).Layout(gtx,
				func(gtx layout.Context) layout.Dimensions {
					macro := op.Record(gtx.Ops)
					dims := layout.UniformInset(constants.DefaultPaddingSmall).
						Layout(gtx, NewLabelWidget(theme, label, textColor))
					call := macro.Stop()
					radius := gtx.Dp(constants.DefaultRoundness)
					rect := image.Rectangle{Max: dims.Size}
					paint.FillShape(gtx.Ops, backgroundColor, clip.UniformRRect(rect, radius).Op(gtx.Ops))
					paint.FillShape(gtx.Ops, borderColor, clip.Stroke{
						Path:  clip.UniformRRect(rect, radius).Path(gtx.Ops),
						Width: float32(gtx.Dp(1)),
					}.Op())
					call.Add(gtx.Ops)
					return dims
				})
		})
	}
}

// NewDropdownRowButtonWidget returns a Widget that renders a dropdown row button with the given text.
// If the button is active, it looks like a selected label widget, otherwise it looks like a normal label widget.
func NewDropdownRowButtonWidget(theme *material.Theme, label string, button *widget.Clickable, active bool) layout.Widget {
	textColor := themes.ColorText
	textFont := font.Font{Weight: font.Normal}
	backgroundColor := themes.ColorInput
	return func(gtx layout.Context) layout.Dimensions {
		if active {
			textColor = themes.ColorAccent
			textFont = font.Font{Weight: font.SemiBold}
			backgroundColor = themes.ColorSelection
		}
		if button.Hovered() {
			backgroundColor = themes.ColorHover
		}

		return material.Clickable(gtx, button, func(gtx layout.Context) layout.Dimensions {
			macro := op.Record(gtx.Ops)
			dims := layout.Inset{Top: unit.Dp(5), Bottom: unit.Dp(5), Left: unit.Dp(8), Right: unit.Dp(8)}.
				Layout(gtx, NewLabelBuilder(theme).WithSizeDefault().
					WithText(label).WithColor(textColor).WithFont(textFont).WithMaxLines(1).Build)
			call := macro.Stop()
			radius := gtx.Dp(constants.DefaultRoundness)
			rect := image.Rectangle{Max: dims.Size}
			paint.FillShape(gtx.Ops, backgroundColor, clip.UniformRRect(rect, radius).Op(gtx.Ops))
			call.Add(gtx.Ops)
			if dims.Size.X < gtx.Constraints.Min.X {
				dims.Size.X = gtx.Constraints.Min.X
			}
			return dims
		})
	}
}

func NewBrightButtonWidget(theme *material.Theme, label string, button *widget.Clickable, disabled bool) layout.Widget {
	textColor := themes.ColorAccentBright
	backgroundColor := themes.ColorPrimaryButton
	borderColor := themes.ColorAccent
	buttonStyle := font.Font{Weight: font.SemiBold}
	return func(gtx layout.Context) layout.Dimensions {
		if disabled {
			gtx = gtx.Disabled()
			textColor = themes.ColorTextDim
			backgroundColor = themes.ColorButtonDisabled
			borderColor = themes.ColorBorderDisabled
		}
		if button.Hovered() {
			borderColor = themes.ColorAccentBright
		}

		return material.Clickable(gtx, button, func(gtx layout.Context) layout.Dimensions {
			macro := op.Record(gtx.Ops)
			dims := newButtonInset().Layout(gtx, NewStyledLabelWidget(theme, label, textColor, buttonStyle))
			call := macro.Stop()
			radius := gtx.Dp(constants.DefaultRoundness)
			rect := image.Rectangle{Max: dims.Size}
			paint.FillShape(gtx.Ops, backgroundColor, clip.UniformRRect(rect, radius).Op(gtx.Ops))
			paint.FillShape(gtx.Ops, borderColor, clip.Stroke{
				Path:  clip.UniformRRect(rect, radius).Path(gtx.Ops),
				Width: float32(gtx.Dp(1)),
			}.Op())
			call.Add(gtx.Ops)
			return dims
		})
	}
}

func NewBrightButtonLargeWidget(theme *material.Theme, label string, button *widget.Clickable, disabled bool) layout.Widget {
	textColor := themes.ColorAccentBright
	backgroundColor := themes.ColorPrimaryButton
	borderColor := themes.ColorAccent
	return func(gtx layout.Context) layout.Dimensions {
		if disabled {
			gtx = gtx.Disabled()
			textColor = themes.ColorTextDim
			backgroundColor = themes.ColorButtonDisabled
			borderColor = themes.ColorBorderDisabled
		} else if button.Hovered() {
			borderColor = themes.ColorAccentBright
		}

		return material.Clickable(gtx, button, func(gtx layout.Context) layout.Dimensions {
			macro := op.Record(gtx.Ops)
			dims := layout.UniformInset(constants.DefaultPaddingLarge).
				Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					label := material.Body2(theme, label)
					label.Font = font.Font{Weight: font.SemiBold}
					label.Color = textColor
					return layout.Center.Layout(gtx, label.Layout)
				})
			call := macro.Stop()
			radius := gtx.Dp(constants.DefaultRoundnessLarge)
			rect := image.Rectangle{Max: dims.Size}
			paint.FillShape(gtx.Ops, backgroundColor, clip.UniformRRect(rect, radius).Op(gtx.Ops))
			paint.FillShape(gtx.Ops, borderColor, clip.Stroke{
				Path:  clip.UniformRRect(rect, radius).Path(gtx.Ops),
				Width: float32(gtx.Dp(1)),
			}.Op())
			call.Add(gtx.Ops)
			return dims
		})
	}
}

func newButtonInset() layout.Inset {
	return layout.Inset{
		Top:    constants.DefaultPaddingSmall,
		Bottom: constants.DefaultPaddingSmall,
		Left:   constants.DefaultPaddingLarge,
		Right:  constants.DefaultPaddingLarge,
	}
}
