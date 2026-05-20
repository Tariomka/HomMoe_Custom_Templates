package content

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/Tariomka/hommoe_custom_templates/internal/gui/components/themes"
)

// segmentButton represents a single selectable option in a row of segments.
type segmentButton struct {
	label  string
	button widget.Clickable
}

func newSegmentButton(label string) *segmentButton {
	return &segmentButton{label: label}
}

// SegmentButtonGroup is a horizontal row of mutually-exclusive segment buttons.
type SegmentButtonGroup struct {
	buttons       []*segmentButton
	selectedIndex int
}

func NewSegmentButtonGroup(labels []string) *SegmentButtonGroup {
	buttons := make([]*segmentButton, len(labels))
	for i, label := range labels {
		buttons[i] = newSegmentButton(label)
	}
	return &SegmentButtonGroup{buttons: buttons}
}

func (this *SegmentButtonGroup) SetSelectedIndex(index int) {
	if index >= 0 && index < len(this.buttons) {
		this.selectedIndex = index
	}
}

// GetSelectedIndex returns the index of the currently-selected segment.
func (this *SegmentButtonGroup) GetSelectedIndex() int {
	return this.selectedIndex
}

// Update returns true if the selection changed this frame.
func (this *SegmentButtonGroup) Update(gtx layout.Context) bool {
	changed := false
	for i, button := range this.buttons {
		if button.button.Clicked(gtx) && this.selectedIndex != i {
			this.selectedIndex = i
			changed = true
		}
	}
	return changed
}

func (this *SegmentButtonGroup) Layout(gtx layout.Context, theme *material.Theme) layout.Dimensions {
	this.Update(gtx)
	children := make([]layout.FlexChild, 0, len(this.buttons))
	for i, button := range this.buttons {
		children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return material.Clickable(gtx, &button.button, func(gtx layout.Context) layout.Dimensions {
				return layout.UniformInset(unit.Dp(2)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return drawSegment(gtx, theme, button.label, this.selectedIndex == i)
				})
			})
		}))
	}
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx, children...)
}
func drawSegment(gtx layout.Context, theme *material.Theme, label string, selected bool) layout.Dimensions {
	bgColor := themes.ColorInput
	fgColor := themes.ColorTextDim
	border := themes.ColorBorder
	if selected {
		bgColor = themes.ColorGenerate
		fgColor = themes.ColorGoldBright
		border = themes.ColorGold
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
