package widgets

import (
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/widget/material"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/constants"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/themes"
)

// NewSectionWidget returns a Widget that renders a group of widgets in a bordered panel under a header.
func NewSectionWidget(theme *material.Theme, title string, rows []layout.Widget) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{Bottom: constants.DefaultPadding}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(newSectionHeaderWidget(theme, title)),
				layout.Rigid(NewPanelWidget(constants.DefaultPadding, func(gtx layout.Context) layout.Dimensions {
					children := []layout.FlexChild{}
					if len(rows) > 0 {
						children = append(children, layout.Rigid(rows[0]))
					}
					for _, rowWidget := range rows[1:] {
						children = append(children, layout.Rigid(NewVerticalSpacerWidget(4)), layout.Rigid(rowWidget))
					}
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx, children...)
				}),
				),
			)
		})
	}
}

// newSectionHeaderWidget returns a Widget that renders a section title.
func newSectionHeaderWidget(theme *material.Theme, title string) layout.Widget {
	label := material.Body2(theme, "◆  "+title)
	label.Color = themes.ColorsBase.Accent
	label.Font = font.Font{Weight: font.SemiBold}
	return label.Layout
}
