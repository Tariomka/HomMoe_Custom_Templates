package components

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"strings"

	"gioui.org/f32"
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/Tariomka/hommoe_custom_templates/internal/constants"
	"github.com/Tariomka/hommoe_custom_templates/internal/gui/components/themes"
	"github.com/Tariomka/hommoe_custom_templates/internal/gui/components/widgets"
	"github.com/Tariomka/hommoe_custom_templates/internal/gui/utils"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/services"
)

// PreviewPanel holds the layout cache + buttons for the preview panel.
type PreviewPanel struct {
	btnSavePNG  widget.Clickable
	btnRefresh  widget.Clickable
	pngStatus   string
	pngStatusOK bool

	state *State
}

func NewPreviewPanel(state *State) *PreviewPanel {
	return &PreviewPanel{state: state}
}

func (this *PreviewPanel) GetPanelWidget(theme *material.Theme) layout.Widget {
	template := this.state.GetLastTemplate()
	return widgets.NewPanelWidget(unit.Dp(8), func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				label := material.H6(theme, "Preview")
				label.Color = themes.ColorGold
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
			layout.Flexed(1, this.getPreviewCanvasWidget(theme, template)),
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
				statusColor := themes.ColorTextDim
				if !this.pngStatusOK {
					statusColor = themes.ColorError
				}
				return layout.Inset{Bottom: unit.Dp(6)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					label := material.Body2(theme, this.pngStatus)
					label.Color = statusColor
					label.TextSize = unit.Sp(11)
					label.MaxLines = 2
					label.Alignment = text.Middle
					return label.Layout(gtx)
				})
			}),
			layout.Rigid(this.getLegendWidget(theme)),
			layout.Rigid(layout.Spacer{Height: unit.Dp(6)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(widgets.NewButtonWidget(theme, "🖼  Save PNG", &this.btnSavePNG, template == nil)),
				)
			}),
		)
	})
}

func (this *PreviewPanel) HandleClicks(gtx layout.Context) {
	if this.btnSavePNG.Clicked(gtx) {
		this.savePreviewPNG()
	}
}

// layoutPreviewCanvas draws the preview area. Always fills the available
// space (so the surrounding panel keeps a consistent size) and renders an
// informational message inside the canvas when there is no template or no
// computed layout yet.
func (this *PreviewPanel) getPreviewCanvasWidget(theme *material.Theme, template *models.RmgTemplate) layout.Widget {
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
		rect := image.Rectangle{Max: canvasSize}
		paint.FillShape(gtx.Ops, themes.ColorPreviewBg, clip.Rect(rect).Op())

		// Frame.
		radius := gtx.Dp(unit.Dp(6))
		frame := image.Rectangle{Min: image.Pt(4, 4), Max: image.Pt(side-4, side-4)}
		paint.FillShape(gtx.Ops, themes.ColorPreviewFrame, clip.Stroke{
			Path:  clip.UniformRRect(frame, radius).Path(gtx.Ops),
			Width: 2,
		}.Op())

		if template == nil {
			return widgets.NewCenteredMessageWidget(theme, "Press \"Generate Template\" to see the map layout.", canvasSize, outerSize)(gtx)
		}

		previewLayout := services.BuildPreviewLayout(template, this.state.GetSettingsFile().Topology, float64(side))
		if len(previewLayout.Positions) == 0 {
			return widgets.NewCenteredMessageWidget(theme, template.Name, canvasSize, outerSize)(gtx)
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

// savePreviewPNG renders the current template into a software bitmap and
// writes it next to the .rmg.json output.
func (this *PreviewPanel) savePreviewPNG() {
	template := this.state.GetLastTemplate()
	if template == nil {
		this.pngStatus = "Generate a template first."
		this.pngStatusOK = false
		return
	}
	dir := strings.TrimSpace(this.state.GetOutputPath())
	if dir == "" {
		this.pngStatus = "Output directory is empty."
		this.pngStatusOK = false
		return
	}
	out, err := services.WritePreviewPNG(dir, template, this.state.GetSettingsFile().Topology, 700)
	if err != nil {
		this.pngStatus = "Save failed: " + err.Error()
		this.pngStatusOK = false
		return
	}
	this.pngStatus = "Saved " + out
	this.pngStatusOK = true
}

func drawConnection(gtx layout.Context, conn services.PreviewConnection, zoneRadius int) {
	// Trim line to circle edges so it doesn't overlap zone fill.
	dx := float64(conn.B.X - conn.A.X)
	dy := float64(conn.B.Y - conn.A.Y)
	distance := math.Hypot(dx, dy)
	if distance < 1 {
		return
	}
	ux := dx / distance
	uy := dy / distance
	radius := float64(zoneRadius)
	ax := float64(conn.A.X) + ux*radius
	ay := float64(conn.A.Y) + uy*radius
	bx := float64(conn.B.X) - ux*radius
	by := float64(conn.B.Y) - uy*radius

	lineColor := themes.ColorPreviewDirectLine
	lineWidth := float32(gtx.Dp(unit.Dp(2.0)))
	if conn.Portal {
		lineColor = themes.ColorPreviewPortalLine
		lineWidth = float32(gtx.Dp(unit.Dp(1.5)))
	}
	drawLine(gtx, image.Pt(int(ax), int(ay)), image.Pt(int(bx), int(by)), lineWidth, lineColor)
}

func drawLine(gtx layout.Context, start, end image.Point, width float32, lineColor color.NRGBA) {
	var path clip.Path
	path.Begin(gtx.Ops)
	path.MoveTo(f32.Point{X: float32(start.X), Y: float32(start.Y)})
	path.LineTo(f32.Point{X: float32(end.X), Y: float32(end.Y)})
	paint.FillShape(gtx.Ops, lineColor, clip.Stroke{
		Path:  path.End(),
		Width: width,
	}.Op())
}

func drawPreviewZone(gtx layout.Context, theme *material.Theme, zone services.PreviewZone, zoneRadius int) {
	radius := zoneRadius
	if zone.IsHub && radius < 28 {
		radius = 28
	}
	cx, cy := zone.Center.X, zone.Center.Y
	rect := image.Rect(cx-radius, cy-radius, cx+radius, cy+radius)

	fill, edge := zoneColors(zone)
	circle := clip.UniformRRect(rect, radius).Op(gtx.Ops)
	paint.FillShape(gtx.Ops, fill, circle)
	paint.FillShape(gtx.Ops, edge, clip.Stroke{
		Path:  clip.UniformRRect(rect, radius).Path(gtx.Ops),
		Width: float32(gtx.Dp(unit.Dp(2))),
	}.Op())

	label := zoneLabel(zone)
	if label != "" {
		utils.DrawOffsetText(gtx, theme, image.Pt(cx, cy), label, 12, color.NRGBA{R: 0xF8, G: 0xE8, B: 0xC0, A: 0xFF})
	}
	if zone.HasCastle && zone.Castles > 0 {
		// Small badge in lower right.
		badgeX := cx + radius/2
		badgeY := cy + radius/2
		utils.DrawOffsetText(gtx, theme, image.Pt(badgeX, badgeY), fmt.Sprintf("⌂%d", zone.Castles), 10, color.NRGBA{R: 0xFF, G: 0xE8, B: 0x90, A: 0xFF})
	}
}

func zoneColors(zone services.PreviewZone) (fill, edge color.NRGBA) {
	switch {
	case zone.IsPlayer:
		return themes.ColorPreviewSpawnFill, themes.ColorPreviewSpawnEdge
	case zone.IsHub:
		return themes.ColorPreviewHubFill, themes.ColorPreviewHubEdge
	}
	switch zone.Tier {
	case 3:
		return themes.ColorPreviewGoldFill, themes.ColorPreviewGoldEdge
	case 2:
		return themes.ColorPreviewSilverFill, themes.ColorPreviewSilverEdge
	default:
		return themes.ColorPreviewBronzeFill, themes.ColorPreviewBronzeEdge
	}
}

func zoneLabel(zone services.PreviewZone) string {
	if zone.IsPlayer {
		if zone.Owner > 0 {
			return fmt.Sprintf("P%d", zone.Owner)
		}
		// Spawn-1 / Spawn-2 → "P1"…
		if strings.HasPrefix(zone.Name, "Spawn-") {
			return "P" + zone.Name[len("Spawn-"):]
		}
		return zone.Letter
	}
	if zone.IsHub {
		return "Hub"
	}
	switch zone.Tier {
	case 3:
		return "G"
	case 2:
		return "S"
	default:
		return "B"
	}
}
