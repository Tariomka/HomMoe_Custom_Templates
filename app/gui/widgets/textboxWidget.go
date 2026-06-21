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

// NewTextboxWidget returns a Widget that renders a text editor box
func NewTextboxWidget(theme *material.Theme, textEditor *widget.Editor, hint string) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		macro := op.Record(gtx.Ops)

		textEditor := material.Editor(theme, textEditor, hint)
		textEditor.Color = themes.ColorText
		textEditor.HintColor = themes.ColorTextDim
		textEditor.TextSize = unit.Sp(12)
		dims := layout.UniformInset(constants.DefaultPadding-2).Layout(gtx, textEditor.Layout)

		call := macro.Stop()
		radius := gtx.Dp(constants.DefaultRoundness)
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
