package services

import (
	"bytes"
	"embed"
	"fmt"
	"image"
	"image/png"
	"sync"
)

// The 15 sprite assets extracted from the official in-game template overview
// images by tools/assetgen: the parchment background, the eight numbered
// player bubbles and the six neutral zone bubbles (low/medium/high quality,
// each with and without the castle glyph).
//
//go:embed previewassets/*.png
var previewAssetsFS embed.FS

// previewSpriteCenter is the bubble centre inside every (96x96) marker
// sprite, matching cropSize/2 in tools/assetgen.
const previewSpriteCenter = 48

// previewSpriteRadius is the bubble outline radius inside the marker
// sprites (the official bubbles are ~21 px in a 700 px canvas).
const previewSpriteRadius = 21.0

type previewAssets struct {
	background image.Image
	players    [8]image.Image // players[i] is the "i+1" bubble
	neutral    map[string]image.Image
}

var (
	previewAssetsOnce sync.Once
	previewAssetsData *previewAssets
)

func loadPreviewAssets() *previewAssets {
	previewAssetsOnce.Do(func() {
		decode := func(name string) image.Image {
			data, err := previewAssetsFS.ReadFile("previewassets/" + name)
			if err != nil {
				panic(fmt.Sprintf("preview asset %s: %v", name, err))
			}
			img, err := png.Decode(bytes.NewReader(data))
			if err != nil {
				panic(fmt.Sprintf("preview asset %s: %v", name, err))
			}
			return img
		}
		a := &previewAssets{neutral: map[string]image.Image{}}
		a.background = decode("background.png")
		for i := range 8 {
			a.players[i] = decode(fmt.Sprintf("player_%d.png", i+1))
		}
		for _, name := range []string{
			"neutral_low", "neutral_low_castle",
			"neutral_medium", "neutral_medium_castle",
			"neutral_high", "neutral_high_castle",
		} {
			a.neutral[name] = decode(name + ".png")
		}
		previewAssetsData = a
	})
	return previewAssetsData
}

// neutralSpriteFor maps a zone to its neutral bubble sprite: the zone tier
// picks the quality fill (gold/silver/none) and the castle glyph marks
// zones holding a city.
func neutralSpriteFor(assets *previewAssets, zone PreviewZone) image.Image {
	quality := "low"
	switch zone.Tier {
	case 3:
		quality = "high"
	case 2:
		quality = "medium"
	}
	name := "neutral_" + quality
	if zone.HasCastle {
		name += "_castle"
	}
	return assets.neutral[name]
}

// playerSpriteFor maps a player zone to its numbered bubble sprite.
func playerSpriteFor(assets *previewAssets, zone PreviewZone) image.Image {
	owner := min(max(zone.Owner, 1), 8)
	return assets.players[owner-1]
}

// drawSpriteScaled alpha-composites a sprite onto dst so that the sprite
// anchor (anchorX, anchorY) lands on dst point (cx, cy), scaled by the given
// factor. Bilinear sampling keeps the artwork smooth at non-integer scales.
func drawSpriteScaled(dst *image.RGBA, sprite image.Image, cx, cy int, anchorX, anchorY, scale float64) {
	sb := sprite.Bounds()
	sw, sh := float64(sb.Dx()), float64(sb.Dy())
	// Destination extent of the scaled sprite.
	x0 := int(float64(cx) - anchorX*scale)
	y0 := int(float64(cy) - anchorY*scale)
	x1 := x0 + int(sw*scale+1)
	y1 := y0 + int(sh*scale+1)
	db := dst.Bounds()
	for y := y0; y <= y1; y++ {
		if y < db.Min.Y || y >= db.Max.Y {
			continue
		}
		for x := x0; x <= x1; x++ {
			if x < db.Min.X || x >= db.Max.X {
				continue
			}
			// Source position (bilinear).
			sx := (float64(x)-float64(cx))/scale + anchorX
			sy := (float64(y)-float64(cy))/scale + anchorY
			r, g, b, a := bilinearSample(sprite, sx, sy)
			if a == 0 {
				continue
			}
			dr, dg, db2, _ := dst.At(x, y).RGBA()
			inv := 255 - a
			outR := (r*a + int(dr>>8)*inv) / 255
			outG := (g*a + int(dg>>8)*inv) / 255
			outB := (b*a + int(db2>>8)*inv) / 255
			dst.Pix[dst.PixOffset(x, y)+0] = uint8(outR)
			dst.Pix[dst.PixOffset(x, y)+1] = uint8(outG)
			dst.Pix[dst.PixOffset(x, y)+2] = uint8(outB)
			dst.Pix[dst.PixOffset(x, y)+3] = 255
		}
	}
}

// bilinearSample reads the sprite at a fractional position with
// premultiplied-alpha-correct bilinear interpolation. Returns straight
// (non-premultiplied) r, g, b and the alpha, all 0-255.
func bilinearSample(img image.Image, fx, fy float64) (int, int, int, int) {
	b := img.Bounds()
	x0 := int(fx)
	y0 := int(fy)
	tx := fx - float64(x0)
	ty := fy - float64(y0)

	var pr, pg, pb, pa float64
	for dy := 0; dy <= 1; dy++ {
		for dx := 0; dx <= 1; dx++ {
			wgt := (1 - tx) * (1 - ty)
			if dx == 1 {
				wgt = tx * (1 - ty)
			}
			if dy == 1 {
				if dx == 1 {
					wgt = tx * ty
				} else {
					wgt = (1 - tx) * ty
				}
			}
			if wgt == 0 {
				continue
			}
			px := x0 + dx
			py := y0 + dy
			if px < b.Min.X || py < b.Min.Y || px >= b.Max.X || py >= b.Max.Y {
				continue // outside sprite = transparent
			}
			r, g, bb, a := img.At(px, py).RGBA() // premultiplied 0-65535
			pr += wgt * float64(r>>8)
			pg += wgt * float64(g>>8)
			pb += wgt * float64(bb>>8)
			pa += wgt * float64(a>>8)
		}
	}
	if pa < 1 {
		return 0, 0, 0, 0
	}
	// Un-premultiply for the src-over blend in drawSpriteScaled.
	return int(pr * 255 / pa), int(pg * 255 / pa), int(pb * 255 / pa), int(pa)
}

// drawBackgroundScaled fills dst with the parchment background asset,
// scaled to the destination size.
func drawBackgroundScaled(dst *image.RGBA, bg image.Image) {
	db := dst.Bounds()
	bb := bg.Bounds()
	sx := float64(bb.Dx()) / float64(db.Dx())
	sy := float64(bb.Dy()) / float64(db.Dy())
	for y := db.Min.Y; y < db.Max.Y; y++ {
		for x := db.Min.X; x < db.Max.X; x++ {
			r, g, b, _ := bilinearSample(bg, float64(x)*sx, float64(y)*sy)
			off := dst.PixOffset(x, y)
			dst.Pix[off+0] = uint8(r)
			dst.Pix[off+1] = uint8(g)
			dst.Pix[off+2] = uint8(b)
			dst.Pix[off+3] = 255
		}
	}
}
