package zones

import (
	"fmt"
	"math"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
)

func (this *ZoneFactory) createNeutralLikeZone(input models.NeutralLikeZoneCreationRequest) entities.Zone {
	mainObjects := this.createNeutralLikeMainObjects(input)
	roadMainObjectCount := len(mainObjects)
	if input.CastleStrategy == models.ZoneCastleStrategyHub {
		roadMainObjectCount = input.CastleCount
	}

	zoneBuilder := this.createNeutralLikeZoneBuilder(input).
		WithMainObjects(mainObjects).
		WithRoads(this.roadFactory.CreateOuterZoneRoads(
			input.ConnectionNames,
			roadMainObjectCount,
			input.FootholdCount,
			input.GenerateRoads,
		))

	if input.MandatoryContentName != "" {
		zoneBuilder = zoneBuilder.WithMandatoryContent(input.MandatoryContentName)
	}
	if input.BiomeMatchPolicy == models.ZoneBiomeMatchPrimaryMainObjectWhenPresent && len(mainObjects) > 0 {
		zoneBuilder = zoneBuilder.WithBiomeMatchMainObject("0")
	} else {
		zoneBuilder = zoneBuilder.WithBiomeMatchZone()
	}

	return zoneBuilder.Build()
}

func (this *ZoneFactory) createNeutralLikeMainObjects(
	input models.NeutralLikeZoneCreationRequest,
) []entities.MainObject {
	if input.CastleStrategy == models.ZoneCastleStrategyHub {
		return this.castleFactory.CreateHubZoneCastles(input.Tuning, input.CastleCount, input.HoldCity)
	}

	mainObjects := this.castleFactory.CreateNeutralZoneCastles(
		input.Profile,
		input.Tuning,
		input.CastleCount,
		input.HoldCity,
	)
	return append(mainObjects, this.castleFactory.createAbandonedOutposts(
		input.Profile,
		input.Tuning,
		input.OutpostCount,
	)...)
}

func (this *ZoneFactory) createNeutralLikeZoneBuilder(
	input models.NeutralLikeZoneCreationRequest,
) *variant_content.ZoneBuilder {
	profile := input.Profile
	tuning := input.Tuning
	return variant_content.NewZoneBuilder().
		WithName(input.Name).
		WithSize(normalizeZoneSize(input.Size)).
		WithLayout(profile.Layout).
		WithGuardCutoffValue(2000).
		WithGuardRandomization(input.GuardRandomization).
		WithGuardMultiplier(tuning.ScaleByNeutralGuardStrengthPrecise(profile.GuardMultiplier)).
		WithGuardWeeklyIncrement(0.20).
		WithGuardReactionDistribution(profile.GuardReactionDistribution).
		WithDiplomacyModifier(-0.5).
		WithGuardedContentPool(profile.GuardedContentPool).
		WithUnguardedContentPool(profile.UnguardedContentPool).
		WithResourcesContentPool(profile.ResourcesContentPool).
		WithContentCountLimits(buildSideContentLimits()).
		WithGuardedContentValue(tuning.ScaleByStructureDensity(
			float64(profile.GuardedContentValue) * tuning.ContentScale)).
		WithGuardedContentValuePerArea(tuning.ScaleByStructureDensity(
			float64(profile.GuardedContentValuePerArea) * math.Sqrt(tuning.ContentScale))).
		WithUnguardedContentValue(tuning.ScaleByStructureDensity(
			float64(profile.UnguardedContentValue) * tuning.ContentScale)).
		WithUnguardedContentValuePerArea(tuning.ScaleByStructureDensity(
			float64(profile.UnguardedContentValuePerArea) * math.Sqrt(tuning.ContentScale))).
		WithResourcesValue(tuning.ScaleByResourceDensity(
			float64(profile.ResourcesValue) * tuning.ContentScale)).
		WithResourcesValuePerArea(tuning.ScaleByResourceDensity(
			float64(profile.ResourcesValuePerArea) * math.Sqrt(tuning.ContentScale))).
		WithCrossroadsPosition(0)
}

func buildSideContentLimits() entities.StringList {
	var limits []string
	for firstIndex := 1; firstIndex <= 5; firstIndex++ {
		for secondIndex := firstIndex + 1; secondIndex <= 6; secondIndex++ {
			limits = append(limits, fmt.Sprintf("content_limits_side_%d_%d", firstIndex, secondIndex))
		}
	}
	return limits
}
