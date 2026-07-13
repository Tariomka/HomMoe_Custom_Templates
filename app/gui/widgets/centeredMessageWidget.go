package widgets

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget/material"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/themes"
)

// NewCenteredMessageWidget renders a Body2 label centered inside the given canvas
// area. Uses the same material.Label approach as the (former) empty-state
// view so text renders reliably (unlike drawCenteredText for longer strings).
func NewCenteredMessageWidget(
	theme *material.Theme,
	message string,
	innerCanvasSize, outerCanvasSize image.Point) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		macro := op.Record(gtx.Ops)
		gtx.Constraints.Min = image.Point{}
		gtx.Constraints.Max = innerCanvasSize
		dims := NewLabelWidget(theme, message, themes.ColorsBase.TextDim)(gtx)
		call := macro.Stop()

		tx := (innerCanvasSize.X - dims.Size.X) / 2
		ty := (innerCanvasSize.Y - dims.Size.Y) / 2
		stack := op.Offset(image.Pt(tx, ty)).Push(gtx.Ops)
		call.Add(gtx.Ops)
		stack.Pop()
		return layout.Dimensions{Size: outerCanvasSize}
	}
}
