package providers

import (
	"fmt"
	"strings"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/provider_interfaces"
)

type ContentLimitProvider struct{}

func NewContentLimitProvider() provider_interfaces.IContentLimitProvider {
	return &ContentLimitProvider{}
}

func (this *ContentLimitProvider) CreateContentCountLimits(
	settings config.GeneratorConfig) []entities.ContentCountLimit {
	sidLimits := this.createDefaultContentLimits()

	// Lift limits when any mandatory-content list (player or neutral or
	// hub) requests more of a given SID than the default cap.
	sidCounts := map[string]int{}
	tally := func(items []entities.MandatoryContentItem) {
		for _, item := range items {
			if item.SID != "" {
				sidCounts[strings.ToLower(item.SID)]++
			}
		}
	}
	tally(settings.PlayerZoneMandatoryContent)
	tally(settings.LowestNeutralMandatoryContent)
	tally(settings.LowNeutralMandatoryContent)
	tally(settings.MediumNeutralMandatoryContent)
	tally(settings.HighNeutralMandatoryContent)
	tally(settings.HubZoneMandatoryContent)
	for i := range sidLimits {
		if count, ok := sidCounts[strings.ToLower(sidLimits[i].SID)]; ok {
			if count > sidLimits[i].MaxCount {
				sidLimits[i].MaxCount = count
			}
		}
	}

	var limits []entities.ContentCountLimit
	limits = append(limits, entities.ContentCountLimit{Name: "content_limits_side", Limits: sidLimits})
	limits = append(limits, entities.ContentCountLimit{Name: "content_limits_side_0_0", Limits: sidLimits})
	for a := 1; a <= 5; a++ {
		for b := a + 1; b <= 6; b++ {
			limits = append(limits, entities.ContentCountLimit{
				Name:   fmt.Sprintf("content_limits_side_%d_%d", a, b),
				Limits: sidLimits,
			})
		}
	}
	return limits
}

func (this *ContentLimitProvider) createDefaultContentLimits() []entities.ContentLimit {
	buildingObjects := registry.GetMapObjectBuildingValues()
	heroBuffBuildings := registry.GetMapObjectHeroBuffBuildingValues()
	magicBuildings := registry.GetMapObjectMagicBuildingValues()
	nonContentObjects := registry.GetMapObjectNonContentValues()
	randomUnitBanks := registry.GetMapObjectRandomUnitBankValues()
	resourceObjects := registry.GetMapObjectResourceValues()
	t1GuardedResourceBanks := registry.GetMapObjectT1GuardedResourceBankValues()
	t1StatsAndSkillsBuildings := registry.GetMapObjectT1StatsAndSkillsBuildingValues()
	t2StatsAndSkillsBuildings := registry.GetMapObjectT2StatsAndSkillsBuildingValues()
	unitBanks := registry.GetMapObjectNamedUnitBankValues()
	visionBuildings := registry.GetMapObjectVisionBuildingValues()
	return []entities.ContentLimit{
		{SID: t1GuardedResourceBanks.BlackTower, MaxCount: 0},
		{SID: heroBuffBuildings.Fountain, MaxCount: 2},
		{SID: heroBuffBuildings.Fountain2, MaxCount: 2},
		{SID: heroBuffBuildings.ManaWell, MaxCount: 2},
		{SID: heroBuffBuildings.BeerFountain, MaxCount: 2},
		{SID: nonContentObjects.Market, MaxCount: 1},
		{SID: nonContentObjects.Forge, MaxCount: 2},
		{SID: heroBuffBuildings.Stables, MaxCount: 1},
		{SID: visionBuildings.Watchtower, MaxCount: 2},
		{SID: visionBuildings.WindRose, MaxCount: 1},
		{SID: heroBuffBuildings.QuixsPath, MaxCount: 2},
		{SID: heroBuffBuildings.CrystalTrail, MaxCount: 3},
		{SID: heroBuffBuildings.MysteriousStone, MaxCount: 2},
		{SID: t2StatsAndSkillsBuildings.University, MaxCount: 2},
		{SID: t1StatsAndSkillsBuildings.WiseOwl, MaxCount: 4},
		{SID: magicBuildings.CelestialSphere, MaxCount: 2},
		{SID: heroBuffBuildings.PileOfBooks, MaxCount: 2},
		{SID: nonContentObjects.InsarasEye, MaxCount: 2},
		{SID: heroBuffBuildings.TearOfTruth, MaxCount: 3},
		{SID: buildingObjects.TreeOfAbundance, MaxCount: 2},
		{SID: buildingObjects.HuntsmansCamp, MaxCount: 2},
		{SID: t1GuardedResourceBanks.ShadyDen, MaxCount: 2},
		{SID: randomUnitBanks.RandomHireTier1, MaxCount: 6},
		{SID: randomUnitBanks.RandomHireTier2, MaxCount: 6},
		{SID: randomUnitBanks.RandomHireTier3, MaxCount: 6},
		{SID: randomUnitBanks.RandomHireTier4, MaxCount: 6},
		{SID: randomUnitBanks.RandomHireTier5, MaxCount: 6},
		{SID: randomUnitBanks.RandomHireTier6, MaxCount: 6},
		{SID: randomUnitBanks.RandomHireTier7, MaxCount: 6},
		{SID: buildingObjects.Arena, MaxCount: 2},
		{SID: buildingObjects.SacrificialShrine, MaxCount: 2},
		{SID: buildingObjects.Chimerologist, MaxCount: 2},
		{SID: t2StatsAndSkillsBuildings.Circus, MaxCount: 2},
		{SID: t2StatsAndSkillsBuildings.InfernalCirque, MaxCount: 2},
		{SID: visionBuildings.FlatteringMirror, MaxCount: 2},
		{SID: buildingObjects.FickleShrine, MaxCount: 1},
		{SID: unitBanks.PointOfBalance, MaxCount: 3},
		{SID: resourceObjects.PandoraBox, MaxCount: 4},
		{SID: unitBanks.RitualPyre, MaxCount: 3},
		{SID: unitBanks.BorealCall, MaxCount: 3},
		{SID: unitBanks.JoustingRange, MaxCount: 1},
		{SID: unitBanks.UnforgottenGrave, MaxCount: 1},
		{SID: unitBanks.PetrifiedMemorial, MaxCount: 1},
		{SID: unitBanks.Gorge, MaxCount: 1},
	}
}
