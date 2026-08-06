package providers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/zone_helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	zone_services "github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
)

// GladiatorArenaProvider stamps the Gladiator Arena onto a generated variant.
//
// Setting the gladiatorArena win-condition flag alone is not enough: the game
// also needs the arena itself somewhere on the map. Which wire form is used
// depends on the topology, mirroring how the shipped templates express it:
//
//   - Topologies with a hub zone (Hub & Spoke, Geometric Hub) get a
//     GladiatorArena main object inside the hub, like Blitz places one in its
//     super-treasure zone.
//   - Every other topology marks the connection between its two richest
//     neutral zones as a GladiatorArena connection, like Helltide's
//     "Win-Connection" and Symmetry's "Arena-Connection".
//   - When no neutral zone touches another (the common alternating ring
//     layouts), the arena falls back to a main object in the richest neutral
//     zone so it is never silently dropped.
type GladiatorArenaProvider struct {
	zoneClassifier *zone_services.ZoneClassifier
}

func NewGladiatorArenaProvider(zoneClassifier *zone_services.ZoneClassifier) *GladiatorArenaProvider {
	return &GladiatorArenaProvider{zoneClassifier: zoneClassifier}
}

// PlaceArena writes the arena into the variant when the configuration asks for
// the Gladiator Arena win condition. Templates without a hub and without any
// neutral zone are left untouched - there is nowhere neutral to put it.
func (this *GladiatorArenaProvider) PlaceArena(configuration config.GeneratorConfig, variant *entities.Variant) {
	if !configuration.IsGladiatorArenaMode() {
		return
	}

	if hubIndex := findHubZoneIndex(variant.Zones); hubIndex >= 0 {
		addArenaMainObject(&variant.Zones[hubIndex])
		return
	}

	if connectionIndex := this.findArenaConnectionIndex(*variant); connectionIndex >= 0 {
		variant.Connections[connectionIndex].ConnectionType = registry.GetConnectionTypeValues().GladiatorArena
		return
	}

	if zoneIndex := this.findRichestNeutralZoneIndex(variant.Zones); zoneIndex >= 0 {
		addArenaMainObject(&variant.Zones[zoneIndex])
	}
}

// findArenaConnectionIndex returns the neutral-to-neutral connection whose two
// endpoints are the richest, or -1 when the variant has none. Ties are broken
// on the connection name so the same configuration always yields the same
// template.
func (this *GladiatorArenaProvider) findArenaConnectionIndex(variant entities.Variant) int {
	qualities := this.mapNeutralZoneQualities(variant.Zones)

	bestIndex, bestScore := -1, 0
	for index, connection := range variant.Connections {
		fromQuality, okFrom := qualities[connection.From]
		toQuality, okTo := qualities[connection.To]
		if !okFrom || !okTo {
			continue
		}

		score := fromQuality.GetIndex() + toQuality.GetIndex()
		if bestIndex < 0 || score > bestScore ||
			(score == bestScore && connection.Name < variant.Connections[bestIndex].Name) {
			bestIndex, bestScore = index, score
		}
	}
	return bestIndex
}

// findRichestNeutralZoneIndex returns the highest-quality neutral zone, or -1
// when the variant has none. Ties are broken on the zone name.
func (this *GladiatorArenaProvider) findRichestNeutralZoneIndex(zones []entities.Zone) int {
	bestIndex, bestQuality := -1, neutral_zone.QualityUnknown
	for index, zone := range zones {
		if !zone_helpers.IsZoneNameNeutral(zone.Name) {
			continue
		}

		quality := this.zoneClassifier.GetQuality(zone)
		if bestIndex < 0 || quality > bestQuality ||
			(quality == bestQuality && zone.Name < zones[bestIndex].Name) {
			bestIndex, bestQuality = index, quality
		}
	}
	return bestIndex
}

func (this *GladiatorArenaProvider) mapNeutralZoneQualities(
	zones []entities.Zone) map[string]neutral_zone.Quality {
	qualities := make(map[string]neutral_zone.Quality, len(zones))
	for _, zone := range zones {
		if zone_helpers.IsZoneNameNeutral(zone.Name) {
			qualities[zone.Name] = this.zoneClassifier.GetQuality(zone)
		}
	}
	return qualities
}

func findHubZoneIndex(zones []entities.Zone) int {
	for index, zone := range zones {
		if zone_helpers.IsZoneNameHub(zone.Name) {
			return index
		}
	}
	return -1
}

// addArenaMainObject appends the arena object using the same placement Blitz
// ships with, so the in-game generator treats it identically.
func addArenaMainObject(zone *entities.Zone) {
	zone.MainObjects = append(zone.MainObjects,
		variant_content.NewObjectBuilder().
			WithTypeGladiatorArena().
			WithPlacementUniform().
			WithPlacementArgs("true", "0", "0").
			Build())
}
