package gui

import (
	"image"
	"image/color"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

// fillBackground paints the given color over gtx.Constraints.Max.
func fillBackground(gtx layout.Context, c color.NRGBA) {
	rect := image.Rectangle{Max: gtx.Constraints.Max}
	paint.FillShape(gtx.Ops, c, clip.Rect(rect).Op())
}

// borderedPanel wraps content with a rounded border on a panel background.
func borderedPanel(gtx layout.Context, padding unit.Dp, w layout.Widget) layout.Dimensions {
	radius := gtx.Dp(4)
	macro := op.Record(gtx.Ops)
	dims := layout.UniformInset(padding).Layout(gtx, w)
	call := macro.Stop()
	rect := image.Rectangle{Max: dims.Size}
	paint.FillShape(gtx.Ops, colPanel, clip.UniformRRect(rect, radius).Op(gtx.Ops))
	paint.FillShape(gtx.Ops, colBorder, clip.Stroke{
		Path:  clip.UniformRRect(rect, radius).Path(gtx.Ops),
		Width: float32(gtx.Dp(1)),
	}.Op())
	call.Add(gtx.Ops)
	return dims
}

// warnBanner draws a yellow warning banner with the given text.
func warnBanner(gtx layout.Context, th *material.Theme, txt string) layout.Dimensions {
	macro := op.Record(gtx.Ops)
	dims := layout.UniformInset(unit.Dp(8)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		lbl := material.Body2(th, txt)
		lbl.Color = colWarnText
		lbl.TextSize = unit.Sp(11)
		return lbl.Layout(gtx)
	})
	call := macro.Stop()
	radius := gtx.Dp(3)
	rect := image.Rectangle{Max: dims.Size}
	paint.FillShape(gtx.Ops, colWarnBg, clip.UniformRRect(rect, radius).Op(gtx.Ops))
	paint.FillShape(gtx.Ops, color.NRGBA{R: 0x6A, G: 0x50, B: 0x20, A: 0xFF}, clip.Stroke{
		Path:  clip.UniformRRect(rect, radius).Path(gtx.Ops),
		Width: float32(gtx.Dp(1)),
	}.Op())
	call.Add(gtx.Ops)
	return dims
}

// sectionHeader draws an accented gold section title (with diamond bullet).
func sectionHeader(gtx layout.Context, th *material.Theme, txt string) layout.Dimensions {
	lbl := material.Body1(th, "◆  "+txt)
	lbl.Color = colGold
	lbl.Font = font.Font{Weight: font.SemiBold}
	lbl.TextSize = unit.Sp(13)
	return lbl.Layout(gtx)
}

// dimLabel renders a small dimmed description line.
func dimLabel(gtx layout.Context, th *material.Theme, txt string) layout.Dimensions {
	lbl := material.Body2(th, txt)
	lbl.Color = colTextDim
	lbl.TextSize = unit.Sp(12)
	return lbl.Layout(gtx)
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

// segmentButton represents a single selectable option in a row of segments.
type segmentButton struct {
	Label string
	click widget.Clickable
}

// segmentGroup is a horizontal row of mutually-exclusive segment buttons.
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

// goldButton renders a Generate-style emphasised button.
type goldButton struct {
	Text     string
	Click    *widget.Clickable
	Disabled bool
}

func (b goldButton) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	if b.Disabled {
		gtx = gtx.Disabled()
	}
	return material.Clickable(gtx, b.Click, func(gtx layout.Context) layout.Dimensions {
		macro := op.Record(gtx.Ops)
		dims := layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				lbl := material.Body1(th, b.Text)
				lbl.Color = colGoldBright
				lbl.TextSize = unit.Sp(14)
				lbl.Font = font.Font{Weight: font.SemiBold}
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
	Text     string
	Click    *widget.Clickable
	Disabled bool
}

func (b toolbarButton) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	if b.Disabled {
		gtx = gtx.Disabled()
	}
	return material.Clickable(gtx, b.Click, func(gtx layout.Context) layout.Dimensions {
		macro := op.Record(gtx.Ops)
		dims := layout.Inset{Top: unit.Dp(5), Bottom: unit.Dp(5), Left: unit.Dp(10), Right: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				lbl := material.Body2(th, b.Text)
				lbl.Color = colText
				lbl.TextSize = unit.Sp(12)
				if b.Disabled {
					lbl.Color = colTextDim
				}
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

// comboBox is an inline-expanding dropdown selector. When closed it shows the
// selected option and a ▾ chevron; clicking expands the list of options below,
// pushing surrounding content downward. Selecting a row collapses it.
type comboBox struct {
	Items    []string
	Selected int
	Open     bool

	toggle widget.Clickable
	rows   []widget.Clickable
}

func newComboBox(items []string) *comboBox {
	return &comboBox{Items: append([]string(nil), items...), rows: make([]widget.Clickable, len(items))}
}

// SetItems replaces the option list and resets selection bounds.
func (c *comboBox) SetItems(items []string) {
	c.Items = append([]string(nil), items...)
	c.rows = make([]widget.Clickable, len(items))
	if c.Selected >= len(items) {
		c.Selected = 0
	}
}

// SelectByName sets Selected to the index whose label matches name.
func (c *comboBox) SelectByName(name string) bool {
	for i, it := range c.Items {
		if it == name {
			c.Selected = i
			return true
		}
	}
	return false
}

// Value returns the currently selected option label, or "" if empty.
func (c *comboBox) Value() string {
	if c.Selected >= 0 && c.Selected < len(c.Items) {
		return c.Items[c.Selected]
	}
	return ""
}

// Update returns true if the selection changed this frame.
func (c *comboBox) Update(gtx layout.Context) bool {
	changed := false
	if c.toggle.Clicked(gtx) {
		c.Open = !c.Open
	}
	for i := range c.rows {
		if i >= len(c.Items) {
			break
		}
		if c.rows[i].Clicked(gtx) {
			if c.Selected != i {
				c.Selected = i
				changed = true
			}
			c.Open = false
		}
	}
	return changed
}

func (c *comboBox) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	c.Update(gtx)
	flex := layout.Flex{Axis: layout.Vertical}
	children := []layout.FlexChild{
		layout.Rigid(func(gtx layout.Context) layout.Dimensions { return c.layoutTrigger(gtx, th) }),
	}
	if c.Open && len(c.Items) > 0 {
		children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Top: unit.Dp(2)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return c.layoutList(gtx, th)
			})
		}))
	}
	return flex.Layout(gtx, children...)
}

func (c *comboBox) layoutTrigger(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return material.Clickable(gtx, &c.toggle, func(gtx layout.Context) layout.Dimensions {
		macro := op.Record(gtx.Ops)
		dims := layout.Inset{Top: unit.Dp(6), Bottom: unit.Dp(6), Left: unit.Dp(8), Right: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					lbl := material.Body1(th, c.Value())
					lbl.Color = colText
					lbl.TextSize = unit.Sp(13)
					lbl.MaxLines = 1
					lbl.Truncator = "…"
					return lbl.Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					arrow := "▾"
					if c.Open {
						arrow = "▴"
					}
					lbl := material.Body1(th, arrow)
					lbl.Color = colGoldDim
					lbl.TextSize = unit.Sp(11)
					return lbl.Layout(gtx)
				}),
			)
		})
		call := macro.Stop()
		radius := gtx.Dp(2)
		rect := image.Rectangle{Max: dims.Size}
		paint.FillShape(gtx.Ops, colInput, clip.UniformRRect(rect, radius).Op(gtx.Ops))
		border := colBorder
		if c.Open {
			border = colGold
		}
		paint.FillShape(gtx.Ops, border, clip.Stroke{
			Path:  clip.UniformRRect(rect, radius).Path(gtx.Ops),
			Width: float32(gtx.Dp(1)),
		}.Op())
		call.Add(gtx.Ops)
		return dims
	})
}

func (c *comboBox) layoutList(gtx layout.Context, th *material.Theme) layout.Dimensions {
	macro := op.Record(gtx.Ops)
	rows := make([]layout.FlexChild, 0, len(c.Items))
	for i := range c.Items {
		i := i
		rows = append(rows, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return material.Clickable(gtx, &c.rows[i], func(gtx layout.Context) layout.Dimensions {
				return drawComboRow(gtx, th, c.Items[i], i == c.Selected)
			})
		}))
	}
	dims := layout.Flex{Axis: layout.Vertical}.Layout(gtx, rows...)
	call := macro.Stop()
	radius := gtx.Dp(2)
	rect := image.Rectangle{Max: dims.Size}
	paint.FillShape(gtx.Ops, color.NRGBA{R: 0x2C, G: 0x26, B: 0x19, A: 0xFF},
		clip.UniformRRect(rect, radius).Op(gtx.Ops))
	paint.FillShape(gtx.Ops, colBorder, clip.Stroke{
		Path:  clip.UniformRRect(rect, radius).Path(gtx.Ops),
		Width: float32(gtx.Dp(1)),
	}.Op())
	call.Add(gtx.Ops)
	return dims
}

func drawComboRow(gtx layout.Context, th *material.Theme, label string, selected bool) layout.Dimensions {
	macro := op.Record(gtx.Ops)
	dims := layout.Inset{Top: unit.Dp(5), Bottom: unit.Dp(5), Left: unit.Dp(8), Right: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		lbl := material.Body1(th, label)
		lbl.Color = colText
		lbl.TextSize = unit.Sp(13)
		lbl.MaxLines = 1
		lbl.Truncator = "…"
		if selected {
			lbl.Color = colGold
			lbl.Font = font.Font{Weight: font.SemiBold}
		}
		return lbl.Layout(gtx)
	})
	call := macro.Stop()
	if selected {
		paint.FillShape(gtx.Ops, color.NRGBA{R: 0x5A, G: 0x4A, B: 0x28, A: 0xFF},
			clip.Rect{Max: dims.Size}.Op())
	}
	call.Add(gtx.Ops)
	if dims.Size.X < gtx.Constraints.Min.X {
		dims.Size.X = gtx.Constraints.Min.X
	}
	return dims
}

// tabs is a horizontal tab strip; the caller draws the active page below.
type tabs struct {
	Labels   []string
	Selected int
	clicks   []widget.Clickable
}

func newTabs(labels []string) *tabs {
	return &tabs{Labels: labels, clicks: make([]widget.Clickable, len(labels))}
}

func (t *tabs) Update(gtx layout.Context) bool {
	for i := range t.clicks {
		if t.clicks[i].Clicked(gtx) && t.Selected != i {
			t.Selected = i
			return true
		}
	}
	return false
}

func (t *tabs) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	t.Update(gtx)
	children := make([]layout.FlexChild, 0, len(t.Labels))
	for i := range t.Labels {
		i := i
		selected := i == t.Selected
		children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return material.Clickable(gtx, &t.clicks[i], func(gtx layout.Context) layout.Dimensions {
				return drawTab(gtx, th, t.Labels[i], selected)
			})
		}))
	}
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx, children...)
}

func drawTab(gtx layout.Context, th *material.Theme, label string, selected bool) layout.Dimensions {
	macro := op.Record(gtx.Ops)
	dims := layout.Inset{Top: unit.Dp(7), Bottom: unit.Dp(7), Left: unit.Dp(20), Right: unit.Dp(20)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		lbl := material.Body1(th, label)
		lbl.TextSize = unit.Sp(13)
		lbl.Alignment = text.Middle
		if selected {
			lbl.Color = colGold
			lbl.Font = font.Font{Weight: font.SemiBold}
		} else {
			lbl.Color = colTextDim
		}
		return lbl.Layout(gtx)
	})
	call := macro.Stop()
	bg := colInput
	border := colBorder
	if selected {
		bg = colPanel
		border = colGold
	}
	rect := image.Rectangle{Max: dims.Size}
	radius := gtx.Dp(4)
	paint.FillShape(gtx.Ops, bg, clip.RRect{Rect: rect, NE: radius, NW: radius}.Op(gtx.Ops))
	paint.FillShape(gtx.Ops, border, clip.Stroke{
		Path:  clip.RRect{Rect: rect, NE: radius, NW: radius}.Path(gtx.Ops),
		Width: float32(gtx.Dp(1)),
	}.Op())
	call.Add(gtx.Ops)
	return dims
}

// snapIntSlider draws a slider snapped to integer steps in [lo, hi] and
// returns the resolved integer value.
func snapIntSlider(gtx layout.Context, th *material.Theme, f *widget.Float, lo, hi int) int {
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
	sl := material.Slider(th, f)
	sl.Color = colGold
	sl.Layout(gtx)
	return rounded
}

func roundHalfAway(v float64) float64 {
	if v < 0 {
		return -roundHalfAway(-v)
	}
	return float64(int(v + 0.5))
}
