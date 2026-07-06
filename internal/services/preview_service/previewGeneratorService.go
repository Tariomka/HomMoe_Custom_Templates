package preview_service

import (
	"image"
	"image/color"
	"math"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/preview"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/asset_provider"
)

const (
	canvasSize         = 700.0
	assetRadius        = 21.0 // Zone bubbles are ~21 px in 700 px canvas
	connectorLineWidth = 4
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
			drawDashedQuadratic(canvas, startPoint, controlPoint, endPoint)
		} else {
			drawThickQuadratic(canvas, startPoint, controlPoint, endPoint)
		}
	}
}

// drawThickQuadratic rasterises a quadratic Bézier (start→ctrl→end) as a chain
// of short thick line segments. With ctrl on the midpoint the samples are
// collinear, so a lone edge is drawn as a straight line.
func drawThickQuadratic(img *image.RGBA, start, ctrl, end image.Point) {
	const segments = 24
	prev := start
	for i := 1; i <= segments; i++ {
		t := float64(i) / float64(segments)
		mt := 1 - t
		x := mt*mt*float64(start.X) + 2*mt*t*float64(ctrl.X) + t*t*float64(end.X)
		y := mt*mt*float64(start.Y) + 2*mt*t*float64(ctrl.Y) + t*t*float64(end.Y)
		curr := image.Pt(int(math.Round(x)), int(math.Round(y)))
		drawThickLine(img, prev, curr)
		prev = curr
	}
}

// drawDashedQuadratic rasterises a quadratic Bézier (start→ctrl→end) as a
// perforated line: it walks the curve in fine steps, accumulating arc length,
// and only strokes the "on" portions of the dashOn/dashOff dash pattern. Used
// for portal connections so they read differently from solid direct lines.
func drawDashedQuadratic(img *image.RGBA, start, ctrl, end image.Point) {
	const dashOn, dashOff = 9.0, 13.0
	const segments = 96
	period := dashOn + dashOff
	traveled := 0.0
	prev := start
	for i := 1; i <= segments; i++ {
		t := float64(i) / float64(segments)
		mt := 1 - t
		x := mt*mt*float64(start.X) + 2*mt*t*float64(ctrl.X) + t*t*float64(end.X)
		y := mt*mt*float64(start.Y) + 2*mt*t*float64(ctrl.Y) + t*t*float64(end.Y)
		curr := image.Pt(int(math.Round(x)), int(math.Round(y)))
		segLen := math.Hypot(float64(curr.X-prev.X), float64(curr.Y-prev.Y))
		// A short segment is fully "on" or "off" based on its midpoint phase;
		// the fine sampling keeps the dashes visually even.
		phase := math.Mod(traveled+segLen/2, period)
		if phase < dashOn {
			drawThickLine(img, prev, curr)
		}
		traveled += segLen
		prev = curr
	}
}

// drawThickLine draws a width-pixel-thick line via DDA with a square brush.
func drawThickLine(img *image.RGBA, a, b image.Point) {
	dx := b.X - a.X
	dy := b.Y - a.Y
	steps := max(-dy, max(dy, max(-dx, dx)))
	if steps <= 0 {
		return
	}
	xinc := float64(dx) / float64(steps)
	yinc := float64(dy) / float64(steps)
	x := float64(a.X)
	y := float64(a.Y)
	half := connectorLineWidth / 2
	for i := 0; i <= steps; i++ {
		px := int(math.Round(x))
		py := int(math.Round(y))
		for oy := -half; oy <= half; oy++ {
			for ox := -half; ox <= half; ox++ {
				xx := px + ox
				yy := py + oy
				if xx < 0 || yy < 0 || xx >= img.Rect.Max.X || yy >= img.Rect.Max.Y {
					continue
				}
				img.SetRGBA(xx, yy, connectorLineColor)
			}
		}
		x += xinc
		y += yinc
	}
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
