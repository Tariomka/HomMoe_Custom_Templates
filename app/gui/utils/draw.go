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
	// Trim both ends back to the circle outlines along the tangent toward the
	// control point so the edge doesn't overlap the zone fill, then stroke the
	// quadratic curve. A lone edge has its control point on the midpoint and so
	// renders as a straight line; parallel edges fan out into distinct curves.
	radius := float64(zoneRadius)
	startX, startY, ok1 := trimToward(conn.A, conn.Ctrl, radius)
	endX, endY, ok2 := trimToward(conn.B, conn.Ctrl, radius)
	if !ok1 || !ok2 {
		return
	}

	lineColor := themes.ColorPreviewDirectLine
	lineWidth := float32(gtx.Dp(unit.Dp(2.0)))
	if conn.Portal {
		lineColor = themes.ColorPreviewPortalLine
		lineWidth = float32(gtx.Dp(unit.Dp(1.5)))
	}
	drawCurve(gtx,
		image.Pt(int(startX), int(startY)),
		conn.Ctrl,
		image.Pt(int(endX), int(endY)),
		lineWidth, lineColor)
}

// trimToward returns the point moved from `from` toward `toward` by `dist`
// pixels. ok is false when the two points are coincident.
func trimToward(from, toward image.Point, dist float64) (x, y float64, ok bool) {
	dx := float64(toward.X - from.X)
	dy := float64(toward.Y - from.Y)
	d := math.Hypot(dx, dy)
	if d < 1 {
		return 0, 0, false
	}
	return float64(from.X) + dx/d*dist, float64(from.Y) + dy/d*dist, true
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

func drawCurve(gtx layout.Context, start, ctrl, end image.Point, width float32, lineColor color.NRGBA) {
	var path clip.Path
	path.Begin(gtx.Ops)
	path.MoveTo(f32.Point{X: float32(start.X), Y: float32(start.Y)})
	path.QuadTo(
		f32.Point{X: float32(ctrl.X), Y: float32(ctrl.Y)},
		f32.Point{X: float32(end.X), Y: float32(end.Y)},
	)
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
		// Spawn-1 / Spawn-2 → "P1"...
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
