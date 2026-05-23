package widgets

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/Tariomka/hommoe_custom_templates/internal/gui/components/themes"
)

// NewCenteredMessageWidget renders a Body2 label centered inside the given canvas
// area. Uses the same material.Label approach as the (former) empty-state
// view so text renders reliably (unlike drawCenteredText for longer strings).
func NewCenteredMessageWidget(theme *material.Theme, message string, innerCanvasSize, outerCanvasSize image.Point) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		macro := op.Record(gtx.Ops)
		gtxLocal := gtx
		gtxLocal.Constraints.Min = image.Point{}
		gtxLocal.Constraints.Max = innerCanvasSize
		label := material.Body2(theme, message)
		label.Color = themes.ColorTextDim
		label.TextSize = unit.Sp(12)
		label.Alignment = text.Middle
		dims := label.Layout(gtxLocal)
		call := macro.Stop()

		tx := (innerCanvasSize.X - dims.Size.X) / 2
		ty := (innerCanvasSize.Y - dims.Size.Y) / 2
		stack := op.Offset(image.Pt(tx, ty)).Push(gtx.Ops)
		call.Add(gtx.Ops)
		stack.Pop()
		return layout.Dimensions{Size: outerCanvasSize}
	}
}
