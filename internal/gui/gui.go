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
	w := new(app.Window)
	w.Option(
		app.Title("Olden Era — Template Generator"),
		app.Size(unit.Dp(1180), unit.Dp(820)),
		app.MinSize(unit.Dp(900), unit.Dp(600)),
	)
	th := newTheme()
	state := newState()

	var ops op.Ops
	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			state.Layout(gtx, th)
			e.Frame(gtx.Ops)
		}
	}
}

// drawEditor renders a text editor with a bordered input background.
func drawEditor(gtx layout.Context, th *material.Theme, ed *widget.Editor, hint string) layout.Dimensions {
	macro := op.Record(gtx.Ops)
	inset := layout.Inset{Top: unit.Dp(6), Bottom: unit.Dp(6), Left: unit.Dp(8), Right: unit.Dp(8)}
	dims := inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		e := material.Editor(th, ed, hint)
		e.Color = colText
		e.HintColor = colTextDim
		e.TextSize = unit.Sp(13)
		return e.Layout(gtx)
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
func sliderLabeled(gtx layout.Context, th *material.Theme, f *widget.Float, value string) layout.Dimensions {
	return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			sl := material.Slider(th, f)
			sl.Color = colGold
			return sl.Layout(gtx)
		}),
		layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min.X = gtx.Dp(64)
			lbl := material.Body1(th, value)
			lbl.Color = colGold
			lbl.TextSize = unit.Sp(13)
			lbl.Alignment = 1
			return lbl.Layout(gtx)
		}),
	)
}

// snapIntSliderLabeled snaps the slider to the [lo, hi] integer range and
// renders the integer value to the right.
func snapIntSliderLabeled(gtx layout.Context, th *material.Theme, f *widget.Float, lo, hi int, suffix string) int {
	v := mapRange(f.Value, float32(lo), float32(hi))
	rounded := int(roundHalfAway(float64(v)))
	if rounded < lo {
		rounded = lo
	}
	if rounded > hi {
		rounded = hi
	}
	target := mapRangeInv(float32(rounded), float32(lo), float32(hi))
	if target != f.Value && !f.Dragging() {
		f.Value = target
	}
	label := itoa(rounded) + suffix
	layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			sl := material.Slider(th, f)
			sl.Color = colGold
			return sl.Layout(gtx)
		}),
		layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min.X = gtx.Dp(64)
			lbl := material.Body1(th, label)
			lbl.Color = colGold
			lbl.TextSize = unit.Sp(13)
			lbl.Alignment = 1
			return lbl.Layout(gtx)
		}),
	)
	return rounded
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}

// checkRow renders one CheckBox + label as a clickable row.
func checkRow(th *material.Theme, b *widget.Bool, label string) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{Top: unit.Dp(3), Bottom: unit.Dp(3)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			cb := material.CheckBox(th, b, label)
			cb.Color = colText
			cb.IconColor = colGold
			cb.TextSize = unit.Sp(13)
			return cb.Layout(gtx)
		})
	}
}

// _ silences unused-color warnings if the file is partially used.
var _ = color.NRGBA{}
