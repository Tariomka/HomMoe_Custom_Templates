package components

import (
	"image"
	"image/color"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
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
				return this.layoutPreviewCanvas(gtx, theme, template)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(6)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions { return this.layoutPreviewLegend(gtx, theme) }),
			layout.Rigid(layout.Spacer{Height: unit.Dp(6)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(widgets.NewButtonWidget(theme, "🖼  Save PNG", &this.btnSavePreview, template == nil)),
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
				)
			}),
		)
	})(gtx)
}

// layoutPreviewCanvas draws the preview area. Always fills the available
// space (so the surrounding panel keeps a consistent size) and renders an
// informational message inside the canvas when there is no template or no
// computed layout yet.
func (this *Window) layoutPreviewCanvas(gtx layout.Context, theme *material.Theme, template *models.RmgTemplate) layout.Dimensions {
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
	rect := image.Rectangle{Max: canvasSize}
	paint.FillShape(gtx.Ops, previewBg, clip.Rect(rect).Op())

	// Frame.
	radius := gtx.Dp(unit.Dp(6))
	frame := image.Rectangle{Min: image.Pt(4, 4), Max: image.Pt(side-4, side-4)}
	paint.FillShape(gtx.Ops, previewFrame, clip.Stroke{
		Path:  clip.UniformRRect(frame, radius).Path(gtx.Ops),
		Width: 2,
	}.Op())

	if template == nil {
		drawCenteredMessage(gtx, theme, canvasSize, "Press \"Generate Template\" to see the map layout.")
		return layout.Dimensions{Size: outerSize}
	}

	previewLayout := services.BuildPreviewLayout(template, this.settingsFile.Topology, float64(side))
	if len(previewLayout.Positions) == 0 {
		drawCenteredMessage(gtx, theme, canvasSize, template.Name)
		return layout.Dimensions{Size: outerSize}
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

	return layout.Dimensions{Size: outerSize}
}

// drawCenteredMessage renders a Body2 label centered inside the given canvas
// area. Uses the same material.Label approach as the (former) empty-state
// view so text renders reliably (unlike drawCenteredText for longer strings).
func drawCenteredMessage(gtx layout.Context, theme *material.Theme, canvasSize image.Point, txt string) {
	macro := op.Record(gtx.Ops)
	gtxLocal := gtx
	gtxLocal.Constraints.Min = image.Point{}
	gtxLocal.Constraints.Max = canvasSize
	label := material.Body2(theme, txt)
	label.Color = colTextDim
	label.TextSize = unit.Sp(12)
	label.Alignment = text.Middle
	dims := label.Layout(gtxLocal)
	call := macro.Stop()

	tx := (canvasSize.X - dims.Size.X) / 2
	ty := (canvasSize.Y - dims.Size.Y) / 2
	stack := op.Offset(image.Pt(tx, ty)).Push(gtx.Ops)
	call.Add(gtx.Ops)
	stack.Pop()
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
