package panels

import (
	"fmt"
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
	"github.com/Tariomka/hommoe_custom_templates/app/gui/constants"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/drivers"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/themes"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/utils"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/widgets"
	"github.com/Tariomka/hommoe_custom_templates/internal/services"
)

// PreviewPanel holds the layout cache + buttons for the preview panel.
type PreviewPanel struct {
	btnRefresh  widget.Clickable
	pngStatus   string
	pngStatusOK bool

	state *drivers.State
}

func NewPreviewPanel(state *drivers.State) *PreviewPanel {
	return &PreviewPanel{state: state}
}

func (this *PreviewPanel) GetPanelWidget(theme *material.Theme) layout.Widget {
	return widgets.NewPanelWidget(constants.DefaultPadding, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
					layout.Rigid(this.getHeaderWidget(theme)),
					layout.Flexed(1, this.getTemplateNameWidget(theme)))
			}),
			layout.Flexed(1, this.getPreviewCanvasWidget(theme)),
			layout.Rigid(widgets.NewVerticalSpacerWidget(6)),
			layout.Rigid(this.getStatusMessageWidget(theme)),
		)
	})
}

func (this *PreviewPanel) HandleClicks(gtx layout.Context) {
	// TODO: add handling for generation and saving after footer is moved to the preview panel
}

func (this *PreviewPanel) getHeaderWidget(theme *material.Theme) layout.Widget {
	label := material.Body1(theme, "Preview")
	label.Color = themes.ColorAccent
	label.Font = font.Font{Weight: font.SemiBold}
	return label.Layout
}

func (this *PreviewPanel) getTemplateNameWidget(theme *material.Theme) layout.Widget {
	template := this.state.GetLastTemplate()

	name := "(no template generated yet)"
	if template != nil {
		name = fmt.Sprintf("Name: '%s'", template.Name)
	}

	label := material.Overline(theme, name)
	label.Color = themes.ColorTextDim
	label.MaxLines = 1
	label.Truncator = "…"
	label.Alignment = text.End
	return label.Layout
}

func (this *PreviewPanel) getStatusMessageWidget(theme *material.Theme) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		// Reserve a fixed-height slot so the canvas above doesn't shift
		// when the status message appears/disappears or grows from 1 to
		// 2 lines. Height is sized for 2 lines of 11sp text + the
		// bottom inset used below the legend separator.
		reserved := gtx.Dp(unit.Dp(34))
		gtx.Constraints.Min.Y = reserved
		gtx.Constraints.Max.Y = reserved
		if this.pngStatus == "" {
			return layout.Dimensions{Size: image.Pt(gtx.Constraints.Max.X, reserved)}
		}

		label := material.Overline(theme, this.pngStatus)
		label.MaxLines = 2
		label.Alignment = text.Middle
		label.Color = themes.ColorTextDim
		if !this.pngStatusOK {
			label.Color = themes.ColorError
		}
		return layout.Inset{Bottom: constants.DefaultPaddingSmall}.Layout(gtx, label.Layout)
	}
}

// layoutPreviewCanvas draws the preview area. Always fills the available
// space (so the surrounding panel keeps a consistent size) and renders an
// informational message inside the canvas when there is no template or no
// computed layout yet.
func (this *PreviewPanel) getPreviewCanvasWidget(theme *material.Theme) layout.Widget {
	getCanvasSizes := func(gtx layout.Context) (outerCanvasSize, innerCanvasSize image.Point) {
		maxX := gtx.Constraints.Max.X
		maxY := gtx.Constraints.Max.Y
		outerCanvasSize = image.Pt(maxX, maxY)
		maxSize := max(min(outerCanvasSize.X, outerCanvasSize.Y), 80)
		innerCanvasSize = image.Pt(maxSize, maxSize)
		return
	}

	renderCanvas := func(gtx layout.Context, canvasSize image.Point) {
		paint.FillShape(gtx.Ops, themes.ColorPreviewBg, clip.Rect(image.Rectangle{Max: canvasSize}).Op())
		paint.FillShape(gtx.Ops, themes.ColorPreviewFrame, clip.Stroke{
			Path: clip.UniformRRect(
				image.Rect(0, 0, canvasSize.X, canvasSize.Y),
				gtx.Dp(constants.DefaultRoundnessLarge)).Path(gtx.Ops),
			Width: 2,
		}.Op())
	}

	renderLegend := func(gtx layout.Context, canvasSize image.Point) {
		legendMacro := op.Record(gtx.Ops)
		contextCopy := gtx
		contextCopy.Constraints.Min = image.Point{}
		contextCopy.Constraints.Max.X = canvasSize.X
		legendWidget := this.getLegendWidget(theme)(contextCopy)
		legendCall := legendMacro.Stop()
		legendOffset := op.Offset(image.Point{
			X: (canvasSize.X - legendWidget.Size.X) / 2,
			Y: canvasSize.Y + gtx.Dp(constants.DefaultPaddingSmall),
		}).Push(gtx.Ops)
		legendCall.Add(gtx.Ops)
		legendOffset.Pop()
	}

	renderTemplate := func(gtx layout.Context, previewLayout services.PreviewLayout) {
		// Connections beneath zones.
		for _, connection := range previewLayout.Connections {
			utils.DrawConnection(gtx, connection, previewLayout.ZoneRadius)
		}
		// Non-spawn zones first, then spawn zones on top.
		for _, zone := range previewLayout.Zones {
			if !zone.IsPlayer {
				utils.DrawPreviewZone(gtx, theme, zone, previewLayout.ZoneRadius)
			}
		}
		for _, zone := range previewLayout.Zones {
			if zone.IsPlayer {
				utils.DrawPreviewZone(gtx, theme, zone, previewLayout.ZoneRadius)
			}
		}
	}

	return func(gtx layout.Context) layout.Dimensions {
		outerCanvasSize, canvasSize := getCanvasSizes(gtx)

		// Center canvas
		defer op.Offset(image.Point{
			X: (gtx.Constraints.Max.X - canvasSize.X) / 2,
			Y: (gtx.Constraints.Max.Y - canvasSize.Y) / 2,
		}).Push(gtx.Ops).Pop()
		renderCanvas(gtx, canvasSize)
		renderLegend(gtx, canvasSize)

		template := this.state.GetLastTemplate()
		if template == nil {
			return widgets.NewCenteredMessageWidget(
				theme, "Adjust the options to generate the map layout.", canvasSize, outerCanvasSize)(gtx)
		}

		previewLayout := services.BuildPreviewLayout(template, this.state.GetStateData().Topology, float64(canvasSize.X))
		if len(previewLayout.Positions) == 0 {
			return widgets.NewCenteredMessageWidget(theme, template.Name, canvasSize, outerCanvasSize)(gtx)
		}

		renderTemplate(gtx, previewLayout)

		return layout.Dimensions{Size: outerCanvasSize}
	}
}

func (this *PreviewPanel) getLegendWidget(theme *material.Theme) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		children := []layout.FlexChild{}
		for i, item := range constants.LegendItems {
			if i > 0 {
				children = append(children, layout.Rigid(widgets.NewHorizontalSpacerWidget(8)))
			}
			children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						side := gtx.Dp(constants.DefaultRoundnessOverlineText)
						rect := image.Rect(0, 0, side, side)
						paint.FillShape(gtx.Ops, item.Color, clip.UniformRRect(rect, side/2).Op(gtx.Ops))
						return layout.Dimensions{Size: rect.Max}
					}),
					layout.Rigid(widgets.NewHorizontalSpacerWidget(4)),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						label := material.Overline(theme, item.Label)
						label.Color = themes.ColorTextDim
						return label.Layout(gtx)
					}),
				)
			}))
		}
		return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx, children...)
	}
}
