package preview_service

import (
	"image"
	"image/color"
	"image/draw"
	"math"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/preview"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/asset_provider"
)

const (
	canvasSize         = 700.0
	assetRadius        = 21.0 // Zone bubbles are ~21 px in 700 px canvas
	connectorLineWidth = 4

	segmentsSolid       = 24
	segmentsDashed      = 96
	dashLength, dashGap = 9.0, 13.0

	// arenaMarkerScale shrinks the swords sprite relative to a zone bubble so it
	// reads as a marker sitting on the connection rather than as another zone.
	arenaMarkerScale = 0.75
)

var connectorLineColor = color.RGBA{R: 0x33, G: 0x18, B: 0x18, A: 0xFF}

type PreviewGeneratorService struct {
	assetProvider *asset_provider.AssetProvider
	layoutService *PreviewLayoutService
}

func NewPreviewGenerator(layoutService *PreviewLayoutService) (*PreviewGeneratorService, error) {
	assetProvider, err := asset_provider.NewAssetProvider()
	if err != nil {
		return nil, err
	}

	return &PreviewGeneratorService{
		assetProvider: assetProvider,
		layoutService: layoutService,
	}, nil
}

func (this *PreviewGeneratorService) CreatePreviewImage(
	template *entities.RmgTemplate,
	topology config.MapTopology) *image.RGBA {
	canvas := image.NewRGBA(image.Rect(0, 0, canvasSize, canvasSize))
	this.assetProvider.DrawBackground(canvas)

	layout := this.layoutService.BuildPreviewLayout(template, topology, canvasSize)
	if len(layout.Positions) == 0 {
		return canvas
	}

	scale := min(float64(layout.ZoneRadius)/assetRadius, 1.15)
	fitterCallback := newAssetFitter(layout.Zones, scale)

	this.drawConnections(canvas, layout.Connections, fitterCallback, scale)
	for _, zone := range layout.Zones {
		if zone.Type != preview.ZoneTypePlayer {
			this.assetProvider.DrawNeutralZone(canvas, zone, fitterCallback(zone.Center), scale)
		}
	}
	for _, zone := range layout.Zones {
		if zone.Type == preview.ZoneTypePlayer {
			this.assetProvider.DrawPlayerZone(canvas, zone, fitterCallback(zone.Center), scale)
		}
	}
	return canvas
}

func (this *PreviewGeneratorService) drawConnections(
	canvas *image.RGBA,
	connections []preview.Connection,
	fitterCallback assetFitter,
	scale float64) {
	zoneRadius := assetRadius * scale
	for _, conn := range connections {
		controlPoint := fitterCallback(conn.Ctrl) // Bézier control point
		startPoint, ok1 := helpers.CalculatePointTowards(fitterCallback(conn.Start), controlPoint, zoneRadius)
		endPoint, ok2 := helpers.CalculatePointTowards(fitterCallback(conn.End), controlPoint, zoneRadius)
		if !ok1 || !ok2 {
			continue
		}

		if conn.IsPortal() {
			this.drawDashedLine(canvas, startPoint, controlPoint, endPoint)
		} else {
			this.drawSolidLine(canvas, startPoint, controlPoint, endPoint)
		}

		if conn.IsGladiatorArena() {
			midPoint := helpers.GetPointOnQuadraticBezierCurve(startPoint, controlPoint, endPoint, 0.5)
			this.assetProvider.DrawArenaMarker(canvas, midPoint, scale*arenaMarkerScale)
		}
	}
}

func (this *PreviewGeneratorService) drawSolidLine(canvas *image.RGBA, start, ctrl, end image.Point) {
	previousPoint := start
	for i := range segmentsSolid {
		t := float64(i+1) / segmentsSolid
		currentPoint := helpers.GetPointOnQuadraticBezierCurve(start, ctrl, end, t)
		this.drawLine(canvas, previousPoint, currentPoint)
		previousPoint = currentPoint
	}
}

func (this *PreviewGeneratorService) drawDashedLine(canvas *image.RGBA, start, ctrl, end image.Point) {
	period := dashLength + dashGap
	traveled := 0.0
	previousPoint := start
	for i := range segmentsDashed {
		t := float64(i+1) / segmentsDashed
		currentPoint := helpers.GetPointOnQuadraticBezierCurve(start, ctrl, end, t)
		segmentLength := math.Hypot(
			float64(currentPoint.X-previousPoint.X),
			float64(currentPoint.Y-previousPoint.Y))
		if math.Mod(traveled+segmentLength/2, period) < dashLength {
			this.drawLine(canvas, previousPoint, currentPoint)
		}
		traveled += segmentLength
		previousPoint = currentPoint
	}
}

func (this *PreviewGeneratorService) drawLine(canvas *image.RGBA, start, end image.Point) {
	delta := end.Sub(start)
	steps := max(math.Abs(float64(delta.X)), math.Abs(float64(delta.Y)))
	if steps <= 0 {
		return
	}

	increment := data.Vec2FromPoint[float64](delta).DivideScalar(steps)
	half := connectorLineWidth / 2
	brushSource := image.NewUniform(connectorLineColor)
	for i := range int(steps) {
		center := data.Vec2FromPoint[float64](start).
			Add(increment.MultiplyScalar(float64(i))).
			ToPointRounded()
		brush := image.Rect(center.X-half, center.Y-half, center.X+half+1, center.Y+half+1).
			Intersect(canvas.Bounds()) // Square brush around the center, clipped to the canvas.
		draw.Draw(canvas, brush, brushSource, image.Point{}, draw.Src)
	}
}
