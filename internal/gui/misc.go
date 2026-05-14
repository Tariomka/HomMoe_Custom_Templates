package gui

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

// sliderLabeled draws a slider in a flex row with a fixed-width gold value
// label on the right.
func sliderLabeled(gtx layout.Context, theme *material.Theme, floatValue *widget.Float, value string) layout.Dimensions {
	return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			slider := material.Slider(theme, floatValue)
			slider.Color = colGold
			return slider.Layout(gtx)
		}),
		layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min.X = gtx.Dp(64)
			label := material.Body1(theme, value)
			label.Color = colGold
			label.TextSize = unit.Sp(13)
			label.Alignment = 1
			return label.Layout(gtx)
		}),
	)
}

// snapIntSliderLabeled snaps the slider to the [low, high] integer range and
// renders the integer value to the right.
func snapIntSliderLabeled(gtx layout.Context, theme *material.Theme, floatValue *widget.Float, low, high int, suffix string) int {
	mapped := mapRange(floatValue.Value, float32(low), float32(high))
	rounded := int(roundHalfAway(float64(mapped)))
	if rounded < low {
		rounded = low
	}
	if rounded > high {
		rounded = high
	}
	target := mapRangeInv(float32(rounded), float32(low), float32(high))
	if target != floatValue.Value && !floatValue.Dragging() {
		floatValue.Value = target
	}
	labelText := itoa(rounded) + suffix
	layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			slider := material.Slider(theme, floatValue)
			slider.Color = colGold
			return slider.Layout(gtx)
		}),
		layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min.X = gtx.Dp(64)
			label := material.Body1(theme, labelText)
			label.Color = colGold
			label.TextSize = unit.Sp(13)
			label.Alignment = 1
			return label.Layout(gtx)
		}),
	)
	return rounded
}

func itoa(number int) string {
	if number == 0 {
		return "0"
	}
	neg := number < 0
	if neg {
		number = -number
	}
	var buf [20]byte
	i := len(buf)
	for number > 0 {
		i--
		buf[i] = byte('0' + number%10)
		number /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}
