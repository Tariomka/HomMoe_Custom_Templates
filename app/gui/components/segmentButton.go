package components

import "gioui.org/widget"

type segmentButton struct {
	label  string
	button widget.Clickable
}

func newSegmentButton(label string) *segmentButton {
	return &segmentButton{label: label}
}
