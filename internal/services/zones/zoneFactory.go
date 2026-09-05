package zones

import (
	"math"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_zones"
	"github.com/Tariomka/hommoe_custom_templates/internal/common/constants"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones/zone_interfaces"
)

type ZoneFactory struct {
	castleFactory zone_interfaces.ICastleFactory
	roadFactory   zone_interfaces.IRoadFactory
}

func NewZoneFactory(
	castleFactory zone_interfaces.ICastleFactory,
	roadFactory zone_interfaces.IRoadFactory) zone_interfaces.IZoneFactory {
	return &ZoneFactory{
		castleFactory: castleFactory,
		roadFactory:   roadFactory,
	}
}

func (this *ZoneFactory) CreateSpawnZone(input models.SpawnZoneCreationRequest) template_model.Zone {
	mainObjects := []template_model.MainObject{
		this.castleFactory.CreatePlayerSpawnCastle(
			input.PlayerName,
			input.Tuning.ScaleByNeutralGuardStrength(5000),
		),
	}
	mainObjects = append(mainObjects,
		this.castleFactory.CreatePlayerOwnedCastles(
			input.MatchFactions, input.PlayerName, input.Tuning.PlayerOwnedCastles)...)
	mainObjects = append(mainObjects,
		this.castleFactory.CreatePlayerUnclaimedCastles(
			input.MatchFactions,
			input.Tuning.ScaleByNeutralGuardStrength(5000),
			input.CastleCount,
		)...)

	roadMainObjectCount := 0
	if input.CastleCount+input.Tuning.PlayerOwnedCastles > 0 {
		roadMainObjectCount = len(mainObjects)
	}

	tuning := input.Tuning

	return variant_content.NewZoneBuilder().
		WithName(constants.GetPlayerZoneNameFor(input.Label)).
		WithSize(normalizeZoneSize(input.Size)).
		WithLayoutSpawns().
		WithGuardCutoffValue(2000).
		WithGuardRandomization(tuning.GuardRandomization).
		WithGuardMultiplier(tuning.ScaleByNeutralGuardStrengthPrecise(1.0)).
		WithGuardWeeklyIncrement(0.20).
		WithGuardReactionDistribution([]int{60, 20, 10, 10, 2, 0}).
		WithDiplomacyModifier(-0.5).
		WithGuardedContentPool(registry.GetGuardedContentPoolT2List()).
		WithUnguardedContentPool(registry.GetUnguardedContentPoolT2List()).
		WithResourcesContentPool([]string{registry.GetResourcesContentPoolValues().StartZonePoor}).
		WithMandatoryContent(constants.GetSideContentNameFor(input.Label)).
		WithContentCountLimits(buildSideContentLimits()).
		WithGuardedContentValue(tuning.ScaleByStructureDensity(200000 * tuning.ContentScale)).
		WithGuardedContentValuePerArea(tuning.ScaleByStructureDensity(2000 * math.Sqrt(tuning.ContentScale))).
		WithUnguardedContentValue(tuning.ScaleByStructureDensity(50000 * tuning.ContentScale)).
		WithUnguardedContentValuePerArea(tuning.ScaleByStructureDensity(400 * math.Sqrt(tuning.ContentScale))).
		WithResourcesValue(tuning.ScaleByResourceDensity(80000 * tuning.ContentScale)).
		WithResourcesValuePerArea(tuning.ScaleByResourceDensity(600 * math.Sqrt(tuning.ContentScale))).
		WithMainObjects(mainObjects).
		WithBiomeMatchMainObject("0").
		WithCrossroadsPosition(0).
		WithRoads(this.roadFactory.CreateOuterZoneRoads(
			input.ConnectionNames,
			roadMainObjectCount,
			input.FootholdCount,
			input.GenerateRoads,
		)).
		Build()
}

func (this *ZoneFactory) CreateNeutralZone(input models.NeutralZoneCreationRequest) template_model.Zone {
	if input.HoldCity && input.CastleCount < 1 {
		input.CastleCount = 1
	}
	zone := this.createNeutralLikeZone(models.NeutralLikeZoneCreationRequest{
		Name:                 input.Name,
		Profile:              common_zones.GetNeutralZoneProfile(input.Quality),
		Size:                 input.Size,
		ConnectionNames:      input.ConnectionNames,
		MandatoryContentName: input.MandatoryContentName,
		CastleStrategy:       models.ZoneCastleStrategyNeutral,
		CastleCount:          input.CastleCount,
		HoldCity:             input.HoldCity,
		OutpostCount:         input.OutpostCount,
		FootholdCount:        input.FootholdCount,
		GuardRandomization:   input.GuardRandomization,
		GenerateRoads:        input.GenerateRoads,
		BiomeMatchPolicy:     models.ZoneBiomeMatchPrimaryMainObjectWhenPresent,
		Tuning:               input.Tuning,
	})
	zone.Quality = &input.Quality
	return zone
}

func (this *ZoneFactory) CreateHubZone(input models.HubZoneCreationRequest) template_model.Zone {
	if input.HoldCity && input.CastleCount < 1 {
		input.CastleCount = 1
	}
	zone := this.createNeutralLikeZone(models.NeutralLikeZoneCreationRequest{
		Name:                 input.Name,
		Profile:              common_zones.GetNeutralZoneProfile(neutral_zone.QualityHighest),
		Size:                 input.Size,
		ConnectionNames:      input.ConnectionNames,
		MandatoryContentName: input.MandatoryContentName,
		CastleStrategy:       models.ZoneCastleStrategyHub,
		CastleCount:          input.CastleCount,
		HoldCity:             input.HoldCity,
		GuardRandomization:   input.GuardRandomization,
		GenerateRoads:        input.GenerateRoads,
		BiomeMatchPolicy:     models.ZoneBiomeMatchPrimaryMainObjectWhenPresent,
		Tuning:               input.Tuning,
	})
	zone.Quality = new(neutral_zone.QualityHighest)
	return zone
}

func normalizeZoneSize(zoneSize float64) float64 {
	if math.IsNaN(zoneSize) || math.IsInf(zoneSize, 0) {
		return 1.0
	}

	return helpers.RoundWithPrecision(math.Max(0.1, math.Min(zoneSize, 2.0)), 2)
}
