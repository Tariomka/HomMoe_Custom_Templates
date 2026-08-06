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
	"github.com/Tariomka/hommoe_custom_templates/app/gui/constants"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/themes"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/utils"
)

// NewButtonWidget returns a Widget that renders a button with the given text.
func NewButtonWidget(theme *material.Theme, label string, button *widget.Clickable, disabled bool) layout.Widget {
	textColor := themes.ColorsBase.Text
	backgroundColor := themes.ColorsBase.Button
	borderColor := themes.ColorsBase.Border
	return func(gtx layout.Context) layout.Dimensions {
		if disabled {
			gtx = gtx.Disabled()
			textColor = themes.ColorsBase.TextDim
		}
		if !disabled && button.Hovered() {
			borderColor = themes.ColorsBase.Hover
		}

		return material.Clickable(gtx, button,
			newBaseButtonWidget(label, NewLabelWidget(theme, label, textColor), backgroundColor, borderColor))
	}
}

// NewToggleButtonWidget returns a Widget that renders a toggle button with the given text.
// If the button is active, it looks like Bright button widget, otherwise it looks like a standard button widget.
func NewToggleButtonWidget(theme *material.Theme, label string, button *widget.Clickable, active bool) layout.Widget {
	textColor := themes.ColorsBase.Text
	backgroundColor := themes.ColorsBase.Button
	borderColor := themes.ColorsBase.Border
	return func(gtx layout.Context) layout.Dimensions {
		if button.Hovered() {
			borderColor = themes.ColorsBase.Hover
		}
		if active {
			textColor = themes.ColorsBase.AccentBright
			backgroundColor = themes.ColorsBase.PrimaryButton
			borderColor = themes.ColorsBase.Accent
		}

		return material.Clickable(gtx, button,
			newBaseButtonWidget(label, NewLabelWidget(theme, label, textColor), backgroundColor, borderColor))
	}
}

// NewSegmentButtonWidget returns a Widget that renders a segment button with the given text.
// If the button is active, it looks like Bright button widget, otherwise it looks like a disabled button widget.
func NewSegmentButtonWidget(theme *material.Theme, label string, button *widget.Clickable, active bool) layout.Widget {
	textColor := themes.ColorsBase.TextDim
	backgroundColor := themes.ColorsBase.Input
	borderColor := themes.ColorsBase.Border
	return func(gtx layout.Context) layout.Dimensions {
		if button.Hovered() && !active {
			borderColor = themes.ColorsBase.Hover
		}
		if active {
			backgroundColor = themes.ColorsBase.PrimaryButton
			textColor = themes.ColorsBase.AccentBright
			borderColor = themes.ColorsBase.Accent
		}

		return material.Clickable(gtx, button, func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(unit.Dp(2)).Layout(gtx,
				func(gtx layout.Context) layout.Dimensions {
					macro := op.Record(gtx.Ops)
					buttonDimensions := layout.UniformInset(constants.DefaultPaddingSmall).
						Layout(gtx, NewLabelWidget(theme, label, textColor))
					call := macro.Stop()
					radius := gtx.Dp(constants.DefaultRoundness)
					rect := image.Rectangle{Max: buttonDimensions.Size}
					paint.FillShape(gtx.Ops, backgroundColor, clip.UniformRRect(rect, radius).Op(gtx.Ops))
					paint.FillShape(gtx.Ops, borderColor, clip.Stroke{
						Path:  clip.UniformRRect(rect, radius).Path(gtx.Ops),
						Width: float32(gtx.Dp(1)),
					}.Op())
					call.Add(gtx.Ops)
					utils.AddButtonSemantics(gtx.Ops, label, buttonDimensions.Size)
					return buttonDimensions
				})
		})
	}
}

// NewDropdownRowButtonWidget returns a Widget that renders a dropdown row button with the given text.
// If the button is active, it looks like a selected label widget, otherwise it looks like a normal label widget.
func NewDropdownRowButtonWidget(
	theme *material.Theme,
	label string,
	button *widget.Clickable,
	active bool) layout.Widget {
	textColor := themes.ColorsBase.Text
	textFont := font.Font{Weight: font.Normal}
	backgroundColor := themes.ColorsBase.Input
	return func(gtx layout.Context) layout.Dimensions {
		if active {
			textColor = themes.ColorsBase.Accent
			textFont = font.Font{Weight: font.SemiBold}
			backgroundColor = themes.ColorsBase.Selection
		}
		if button.Hovered() {
			backgroundColor = themes.ColorsBase.Hover
		}

		return material.Clickable(gtx, button, func(gtx layout.Context) layout.Dimensions {
			macro := op.Record(gtx.Ops)
			buttonDimensions := layout.Inset{Top: unit.Dp(5), Bottom: unit.Dp(5), Left: unit.Dp(8), Right: unit.Dp(8)}.
				Layout(gtx, NewLabelBuilder(theme).WithSizeDefault().
					WithText(label).WithColor(textColor).WithFont(textFont).WithMaxLines(1).Build)
			call := macro.Stop()
			radius := gtx.Dp(constants.DefaultRoundness)
			rect := image.Rectangle{Max: buttonDimensions.Size}
			paint.FillShape(gtx.Ops, backgroundColor, clip.UniformRRect(rect, radius).Op(gtx.Ops))
			call.Add(gtx.Ops)
			utils.AddButtonSemantics(gtx.Ops, label, buttonDimensions.Size)
			if buttonDimensions.Size.X < gtx.Constraints.Min.X {
				buttonDimensions.Size.X = gtx.Constraints.Min.X
			}
			return buttonDimensions
		})
	}
}

func NewBrightButtonWidget(theme *material.Theme, label string, button *widget.Clickable, disabled bool) layout.Widget {
	textColor := themes.ColorsBase.AccentBright
	backgroundColor := themes.ColorsBase.PrimaryButton
	borderColor := themes.ColorsBase.Accent
	buttonStyle := font.Font{Weight: font.SemiBold}
	return func(gtx layout.Context) layout.Dimensions {
		if disabled {
			gtx = gtx.Disabled()
			textColor = themes.ColorsBase.TextDim
			backgroundColor = themes.ColorsBase.ButtonDisabled
			borderColor = themes.ColorsBase.BorderDisabled
		}
		if button.Hovered() {
			borderColor = themes.ColorsBase.AccentBright
		}

		return material.Clickable(gtx, button,
			newBaseButtonWidget(
				label, NewStyledLabelWidget(theme, label, textColor, buttonStyle), backgroundColor, borderColor),
		)
	}
}

func NewBrightButtonLargeWidget(
	theme *material.Theme,
	label string,
	button *widget.Clickable,
	disabled bool) layout.Widget {
	textColor := themes.ColorsBase.AccentBright
	backgroundColor := themes.ColorsBase.PrimaryButton
	borderColor := themes.ColorsBase.Accent
	return func(gtx layout.Context) layout.Dimensions {
		if disabled {
			gtx = gtx.Disabled()
			textColor = themes.ColorsBase.TextDim
			backgroundColor = themes.ColorsBase.ButtonDisabled
			borderColor = themes.ColorsBase.BorderDisabled
		} else if button.Hovered() {
			borderColor = themes.ColorsBase.AccentBright
		}

		return material.Clickable(gtx, button, func(gtx layout.Context) layout.Dimensions {
			macro := op.Record(gtx.Ops)
			buttonDimensions := layout.UniformInset(constants.DefaultPaddingLarge).
				Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					label := material.Body2(theme, label)
					label.Font = font.Font{Weight: font.SemiBold}
					label.Color = textColor
					return layout.Center.Layout(gtx, label.Layout)
				})
			call := macro.Stop()
			radius := gtx.Dp(constants.DefaultRoundnessLarge)
			rect := image.Rectangle{Max: buttonDimensions.Size}
			paint.FillShape(gtx.Ops, backgroundColor, clip.UniformRRect(rect, radius).Op(gtx.Ops))
			paint.FillShape(gtx.Ops, borderColor, clip.Stroke{
				Path:  clip.UniformRRect(rect, radius).Path(gtx.Ops),
				Width: float32(gtx.Dp(1)),
			}.Op())
			call.Add(gtx.Ops)
			utils.AddButtonSemantics(gtx.Ops, label, buttonDimensions.Size)
			return buttonDimensions
		})
	}
}

func newBaseButtonWidget(
	label string,
	labelWidget layout.Widget,
	backgroundColor, borderColor color.NRGBA) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		macro := op.Record(gtx.Ops)
		buttonDimensions := newButtonInset().Layout(gtx, labelWidget)
		call := macro.Stop()
		radius := gtx.Dp(constants.DefaultRoundness)
		rect := image.Rectangle{Max: buttonDimensions.Size}
		paint.FillShape(gtx.Ops, backgroundColor, clip.UniformRRect(rect, radius).Op(gtx.Ops))
		paint.FillShape(gtx.Ops, borderColor, clip.Stroke{
			Path:  clip.UniformRRect(rect, radius).Path(gtx.Ops),
			Width: float32(gtx.Dp(1)),
		}.Op())
		call.Add(gtx.Ops)
		utils.AddButtonSemantics(gtx.Ops, label, buttonDimensions.Size)
		return buttonDimensions
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
