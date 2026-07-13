package widgets

import (
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/widget/material"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/themes"
)

func NewTitleBarWidget(theme *material.Theme, title string) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Horizontal, Alignment: layout.End}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				label := material.H6(theme, title)
				label.Color = themes.ColorsBase.Accent
				label.Font = font.Font{Weight: font.SemiBold}
				return label.Layout(gtx)
			}))
	}
}
