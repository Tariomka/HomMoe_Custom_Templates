package connection_editor

// Manual-edit reapply support: pushing changed castle-count options into the
// manually edited zone snapshot. The snapshot itself is authoritative on
// regeneration - castle counts are the ONLY generator options that override
// manual edits, and only when they changed since the last generation.

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_zones"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/zone_helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/preview"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/generation_tuning"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones/zone_interfaces"
)

type ManualReapplyService struct {
	zoneEditor    IZoneEditorService
	castleFactory zone_interfaces.ICastleFactory
	tierService   zone_interfaces.IZoneTierService
	tuningFactory generation_tuning.IGenerationTuningFactory
}

func NewManualReapplyService(
	zoneEditor IZoneEditorService,
	castleFactory zone_interfaces.ICastleFactory,
	tierService zone_interfaces.IZoneTierService,
	tuningFactory generation_tuning.IGenerationTuningFactory) IManualReapplyService {
	return &ManualReapplyService{
		zoneEditor:    zoneEditor,
		castleFactory: castleFactory,
		tierService:   tierService,
		tuningFactory: tuningFactory,
	}
}

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
func (this *ManualReapplyService) ApplyCastleSettingChanges(
	zones []template_model.Zone,
	changes editor_state_model.CastleSettingChanges,
	configuration *config.GeneratorConfig) {
	if !changes.Any() {
		return
	}
	tuning := this.tuningFactory.Create(configuration, len(zones))
	for i := range zones {
		zone := &zones[i]
		switch zone_helpers.GetZoneTypeFromName(zone.Name) {
		case preview.ZoneTypePlayer:
			if changes.PlayerCastles {
				this.rebuildSpawnZoneCastles(zone, configuration, tuning)
			}
		case preview.ZoneTypeHub:
			if changes.Hub {
				this.rebuildHubZoneCastles(
					zone,
					configuration.ZoneConfiguration.Advanced.HubZoneCastles,
					tuning,
				)
			}
		case preview.ZoneTypeNeutral:
			if count, update := this.neutralCastleTarget(*zone, changes, configuration); update {
				this.SetNeutralZoneCastleCount(zone, count, tuning)
			}
		case preview.ZoneTypeUnknown:
		}
	}
}

// SetNeutralZoneCastleCount rebuilds only the zone's City castles for the new
// count, keeping the quality profile, guard values, content pools and any
// non-castle main objects (abandoned outposts) untouched - unlike
// ApplyNeutralZoneQuality, which re-profiles the whole zone.
//
// The profile follows the zone's RECORDED tier where it has one. Inferring it
// instead used to hand an unclassifiable zone QualityUnknown, whose profile is
// the Lowest one - a silent down-tier to Plastic stats on a zone the user may
// have set to Gold.
func (this *ManualReapplyService) SetNeutralZoneCastleCount(
	zone *template_model.Zone,
	castleCount int,
	tuning models.GenerationTuning) {
	quality := this.tierService.ResolveQuality(*zone)
	profile := common_zones.GetNeutralZoneProfile(quality)
	preserved, isHoldCity := splitOutNonCastles(zone.MainObjects)
	zone.MainObjects = append(
		this.castleFactory.CreateNeutralZoneCastles(profile, tuning, castleCount, isHoldCity),
		preserved...)
	this.zoneEditor.RebuildCastleRoads(zone)
}

// neutralCastleTarget decides whether a changed option drives the neutral
// zone's castle count, and with which count. Simple mode has a single count
// for every neutral zone; advanced mode has per-tier counts that apply only
// to with-castle zones of the matching quality.
func (this *ManualReapplyService) neutralCastleTarget(
	zone template_model.Zone,
	changes editor_state_model.CastleSettingChanges,
	configuration *config.GeneratorConfig) (int, bool) {
	zoneConfiguration := configuration.ZoneConfiguration
	if changes.NeutralSimple {
		return helpers.Clamp(zoneConfiguration.NeutralZoneCastles, 0, 4), true
	}

	if this.zoneEditor.CountZoneCastles(zone) == 0 {
		return 0, false
	}

	switch this.tierService.ResolveQuality(zone) {
	case neutral_zone.QualityHighest:
		if changes.Hub {
			return helpers.Clamp(zoneConfiguration.Advanced.HubZoneCastles, 0, 4), true
		}
	case neutral_zone.QualityHigh:
		if changes.NeutralHigh {
			return helpers.Clamp(zoneConfiguration.Advanced.NeutralHighCastlesPerZone, 0, 4), true
		}
	case neutral_zone.QualityMedium:
		if changes.NeutralMedium {
			return helpers.Clamp(zoneConfiguration.Advanced.NeutralMediumCastlesPerZone, 0, 4), true
		}
	case neutral_zone.QualityLow:
		if changes.NeutralLow {
			return helpers.Clamp(zoneConfiguration.Advanced.NeutralLowCastlesPerZone, 0, 4), true
		}
	case neutral_zone.QualityLowest:
		if changes.NeutralLowest {
			return helpers.Clamp(zoneConfiguration.Advanced.NeutralLowestCastlesPerZone, 0, 4), true
		}
	case neutral_zone.QualityUnknown:
	}
	return 0, false
}

// rebuildSpawnZoneCastles rebuilds a player spawn zone's extra castles for the
// current player-castle options. The spawn castle (main object 0) is kept
// verbatim so the player assignment and faction survive.
func (this *ManualReapplyService) rebuildSpawnZoneCastles(
	zone *template_model.Zone,
	configuration *config.GeneratorConfig,
	tuning models.GenerationTuning) {
	if len(zone.MainObjects) == 0 || zone.MainObjects[0].Type != registry.GetMainObjectTypeValues().Spawn {
		return
	}

	spawnCastle := zone.MainObjects[0]
	matchFactions := configuration.MatchPlayerCastleFactions
	mainObjects := []template_model.MainObject{spawnCastle}
	mainObjects = append(mainObjects,
		this.castleFactory.CreatePlayerOwnedCastles(
			matchFactions, spawnCastle.Spawn, tuning.PlayerOwnedCastles)...)
	mainObjects = append(mainObjects,
		this.castleFactory.CreatePlayerUnclaimedCastles(
			matchFactions,
			tuning.ScaleByNeutralGuardStrength(5000),
			configuration.ZoneConfiguration.PlayerZoneCastles)...)
	zone.MainObjects = mainObjects
	this.zoneEditor.RebuildCastleRoads(zone)
}

// rebuildHubZoneCastles rebuilds a hub zone's castles for the current
// hub-castle option, keeping any non-castle main objects.
func (this *ManualReapplyService) rebuildHubZoneCastles(
	zone *template_model.Zone,
	castleCount int,
	tuning models.GenerationTuning) {
	preserved, isHoldCity := splitOutNonCastles(zone.MainObjects)
	zone.MainObjects = append(
		this.castleFactory.CreateHubZoneCastles(tuning, castleCount, isHoldCity),
		preserved...)
	this.zoneEditor.RebuildCastleRoads(zone)
}

// splitOutNonCastles returns the zone's non-City main objects and whether any
// of its City castles carries the hold-city win condition, so a rebuild can
// preserve both.
func splitOutNonCastles(
	mainObjects []template_model.MainObject) (preserved []template_model.MainObject, isHoldCity bool) {
	mainObjectType := registry.GetMainObjectTypeValues().City
	for _, mainObject := range mainObjects {
		if mainObject.Type == mainObjectType {
			isHoldCity = isHoldCity || mainObject.HoldCityWinCon
			continue
		}
		preserved = append(preserved, mainObject)
	}
	return preserved, isHoldCity
}
