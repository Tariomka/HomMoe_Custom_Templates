package components

import (
	"image"
	"image/color"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/Tariomka/hommoe_custom_templates/internal/gui/components/widgets"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/services"
)

// Preview holds the layout cache + buttons for the preview panel.
type Preview struct {
	btnSavePNG  widget.Clickable
	btnRefresh  widget.Clickable
	pngStatus   string
	pngStatusOK bool
}

// — Layout —

// layoutPreviewPanel renders the right-hand preview area. Returns empty
// dimensions when there's nothing to show (so the caller can omit it).
func (this *Window) layoutPreviewPanel(gtx layout.Context, theme *material.Theme) layout.Dimensions {
	if this.btnSavePreview.Clicked(gtx) {
		this.savePreviewPNG()
	}
	template := this.lastTemplate

	return widgets.NewPanelWidget(unit.Dp(8), func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				label := material.H6(theme, "Preview")
				label.Color = colGold
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
				label.Color = colTextDim
				label.TextSize = unit.Sp(11)
				label.MaxLines = 1
				label.Truncator = "…"
				return layout.Inset{Top: unit.Dp(2), Bottom: unit.Dp(6)}.Layout(gtx, label.Layout)
			}),
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				if template == nil {
					return this.layoutPreviewEmpty(gtx, theme)
				}
				return this.layoutPreviewCanvas(gtx, theme, template)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(6)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions { return this.layoutPreviewLegend(gtx, theme) }),
			layout.Rigid(layout.Spacer{Height: unit.Dp(6)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						if this.preview.pngStatus == "" {
							return layout.Dimensions{}
						}
						statusColor := colTextDim
						if !this.preview.pngStatusOK {
							statusColor = colError
						}
						label := material.Body2(theme, this.preview.pngStatus)
						label.Color = statusColor
						label.TextSize = unit.Sp(11)
						label.MaxLines = 2
						return label.Layout(gtx)
					}),
					layout.Rigid(widgets.NewButtonWidget(theme, "🖼  Save PNG", &this.btnSavePreview, template == nil)),
				)
			}),
		)
	})(gtx)
}

func (this *Window) layoutPreviewEmpty(gtx layout.Context, theme *material.Theme) layout.Dimensions {
	return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		label := material.Body2(theme, "Press \"Generate Template\" to see the map layout.")
		label.Color = colTextDim
		label.TextSize = unit.Sp(12)
		label.Alignment = text.Middle
		return label.Layout(gtx)
	})
}

// layoutPreviewCanvas draws the actual map preview at the largest possible
// square fitting inside the available area.
func (this *Window) layoutPreviewCanvas(gtx layout.Context, theme *material.Theme, template *models.RmgTemplate) layout.Dimensions {
	maxX := gtx.Constraints.Max.X
	maxY := gtx.Constraints.Max.Y
	side := max(min(maxY, maxX), 80)
	canvasSize := image.Pt(side, side)

	// Background.
	rect := image.Rectangle{Max: canvasSize}
	paint.FillShape(gtx.Ops, previewBg, clip.Rect(rect).Op())

	// Frame.
	radius := gtx.Dp(unit.Dp(6))
	frame := image.Rectangle{Min: image.Pt(4, 4), Max: image.Pt(side-4, side-4)}
	paint.FillShape(gtx.Ops, previewFrame, clip.Stroke{
		Path:  clip.UniformRRect(frame, radius).Path(gtx.Ops),
		Width: 2,
	}.Op())

	previewLayout := services.BuildPreviewLayout(template, this.settingsFile.Topology, float64(side))
	if len(previewLayout.Positions) == 0 {
		drawCenteredText(gtx, theme, canvasSize, template.Name, 18, colText)
		return layout.Dimensions{Size: canvasSize}
	}

	// Connections beneath zones.
	for _, conn := range previewLayout.Connections {
		drawConnection(gtx, conn, previewLayout.ZoneRadius)
	}
	// Non-spawn zones first, then spawn zones on top.
	for _, zone := range previewLayout.Zones {
		if zone.IsPlayer {
			continue
		}
		drawPreviewZone(gtx, theme, zone, previewLayout.ZoneRadius)
	}
	for _, zone := range previewLayout.Zones {
		if !zone.IsPlayer {
			continue
		}
		drawPreviewZone(gtx, theme, zone, previewLayout.ZoneRadius)
	}

	return layout.Dimensions{Size: canvasSize}
}

func (this *Window) layoutPreviewLegend(gtx layout.Context, theme *material.Theme) layout.Dimensions {
	type legendItem struct {
		Color color.NRGBA
		Label string
	}
	items := []legendItem{
		{previewSpawnEdge, "Player"},
		{previewBronzeEdge, "Bronze"},
		{previewSilverEdge, "Silver"},
		{previewGoldEdge, "Gold"},
		{previewHubEdge, "Hub"},
	}
	children := make([]layout.FlexChild, 0, len(items)*2)
	for i, item := range items {
		item := item
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
					label.Color = colTextDim
					label.TextSize = unit.Sp(10)
					return label.Layout(gtx)
				}),
			)
		}))
	}
	return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx, children...)
}
