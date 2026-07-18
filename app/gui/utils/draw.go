package utils

import (
	"fmt"
	"image"
	"image/color"

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
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/zone_helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutralZone"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/preview"
)

func DrawConnection(gtx layout.Context, conn preview.Connection, zoneRadius int) {
	radius := float64(zoneRadius)
	start, ok1 := helpers.CalculatePointTowards(conn.Start, conn.Ctrl, radius)
	end, ok2 := helpers.CalculatePointTowards(conn.End, conn.Ctrl, radius)
	if !ok1 || !ok2 {
		return
	}

	lineColor := themes.ColorsPreview.DirectLine
	lineWidth := float32(gtx.Dp(unit.Dp(2.0)))
	if conn.Portal {
		lineColor = themes.ColorsPreview.PortalLine
		lineWidth = float32(gtx.Dp(unit.Dp(1.5)))
	}
	drawCurve(gtx, start, conn.Ctrl, end, lineWidth, lineColor)
}

func DrawPreviewZone(gtx layout.Context, theme *material.Theme, zone preview.Zone, zoneRadius int) {
	radius := zoneRadius
	if zone.Type == preview.ZoneTypeHub {
		radius = max(radius, 28)
	}
	rect := image.Rect(zone.Center.X-radius, zone.Center.Y-radius, zone.Center.X+radius, zone.Center.Y+radius)

	fill, edge := zoneColors(zone)
	circle := clip.UniformRRect(rect, radius).Op(gtx.Ops)
	paint.FillShape(gtx.Ops, fill, circle)
	paint.FillShape(gtx.Ops, edge, clip.Stroke{
		Path:  clip.UniformRRect(rect, radius).Path(gtx.Ops),
		Width: float32(gtx.Dp(unit.Dp(2))),
	}.Op())

	label := zoneLabel(zone)
	if label != "" {
		drawOffsetText(gtx, theme, zone.Center, label, 12, themes.ColorsPreview.ZoneLabel)
	}
	if zone.HasCastles() {
		drawOffsetText(
			gtx, theme, zone.Center.Add(image.Pt(radius/2, radius/2)),
			fmt.Sprintf("⌂%d", zone.Castles), 10, themes.ColorsPreview.CastleBadge)
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
	switch zone.Type {
	case preview.ZoneTypePlayer:
		return themes.ColorsPreview.SpawnFill, themes.ColorsPreview.SpawnEdge
	case preview.ZoneTypeHub:
		return themes.ColorsPreview.HubFill, themes.ColorsPreview.HubEdge
	case preview.ZoneTypeNeutralZone, preview.ZoneTypeUnknown:
	}

	switch zone.Quality {
	case neutralZone.QualityHighest:
		return themes.ColorsPreview.PlatinumFill, themes.ColorsPreview.PlatinumEdge
	case neutralZone.QualityHigh:
		return themes.ColorsPreview.GoldFill, themes.ColorsPreview.GoldEdge
	case neutralZone.QualityMedium:
		return themes.ColorsPreview.SilverFill, themes.ColorsPreview.SilverEdge
	case neutralZone.QualityLow:
		return themes.ColorsPreview.BronzeFill, themes.ColorsPreview.BronzeEdge
	case neutralZone.QualityLowest:
		return themes.ColorsPreview.PlasticFill, themes.ColorsPreview.PlasticEdge
	case neutralZone.QualityUnknown:
		fallthrough
	default:
		// shouldn't happen, but if it does, at least it's visible
		return color.NRGBA{A: 255}, color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	}
}

func zoneLabel(zone preview.Zone) string {
	if zone.Type == preview.ZoneTypePlayer {
		if zone.Owner > 0 {
			return fmt.Sprintf("P%d", zone.Owner)
		}

		// Spawn-1 / Spawn-2 → "P1"...
		if zone_helpers.IsZoneNamePlayer(zone.Name) {
			return "P" + zone.Name[len("Spawn-"):]
		}

		return zone.Label
	}

	if zone.Type == preview.ZoneTypeHub {
		return "Hub"
	}

	switch zone.Quality {
	case neutralZone.QualityHighest:
		return "Pt"
	case neutralZone.QualityHigh:
		return "G"
	case neutralZone.QualityMedium:
		return "S"
	case neutralZone.QualityLow:
		return "B"
	case neutralZone.QualityLowest:
		return "p"
	case neutralZone.QualityUnknown:
		fallthrough
	default:
		return "?"
	}
}
