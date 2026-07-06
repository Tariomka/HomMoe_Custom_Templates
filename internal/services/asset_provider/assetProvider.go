package asset_provider

import (
	"bytes"
	"embed"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"sync"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/preview"
)

const (
	assetCenter     = 48        // all assets are 96x96
	assetFolder     = "assets/" // all asserts are stored inside ./assets folder
	backgroundAsset = "background.png"
)

var (
	//go:embed assets/*.png
	assetFileSystem embed.FS

	providerOnce      sync.Once
	providerSingleton *AssetProvider
	providerErr       error

	neutralAssetNames = []string{
		"neutral_low", "neutral_low_castle",
		"neutral_medium", "neutral_medium_castle",
		"neutral_high", "neutral_high_castle",
	}
)

type AssetProvider struct {
	background   image.Image
	players      [8]image.Image
	neutralZones map[string]image.Image
}

func NewAssetProvider() (*AssetProvider, error) {
	if err := providerSingleton.loadAssets(); err != nil {
		return nil, fmt.Errorf("failed to load assets: %v", err)
	}

	return providerSingleton, nil
}

func (this *AssetProvider) DrawBackground(canvas *image.RGBA) {
	this.ensureAssetsAreLoaded()

	canvasBounds := canvas.Bounds()
	backgroundBounds := this.background.Bounds()
	scaledPosition := data.NewVec2(
		float64(backgroundBounds.Dx())/float64(canvasBounds.Dx()),
		float64(backgroundBounds.Dy())/float64(canvasBounds.Dy()))

	for y := canvasBounds.Min.Y; y < canvasBounds.Max.Y; y++ {
		for x := canvasBounds.Min.X; x < canvasBounds.Max.X; x++ {
			interpolatedColor := this.calculateBilinearInterpolation(this.background,
				scaledPosition.MultiplyComponent(data.NewVec2(float64(x), float64(y))))
			off := canvas.PixOffset(x, y)
			canvas.Pix[off+0] = interpolatedColor.R
			canvas.Pix[off+1] = interpolatedColor.G
			canvas.Pix[off+2] = interpolatedColor.B
			canvas.Pix[off+3] = 255
		}
	}
}

func (this *AssetProvider) DrawPlayerZone(
	canvas *image.RGBA,
	zone preview.PreviewZone,
	center image.Point,
	scale float64) {
	this.ensureAssetsAreLoaded()

	sprite := this.getPlayerAsset(zone)
	this.drawSpriteScaled(canvas, sprite, center, scale)
}

func (this *AssetProvider) DrawNeutralZone(
	canvas *image.RGBA,
	zone preview.PreviewZone,
	center image.Point,
	scale float64) {
	this.ensureAssetsAreLoaded()

	sprite := this.getNeutralZoneAsset(zone)
	this.drawSpriteScaled(canvas, sprite, center, scale)
}

// drawSpriteScaled alpha-composites a sprite onto dst so that the sprite
// anchor (anchorX, anchorY) lands on dst point (cx, cy), scaled by the given
// factor. Bilinear sampling keeps the artwork smooth at non-integer scales.
func (this *AssetProvider) drawSpriteScaled(
	canvas *image.RGBA,
	asset image.Image,
	center image.Point,
	scale float64,
) {
	assetBounds := asset.Bounds()
	const roundness = 2

	// Destination rectangle covered by the scaled sprite, clipped to dst.
	left := int(float64(center.X) - assetCenter*scale)
	top := int(float64(center.Y) - assetCenter*scale)
	width := int(float64(assetBounds.Dx())*scale) + roundness // +2 covers rounding at both edges
	height := int(float64(assetBounds.Dy())*scale) + roundness
	rect := image.Rect(left, top, left+width, top+height).Intersect(canvas.Bounds())

	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			normalizedPosition := data.NewVec2(
				(float64(x)-float64(center.X))/scale+assetCenter,
				(float64(y)-float64(center.Y))/scale+assetCenter)
			interpolatedColor := this.calculateBilinearInterpolation(asset, normalizedPosition)
			if interpolatedColor.A == 0 {
				continue
			}

			// Src-over blend: out = src*alpha + dst*(1-alpha).
			offset := canvas.PixOffset(x, y)
			alpha, keep := int(interpolatedColor.A), 255-int(interpolatedColor.A)
			canvas.Pix[offset+0] = uint8((int(interpolatedColor.R)*alpha + int(canvas.Pix[offset+0])*keep) / 255)
			canvas.Pix[offset+1] = uint8((int(interpolatedColor.G)*alpha + int(canvas.Pix[offset+1])*keep) / 255)
			canvas.Pix[offset+2] = uint8((int(interpolatedColor.B)*alpha + int(canvas.Pix[offset+2])*keep) / 255)
			canvas.Pix[offset+3] = 255
		}
	}
}

// calculateBilinearInterpolation reads the asset at a fractional position and returns interpolated color value.
func (this *AssetProvider) calculateBilinearInterpolation(asset image.Image, position data.Vec2[float64]) color.NRGBA {
	calculateTextureElement := func(pixel data.Vec2[int]) data.Vec4[float64] {
		if !image.Pt(pixel.X, pixel.Y).In(asset.Bounds()) {
			return data.Vec4[float64]{}
		}
		red, green, blue, alpha := asset.At(pixel.X, pixel.Y).RGBA()
		return data.NewVec4(float64(red>>8), float64(green>>8), float64(blue>>8), float64(alpha>>8))
	}

	cornerPixel := data.Transform[float64, int](position)
	pixelOffset := data.NewVec2(position.X-float64(cornerPixel.X), position.Y-float64(cornerPixel.Y))
	premultipliedColor := data.Vec4[float64]{}
	pixelNeighbors := [4]struct {
		pixel  data.Vec2[int]
		weight float64
	}{
		{cornerPixel, (1 - pixelOffset.X) * (1 - pixelOffset.Y)},
		{data.NewVec2(cornerPixel.X+1, cornerPixel.Y), pixelOffset.X * (1 - pixelOffset.Y)},
		{data.NewVec2(cornerPixel.X, cornerPixel.Y+1), (1 - pixelOffset.X) * pixelOffset.Y},
		{data.NewVec2(cornerPixel.X+1, cornerPixel.Y+1), pixelOffset.X * pixelOffset.Y},
	}

	for _, tap := range pixelNeighbors {
		texel := calculateTextureElement(tap.pixel)
		premultipliedColor.X += tap.weight * texel.X
		premultipliedColor.Y += tap.weight * texel.Y
		premultipliedColor.Z += tap.weight * texel.Z
		premultipliedColor.W += tap.weight * texel.W
	}

	if premultipliedColor.W < 1 {
		return color.NRGBA{}
	}

	return color.NRGBA{
		R: uint8(premultipliedColor.X * 255 / premultipliedColor.W),
		G: uint8(premultipliedColor.Y * 255 / premultipliedColor.W),
		B: uint8(premultipliedColor.Z * 255 / premultipliedColor.W),
		A: uint8(premultipliedColor.W),
	}
}

func (this *AssetProvider) getNeutralZoneAsset(zone preview.PreviewZone) image.Image {
	quality := "low"
	switch zone.Tier {
	case 3:
		quality = "high"
	case 2:
		quality = "medium"
	case 1:
		quality = "low"
	}
	name := "neutral_" + quality
	if zone.HasCastle {
		name += "_castle"
	}
	return this.neutralZones[name]
}

func (this *AssetProvider) getPlayerAsset(zone preview.PreviewZone) image.Image {
	owner := min(max(zone.Owner, 1), 8)
	return this.players[owner-1]
}

func (this *AssetProvider) loadAssets() error {
	providerOnce.Do(func() {
		decode := func(name string) (image.Image, error) {
			data, err := assetFileSystem.ReadFile(assetFolder + name)
			if err != nil {
				return nil, fmt.Errorf("preview asset %s: %v", name, err)
			}

			img, err := png.Decode(bytes.NewReader(data))
			if err != nil {
				return nil, fmt.Errorf("preview asset %s: %v", name, err)
			}

			return img, nil
		}

		assetProvider := &AssetProvider{neutralZones: map[string]image.Image{}}
		var err error

		if assetProvider.background, err = decode(backgroundAsset); err != nil {
			providerErr = fmt.Errorf("failed to load background: %v", err)
			return
		}

		for i := range 8 {
			if assetProvider.players[i], err = decode(fmt.Sprintf("player_%d.png", i+1)); err != nil {
				providerErr = fmt.Errorf("failed to load player %d: %v", i+1, err)
				return
			}
		}

		for _, name := range neutralAssetNames {
			if assetProvider.neutralZones[name], err = decode(name + ".png"); err != nil {
				providerErr = fmt.Errorf("failed to load neutral asset %s: %v", name, err)
				return
			}
		}

		providerSingleton = assetProvider
	})
	return providerErr
}

func (this *AssetProvider) ensureAssetsAreLoaded() {
	if providerSingleton == nil && this.loadAssets() != nil {
		panic("failed to load assets: " + providerErr.Error())
	}
}
