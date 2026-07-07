package preview_service

import (
	"image"
	"image/color"
	"math"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
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
)

var connectorLineColor = color.RGBA{R: 0x39, G: 0x11, B: 0x14, A: 0xFF}

type PreviewGeneratorService struct {
	assetProvider *asset_provider.AssetProvider
}

func NewPreviewGenerator() (*PreviewGeneratorService, error) {
	assetProvider, err := asset_provider.NewAssetProvider()
	if err != nil {
		return nil, err
	}

	return &PreviewGeneratorService{assetProvider: assetProvider}, nil
}

func (this *PreviewGeneratorService) CreatePreviewImage(
	template *entities.RmgTemplate,
	topology config.MapTopology) *image.RGBA {
	canvas := image.NewRGBA(image.Rect(0, 0, canvasSize, canvasSize))
	this.assetProvider.DrawBackground(canvas)

	layout := BuildPreviewLayout(template, topology, canvasSize)
	if len(layout.Positions) == 0 {
		return canvas
	}

	scale := min(float64(layout.ZoneRadius)/assetRadius, 1.15)
	fitterCallback := newAssetFitter(layout.Zones, scale)

	this.drawConnections(canvas, layout.Connections, fitterCallback, assetRadius*scale)
	for _, zone := range layout.Zones {
		if !zone.IsPlayer {
			this.assetProvider.DrawNeutralZone(canvas, zone, fitterCallback(zone.Center), scale)
		}
	}
	for _, zone := range layout.Zones {
		if zone.IsPlayer {
			this.assetProvider.DrawPlayerZone(canvas, zone, fitterCallback(zone.Center), scale)
		}
	}
	return canvas
}

func (this *PreviewGeneratorService) drawConnections(
	canvas *image.RGBA,
	connections []preview.PreviewConnection,
	fitterCallback assetFitter,
	zoneRadius float64) {
	for _, conn := range connections {
		controlPoint := fitterCallback(conn.Ctrl) // Bézier control point
		startPoint, ok1 := this.movePointTowards(fitterCallback(conn.A), controlPoint, zoneRadius)
		endPoint, ok2 := this.movePointTowards(fitterCallback(conn.B), controlPoint, zoneRadius)
		if !ok1 || !ok2 {
			continue
		}

		if conn.Portal {
			this.drawDashedLine(canvas, startPoint, controlPoint, endPoint)
		} else {
			this.drawSolidLine(canvas, startPoint, controlPoint, endPoint)
		}
	}
}

func (this *PreviewGeneratorService) drawSolidLine(canvas *image.RGBA, start, ctrl, end image.Point) {
	previousPoint := start
	for i := range segmentsSolid {
		t := float64(i+1) / segmentsSolid
		currentPoint := this.getPointOnQuadraticBezierCurve(start, ctrl, end, t)
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
		currentPoint := this.getPointOnQuadraticBezierCurve(start, ctrl, end, t)
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

func (this *PreviewGeneratorService) drawLine(canvas *image.RGBA, a, b image.Point) {
	delta := b.Sub(a)
	steps := max(math.Abs(float64(delta.X)), math.Abs(float64(delta.Y)))
	if steps <= 0 {
		return
	}

	increment := data.Vec2FromPoint[float64](delta).DivideScalar(steps)
	half := connectorLineWidth / 2
	for i := range int(steps) {
		center := data.Vec2FromPoint[float64](a).
			Add(increment.MultiplyScalar(float64(i))).
			ToPointRounded()
		brush := image.Rect(center.X-half, center.Y-half, center.X+half+1, center.Y+half+1).
			Intersect(canvas.Bounds()) // Square brush around the centre, clipped to the canvas.
		for y := brush.Min.Y; y < brush.Max.Y; y++ {
			for x := brush.Min.X; x < brush.Max.X; x++ {
				canvas.SetRGBA(x, y, connectorLineColor)
			}
		}
	}
}

// getPointOnQuadraticBezierCurve evaluates the quadratic Bézier (start(P0)→ctrl(P1)→end(P2)) at
// the point along the curve t in [0,1] by applying this formula:
//
// B(t) = (1-t)²*P0 + 2*(1-t)*t*P1 + t²*P2.
func (this *PreviewGeneratorService) getPointOnQuadraticBezierCurve(
	start, ctrl, end image.Point,
	t float64) image.Point {
	mt := 1 - t
	point := data.Vec2FromPoint[float64](start).MultiplyScalar(mt * mt).
		Add(data.Vec2FromPoint[float64](ctrl).MultiplyScalar(2 * mt * t)).
		Add(data.Vec2FromPoint[float64](end).MultiplyScalar(t * t))
	return image.Pt(int(math.Round(point.X)), int(math.Round(point.Y)))
}

// movePointTowards returns `from` moved toward `toward` by `distance` pixels.
// ok is false when the two points are coincident.
func (this *PreviewGeneratorService) movePointTowards(
	from, toward image.Point,
	distance float64) (movedPoint image.Point, ok bool) {
	deltaX := float64(toward.X - from.X)
	deltaY := float64(toward.Y - from.Y)
	deltaDistance := math.Hypot(deltaX, deltaY)
	if deltaDistance < 1 {
		return movedPoint, false
	}

	movedPoint.X = int(float64(from.X) + deltaX/deltaDistance*distance)
	movedPoint.Y = int(float64(from.Y) + deltaY/deltaDistance*distance)
	return movedPoint, true
}
