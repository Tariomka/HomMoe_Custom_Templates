package utils

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
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/themes"
	"github.com/Tariomka/hommoe_custom_templates/internal/services"
)

func DrawConnection(gtx layout.Context, conn services.PreviewConnection, zoneRadius int) {
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

func DrawPreviewZone(gtx layout.Context, theme *material.Theme, zone services.PreviewZone, zoneRadius int) {
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
		drawOffsetText(gtx, theme, image.Pt(cx, cy), label, 12, themes.ColorPreviewZoneLabel)
	}
	if zone.HasCastle && zone.Castles > 0 {
		// Small badge in lower right.
		badgeX := cx + radius/2
		badgeY := cy + radius/2
		drawOffsetText(gtx, theme, image.Pt(badgeX, badgeY), fmt.Sprintf("⌂%d", zone.Castles), 10, themes.ColorPreviewCastleBadge)
	}
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

func drawOffsetText(gtx layout.Context, theme *material.Theme, offset image.Point, text string, sizeSp int, textColor color.NRGBA) {
	macro := op.Record(gtx.Ops)
	dims := func() layout.Dimensions {
		gtxLocal := gtx
		gtxLocal.Constraints.Min = image.Point{}
		gtxLocal.Constraints.Max = image.Pt(1<<14, 1<<14)
		label := material.Label(theme, unit.Sp(float32(sizeSp)), text)
		label.Color = textColor
		label.Font = font.Font{Weight: font.SemiBold}
		return label.Layout(gtxLocal)
	}()
	call := macro.Stop()

	tx := offset.X - dims.Size.X/2
	ty := offset.Y - dims.Size.Y/2
	stack := op.Offset(image.Pt(tx, ty)).Push(gtx.Ops)
	call.Add(gtx.Ops)
	stack.Pop()
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
