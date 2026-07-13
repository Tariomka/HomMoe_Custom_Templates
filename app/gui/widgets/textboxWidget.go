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
	"github.com/Tariomka/hommoe_custom_templates/app/gui/constants"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/themes"
)

// NewTextboxWidget returns a Widget that renders a text editor box.
func NewTextboxWidget(theme *material.Theme, textEditor *widget.Editor, hint string, readonly bool) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		macro := op.Record(gtx.Ops)

		textEditor := material.Editor(theme, textEditor, hint)
		textEditor.Editor.ReadOnly = readonly
		textEditor.Color = themes.ColorsBase.Text
		textEditor.HintColor = themes.ColorsBase.TextDim
		textEditor.TextSize = unit.Sp(12)
		dims := layout.UniformInset(constants.DefaultPaddingSmall).Layout(gtx, textEditor.Layout)

		call := macro.Stop()
		radius := gtx.Dp(constants.DefaultRoundness)
		rect := image.Rectangle{Max: dims.Size}
		paint.FillShape(gtx.Ops, themes.ColorsBase.Input, clip.UniformRRect(rect, radius).Op(gtx.Ops))
		paint.FillShape(gtx.Ops, themes.ColorsBase.Border, clip.Stroke{
			Path:  clip.UniformRRect(rect, radius).Path(gtx.Ops),
			Width: float32(gtx.Dp(1)),
		}.Op())
		call.Add(gtx.Ops)
		return dims
	}
}
