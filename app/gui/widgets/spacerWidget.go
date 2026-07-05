package widgets

import (
	"gioui.org/layout"
	"gioui.org/unit"
)

// NewVerticalSpacerWidget returns a Widget that renders a vertical spacer.
// Height is device independent pixel count.
func NewVerticalSpacerWidget(height float32) layout.Widget {
	return layout.Spacer{Height: unit.Dp(height)}.Layout
}

// NewHorizontalSpacerWidget returns a Widget that renders a horizontal spacer.
// Width is device independent pixel count.
func NewHorizontalSpacerWidget(width float32) layout.Widget {
	return layout.Spacer{Width: unit.Dp(width)}.Layout
}

// NewDefaultWidgetSpacer returns a [layout.Rigid] for horizontal widget separation.
// The separation is 12 pixels wide.
func NewDefaultWidgetSpacer() layout.FlexChild {
	return layout.Rigid(NewHorizontalSpacerWidget(12))
}

// NewDefaultComponentSpacer returns a [layout.Rigid] for horizontal component separation.
// The separation is 8 pixels wide.
func NewDefaultComponentSpacer() layout.FlexChild {
	return layout.Rigid(NewHorizontalSpacerWidget(8))
}
