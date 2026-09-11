package widgets

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/constants"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/themes"
)

// NewLegendWidget stacks the given legend rows, each a horizontal run of
// swatch-and-label pairs.
func NewLegendWidget(theme *material.Theme, rows [][]constants.LegendItem) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		children := []layout.FlexChild{}
		for index, row := range rows {
			if index > 0 {
				children = append(children, layout.Rigid(NewVerticalSpacerWidget(4)))
			}
			children = append(children, layout.Rigid(NewLegendRowWidget(theme, row)))
		}
		return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx, children...)
	}
}

// NewLegendRowWidget draws one row of legend entries: a line entry gets a short
// stroke swatch, anything else a round dot.
func NewLegendRowWidget(theme *material.Theme, items []constants.LegendItem) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		children := []layout.FlexChild{}
		for index, item := range items {
			if index > 0 {
				children = append(children, NewDefaultComponentSpacer())
			}
			children = append(children, layout.Rigid(newLegendItemWidget(theme, item)))
		}
		return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx, children...)
	}
}

func newLegendItemWidget(theme *material.Theme, item constants.LegendItem) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
			layout.Rigid(newLegendSwatchWidget(item)),
			layout.Rigid(NewHorizontalSpacerWidget(4)),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				label := material.Overline(theme, item.Label)
				label.Color = themes.ColorsBase.TextDim
				return label.Layout(gtx)
			}),
		)
	}
}

func newLegendSwatchWidget(item constants.LegendItem) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		if item.Line {
			rect := image.Rect(
				0, 0,
				gtx.Dp(unit.Dp(20)), // width
				gtx.Dp(unit.Dp(2)))  // height
			paint.FillShape(gtx.Ops, item.Color, clip.Rect(rect).Op())
			return layout.Dimensions{Size: rect.Max}
		}

		side := gtx.Dp(unit.Dp(10))
		rect := image.Rect(0, 0, side, side)
		paint.FillShape(gtx.Ops, item.Color, clip.UniformRRect(rect, side/2).Op(gtx.Ops))
		return layout.Dimensions{Size: rect.Max}
	}
}
