package widgets

import (
	"gioui.org/layout"
	"gioui.org/unit"
)

// NewVerticalSpacerWidget returns a Widget that renders a vertical spacer.
// Height is device independent pixel count
func NewVerticalSpacerWidget(height float32) layout.Widget {
	return layout.Spacer{Height: unit.Dp(height)}.Layout
}

// NewHorizontalSpacerWidget returns a Widget that renders a horizontal spacer.
// Width is device independent pixel count
func NewHorizontalSpacerWidget(width float32) layout.Widget {
	return layout.Spacer{Width: unit.Dp(width)}.Layout
}
