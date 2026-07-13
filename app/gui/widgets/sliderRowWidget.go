package widgets

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

// NewSliderRowWidget returns a Widget that renders a labeled row containing a
// slider followed by its formatted value, combining [NewLabeledRowWidget] and
// [NewLabeledSliderWidget]. The format function receives the slider's current
// value at layout time.
func NewSliderRowWidget(
	theme *material.Theme,
	label string,
	labelWidthPixels int,
	slider *widget.Float,
	format func(value float32) string) layout.Widget {
	return NewLabeledRowWidget(theme, label, labelWidthPixels,
		NewLabeledSliderWidget(theme, slider, format(slider.Value)))
}
