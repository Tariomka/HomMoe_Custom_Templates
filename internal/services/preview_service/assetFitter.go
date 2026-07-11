package preview_service

import (
	"image"
	"math"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/preview"
)

type assetFitter func(image.Point) image.Point

// newAssetFitter returns a function that pulls canvas points uniformly toward the centre
// just enough that no zone asset crosses the border line painted on the background
// (inset ≈16 px on the 700 px canvas). Player emblems extend well past the bubble outline
// (radius ≈36 in sprite space, plus the drop shadow), so they need more headroom.
func newAssetFitter(zones []preview.PreviewZone, scale float64) assetFitter {
	const borderInset = 19.0
	const centre = canvasSize / 2

	fit := 1.0
	for _, zone := range zones {
		extendedFromCenter := 28.0 * scale // artwork overhang beyond the zone centre
		if zone.IsPlayer {
			extendedFromCenter = 41.0 * scale
		}
		allowedDeviation := centre - borderInset - extendedFromCenter // max centre-to-centre deviation
		if allowedDeviation < 1 {
			continue
		}

		// Chebyshev actualDeviation: the worse of the two axes.
		actualDeviation := math.Max(
			math.Abs(float64(zone.Center.X)-centre),
			math.Abs(float64(zone.Center.Y)-centre))
		if actualDeviation > allowedDeviation {
			fit = min(fit, allowedDeviation/actualDeviation)
		}
	}

	if fit >= 1 {
		return func(p image.Point) image.Point { return p }
	}

	return func(p image.Point) image.Point {
		return image.Pt(
			int(math.Round(centre+(float64(p.X)-centre)*fit)),
			int(math.Round(centre+(float64(p.Y)-centre)*fit)),
		)
	}
}
