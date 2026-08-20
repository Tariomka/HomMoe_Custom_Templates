package preview_service

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/zone_helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
)

// layoutFixedPositions places zones at their exact GeneratorPosition stamps,
// preserving the deterministic geometric figure built by the Square, Geometric,
// Cross, Fractal and Geometric Hub topologies. The normalized positions are
// centered and uniformly scaled to fill the padded draw area (never relaxed),
// then the zone radius is shrunk just enough to keep the closest pair from
// overlapping. extraInset (canvas-reference pixels, see csGeoHubEdgeInset)
// widens the padding for topologies whose figure should keep clear of the
// border.
func (this *PreviewLayoutService) layoutFixedPositions(
	zones []entities.Zone, side float64, extraInset float64) {
	metrics := newCanvasMetrics(side)
	if this.placeTrivial(zones, metrics) {
		return
	}

	positions := getGeneratorCoordinates(zones)
	// The margin reserves room for the largest possible zone radius, so the eventual radius (never larger) always fits.
	fitToCanvas(positions, metrics, metrics.margin+metrics.zoneRadiusMax+extraInset*metrics.scale, true)
	radius := radiusFromClosestPair(positions, metrics.zoneRadiusMax, metrics.minGap)
	this.commitPositions(zones, positions, radius)
}

// fixedGeometryEdgeInset returns the extra border padding for a
// fixed-geometry topology's preview figure. Crowded Geometric Hub figures
// (csGeoHubCrowdedMinPlayers+ players) keep a smaller inset so the players
// sit closer to the border and further from the central hub.
func fixedGeometryEdgeInset(topology config.MapTopology, zones []entities.Zone) float64 {
	if topology != config.TopologyGeometricHub {
		return 0
	}

	if countPlayerZones(zones) >= csGeoHubCrowdedMinPlayers {
		return csGeoHubEdgeInsetCrowded
	}

	return csGeoHubEdgeInset
}

// countPlayerZones counts the player (spawn) zones in the cluster.
func countPlayerZones(zones []entities.Zone) int {
	playerCount := 0
	for _, zone := range zones {
		if zone_helpers.IsZoneNamePlayer(zone.Name) {
			playerCount++
		}
	}
	return playerCount
}
