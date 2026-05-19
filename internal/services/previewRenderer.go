package services

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"path/filepath"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
)

// Preview palette — ported from TemplatePreviewPngWriter colours.
var (
	previewBg         = color.NRGBA{R: 0x1C, G: 0x16, B: 0x10, A: 0xFF}
	previewFrame      = color.NRGBA{R: 0x8F, G: 0x73, B: 0x3F, A: 0xFF}
	previewBronzeFill = color.NRGBA{R: 0x65, G: 0x43, B: 0x21, A: 0xFF}
	previewBronzeEdge = color.NRGBA{R: 0xCD, G: 0x7F, B: 0x32, A: 0xFF}
	previewSilverFill = color.NRGBA{R: 0x48, G: 0x4C, B: 0x50, A: 0xFF}
	previewSilverEdge = color.NRGBA{R: 0xC0, G: 0xC0, B: 0xC0, A: 0xFF}
	previewGoldFill   = color.NRGBA{R: 0x78, G: 0x5A, B: 0x14, A: 0xFF}
	previewGoldEdge   = color.NRGBA{R: 0xFF, G: 0xD2, B: 0x32, A: 0xFF}
	previewSpawnFill  = color.NRGBA{R: 0x2A, G: 0x5A, B: 0x32, A: 0xFF}
	previewSpawnEdge  = color.NRGBA{R: 0x64, G: 0xC8, B: 0x78, A: 0xFF}
	previewHubFill    = color.NRGBA{R: 0x37, G: 0x50, B: 0x5F, A: 0xFF}
	previewHubEdge    = color.NRGBA{R: 0x82, G: 0xB4, B: 0xC8, A: 0xFF}
	previewDirectLine = color.NRGBA{R: 0xB4, G: 0x91, B: 0x3C, A: 0xFF}
	previewPortalLine = color.NRGBA{R: 0x5A, G: 0xAA, B: 0xD2, A: 0xB4}
)

// WritePreviewPNG rasterises the given template and writes it as a PNG into
// dir/<safeName>.png at the requested side length. The directory is created
// if missing. Returns the final path on success.
func WritePreviewPNG(dir string, template *models.RmgTemplate, topology models.MapTopology, side int) (string, error) {
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

// RenderPreviewImage rasterises the layout into an *image.RGBA. It uses only
// the standard library — circles are filled scanline-by-scanline, lines via
// a simple DDA with a tiny brush.
func RenderPreviewImage(template *models.RmgTemplate, topology models.MapTopology, side int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, side, side))
	fillRect(img, img.Bounds(), previewBg)

	layout := BuildPreviewLayout(template, topology, float64(side))
	if len(layout.Positions) == 0 {
		return img
	}
	radius := layout.ZoneRadius

	// Frame.
	strokeRect(img, image.Rect(4, 4, side-4, side-4), 2, previewFrame)

	// Connections.
	for _, conn := range layout.Connections {
		dx := float64(conn.B.X - conn.A.X)
		dy := float64(conn.B.Y - conn.A.Y)
		distance := math.Hypot(dx, dy)
		if distance < 1 {
			continue
		}
		ux := dx / distance
		uy := dy / distance
		ax := image.Pt(int(float64(conn.A.X)+ux*float64(radius)), int(float64(conn.A.Y)+uy*float64(radius)))
		bx := image.Pt(int(float64(conn.B.X)-ux*float64(radius)), int(float64(conn.B.Y)-uy*float64(radius)))
		lineColor := previewDirectLine
		lineWidth := 3
		if conn.Portal {
			lineColor = previewPortalLine
			lineWidth = 2
		}
		drawThickLine(img, ax, bx, lineWidth, lineColor)
	}
	// Zones — non-player first, then player on top.
	labelColor := color.NRGBA{R: 0xF8, G: 0xE8, B: 0xC0, A: 0xFF}
	badgeColor := color.NRGBA{R: 0xFF, G: 0xE8, B: 0x90, A: 0xFF}
	for _, zone := range layout.Zones {
		if zone.IsPlayer {
			continue
		}
		fill, edge := pngZoneColors(zone)
		zoneRadius := radius
		if zone.IsHub && zoneRadius < 28 {
			zoneRadius = 28
		}
		fillCircle(img, zone.Center, zoneRadius, fill)
		strokeCircle(img, zone.Center, zoneRadius, 2, edge)
		drawBitmapTextCentered(img, zone.Center, pngZoneLabel(zone), 2, labelColor)
		if zone.HasCastle && zone.Castles > 0 {
			badgePos := image.Pt(zone.Center.X+zoneRadius/2, zone.Center.Y+zoneRadius/2)
			drawBitmapTextCentered(img, badgePos, intToString(zone.Castles), 1, badgeColor)
		}
	}
	for _, zone := range layout.Zones {
		if !zone.IsPlayer {
			continue
		}
		fill, edge := pngZoneColors(zone)
		fillCircle(img, zone.Center, radius, fill)
		strokeCircle(img, zone.Center, radius, 2, edge)
		drawBitmapTextCentered(img, zone.Center, pngZoneLabel(zone), 2, labelColor)
		if zone.HasCastle && zone.Castles > 0 {
			badgePos := image.Pt(zone.Center.X+radius/2, zone.Center.Y+radius/2)
			drawBitmapTextCentered(img, badgePos, intToString(zone.Castles), 1, badgeColor)
		}
	}
	return img
}

func pngZoneColors(zone PreviewZone) (fill, edge color.NRGBA) {
	switch {
	case zone.IsPlayer:
		return previewSpawnFill, previewSpawnEdge
	case zone.IsHub:
		return previewHubFill, previewHubEdge
	}
	switch zone.Tier {
	case 3:
		return previewGoldFill, previewGoldEdge
	case 2:
		return previewSilverFill, previewSilverEdge
	default:
		return previewBronzeFill, previewBronzeEdge
	}
}

func pngZoneLabel(zone PreviewZone) string {
	if zone.IsPlayer {
		if zone.Owner > 0 {
			return "P" + intToString(zone.Owner)
		}
		// Spawn-1 / Spawn-2 → "P1"…
		if len(zone.Letter) > 0 {
			return "P" + zone.Letter
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

func intToString(value int) string {
	if value == 0 {
		return "0"
	}
	negative := value < 0
	if negative {
		value = -value
	}
	var buf [20]byte
	i := len(buf)
	for value > 0 {
		i--
		buf[i] = byte('0' + value%10)
		value /= 10
	}
	if negative {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}

// ── Pixel-pushing primitives ────────────────────────────────────────

func fillRect(img *image.RGBA, rect image.Rectangle, fillColor color.NRGBA) {
	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			img.SetRGBA(x, y, color.RGBA{R: fillColor.R, G: fillColor.G, B: fillColor.B, A: fillColor.A})
		}
	}
}

func strokeRect(img *image.RGBA, rect image.Rectangle, width int, strokeColor color.NRGBA) {
	fillRect(img, image.Rect(rect.Min.X, rect.Min.Y, rect.Max.X, rect.Min.Y+width), strokeColor)
	fillRect(img, image.Rect(rect.Min.X, rect.Max.Y-width, rect.Max.X, rect.Max.Y), strokeColor)
	fillRect(img, image.Rect(rect.Min.X, rect.Min.Y, rect.Min.X+width, rect.Max.Y), strokeColor)
	fillRect(img, image.Rect(rect.Max.X-width, rect.Min.Y, rect.Max.X, rect.Max.Y), strokeColor)
}

func fillCircle(img *image.RGBA, center image.Point, radius int, fillColor color.NRGBA) {
	radiusSq := radius * radius
	for dy := -radius; dy <= radius; dy++ {
		for dx := -radius; dx <= radius; dx++ {
			if dx*dx+dy*dy <= radiusSq {
				x := center.X + dx
				y := center.Y + dy
				if x < 0 || y < 0 || x >= img.Rect.Max.X || y >= img.Rect.Max.Y {
					continue
				}
				img.SetRGBA(x, y, color.RGBA{R: fillColor.R, G: fillColor.G, B: fillColor.B, A: fillColor.A})
			}
		}
	}
}

func strokeCircle(img *image.RGBA, center image.Point, radius, width int, strokeColor color.NRGBA) {
	outerSq := radius * radius
	innerRadius := radius - width
	innerSq := innerRadius * innerRadius
	if innerSq < 0 {
		innerSq = 0
	}
	for dy := -radius; dy <= radius; dy++ {
		for dx := -radius; dx <= radius; dx++ {
			distSq := dx*dx + dy*dy
			if distSq <= outerSq && distSq >= innerSq {
				x := center.X + dx
				y := center.Y + dy
				if x < 0 || y < 0 || x >= img.Rect.Max.X || y >= img.Rect.Max.Y {
					continue
				}
				img.SetRGBA(x, y, color.RGBA{R: strokeColor.R, G: strokeColor.G, B: strokeColor.B, A: strokeColor.A})
			}
		}
	}
}

// drawThickLine draws a width-pixel-thick line via DDA with a square brush.
func drawThickLine(img *image.RGBA, a, b image.Point, width int, lineColor color.NRGBA) {
	dx := b.X - a.X
	dy := b.Y - a.Y
	steps := dx
	if -dx > steps {
		steps = -dx
	}
	if dy > steps {
		steps = dy
	}
	if -dy > steps {
		steps = -dy
	}
	if steps <= 0 {
		fillCircle(img, a, width, lineColor)
		return
	}
	xinc := float64(dx) / float64(steps)
	yinc := float64(dy) / float64(steps)
	x := float64(a.X)
	y := float64(a.Y)
	half := width / 2
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
				img.SetRGBA(xx, yy, color.RGBA{R: lineColor.R, G: lineColor.G, B: lineColor.B, A: lineColor.A})
			}
		}
		x += xinc
		y += yinc
	}
}

// ── Bitmap font for PNG annotations ──────────────────────────────────
// A tiny 5×7 pixel font covering 0-9, A-Z, and a few symbols.
// Each glyph is 5 columns wide; each column is a byte with 7 bit-rows
// (bit 0 = top row).

var bitmapGlyphs = map[byte][5]byte{
	'0': {0x3E, 0x51, 0x49, 0x45, 0x3E},
	'1': {0x00, 0x42, 0x7F, 0x40, 0x00},
	'2': {0x42, 0x61, 0x51, 0x49, 0x46},
	'3': {0x21, 0x41, 0x45, 0x4B, 0x31},
	'4': {0x18, 0x14, 0x12, 0x7F, 0x10},
	'5': {0x27, 0x45, 0x45, 0x45, 0x39},
	'6': {0x3C, 0x4A, 0x49, 0x49, 0x30},
	'7': {0x01, 0x71, 0x09, 0x05, 0x03},
	'8': {0x36, 0x49, 0x49, 0x49, 0x36},
	'9': {0x06, 0x49, 0x49, 0x29, 0x1E},
	'A': {0x7E, 0x11, 0x11, 0x11, 0x7E},
	'B': {0x7F, 0x49, 0x49, 0x49, 0x36},
	'C': {0x3E, 0x41, 0x41, 0x41, 0x22},
	'D': {0x7F, 0x41, 0x41, 0x22, 0x1C},
	'E': {0x7F, 0x49, 0x49, 0x49, 0x41},
	'F': {0x7F, 0x09, 0x09, 0x09, 0x01},
	'G': {0x3E, 0x41, 0x49, 0x49, 0x7A},
	'H': {0x7F, 0x08, 0x08, 0x08, 0x7F},
	'I': {0x00, 0x41, 0x7F, 0x41, 0x00},
	'J': {0x20, 0x40, 0x41, 0x3F, 0x01},
	'K': {0x7F, 0x08, 0x14, 0x22, 0x41},
	'L': {0x7F, 0x40, 0x40, 0x40, 0x40},
	'M': {0x7F, 0x02, 0x0C, 0x02, 0x7F},
	'N': {0x7F, 0x04, 0x08, 0x10, 0x7F},
	'O': {0x3E, 0x41, 0x41, 0x41, 0x3E},
	'P': {0x7F, 0x09, 0x09, 0x09, 0x06},
	'Q': {0x3E, 0x41, 0x51, 0x21, 0x5E},
	'R': {0x7F, 0x09, 0x19, 0x29, 0x46},
	'S': {0x46, 0x49, 0x49, 0x49, 0x31},
	'T': {0x01, 0x01, 0x7F, 0x01, 0x01},
	'U': {0x3F, 0x40, 0x40, 0x40, 0x3F},
	'V': {0x1F, 0x20, 0x40, 0x20, 0x1F},
	'W': {0x3F, 0x40, 0x38, 0x40, 0x3F},
	'X': {0x63, 0x14, 0x08, 0x14, 0x63},
	'Y': {0x07, 0x08, 0x70, 0x08, 0x07},
	'Z': {0x61, 0x51, 0x49, 0x45, 0x43},
}

// drawBitmapText draws a string onto img at (x, y) with the given scale.
func drawBitmapText(img *image.RGBA, x, y int, text string, scale int, textColor color.NRGBA) {
	if scale < 1 {
		scale = 1
	}
	rgba := color.RGBA{R: textColor.R, G: textColor.G, B: textColor.B, A: textColor.A}
	cursorX := x
	for i := 0; i < len(text); i++ {
		ch := text[i]
		if ch >= 'a' && ch <= 'z' {
			ch -= 32
		}
		glyph, ok := bitmapGlyphs[ch]
		if !ok {
			cursorX += 4 * scale // space for unknown chars
			continue
		}
		for col := 0; col < 5; col++ {
			bits := glyph[col]
			for row := 0; row < 7; row++ {
				if bits&(1<<uint(row)) != 0 {
					for sy := 0; sy < scale; sy++ {
						for sx := 0; sx < scale; sx++ {
							px := cursorX + col*scale + sx
							py := y + row*scale + sy
							if px >= 0 && py >= 0 && px < img.Rect.Max.X && py < img.Rect.Max.Y {
								img.SetRGBA(px, py, rgba)
							}
						}
					}
				}
			}
		}
		cursorX += 6 * scale // 5 pixel cols + 1 pixel gap
	}
}

// bitmapTextWidth returns the pixel width of a string at the given scale.
func bitmapTextWidth(text string, scale int) int {
	if len(text) == 0 {
		return 0
	}
	return len(text)*6*scale - scale // subtract trailing gap
}

// drawBitmapTextCentered draws a string centered on the given point.
func drawBitmapTextCentered(img *image.RGBA, center image.Point, text string, scale int, textColor color.NRGBA) {
	if text == "" {
		return
	}
	textWidth := bitmapTextWidth(text, scale)
	textHeight := 7 * scale
	x := center.X - textWidth/2
	y := center.Y - textHeight/2
	drawBitmapText(img, x, y, text, scale, textColor)
}
