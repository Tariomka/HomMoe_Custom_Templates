package services

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"path/filepath"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/asset_provider"
)

// previewLineColor is the connection-line color sampled from the official
// in-game template overview images (a dark warm brown, drawn under the
// zone bubbles).
var previewLineColor = color.NRGBA{R: 0x39, G: 0x11, B: 0x14, A: 0xFF}

// WritePreviewPNG rasterize the given template and writes it as a PNG into
// dir/<safeName>.png at the requested side length. The directory is created
// if missing. Returns the final path on success.
func WritePreviewPNG(dir string, template *entities.RmgTemplate, topology config.MapTopology) (string, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	safeName := helpers.SanitizeFilename(template.Name)
	if safeName == "" {
		safeName = "Generated_Template"
	}
	out := filepath.Join(dir, safeName+".png")
	img := RenderPreviewImage(template, topology, 700)
	file, err := os.Create(out)
	if err != nil {
		return "", err
	}
	defer file.Close()
	if err := png.Encode(file, img); err != nil {
		return "", err
	}
	return out, nil
}

// RenderPreviewImage rasterize the layout in the style of the official
// in-game template overview images: the parchment background, the zone
// bubbles and the connection lines are composited from sprite assets
// extracted from those images.
func RenderPreviewImage(template *entities.RmgTemplate, topology config.MapTopology, side int) *image.RGBA {
	assetProvider, err := asset_provider.NewAssetProvider()
	if err != nil {
		return nil
	}
	side = 700
	canvas := image.NewRGBA(image.Rect(0, 0, side, side))
	assetProvider.DrawBackground(canvas)

	layout := BuildPreviewLayout(template, topology, float64(side))
	if len(layout.Positions) == 0 {
		return canvas
	}

	// Sprite scale: follow the layout's zone radius for dense maps, but cap
	// near the official proportions (bubble radius ≈21 px on a 700 px
	// canvas) so sparse maps keep the in-game look instead of ballooning.
	maxScale := 1.15 * float64(side) / 700.0
	scale := float64(layout.ZoneRadius) / previewSpriteRadius
	if scale > maxScale {
		scale = maxScale
	}
	drawnRadius := previewSpriteRadius * scale

	// Keep every sprite inside the border line painted on the parchment
	// background (inset ≈16 px on the 700 px canvas): if any zone artwork
	// would cross it, pull all positions uniformly towards the centre.
	// Player emblems extend well past the bubble outline (radius ≈36 in
	// sprite space, plus the drop shadow).
	borderInset := 19.0 * float64(side) / 700.0
	centre := float64(side) / 2
	fit := 1.0
	for _, zone := range layout.Zones {
		extent := 28.0 * scale
		if zone.IsPlayer {
			extent = 41.0 * scale
		}
		allowed := centre - borderInset - extent
		if allowed < 1 {
			continue
		}
		dev := math.Max(math.Abs(float64(zone.Center.X)-centre), math.Abs(float64(zone.Center.Y)-centre))
		if dev > allowed {
			if k := allowed / dev; k < fit {
				fit = k
			}
		}
	}
	fitPt := func(p image.Point) image.Point {
		if fit >= 1 {
			return p
		}
		return image.Pt(
			int(math.Round(centre+(float64(p.X)-centre)*fit)),
			int(math.Round(centre+(float64(p.Y)-centre)*fit)),
		)
	}

	// Connections: solid dark lines for direct connections, dashed
	// (perforated) lines for portal connections, ended at the bubble
	// outlines so they do not show through the open (low-quality) rings.
	lineWidth := max(int(math.Round(4.0*float64(side)/700.0)), 2)

	// Dash geometry for portal connections, scaled with the canvas.
	dashOn := 9.0 * float64(side) / 700.0
	dashOff := 13.0 * float64(side) / 700.0
	for _, conn := range layout.Connections {
		a := fitPt(conn.A)
		b := fitPt(conn.B)
		ctrl := fitPt(conn.Ctrl)
		// Trim both ends back to the bubble outlines along the tangent toward
		// the control point, then stroke the quadratic curve. A lone edge has
		// its control point on the midpoint, so it renders as a straight line;
		// parallel edges fan out into distinct curves.
		ax, ay, ok1 := trimTowardPoint(a, ctrl, drawnRadius)
		bx, by, ok2 := trimTowardPoint(b, ctrl, drawnRadius)
		if !ok1 || !ok2 {
			continue
		}
		if conn.Portal {
			drawDashedQuadratic(canvas,
				image.Pt(ax, ay), ctrl, image.Pt(bx, by), lineWidth, previewLineColor, dashOn, dashOff)
		} else {
			drawThickQuadratic(canvas, image.Pt(ax, ay), ctrl, image.Pt(bx, by), lineWidth, previewLineColor)
		}
	}

	// Zones - non-player bubbles first, then the player bubbles on top,
	// exactly like the official overview images layer them.
	for _, zone := range layout.Zones {
		if zone.IsPlayer {
			continue
		}

		assetProvider.DrawNeutralZone(canvas, zone, fitPt(zone.Center), scale)
	}
	for _, zone := range layout.Zones {
		if !zone.IsPlayer {
			continue
		}

		assetProvider.DrawPlayerZone(canvas, zone, fitPt(zone.Center), scale)
	}
	return canvas
}

// trimTowardPoint returns `from` moved toward `toward` by `dist` pixels. ok is
// false when the two points are coincident.
func trimTowardPoint(from, toward image.Point, dist float64) (x, y int, ok bool) {
	dx := float64(toward.X - from.X)
	dy := float64(toward.Y - from.Y)
	d := math.Hypot(dx, dy)
	if d < 1 {
		return 0, 0, false
	}
	return int(float64(from.X) + dx/d*dist), int(float64(from.Y) + dy/d*dist), true
}

// drawThickQuadratic rasterises a quadratic Bézier (start→ctrl→end) as a chain
// of short thick line segments. With ctrl on the midpoint the samples are
// collinear, so a lone edge is drawn as a straight line.
func drawThickQuadratic(img *image.RGBA, start, ctrl, end image.Point, width int, lineColor color.NRGBA) {
	const segments = 24
	prev := start
	for i := 1; i <= segments; i++ {
		t := float64(i) / float64(segments)
		mt := 1 - t
		x := mt*mt*float64(start.X) + 2*mt*t*float64(ctrl.X) + t*t*float64(end.X)
		y := mt*mt*float64(start.Y) + 2*mt*t*float64(ctrl.Y) + t*t*float64(end.Y)
		curr := image.Pt(int(math.Round(x)), int(math.Round(y)))
		drawThickLine(img, prev, curr, width, lineColor)
		prev = curr
	}
}

// drawDashedQuadratic rasterises a quadratic Bézier (start→ctrl→end) as a
// perforated line: it walks the curve in fine steps, accumulating arc length,
// and only strokes the "on" portions of the dashOn/dashOff dash pattern. Used
// for portal connections so they read differently from solid direct lines.
func drawDashedQuadratic(
	img *image.RGBA,
	start, ctrl, end image.Point,
	width int,
	lineColor color.NRGBA,
	dashOn, dashOff float64,
) {
	if dashOn <= 0 {
		drawThickQuadratic(img, start, ctrl, end, width, lineColor)
		return
	}
	if dashOff < 0 {
		dashOff = 0
	}
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
			drawThickLine(img, prev, curr, width, lineColor)
		}
		traveled += segLen
		prev = curr
	}
}

// drawThickLine draws a width-pixel-thick line via DDA with a square brush.
func drawThickLine(img *image.RGBA, a, b image.Point, width int, lineColor color.NRGBA) {
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
	half := width / 2
	rgba := color.RGBA(lineColor)
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
				img.SetRGBA(xx, yy, rgba)
			}
		}
		x += xinc
		y += yinc
	}
}
