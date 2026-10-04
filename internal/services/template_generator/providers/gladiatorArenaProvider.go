package providers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/zone_helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/provider_interfaces"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones/zone_interfaces"
)

// GladiatorArenaProvider stamps the Gladiator Arena onto a generated variant.
//
// Setting the gladiatorArena win-condition flag alone is not enough: the game
// also needs the arena itself somewhere on the map. Which wire form is used
// depends on the topology, mirroring how the shipped templates express it:
//
//   - Topologies with a hub zone (Geometric Hub) get GladiatorArena main object
//     inside the hub, like Blitz places one in its super-treasure zone.
//   - Every other topology marks the connection between its two richest
//     neutral zones as a GladiatorArena connection, like Helltide's
//     "Win-Connection" and Symmetry's "Arena-Connection".
//   - When no neutral zone touches another (the common alternating ring
//     layouts), the arena falls back to a main object in the richest neutral
//     zone so it is never silently dropped.
type GladiatorArenaProvider struct {
	tierService zone_interfaces.IZoneTierService
}

func NewGladiatorArenaProvider(
	tierService zone_interfaces.IZoneTierService) provider_interfaces.IGladiatorArenaProvider {
	return &GladiatorArenaProvider{tierService: tierService}
}

// PlaceArena writes the arena into the variant when the configuration asks for the Gladiator Arena win condition.
func (this *GladiatorArenaProvider) PlaceArena(configuration config.GeneratorConfig, variant *template_model.Variant) {
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
// endpoints are the richest, or -1 when the variant has none.
func (this *GladiatorArenaProvider) findArenaConnectionIndex(variant template_model.Variant) int {
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
// when the variant has none.
func (this *GladiatorArenaProvider) findRichestNeutralZoneIndex(zones []template_model.Zone) int {
	bestIndex, bestQuality := -1, neutral_zone.QualityUnknown
	for index, zone := range zones {
		if !zone_helpers.IsZoneNameNeutral(zone.Name) {
			continue
		}

		quality := this.tierService.ResolveQuality(zone)
		if bestIndex < 0 || quality > bestQuality ||
			(quality == bestQuality && zone.Name < zones[bestIndex].Name) {
			bestIndex, bestQuality = index, quality
		}
	}
	return bestIndex
}

func (this *GladiatorArenaProvider) mapNeutralZoneQualities(
	zones []template_model.Zone) map[string]neutral_zone.Quality {
	qualities := make(map[string]neutral_zone.Quality, len(zones))
	for _, zone := range zones {
		if zone_helpers.IsZoneNameNeutral(zone.Name) {
			qualities[zone.Name] = this.tierService.ResolveQuality(zone)
		}
	}
	return qualities
}

func findHubZoneIndex(zones []template_model.Zone) int {
	for index, zone := range zones {
		if zone_helpers.IsZoneNameHub(zone.Name) {
			return index
		}
	}
	return -1
}

func addArenaMainObject(zone *template_model.Zone) {
	zone.MainObjects = append(zone.MainObjects,
		variant_content.NewObjectBuilder().
			WithTypeGladiatorArena().
			WithPlacementUniform().
			WithPlacementArgs("true", "0", "0").
			Build())
}
