package utils

import (
	"fmt"
	"image"
	"image/color"
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
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/preview"
)

func DrawConnection(gtx layout.Context, conn preview.Connection, zoneRadius int) {
	radius := float64(zoneRadius)
	start, ok1 := helpers.CalculatePointTowards(conn.Start, conn.Ctrl, radius)
	end, ok2 := helpers.CalculatePointTowards(conn.End, conn.Ctrl, radius)
	if !ok1 || !ok2 {
		return
	}

	lineColor := themes.ColorPreviewDirectLine
	lineWidth := float32(gtx.Dp(unit.Dp(2.0)))
	if conn.Portal {
		lineColor = themes.ColorPreviewPortalLine
		lineWidth = float32(gtx.Dp(unit.Dp(1.5)))
	}
	drawCurve(gtx, start, conn.Ctrl, end, lineWidth, lineColor)
}

func DrawPreviewZone(gtx layout.Context, theme *material.Theme, zone preview.Zone, zoneRadius int) {
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
		drawOffsetText(
			gtx,
			theme,
			image.Pt(badgeX, badgeY),
			fmt.Sprintf("⌂%d", zone.Castles),
			10,
			themes.ColorPreviewCastleBadge,
		)
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

func drawOffsetText(
	gtx layout.Context,
	theme *material.Theme,
	offset image.Point,
	text string,
	sizeSp int,
	textColor color.NRGBA) {
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

func zoneColors(zone preview.Zone) (fill, edge color.NRGBA) {
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

func zoneLabel(zone preview.Zone) string {
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
