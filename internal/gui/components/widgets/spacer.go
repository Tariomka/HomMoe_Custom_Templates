package widgets

import (
	"gioui.org/layout"
	"gioui.org/unit"
)

func NewVerticalSpacerWidget(height float32) layout.Widget {
	return layout.Spacer{Height: unit.Dp(height)}.Layout
}

func NewHorizontalSpacerWidget(width float32) layout.Widget {
	return layout.Spacer{Width: unit.Dp(width)}.Layout
}
