package gui

import (
	"image"
	"image/color"

	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

// fillBackground paints the given color over gtx.Constraints.Max.
func fillBackground(gtx layout.Context, c color.NRGBA) {
	rect := image.Rectangle{Max: gtx.Constraints.Max}
	paint.FillShape(gtx.Ops, c, clip.Rect(rect).Op())
}

// roundedRect paints a rounded rect of the given color, filling the constraints.
func roundedRect(gtx layout.Context, c color.NRGBA, radius int) {
	rect := image.Rectangle{Max: gtx.Constraints.Max}
	paint.FillShape(gtx.Ops, c, clip.UniformRRect(rect, radius).Op(gtx.Ops))
}

// borderedPanel wraps content with a rounded border on a panel background.
func borderedPanel(gtx layout.Context, padding unit.Dp, w layout.Widget) layout.Dimensions {
	radius := gtx.Dp(4)
	macro := op.Record(gtx.Ops)
	dims := layout.UniformInset(padding).Layout(gtx, w)
	call := macro.Stop()
	rect := image.Rectangle{Max: dims.Size}
	// Fill background.
	paint.FillShape(gtx.Ops, colPanel, clip.UniformRRect(rect, radius).Op(gtx.Ops))
	// Border (1dp).
	border := clip.Stroke{
		Path:  clip.UniformRRect(rect, radius).Path(gtx.Ops),
		Width: float32(gtx.Dp(1)),
	}.Op()
	paint.FillShape(gtx.Ops, colBorder, border)
	// Replay content.
	call.Add(gtx.Ops)
	return dims
}

// sectionHeader draws an accented gold section title.
func sectionHeader(gtx layout.Context, th *material.Theme, text string) layout.Dimensions {
	lbl := material.Body1(th, "◆  "+text)
	lbl.Color = colGold
	lbl.TextSize = unit.Sp(13)
	return lbl.Layout(gtx)
}

// dimLabel renders a small dimmed description line.
func dimLabel(gtx layout.Context, th *material.Theme, text string) layout.Dimensions {
	lbl := material.Body2(th, text)
	lbl.Color = colTextDim
	lbl.TextSize = unit.Sp(12)
	return lbl.Layout(gtx)
}

// segmentButton represents a single selectable option in a row of segments.
type segmentButton struct {
	Label string
	click widget.Clickable
}

// segmentGroup is a horizontal row of mutually-exclusive segment buttons,
// used in place of a ComboBox for small enums.
type segmentGroup struct {
	buttons  []segmentButton
	Selected int
}

func newSegmentGroup(labels []string) *segmentGroup {
	g := &segmentGroup{buttons: make([]segmentButton, len(labels))}
	for i, l := range labels {
		g.buttons[i].Label = l
	}
	return g
}

// Update returns true if the selection changed this frame.
func (g *segmentGroup) Update(gtx layout.Context) bool {
	changed := false
	for i := range g.buttons {
		if g.buttons[i].click.Clicked(gtx) && g.Selected != i {
			g.Selected = i
			changed = true
		}
	}
	return changed
}

func (g *segmentGroup) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	g.Update(gtx)
	children := make([]layout.FlexChild, 0, len(g.buttons))
	for i := range g.buttons {
		i := i
		b := &g.buttons[i]
		selected := i == g.Selected
		children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return material.Clickable(gtx, &b.click, func(gtx layout.Context) layout.Dimensions {
				return layout.UniformInset(unit.Dp(2)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return drawSegment(gtx, th, b.Label, selected)
				})
			})
		}))
	}
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx, children...)
}

func drawSegment(gtx layout.Context, th *material.Theme, label string, selected bool) layout.Dimensions {
	bg := colInput
	fg := colTextDim
	border := colBorder
	if selected {
		bg = colGenerate
		fg = colGoldBright
		border = colGold
	}
	macro := op.Record(gtx.Ops)
	dims := layout.UniformInset(unit.Dp(8)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		lbl := material.Body2(th, label)
		lbl.Color = fg
		lbl.TextSize = unit.Sp(12)
		return lbl.Layout(gtx)
	})
	call := macro.Stop()
	radius := gtx.Dp(3)
	rect := image.Rectangle{Max: dims.Size}
	paint.FillShape(gtx.Ops, bg, clip.UniformRRect(rect, radius).Op(gtx.Ops))
	paint.FillShape(gtx.Ops, border, clip.Stroke{
		Path:  clip.UniformRRect(rect, radius).Path(gtx.Ops),
		Width: float32(gtx.Dp(1)),
	}.Op())
	call.Add(gtx.Ops)
	return dims
}

// labeledRow lays out a label on the left and the supplied widget on the right.
func labeledRow(gtx layout.Context, th *material.Theme, label string, labelW unit.Dp, control layout.Widget) layout.Dimensions {
	return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min.X = gtx.Dp(labelW)
			lbl := material.Body1(th, label)
			lbl.Color = colText
			lbl.TextSize = unit.Sp(13)
			return lbl.Layout(gtx)
		}),
		layout.Flexed(1, control),
	)
}

// goldButton renders a Generate-style emphasised button.
type goldButton struct {
	Text     string
	Click    *widget.Clickable
	Disabled bool
}

func (b goldButton) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	if b.Disabled {
		// Block pointer events when disabled.
		gtx = gtx.Disabled()
	}
	return material.Clickable(gtx, b.Click, func(gtx layout.Context) layout.Dimensions {
		macro := op.Record(gtx.Ops)
		dims := layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				lbl := material.Body1(th, b.Text)
				lbl.Color = colGoldBright
				lbl.TextSize = unit.Sp(14)
				if b.Disabled {
					lbl.Color = colTextDim
				}
				return lbl.Layout(gtx)
			})
		})
		call := macro.Stop()
		radius := gtx.Dp(3)
		rect := image.Rectangle{Max: dims.Size}
		bg := colGenerate
		border := colGold
		if b.Disabled {
			bg = color.NRGBA{R: 0x3A, G: 0x30, B: 0x20, A: 0xFF}
			border = color.NRGBA{R: 0x4A, G: 0x40, B: 0x30, A: 0xFF}
		}
		paint.FillShape(gtx.Ops, bg, clip.UniformRRect(rect, radius).Op(gtx.Ops))
		paint.FillShape(gtx.Ops, border, clip.Stroke{
			Path:  clip.UniformRRect(rect, radius).Path(gtx.Ops),
			Width: float32(gtx.Dp(1)),
		}.Op())
		call.Add(gtx.Ops)
		return dims
	})
}

// toolbarButton renders a small dark toolbar-style button.
type toolbarButton struct {
	Text  string
	Click *widget.Clickable
}

func (b toolbarButton) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return material.Clickable(gtx, b.Click, func(gtx layout.Context) layout.Dimensions {
		macro := op.Record(gtx.Ops)
		dims := layout.UniformInset(unit.Dp(6)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				lbl := material.Body2(th, b.Text)
				lbl.Color = colText
				lbl.TextSize = unit.Sp(12)
				return lbl.Layout(gtx)
			})
		})
		call := macro.Stop()
		radius := gtx.Dp(3)
		rect := image.Rectangle{Max: dims.Size}
		paint.FillShape(gtx.Ops, color.NRGBA{R: 0x2A, G: 0x2A, B: 0x2A, A: 0xFF},
			clip.UniformRRect(rect, radius).Op(gtx.Ops))
		paint.FillShape(gtx.Ops, colBorder, clip.Stroke{
			Path:  clip.UniformRRect(rect, radius).Path(gtx.Ops),
			Width: float32(gtx.Dp(1)),
		}.Op())
		call.Add(gtx.Ops)
		return dims
	})
}

// _ silences unused import warning for pointer when the file is partially used.
var _ = pointer.CursorPointer
