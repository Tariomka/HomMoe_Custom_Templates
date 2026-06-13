package widgets

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/themes"
)

// NewTextboxWidget returns a Widget that renders a text editor box
func NewTextboxWidget(theme *material.Theme, textEditor *widget.Editor, hint string) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		macro := op.Record(gtx.Ops)
		inset := layout.Inset{Top: unit.Dp(6), Bottom: unit.Dp(6), Left: unit.Dp(8), Right: unit.Dp(8)}
		dims := inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			textEditorStyle := material.Editor(theme, textEditor, hint)
			textEditorStyle.Color = themes.ColorText
			textEditorStyle.HintColor = themes.ColorTextDim
			textEditorStyle.TextSize = unit.Sp(13)
			return textEditorStyle.Layout(gtx)
		})
		call := macro.Stop()
		radius := gtx.Dp(2)
		rect := image.Rectangle{Max: dims.Size}
		paint.FillShape(gtx.Ops, themes.ColorInput, clip.UniformRRect(rect, radius).Op(gtx.Ops))
		paint.FillShape(gtx.Ops, themes.ColorBorder, clip.Stroke{
			Path:  clip.UniformRRect(rect, radius).Path(gtx.Ops),
			Width: float32(gtx.Dp(1)),
		}.Op())
		call.Add(gtx.Ops)
		return dims
	}
}
