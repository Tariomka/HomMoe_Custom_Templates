package gui

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"strings"

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

var (
	previewBg         = color.NRGBA{R: 0x1C, G: 0x16, B: 0x10, A: 0xFF}
	previewFrame      = color.NRGBA{R: 0x8F, G: 0x73, B: 0x3F, A: 0xFF}
	previewBronzeFill = color.NRGBA{R: 0x65, G: 0x43, B: 0x21, A: 0xFF}
	previewBronzeEdge = color.NRGBA{R: 0xCD, G: 0x7F, B: 0x32, A: 0xFF}
	previewSilverFill = color.NRGBA{R: 0x48, G: 0x4C, B: 0x50, A: 0xFF}
	previewSilverEdge = color.NRGBA{R: 0xC0, G: 0xC0, B: 0xC0, A: 0xFF}
	previewGoldFill   = color.NRGBA{R: 0x78, G: 0x5A, B: 0x14, A: 0xFF}
	previewGoldEdge   = color.NRGBA{R: 0xFF, G: 0xD2, B: 0x32, A: 0xFF}
	previewSpawnFill  = color.NRGBA{R: 0x2A, G: 0x5A, B: 0x32, A: 0xFF}
	previewSpawnEdge  = color.NRGBA{R: 0x64, G: 0xC8, B: 0x78, A: 0xFF}
	previewHubFill    = color.NRGBA{R: 0x37, G: 0x50, B: 0x5F, A: 0xFF}
	previewHubEdge    = color.NRGBA{R: 0x82, G: 0xB4, B: 0xC8, A: 0xFF}
	previewDirectLine = color.NRGBA{R: 0xB4, G: 0x91, B: 0x3C, A: 0xFF}
	previewPortalLine = color.NRGBA{R: 0x5A, G: 0xAA, B: 0xD2, A: 0xB4}
)

// previewState holds the layout cache + buttons for the preview panel.
type previewState struct {
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

// — Drawing primitives (Gio canvas) —

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

	lineColor := previewDirectLine
	lineWidth := float32(gtx.Dp(unit.Dp(2.0)))
	if conn.Portal {
		lineColor = previewPortalLine
		lineWidth = float32(gtx.Dp(unit.Dp(1.5)))
	}
	drawLine(gtx, image.Pt(int(ax), int(ay)), image.Pt(int(bx), int(by)), lineWidth, lineColor)
}

func drawLine(gtx layout.Context, start, end image.Point, width float32, lineColor color.NRGBA) {
	var path clip.Path
	path.Begin(gtx.Ops)
	path.MoveTo(f32Pt(start))
	path.LineTo(f32Pt(end))
	paint.FillShape(gtx.Ops, lineColor, clip.Stroke{
		Path:  path.End(),
		Width: width,
	}.Op())
}

func f32Pt(point image.Point) (out f32Point) {
	out.X = float32(point.X)
	out.Y = float32(point.Y)
	return
}

// f32Point is a tiny shim so we don't import gioui.org/f32 just for two fields.
type f32Point = struct {
	X, Y float32
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
		drawCenteredText(gtx, theme, image.Pt(cx, cy), label, 12, color.NRGBA{R: 0xF8, G: 0xE8, B: 0xC0, A: 0xFF})
	}
	if zone.HasCastle && zone.Castles > 0 {
		// Small badge in lower right.
		badgeX := cx + radius/2
		badgeY := cy + radius/2
		drawCenteredText(gtx, theme, image.Pt(badgeX, badgeY), fmt.Sprintf("⌂%d", zone.Castles), 10, color.NRGBA{R: 0xFF, G: 0xE8, B: 0x90, A: 0xFF})
	}
}

func zoneColors(zone services.PreviewZone) (fill, edge color.NRGBA) {
	switch {
	case zone.IsPlayer:
		return previewSpawnFill, previewSpawnEdge
	case zone.IsHub:
		return previewHubFill, previewHubEdge
	}
	switch zone.Tier {
	case 3:
		return previewGoldFill, previewGoldEdge
	case 2:
		return previewSilverFill, previewSilverEdge
	default:
		return previewBronzeFill, previewBronzeEdge
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

// drawCenteredText draws text centered on the given canvas point.
func drawCenteredText(gtx layout.Context, theme *material.Theme, center image.Point, txt string, sizeSp int, textColor color.NRGBA) {
	macro := op.Record(gtx.Ops)
	dims := func() layout.Dimensions {
		gtxLocal := gtx
		gtxLocal.Constraints.Min = image.Point{}
		gtxLocal.Constraints.Max = image.Pt(1<<14, 1<<14)
		label := material.Label(theme, unit.Sp(float32(sizeSp)), txt)
		label.Color = textColor
		label.Font = font.Font{Weight: font.SemiBold}
		return label.Layout(gtxLocal)
	}()
	call := macro.Stop()

	tx := center.X - dims.Size.X/2
	ty := center.Y - dims.Size.Y/2
	stack := op.Offset(image.Pt(tx, ty)).Push(gtx.Ops)
	call.Add(gtx.Ops)
	stack.Pop()
}

// — PNG export —

// savePreviewPNG renders the current template into a software bitmap and
// writes it next to the .rmg.json output.
func (this *Window) savePreviewPNG() {
	template := this.lastTemplate
	if template == nil {
		this.preview.pngStatus = "Generate a template first."
		this.preview.pngStatusOK = false
		return
	}
	dir := strings.TrimSpace(this.outputPath.Text())
	if dir == "" {
		this.preview.pngStatus = "Output directory is empty."
		this.preview.pngStatusOK = false
		return
	}
	out, err := services.WritePreviewPNG(dir, template, this.settingsFile.Topology, 700)
	if err != nil {
		this.preview.pngStatus = "Save failed: " + err.Error()
		this.preview.pngStatusOK = false
		return
	}
	this.preview.pngStatus = "Saved " + out
	this.preview.pngStatusOK = true
}
