package preview_service

import (
	"image"
	"image/color"
	"image/draw"
	"math"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/preview"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
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

type PreviewGeneratorService struct {
	assetProvider *asset_provider.AssetProvider
	layoutService IPreviewLayoutService
}

func NewPreviewGenerator(layoutService IPreviewLayoutService) (IPreviewGeneratorService, error) {
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
	template *template_model.Template,
	topology config.MapTopology) *image.RGBA {
	canvas := image.NewRGBA(image.Rect(0, 0, canvasSize, canvasSize))
	this.assetProvider.DrawBackground(canvas)

	layout := this.layoutService.BuildPreviewLayout(template, topology, canvasSize)
	if len(layout.Positions) == 0 {
		return canvas
	}

	scale := min(layout.ZoneRadius/assetRadius, 1.15)
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
	brushSource := image.NewUniform(color.RGBA{R: 0x33, G: 0x18, B: 0x18, A: 0xFF})
	maskSource := image.NewUniform(color.Alpha{A: 128})
	var strokeMask *image.Alpha
	for _, connection := range connections {
		controlPoint := fitterCallback(connection.Ctrl) // Bézier control point
		startPoint, startValid := helpers.CalculatePointTowards(
			fitterCallback(connection.Start), controlPoint, zoneRadius)
		endPoint, endValid := helpers.CalculatePointTowards(fitterCallback(connection.End), controlPoint, zoneRadius)
		if !startValid || !endValid {
			continue
		}

		var target draw.Image = canvas
		source := brushSource
		if !connection.HasRoad {
			if strokeMask == nil {
				strokeMask = image.NewAlpha(canvas.Bounds())
			}
			target, source = strokeMask, maskSource
		}

		var painted image.Rectangle
		if connection.IsPortal() {
			painted = this.drawDashedLine(target, source, startPoint, controlPoint, endPoint)
		} else {
			painted = this.drawSolidLine(target, source, startPoint, controlPoint, endPoint)
		}

		if !connection.HasRoad && !painted.Empty() {
			// Src stamps overwrite the mask, so an edge contributes opacity only once.
			draw.DrawMask(canvas, painted, brushSource, image.Point{}, strokeMask, painted.Min, draw.Over)
			for pixelY := painted.Min.Y; pixelY < painted.Max.Y; pixelY++ {
				offset := strokeMask.PixOffset(painted.Min.X, pixelY)
				clear(strokeMask.Pix[offset : offset+painted.Dx()])
			}
		}

		if connection.IsGladiatorArena() {
			midPoint := helpers.GetVectorOnQuadraticBezierCurve(startPoint, controlPoint, endPoint, 0.5)
			this.assetProvider.DrawArenaMarker(canvas, midPoint, scale*arenaMarkerScale)
		}
	}
}

func (this *PreviewGeneratorService) drawSolidLine(
	canvas draw.Image,
	source *image.Uniform,
	start, control, end data.Vec2[float64]) image.Rectangle {
	painted := image.Rectangle{}
	previousPoint := start
	for index := range segmentsSolid {
		progress := float64(index+1) / segmentsSolid
		currentPoint := helpers.GetVectorOnQuadraticBezierCurve(start, control, end, progress)
		painted = painted.Union(this.drawLine(canvas, source, previousPoint, currentPoint))
		previousPoint = currentPoint
	}
	return painted
}

func (this *PreviewGeneratorService) drawDashedLine(
	canvas draw.Image,
	source *image.Uniform,
	start, control, end data.Vec2[float64]) image.Rectangle {
	painted := image.Rectangle{}
	period := dashLength + dashGap
	traveled := 0.0
	previousPoint := start
	for index := range segmentsDashed {
		progress := float64(index+1) / segmentsDashed
		currentPoint := helpers.GetVectorOnQuadraticBezierCurve(start, control, end, progress)
		segmentLength := math.Hypot(
			currentPoint.X-previousPoint.X,
			currentPoint.Y-previousPoint.Y)
		if math.Mod(traveled+segmentLength/2, period) < dashLength {
			painted = painted.Union(this.drawLine(canvas, source, previousPoint, currentPoint))
		}
		traveled += segmentLength
		previousPoint = currentPoint
	}
	return painted
}

func (this *PreviewGeneratorService) drawLine(
	canvas draw.Image,
	source *image.Uniform,
	start, end data.Vec2[float64]) image.Rectangle {
	painted := image.Rectangle{}
	delta := end.Subtract(start)
	steps := max(math.Abs(delta.X), math.Abs(delta.Y))
	if steps <= 0 {
		return painted
	}

	increment := delta.DivideScalar(steps)
	half := connectorLineWidth / 2
	for index := range int(math.Ceil(steps)) {
		center := start.Add(increment.MultiplyScalar(float64(index))).ToPointRounded()
		brush := image.Rect(center.X-half, center.Y-half, center.X+half+1, center.Y+half+1).
			Intersect(canvas.Bounds()) // Square brush around the center, clipped to the canvas.
		draw.Draw(canvas, brush, source, image.Point{}, draw.Src)
		painted = painted.Union(brush)
	}
	return painted
}
