package preview_service

import (
	"image"
	"math"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/preview"
)

type assetFitter func(image.Point) image.Point

// newAssetFitter returns a function that pulls canvas points uniformly toward the center
// just enough that no zone asset crosses the border line painted on the background
// (inset ≈16 px on the 700 px canvas). Player emblems extend well past the bubble outline
// (radius ≈36 in sprite space, plus the drop shadow), so they need more headroom.
func newAssetFitter(zones []preview.Zone, scale float64) assetFitter {
	const borderInset = 19.0
	const center = canvasSize / 2

	fit := 1.0
	for _, zone := range zones {
		extendedFromCenter := 28.0 * scale // artwork overhang beyond the zone center
		if zone.IsPlayer {
			extendedFromCenter = 41.0 * scale
		}
		allowedDeviation := center - borderInset - extendedFromCenter // max center-to-center deviation
		if allowedDeviation < 1 {
			continue
		}

		// Chebyshev actualDeviation: the worse of the two axes.
		actualDeviation := math.Max(
			math.Abs(float64(zone.Center.X)-center),
			math.Abs(float64(zone.Center.Y)-center))
		if actualDeviation > allowedDeviation {
			fit = min(fit, allowedDeviation/actualDeviation)
		}
	}

	if fit >= 1 {
		return func(p image.Point) image.Point { return p }
	}

	return func(p image.Point) image.Point {
		return image.Pt(
			int(math.Round(center+(float64(p.X)-center)*fit)),
			int(math.Round(center+(float64(p.Y)-center)*fit)),
		)
	}
}
