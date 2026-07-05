package components

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/widgets"
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

func (this *SegmentButtonGroup) GetWidget(theme *material.Theme) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		this.Update(gtx)
		children := make([]layout.FlexChild, 0, len(this.buttons))
		for i, button := range this.buttons {
			children = append(children, layout.Rigid(
				widgets.NewSegmentButtonWidget(theme, button.label, &button.button, this.selectedIndex == i)))
		}
		return layout.Flex{Axis: layout.Horizontal}.Layout(gtx, children...)
	}
}
