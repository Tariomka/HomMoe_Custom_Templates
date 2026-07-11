package widgets

import "gioui.org/layout"

func NewEmptyWidget() layout.Widget {
	return func(layout.Context) layout.Dimensions { return layout.Dimensions{} }
}

func NewMinimumConstraintsWidget() layout.Widget {
	return func(gtx layout.Context) layout.Dimensions { return layout.Dimensions{Size: gtx.Constraints.Min} }
}
