package preview_service

import (
	"math"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/preview"
)

type assetFitter func(data.Vec2[float64]) data.Vec2[float64]

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
		if zone.Type == preview.ZoneTypePlayer {
			extendedFromCenter = 41.0 * scale
		}
		allowedDeviation := center - borderInset - extendedFromCenter // max center-to-center deviation
		if allowedDeviation < 1 {
			continue
		}

		// Chebyshev actualDeviation: the worse of the two axes.
		actualDeviation := math.Max(math.Abs(zone.Center.X-center), math.Abs(zone.Center.Y-center))
		if actualDeviation > allowedDeviation {
			fit = min(fit, allowedDeviation/actualDeviation)
		}
	}

	if fit >= 1 {
		return func(point data.Vec2[float64]) data.Vec2[float64] { return point }
	}

	return func(point data.Vec2[float64]) data.Vec2[float64] {
		return data.NewVec2(center+(point.X-center)*fit, center+(point.Y-center)*fit)
	}
}
