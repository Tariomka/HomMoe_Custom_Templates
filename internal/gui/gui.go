package gui

import (
	"image"
	"image/color"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

// Run opens the GUI window and runs the event loop until close.
func Run() error {
	window := new(app.Window)
	window.Option(
		app.Title("Olden Era — Template Generator"),
		app.Size(unit.Dp(1180), unit.Dp(820)),
		app.MinSize(unit.Dp(900), unit.Dp(600)),
	)
	theme := newTheme()
	state := newState()

	var ops op.Ops
	for {
		switch event := window.Event().(type) {
		case app.DestroyEvent:
			return event.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, event)
			state.Layout(gtx, theme)
			event.Frame(gtx.Ops)
		}
	}
}

// drawEditor renders a text editor with a bordered input background.
func drawEditor(gtx layout.Context, theme *material.Theme, editor *widget.Editor, hint string) layout.Dimensions {
	macro := op.Record(gtx.Ops)
	inset := layout.Inset{Top: unit.Dp(6), Bottom: unit.Dp(6), Left: unit.Dp(8), Right: unit.Dp(8)}
	dims := inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		editorWidget := material.Editor(theme, editor, hint)
		editorWidget.Color = colText
		editorWidget.HintColor = colTextDim
		editorWidget.TextSize = unit.Sp(13)
		return editorWidget.Layout(gtx)
	})
	call := macro.Stop()
	radius := gtx.Dp(2)
	rect := image.Rectangle{Max: dims.Size}
	paint.FillShape(gtx.Ops, colInput, clip.UniformRRect(rect, radius).Op(gtx.Ops))
	paint.FillShape(gtx.Ops, colBorder, clip.Stroke{
		Path:  clip.UniformRRect(rect, radius).Path(gtx.Ops),
		Width: float32(gtx.Dp(1)),
	}.Op())
	call.Add(gtx.Ops)
	return dims
}

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

// checkRow renders one CheckBox + label as a clickable row.
func checkRow(theme *material.Theme, boolValue *widget.Bool, label string) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{Top: unit.Dp(3), Bottom: unit.Dp(3)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			checkbox := material.CheckBox(theme, boolValue, label)
			checkbox.Color = colText
			checkbox.IconColor = colGold
			checkbox.TextSize = unit.Sp(13)
			return checkbox.Layout(gtx)
		})
	}
}

// _ silences unused-color warnings if the file is partially used.
var _ = color.NRGBA{}
