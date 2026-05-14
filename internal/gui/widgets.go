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
	group := &segmentGroup{buttons: make([]segmentButton, len(labels))}
	for i, label := range labels {
		group.buttons[i].Label = label
	}
	return group
}

// Update returns true if the selection changed this frame.
func (this *segmentGroup) Update(gtx layout.Context) bool {
	changed := false
	for i := range this.buttons {
		if this.buttons[i].click.Clicked(gtx) && this.Selected != i {
			this.Selected = i
			changed = true
		}
	}
	return changed
}

func (this *segmentGroup) Layout(gtx layout.Context, theme *material.Theme) layout.Dimensions {
	this.Update(gtx)
	children := make([]layout.FlexChild, 0, len(this.buttons))
	for i := range this.buttons {
		i := i
		button := &this.buttons[i]
		selected := i == this.Selected
		children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return material.Clickable(gtx, &button.click, func(gtx layout.Context) layout.Dimensions {
				return layout.UniformInset(unit.Dp(2)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return drawSegment(gtx, theme, button.Label, selected)
				})
			})
		}))
	}
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx, children...)
}

func drawSegment(gtx layout.Context, theme *material.Theme, label string, selected bool) layout.Dimensions {
	bgColor := colInput
	fgColor := colTextDim
	border := colBorder
	if selected {
		bgColor = colGenerate
		fgColor = colGoldBright
		border = colGold
	}
	macro := op.Record(gtx.Ops)
	dims := layout.UniformInset(unit.Dp(8)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		label := material.Body2(theme, label)
		label.Color = fgColor
		label.TextSize = unit.Sp(12)
		return label.Layout(gtx)
	})
	call := macro.Stop()
	radius := gtx.Dp(3)
	rect := image.Rectangle{Max: dims.Size}
	paint.FillShape(gtx.Ops, bgColor, clip.UniformRRect(rect, radius).Op(gtx.Ops))
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

func (this goldButton) Layout(gtx layout.Context, theme *material.Theme) layout.Dimensions {
	if this.Disabled {
		gtx = gtx.Disabled()
	}
	return material.Clickable(gtx, this.Click, func(gtx layout.Context) layout.Dimensions {
		macro := op.Record(gtx.Ops)
		dims := layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				label := material.Body1(theme, this.Text)
				label.Color = colGoldBright
				label.TextSize = unit.Sp(14)
				label.Font = font.Font{Weight: font.SemiBold}
				if this.Disabled {
					label.Color = colTextDim
				}
				return label.Layout(gtx)
			})
		})
		call := macro.Stop()
		radius := gtx.Dp(3)
		rect := image.Rectangle{Max: dims.Size}
		bgColor := colGenerate
		border := colGold
		if this.Disabled {
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

// toolbarButton renders a small dark toolbar-style button.
type toolbarButton struct {
	Text     string
	Click    *widget.Clickable
	Disabled bool
}

func (this toolbarButton) Layout(gtx layout.Context, theme *material.Theme) layout.Dimensions {
	if this.Disabled {
		gtx = gtx.Disabled()
	}
	return material.Clickable(gtx, this.Click, func(gtx layout.Context) layout.Dimensions {
		macro := op.Record(gtx.Ops)
		dims := layout.Inset{Top: unit.Dp(5), Bottom: unit.Dp(5), Left: unit.Dp(10), Right: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				label := material.Body2(theme, this.Text)
				label.Color = colText
				label.TextSize = unit.Sp(12)
				if this.Disabled {
					label.Color = colTextDim
				}
				return label.Layout(gtx)
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
func (this *comboBox) SetItems(items []string) {
	this.Items = append([]string(nil), items...)
	this.rows = make([]widget.Clickable, len(items))
	if this.Selected >= len(items) {
		this.Selected = 0
	}
}

// SelectByName sets Selected to the index whose label matches name.
func (this *comboBox) SelectByName(name string) bool {
	for i, it := range this.Items {
		if it == name {
			this.Selected = i
			return true
		}
	}
	return false
}

// Value returns the currently selected option label, or "" if empty.
func (this *comboBox) Value() string {
	if this.Selected >= 0 && this.Selected < len(this.Items) {
		return this.Items[this.Selected]
	}
	return ""
}

// Update returns true if the selection changed this frame.
func (this *comboBox) Update(gtx layout.Context) bool {
	changed := false
	if this.toggle.Clicked(gtx) {
		this.Open = !this.Open
	}
	for i := range this.rows {
		if i >= len(this.Items) {
			break
		}
		if this.rows[i].Clicked(gtx) {
			if this.Selected != i {
				this.Selected = i
				changed = true
			}
			this.Open = false
		}
	}
	return changed
}

func (this *comboBox) Layout(gtx layout.Context, theme *material.Theme) layout.Dimensions {
	this.Update(gtx)
	flex := layout.Flex{Axis: layout.Vertical}
	children := []layout.FlexChild{
		layout.Rigid(func(gtx layout.Context) layout.Dimensions { return this.layoutTrigger(gtx, theme) }),
	}
	if this.Open && len(this.Items) > 0 {
		children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Top: unit.Dp(2)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return this.layoutList(gtx, theme)
			})
		}))
	}
	return flex.Layout(gtx, children...)
}

func (this *comboBox) layoutTrigger(gtx layout.Context, theme *material.Theme) layout.Dimensions {
	return material.Clickable(gtx, &this.toggle, func(gtx layout.Context) layout.Dimensions {
		macro := op.Record(gtx.Ops)
		dims := layout.Inset{Top: unit.Dp(6), Bottom: unit.Dp(6), Left: unit.Dp(8), Right: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					label := material.Body1(theme, this.Value())
					label.Color = colText
					label.TextSize = unit.Sp(13)
					label.MaxLines = 1
					label.Truncator = "…"
					return label.Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					arrow := "▾"
					if this.Open {
						arrow = "▴"
					}
					label := material.Body1(theme, arrow)
					label.Color = colGoldDim
					label.TextSize = unit.Sp(11)
					return label.Layout(gtx)
				}),
			)
		})
		call := macro.Stop()
		radius := gtx.Dp(2)
		rect := image.Rectangle{Max: dims.Size}
		paint.FillShape(gtx.Ops, colInput, clip.UniformRRect(rect, radius).Op(gtx.Ops))
		border := colBorder
		if this.Open {
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

func (this *comboBox) layoutList(gtx layout.Context, theme *material.Theme) layout.Dimensions {
	macro := op.Record(gtx.Ops)
	rows := make([]layout.FlexChild, 0, len(this.Items))
	for i := range this.Items {
		i := i
		rows = append(rows, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return material.Clickable(gtx, &this.rows[i], func(gtx layout.Context) layout.Dimensions {
				return drawComboRow(gtx, theme, this.Items[i], i == this.Selected)
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

func drawComboRow(gtx layout.Context, theme *material.Theme, label string, selected bool) layout.Dimensions {
	macro := op.Record(gtx.Ops)
	dims := layout.Inset{Top: unit.Dp(5), Bottom: unit.Dp(5), Left: unit.Dp(8), Right: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		label := material.Body1(theme, label)
		label.Color = colText
		label.TextSize = unit.Sp(13)
		label.MaxLines = 1
		label.Truncator = "…"
		if selected {
			label.Color = colGold
			label.Font = font.Font{Weight: font.SemiBold}
		}
		return label.Layout(gtx)
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

func (this *tabs) Update(gtx layout.Context) bool {
	for i := range this.clicks {
		if this.clicks[i].Clicked(gtx) && this.Selected != i {
			this.Selected = i
			return true
		}
	}
	return false
}

func (this *tabs) Layout(gtx layout.Context, theme *material.Theme) layout.Dimensions {
	this.Update(gtx)
	children := make([]layout.FlexChild, 0, len(this.Labels))
	for i := range this.Labels {
		i := i
		selected := i == this.Selected
		children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return material.Clickable(gtx, &this.clicks[i], func(gtx layout.Context) layout.Dimensions {
				return drawTab(gtx, theme, this.Labels[i], selected)
			})
		}))
	}
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx, children...)
}

func drawTab(gtx layout.Context, theme *material.Theme, label string, selected bool) layout.Dimensions {
	macro := op.Record(gtx.Ops)
	dims := layout.Inset{Top: unit.Dp(7), Bottom: unit.Dp(7), Left: unit.Dp(20), Right: unit.Dp(20)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		label := material.Body1(theme, label)
		label.TextSize = unit.Sp(13)
		label.Alignment = text.Middle
		if selected {
			label.Color = colGold
			label.Font = font.Font{Weight: font.SemiBold}
		} else {
			label.Color = colTextDim
		}
		return label.Layout(gtx)
	})
	call := macro.Stop()
	bgColor := colInput
	border := colBorder
	if selected {
		bgColor = colPanel
		border = colGold
	}
	rect := image.Rectangle{Max: dims.Size}
	radius := gtx.Dp(4)
	paint.FillShape(gtx.Ops, bgColor, clip.RRect{Rect: rect, NE: radius, NW: radius}.Op(gtx.Ops))
	paint.FillShape(gtx.Ops, border, clip.Stroke{
		Path:  clip.RRect{Rect: rect, NE: radius, NW: radius}.Path(gtx.Ops),
		Width: float32(gtx.Dp(1)),
	}.Op())
	call.Add(gtx.Ops)
	return dims
}

// snapIntSlider draws a slider snapped to integer steps in [low, high] and
// returns the resolved integer value.
func snapIntSlider(gtx layout.Context, theme *material.Theme, f *widget.Float, low, high int) int {
	value := mapRange(f.Value, float32(low), float32(high))
	rounded := int(roundHalfAway(float64(value)))
	if rounded < low {
		rounded = low
	}
	if rounded > high {
		rounded = high
	}
	target := mapRangeInv(float32(rounded), float32(low), float32(high))
	if target != f.Value && !f.Dragging() {
		f.Value = target
	}
	slider := material.Slider(theme, f)
	slider.Color = colGold
	slider.Layout(gtx)
	return rounded
}

func roundHalfAway(value float64) float64 {
	if value < 0 {
		return -roundHalfAway(-value)
	}
	return float64(int(value + 0.5))
}
