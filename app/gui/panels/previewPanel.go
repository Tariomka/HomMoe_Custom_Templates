package panels

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
	template := this.state.GetLastTemplate()
	return widgets.NewPanelWidget(unit.Dp(8), func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				label := material.H6(theme, "Preview")
				label.Color = themes.ColorAccent
				label.Font = font.Font{Weight: font.SemiBold}
				label.TextSize = unit.Sp(15)
				return label.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				name := "(no template generated yet)"
				if template != nil {
					name = template.Name
				}
				label := material.Body2(theme, name)
				label.Color = themes.ColorTextDim
				label.TextSize = unit.Sp(11)
				label.MaxLines = 1
				label.Truncator = "…"
				return layout.Inset{Top: unit.Dp(2), Bottom: unit.Dp(6)}.Layout(gtx, label.Layout)
			}),
			layout.Flexed(1, this.getPreviewCanvasWidget(theme)),
			layout.Rigid(this.getLegendWidget(theme)),
			layout.Rigid(layout.Spacer{Height: unit.Dp(6)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
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

				label := material.Body2(theme, this.pngStatus)
				label.TextSize = unit.Sp(11)
				label.MaxLines = 2
				label.Alignment = text.Middle
				label.Color = themes.ColorTextDim
				if !this.pngStatusOK {
					label.Color = themes.ColorError
				}
				return layout.Inset{Bottom: unit.Dp(6)}.Layout(gtx, label.Layout)
			}),
		)
	})
}

func (this *PreviewPanel) HandleClicks(gtx layout.Context) {
	// TODO: add handling for generation and saving after footer is moved to the preview panel
}

// layoutPreviewCanvas draws the preview area. Always fills the available
// space (so the surrounding panel keeps a consistent size) and renders an
// informational message inside the canvas when there is no template or no
// computed layout yet.
func (this *PreviewPanel) getPreviewCanvasWidget(theme *material.Theme) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		maxX := gtx.Constraints.Max.X
		maxY := gtx.Constraints.Max.Y
		outerSize := image.Pt(maxX, maxY)
		side := max(min(maxY, maxX), 80)
		canvasSize := image.Pt(side, side)

		// Center the square canvas inside the available area.
		offsetX := (maxX - side) / 2
		offsetY := (maxY - side) / 2
		defer op.Offset(image.Pt(offsetX, offsetY)).Push(gtx.Ops).Pop()

		// Background.
		paint.FillShape(gtx.Ops, themes.ColorPreviewBg, clip.Rect(image.Rectangle{Max: canvasSize}).Op())

		// Frame.
		radius := gtx.Dp(unit.Dp(6))
		frame := image.Rectangle{Min: image.Pt(4, 4), Max: image.Pt(side-4, side-4)}
		paint.FillShape(gtx.Ops, themes.ColorPreviewFrame, clip.Stroke{
			Path:  clip.UniformRRect(frame, radius).Path(gtx.Ops),
			Width: 2,
		}.Op())

		template := this.state.GetLastTemplate()
		if template == nil {
			return widgets.NewCenteredMessageWidget(theme, "Adjust the options to generate the map layout.", canvasSize, outerSize)(gtx)
		}

		previewLayout := services.BuildPreviewLayout(template, this.state.GetStateData().Topology, float64(side))
		if len(previewLayout.Positions) == 0 {
			return widgets.NewCenteredMessageWidget(theme, template.Name, canvasSize, outerSize)(gtx)
		}

		// Connections beneath zones.
		for _, connection := range previewLayout.Connections {
			utils.DrawConnection(gtx, connection, previewLayout.ZoneRadius)
		}
		// Non-spawn zones first, then spawn zones on top.
		for _, zone := range previewLayout.Zones {
			if zone.IsPlayer {
				continue
			}

			utils.DrawPreviewZone(gtx, theme, zone, previewLayout.ZoneRadius)
		}
		for _, zone := range previewLayout.Zones {
			if !zone.IsPlayer {
				continue
			}

			utils.DrawPreviewZone(gtx, theme, zone, previewLayout.ZoneRadius)
		}

		return layout.Dimensions{Size: outerSize}
	}
}

func (this *PreviewPanel) getLegendWidget(theme *material.Theme) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		children := make([]layout.FlexChild, 0, len(constants.LegendItems)*2)
		for i, item := range constants.LegendItems {
			if i > 0 {
				children = append(children, layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout))
			}
			children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						side := gtx.Dp(unit.Dp(10))
						rect := image.Rect(0, 0, side, side)
						paint.FillShape(gtx.Ops, item.Color, clip.UniformRRect(rect, side/2).Op(gtx.Ops))
						return layout.Dimensions{Size: rect.Max}
					}),
					layout.Rigid(layout.Spacer{Width: unit.Dp(4)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						label := material.Body2(theme, item.Label)
						label.Color = themes.ColorTextDim
						label.TextSize = unit.Sp(10)
						return label.Layout(gtx)
					}),
				)
			}))
		}
		return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx, children...)
	}
}
