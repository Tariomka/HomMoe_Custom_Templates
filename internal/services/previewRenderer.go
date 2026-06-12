package services

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"path/filepath"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template"
)

// previewLineColor is the connection-line colour sampled from the official
// in-game template overview images (a dark warm brown, drawn under the
// zone bubbles).
var previewLineColor = color.NRGBA{R: 0x39, G: 0x11, B: 0x14, A: 0xFF}

// WritePreviewPNG rasterises the given template and writes it as a PNG into
// dir/<safeName>.png at the requested side length. The directory is created
// if missing. Returns the final path on success.
func WritePreviewPNG(dir string, template *template.RmgTemplateModel, topology config.MapTopology, side int) (string, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	safeName := SanitizeFilename(template.Name)
	if safeName == "" {
		safeName = "Generated_Template"
	}
	out := filepath.Join(dir, safeName+".png")
	img := RenderPreviewImage(template, topology, side)
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

// RenderPreviewImage rasterises the layout in the style of the official
// in-game template overview images: the parchment background, the zone
// bubbles and the connection lines are composited from sprite assets
// extracted from those images (see tools/assetgen).
func RenderPreviewImage(template *template.RmgTemplateModel, topology config.MapTopology, side int) *image.RGBA {
	assets := loadPreviewAssets()
	img := image.NewRGBA(image.Rect(0, 0, side, side))
	drawBackgroundScaled(img, assets.background)

	layout := BuildPreviewLayout(template, topology, float64(side))
	if len(layout.Positions) == 0 {
		return img
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

	// Connections: solid dark lines, ended at the bubble outlines so they
	// do not show through the open (low-quality) rings.
	lineWidth := int(math.Round(4.0 * float64(side) / 700.0))
	if lineWidth < 2 {
		lineWidth = 2
	}
	for _, conn := range layout.Connections {
		a := fitPt(conn.A)
		b := fitPt(conn.B)
		dx := float64(b.X - a.X)
		dy := float64(b.Y - a.Y)
		distance := math.Hypot(dx, dy)
		if distance < 1 {
			continue
		}
		ux := dx / distance
		uy := dy / distance
		ax := image.Pt(int(float64(a.X)+ux*drawnRadius), int(float64(a.Y)+uy*drawnRadius))
		bx := image.Pt(int(float64(b.X)-ux*drawnRadius), int(float64(b.Y)-uy*drawnRadius))
		drawThickLine(img, ax, bx, lineWidth, previewLineColor)
	}

	// Zones — non-player bubbles first, then the player bubbles on top,
	// exactly like the official overview images layer them.
	for _, zone := range layout.Zones {
		if zone.IsPlayer {
			continue
		}
		sprite := neutralSpriteFor(assets, zone)
		c := fitPt(zone.Center)
		drawSpriteScaled(img, sprite, c.X, c.Y, previewSpriteCenter, previewSpriteCenter, scale)
	}
	for _, zone := range layout.Zones {
		if !zone.IsPlayer {
			continue
		}
		sprite := playerSpriteFor(assets, zone)
		c := fitPt(zone.Center)
		drawSpriteScaled(img, sprite, c.X, c.Y, previewSpriteCenter, previewSpriteCenter, scale)
	}
	return img
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
