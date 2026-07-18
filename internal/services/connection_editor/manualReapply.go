package connection_editor

// Manual-edit reapply support: pushing changed castle-count options into the
// manually edited zone snapshot. The snapshot itself is authoritative on
// regeneration - castle counts are the ONLY generator options that override
// manual edits, and only when they changed since the last generation.

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/zone_helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutralZone"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/preview"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/base"
)

// ApplyCastleSettingChanges rewrites the castles of the manually edited zones
// whose castle-count option changed, in place. Everything else about the
// zones - quality, guard values, pools, sizes, positions - stays untouched:
//   - spawn zones follow the player-castle options (the spawn castle itself,
//     and with it the player assignment, is preserved);
//   - hub zones follow the hub-castle option;
//   - neutral zones follow the simple-mode count (every neutral zone) or, in
//     advanced mode, the per-tier count matching the zone's CURRENT (possibly
//     manually re-tiered) quality - castle-less zones then keep their
//     no-castle plan.
func ApplyCastleSettingChanges(
	zones []entities.Zone,
	changes editor_state_dto.CastleSettingChanges,
	configuration *config.GeneratorConfig) {
	if !changes.Any() {
		return
	}
	tuning := models.NewGenerationTuning(configuration, len(zones))
	topology := base.NewTopologyBase()
	for i := range zones {
		zone := &zones[i]
		switch zone_helpers.GetZoneTypeFromName(zone.Name) {
		case preview.ZoneTypePlayer:
			if changes.PlayerCastles {
				rebuildSpawnZoneCastles(zone, configuration, tuning, &topology)
			}
		case preview.ZoneTypeHub:
			if changes.Hub {
				rebuildHubZoneCastles(zone, configuration.ZoneConfiguration.Advanced.HubZoneCastles, tuning, &topology)
			}
		case preview.ZoneTypeNeutral:
			if count, update := neutralCastleTarget(*zone, changes, configuration); update {
				SetNeutralZoneCastleCount(zone, count, tuning)
			}
		case preview.ZoneTypeUnknown:
		}
	}
}

// neutralCastleTarget decides whether a changed option drives the neutral
// zone's castle count, and with which count. Simple mode has a single count
// for every neutral zone; advanced mode has per-tier counts that apply only
// to with-castle zones of the matching quality.
func neutralCastleTarget(
	zone entities.Zone,
	changes editor_state_dto.CastleSettingChanges,
	configuration *config.GeneratorConfig) (int, bool) {
	zoneConfiguration := configuration.ZoneConfiguration
	if changes.NeutralSimple {
		return helpers.Clamp(zoneConfiguration.NeutralZoneCastles, 0, 4), true
	}

	if CountZoneCastles(zone) == 0 {
		return 0, false
	}

	switch neutralZone.GetQualityFrom(zone) {
	case neutralZone.QualityHighest:
		if changes.Hub {
			return helpers.Clamp(zoneConfiguration.Advanced.HubZoneCastles, 0, 4), true
		}
	case neutralZone.QualityHigh:
		if changes.NeutralHigh {
			return helpers.Clamp(zoneConfiguration.Advanced.NeutralHighCastlesPerZone, 0, 4), true
		}
	case neutralZone.QualityMedium:
		if changes.NeutralMedium {
			return helpers.Clamp(zoneConfiguration.Advanced.NeutralMediumCastlesPerZone, 0, 4), true
		}
	case neutralZone.QualityLow:
		if changes.NeutralLow {
			return helpers.Clamp(zoneConfiguration.Advanced.NeutralLowCastlesPerZone, 0, 4), true
		}
	case neutralZone.QualityLowest:
		if changes.NeutralLowest {
			return helpers.Clamp(zoneConfiguration.Advanced.NeutralLowestCastlesPerZone, 0, 4), true
		}
	case neutralZone.QualityUnknown:
	}
	return 0, false
}

// SetNeutralZoneCastleCount rebuilds only the zone's City castles for the new
// count, keeping the quality profile, guard values, content pools and any
// non-castle main objects (abandoned outposts) untouched - unlike
// ApplyNeutralZoneQuality, which re-profiles the whole zone.
func SetNeutralZoneCastleCount(zone *entities.Zone, castleCount int, tuning models.GenerationTuning) {
	profile := neutralZone.NewNeutralZoneProfile(neutralZone.GetQualityFrom(*zone))
	preserved, isHoldCity := splitOutNonCastles(zone.MainObjects)
	zone.MainObjects = append(
		base.CreateNeutralZoneCastles(profile, tuning, castleCount, isHoldCity),
		preserved...)
	rebuildCastleRoads(zone)
}

// rebuildSpawnZoneCastles rebuilds a player spawn zone's extra castles for the
// current player-castle options. The spawn castle (main object 0) is kept
// verbatim so the player assignment and faction survive.
func rebuildSpawnZoneCastles(
	zone *entities.Zone,
	configuration *config.GeneratorConfig,
	tuning models.GenerationTuning,
	topology *base.TopologyBase) {
	if len(zone.MainObjects) == 0 || zone.MainObjects[0].Type != registry.GetMainObjectTypeValues().Spawn {
		return
	}

	spawnCastle := zone.MainObjects[0]
	matchFactions := configuration.MatchPlayerCastleFactions
	mainObjects := []entities.MainObject{spawnCastle}
	mainObjects = append(mainObjects,
		topology.CreatePlayerOwnedCastles(matchFactions, spawnCastle.Spawn, tuning.PlayerOwnedCastles)...)
	mainObjects = append(mainObjects,
		topology.CreatePlayerUnclaimedCastles(
			matchFactions,
			tuning.ScaleByNeutralGuardStrength(5000),
			configuration.ZoneConfiguration.PlayerZoneCastles)...)
	zone.MainObjects = mainObjects
	rebuildCastleRoads(zone)
}

// rebuildHubZoneCastles rebuilds a hub zone's castles for the current
// hub-castle option, keeping any non-castle main objects.
func rebuildHubZoneCastles(
	zone *entities.Zone,
	castleCount int,
	tuning models.GenerationTuning,
	topology *base.TopologyBase) {
	preserved, isHoldCity := splitOutNonCastles(zone.MainObjects)
	zone.MainObjects = append(
		topology.CreateHubZoneCastles(tuning, castleCount, isHoldCity),
		preserved...)
	rebuildCastleRoads(zone)
}

// splitOutNonCastles returns the zone's non-City main objects and whether any
// of its City castles carries the hold-city win condition, so a rebuild can
// preserve both.
func splitOutNonCastles(mainObjects []entities.MainObject) (preserved []entities.MainObject, isHoldCity bool) {
	for _, mainObject := range mainObjects {
		if mainObject.Type == registry.GetMainObjectTypeValues().City {
			isHoldCity = isHoldCity || mainObject.HoldCityWinCon
			continue
		}
		preserved = append(preserved, mainObject)
	}
	return preserved, isHoldCity
}

// rebuildCastleRoads regenerates the stone castle<->castle roads for the
// zone's current main objects, keeping all other roads.
func rebuildCastleRoads(zone *entities.Zone) {
	kept := make([]entities.Road, 0, len(zone.Roads))
	for _, road := range zone.Roads {
		if isCastleRoad(road) {
			continue
		}
		kept = append(kept, road)
	}
	zone.Roads = append(buildCastleRoads(len(zone.MainObjects)), kept...)
}
