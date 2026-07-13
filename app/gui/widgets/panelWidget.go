package widgets

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/constants"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/themes"
)

// NewPanelWidget return a Widget that renders content in a rounded border panel.
func NewPanelWidget(padding unit.Dp, content layout.Widget) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		radius := gtx.Dp(constants.DefaultRoundness)
		macro := op.Record(gtx.Ops)
		dims := layout.UniformInset(padding).Layout(gtx, content)
		call := macro.Stop()
		rect := image.Rectangle{Max: dims.Size}
		paint.FillShape(gtx.Ops, themes.ColorsBase.Panel, clip.UniformRRect(rect, radius).Op(gtx.Ops))
		paint.FillShape(gtx.Ops, themes.ColorsBase.Border, clip.Stroke{
			Path:  clip.UniformRRect(rect, radius).Path(gtx.Ops),
			Width: float32(gtx.Dp(1)),
		}.Op())
		call.Add(gtx.Ops)
		return dims
	}
}
