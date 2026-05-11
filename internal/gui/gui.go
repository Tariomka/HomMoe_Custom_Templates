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
		app.Size(unit.Dp(1000), unit.Dp(780)),
		app.MinSize(unit.Dp(720), unit.Dp(560)),
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

// sliderWithValue draws a Material slider with a gold value label on the right.
func sliderWithValue(gtx layout.Context, th *material.Theme, f *widget.Float, value string) layout.Dimensions {
	return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			sl := material.Slider(th, f)
			sl.Color = colGold
			return sl.Layout(gtx)
		}),
		layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min.X = gtx.Dp(48)
			lbl := material.Body1(th, value)
			lbl.Color = colGold
			lbl.TextSize = unit.Sp(13)
			lbl.Alignment = 1 // text.End
			return lbl.Layout(gtx)
		}),
	)
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

// _ silences unused color import when builds remove specific consts.
var _ = color.NRGBA{}
