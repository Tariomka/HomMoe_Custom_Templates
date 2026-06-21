package widgets

import (
	"gioui.org/layout"
	"gioui.org/widget/material"
)

type ThemelessWidget = func(theme *material.Theme) layout.Widget

// NewHorizontallySplitWidget returns a Widget that renders multiple widgets in a horizontal split layout
func NewHorizontallySplitWidget(theme *material.Theme, columns ...ThemelessWidget) layout.Widget {
	flexBase := layout.Flex{Axis: layout.Horizontal}
	if len(columns) == 0 {
		return func(gtx layout.Context) layout.Dimensions { return flexBase.Layout(gtx) }
	}

	if len(columns) == 1 {
		return func(gtx layout.Context) layout.Dimensions {
			return flexBase.Layout(gtx, layout.Rigid(columns[0](theme)))
		}
	}

	splitRatio := 1.0 / float32(len(columns))
	children := []layout.FlexChild{layout.Flexed(splitRatio, columns[0](theme))}

	for _, column := range columns[1:] {
		children = append(children, NewDefaultWidgetSpacer(), layout.Flexed(splitRatio, column(theme)))
	}

	return func(gtx layout.Context) layout.Dimensions { return flexBase.Layout(gtx, children...) }
}
