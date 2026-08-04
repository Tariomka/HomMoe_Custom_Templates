package zones

import (
	"math"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_zones"
	"github.com/Tariomka/hommoe_custom_templates/internal/common/constants"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
)

type ZoneFactory struct {
	castleFactory *CastleFactory
	roadFactory   *RoadFactory
}

func NewZoneFactory(castleFactory *CastleFactory, roadFactory *RoadFactory) *ZoneFactory {
	return &ZoneFactory{
		castleFactory: castleFactory,
		roadFactory:   roadFactory,
	}
}

func (this *ZoneFactory) CreateSpawnZone(
	label, playerName string,
	connectionNames []string,
	castleCount int,
	matchFactions bool,
	zoneSize float64,
	footholdCount int,
	generateRoads bool,
	tuning models.GenerationTuning,
) entities.Zone {
	mainObjects := []entities.MainObject{
		this.castleFactory.createPlayerSpawnCastle(
			playerName,
			tuning.ScaleByNeutralGuardStrength(5000),
		),
	}
	mainObjects = append(mainObjects,
		this.castleFactory.CreatePlayerOwnedCastles(matchFactions, playerName, tuning.PlayerOwnedCastles)...)
	mainObjects = append(mainObjects,
		this.castleFactory.CreatePlayerUnclaimedCastles(
			matchFactions,
			tuning.ScaleByNeutralGuardStrength(5000),
			castleCount,
		)...)

	roadMainObjectCount := 0
	if castleCount+tuning.PlayerOwnedCastles > 0 {
		roadMainObjectCount = len(mainObjects)
	}

	return variant_content.NewZoneBuilder().
		WithName(constants.PlayerZonePrefix + label).
		WithSize(normalizeZoneSize(zoneSize)).
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
		WithMandatoryContent("mandatory_content_side_" + label).
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
			connectionNames,
			roadMainObjectCount,
			footholdCount,
			generateRoads,
		)).
		Build()
}

func (this *ZoneFactory) CreateNeutralZone(input models.NeutralZoneCreation) entities.Zone {
	if input.HoldCity && input.CastleCount < 1 {
		input.CastleCount = 1
	}
	return this.createNeutralLikeZone(models.NeutralLikeZoneCreation{
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
}

func (this *ZoneFactory) CreateHubZone(input models.HubZoneCreation) entities.Zone {
	if input.HoldCity && input.CastleCount < 1 {
		input.CastleCount = 1
	}
	return this.createNeutralLikeZone(models.NeutralLikeZoneCreation{
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
}

func normalizeZoneSize(zoneSize float64) float64 {
	if math.IsNaN(zoneSize) || math.IsInf(zoneSize, 0) {
		return 1.0
	}
	return helpers.RoundWithPrecision(math.Max(0.1, math.Min(zoneSize, 2.0)), 2)
}
