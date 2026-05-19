package components

import (
	"image"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/Tariomka/hommoe_custom_templates/internal/gui/components/themes"
)

type Tab struct {
	name   string
	button widget.Clickable

	panel      IPanel
	isSelected bool
}

func NewTab(name string, panel IPanel) *Tab {
	return &Tab{name: name, panel: panel}
}

func (this *Tab) SetSelected(selected bool) {
	this.isSelected = selected
}

func (this *Tab) IsTabClicked(gtx layout.Context) bool {
	return this.button.Clicked(gtx)
}

func (this *Tab) GetWidget(theme *material.Theme) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		return material.Clickable(gtx, &this.button, func(gtx layout.Context) layout.Dimensions {
			macro := op.Record(gtx.Ops)
			dims := layout.Inset{Top: unit.Dp(7), Bottom: unit.Dp(7), Left: unit.Dp(20), Right: unit.Dp(20)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				label := material.Body1(theme, this.name)
				label.TextSize = unit.Sp(13)
				label.Alignment = text.Middle
				if this.isSelected {
					label.Color = themes.ColorGold
					label.Font = font.Font{Weight: font.SemiBold}
				} else {
					label.Color = themes.ColorTextDim
				}
				return label.Layout(gtx)
			})
			call := macro.Stop()
			bgColor := themes.ColorInput
			border := themes.ColorBorder
			if this.isSelected {
				bgColor = themes.ColorPanel
				border = themes.ColorGold
			}
			rect := image.Rectangle{Max: dims.Size}
			radius := gtx.Dp(4)
			paint.FillShape(gtx.Ops, bgColor, clip.RRect{Rect: rect, NE: radius, NW: radius}.Op(gtx.Ops))
			paint.FillShape(gtx.Ops, border, clip.Stroke{
				Path:  clip.RRect{Rect: rect, NE: radius, NW: radius}.Path(gtx.Ops),
				Width: float32(gtx.Dp(1)),
			}.Op())
			call.Add(gtx.Ops)
			return dims
		})
	}
}

func (this *Tab) GetPanelWidget(theme *material.Theme) layout.Widget {
	if this.isSelected {
		return this.panel.GetPanelWidget(theme)
	}

	return func(gtx layout.Context) layout.Dimensions { return layout.Dimensions{} }
}

func (this *Tab) LoadFromState() {
	this.panel.LoadFromState()
}

func (this *Tab) SaveToState() {
	this.panel.SaveToState()
}
