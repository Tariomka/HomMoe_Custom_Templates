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
	"github.com/Tariomka/hommoe_custom_templates/internal/models/preview"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/preview_service"
)

// PreviewPanel holds the layout cache + buttons for the preview panel.
type PreviewPanel struct {
	btnGenerate     widget.Clickable
	btnSaveTemplate widget.Clickable
	btnPickOutput   widget.Clickable
	btnRevealOutput widget.Clickable

	state         *drivers.State
	layoutService *preview_service.PreviewLayoutService
}

func NewPreviewPanel(state *drivers.State) *PreviewPanel {
	return &PreviewPanel{state: state, layoutService: preview_service.NewPreviewLayoutService()}
}

func (this *PreviewPanel) GetPanelWidget(theme *material.Theme) layout.Widget {
	return widgets.NewPanelWidget(constants.DefaultPadding, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
					layout.Flexed(0.5, this.getHeaderWidget(theme)),
					layout.Flexed(0.5, this.getTemplateNameWidget(theme)))
			}),
			layout.Flexed(1, this.getPreviewCanvasWidget(theme)),
			layout.Rigid(this.getStatusMessageWidget(theme)),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
					layout.Flexed(1, widgets.NewLabelBigWidget(theme, "Output directory:", themes.ColorText)),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
							layout.Rigid(widgets.NewButtonWidget(theme, "Browse", &this.btnPickOutput, false)),
							layout.Rigid(widgets.NewHorizontalSpacerWidget(6)),
							layout.Rigid(widgets.NewButtonWidget(theme, "Reveal", &this.btnRevealOutput, false)),
						)
					}),
				)
			}),
			layout.Rigid(widgets.NewVerticalSpacerWidget(8)),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Flexed(1, this.state.GetOutputPathWidget(theme)))
			}),
			layout.Rigid(widgets.NewVerticalSpacerWidget(8)),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Flexed(0.5, widgets.NewBrightButtonLargeWidget(theme, "Generate", &this.btnGenerate, false)),
					widgets.NewDefaultComponentSpacer(),
					layout.Flexed(0.5, widgets.NewBrightButtonLargeWidget(
						theme, "Save Template", &this.btnSaveTemplate, this.state.GetLastTemplate() == nil)))
			}),
		)
	})
}

func (this *PreviewPanel) HandleClicks(gtx layout.Context) {
	if this.btnGenerate.Clicked(gtx) {
		this.state.Generate()
	}
	if this.btnSaveTemplate.Clicked(gtx) {
		this.state.SaveTemplate()
	}
	if this.btnPickOutput.Clicked(gtx) {
		this.state.PickOutputDir()
	}
	if this.btnRevealOutput.Clicked(gtx) {
		this.state.RevealOutputDir()
	}
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
	label.Truncator = "..."
	label.Alignment = text.End
	return label.Layout
}

func (this *PreviewPanel) getStatusMessageWidget(theme *material.Theme) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		reserved := gtx.Dp(unit.Dp(48))
		gtx.Constraints.Min.Y = reserved
		gtx.Constraints.Max.Y = reserved

		message, isError := this.state.GetStatus()
		if message == "" {
			return layout.Dimensions{Size: image.Pt(gtx.Constraints.Max.X, reserved)}
		}

		label := material.Caption(theme, message)
		label.MaxLines = 3
		label.Alignment = text.Middle
		label.Color = themes.ColorTextDim
		if isError {
			label.Color = themes.ColorError
		}

		return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{Bottom: constants.DefaultPaddingSmall}.Layout(gtx, label.Layout)
			}),
		)
	}
}

// getPreviewCanvasWidget draws the preview area. Always fills the available
// space (so the surrounding panel keeps a consistent size) and renders an
// informational message inside the canvas when there is no template or no
// computed layout yet.
func (this *PreviewPanel) getPreviewCanvasWidget(theme *material.Theme) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		outerCanvasSize, canvasSize := getPreviewCanvasSizes(gtx)

		defer op.Offset(image.Pt((gtx.Constraints.Max.X-canvasSize.X)/2, (gtx.Constraints.Max.Y-canvasSize.Y)/2)).
			Push(gtx.Ops).Pop() // Center canvas
		renderPreviewCanvas(gtx, canvasSize)
		this.renderPreviewLegend(gtx, theme, canvasSize)

		template := this.state.GetLastTemplate()
		if template == nil {
			return widgets.NewCenteredMessageWidget(
				theme, "Adjust the options to generate the map layout.", canvasSize, outerCanvasSize)(gtx)
		}

		previewLayout := this.layoutService.BuildPreviewLayout(
			template, this.state.GetStateData().Topology, float64(canvasSize.X))
		if len(previewLayout.Positions) == 0 {
			return widgets.NewCenteredMessageWidget(theme, template.Name, canvasSize, outerCanvasSize)(gtx)
		}

		renderPreviewTemplate(gtx, theme, previewLayout)

		return layout.Dimensions{Size: outerCanvasSize}
	}
}

// getPreviewCanvasSizes returns the full available area and the centered square
// canvas (clamped to a small minimum) that the preview is drawn into.
func getPreviewCanvasSizes(gtx layout.Context) (outerCanvasSize, innerCanvasSize image.Point) {
	outerCanvasSize = gtx.Constraints.Max
	maxSize := max(min(outerCanvasSize.X, outerCanvasSize.Y), 80)
	innerCanvasSize = image.Pt(maxSize, maxSize)
	return outerCanvasSize, innerCanvasSize
}

// renderPreviewCanvas paints the canvas backdrop and its rounded frame.
func renderPreviewCanvas(gtx layout.Context, canvasSize image.Point) {
	paint.FillShape(gtx.Ops, themes.ColorPreviewBg, clip.Rect(image.Rectangle{Max: canvasSize}).Op())
	paint.FillShape(gtx.Ops, themes.ColorPreviewFrame, clip.Stroke{
		Path: clip.UniformRRect(
			image.Rect(0, 0, canvasSize.X, canvasSize.Y),
			gtx.Dp(constants.DefaultRoundnessLarge)).Path(gtx.Ops),
		Width: 2,
	}.Op())
}

// renderPreviewLegend draws the colour legend centered underneath the canvas.
func (this *PreviewPanel) renderPreviewLegend(gtx layout.Context, theme *material.Theme, canvasSize image.Point) {
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

// renderPreviewTemplate draws the computed preview layout: connections beneath
// zones, non-spawn zones first, then spawn zones on top.
func renderPreviewTemplate(gtx layout.Context, theme *material.Theme, previewLayout preview.Layout) {
	for _, connection := range previewLayout.Connections {
		utils.DrawConnection(gtx, connection, previewLayout.ZoneRadius)
	}
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

func (this *PreviewPanel) getLegendWidget(theme *material.Theme) layout.Widget {
	renderRow := func(items []constants.LegendItem) layout.Widget {
		return func(gtx layout.Context) layout.Dimensions {
			children := []layout.FlexChild{}
			for i, item := range items {
				if i > 0 {
					children = append(children, widgets.NewDefaultComponentSpacer())
				}
				children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							if item.Line {
								width := gtx.Dp(constants.DefaultRoundnessOverlineText) * 2
								height := gtx.Dp(unit.Dp(2))
								rect := image.Rect(0, 0, width, height)
								paint.FillShape(gtx.Ops, item.Color, clip.Rect(rect).Op())
								return layout.Dimensions{Size: rect.Max}
							}

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

	return func(gtx layout.Context) layout.Dimensions {
		rows := []layout.FlexChild{}
		for i, row := range constants.LegendRows {
			if i > 0 {
				rows = append(rows, layout.Rigid(widgets.NewVerticalSpacerWidget(4)))
			}
			rows = append(rows, layout.Rigid(renderRow(row)))
		}
		return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx, rows...)
	}
}
