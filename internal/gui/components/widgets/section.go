package widgets

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

// NewSectionWidget returns a Widget that renders a group of widgets in a bordered panel under a header
func NewSectionWidget(theme *material.Theme, title string, rows []layout.Widget) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{Bottom: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(NewSectionHeaderWidget(theme, title)),
				layout.Rigid(NewPanelWidget(unit.Dp(8), func(gtx layout.Context) layout.Dimensions {
					children := make([]layout.FlexChild, 0, len(rows)*2)
					for i, rowWidget := range rows {
						if i > 0 {
							children = append(children, layout.Rigid(layout.Spacer{Height: unit.Dp(4)}.Layout))
						}
						children = append(children, layout.Rigid(rowWidget))
					}
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx, children...)
				}),
				),
			)
		})
	}
}
