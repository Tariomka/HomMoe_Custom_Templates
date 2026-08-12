package common_zones

import (
	"slices"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
)

func GetNeutralZoneProfile(quality neutral_zone.Quality) neutral_zone.Profile {
	switch quality {
	case neutral_zone.QualityHighest:
		return getHighestQualityProfile()
	case neutral_zone.QualityHigh:
		return getHighQualityProfile()
	case neutral_zone.QualityMedium:
		return getMediumQualityProfile()
	case neutral_zone.QualityLow:
		return getLowQualityProfile()
	case neutral_zone.QualityLowest, neutral_zone.QualityUnknown:
		fallthrough
	default:
		return getLowestQualityProfile()
	}
}

func getLowestQualityProfile() neutral_zone.Profile {
	layoutValues := registry.GetLayoutValues()
	constructionValues := registry.GetBuildingsConstructionSidValues()
	return neutral_zone.Profile{
		Layout:                       layoutValues.Sides,
		GuardReactionDistribution:    []int{0, 10, 10, 10, 10, 0},
		GuardMultiplier:              0.9,
		GuardedContentPool:           slices.Clone(registry.GetGuardedContentPoolT1List()),
		UnguardedContentPool:         slices.Clone(registry.GetUnguardedContentPoolT1List()),
		ResourcesContentPool:         []string{registry.GetResourcesContentPoolValues().StartZoneVeryPoor},
		GuardedContentValue:          60000,
		GuardedContentValuePerArea:   500,
		UnguardedContentValue:        16000,
		UnguardedContentValuePerArea: 130,
		ResourcesValue:               16000,
		ResourcesValuePerArea:        140,
		PrimaryCityGuardValue:        2000,
		ExtraCityGuardValue:          1000,
		PrimaryBuildingsSid:          constructionValues.ExtraPoor,
		ExtraBuildingsSid:            constructionValues.ExtraPoor,
	}
}

func getLowQualityProfile() neutral_zone.Profile {
	layoutValues := registry.GetLayoutValues()
	constructionValues := registry.GetBuildingsConstructionSidValues()
	return neutral_zone.Profile{
		Layout:                       layoutValues.Sides,
		GuardReactionDistribution:    []int{0, 10, 10, 10, 10, 0},
		GuardMultiplier:              1.1,
		GuardedContentPool:           slices.Clone(registry.GetGuardedContentPoolT2List()),
		UnguardedContentPool:         slices.Clone(registry.GetUnguardedContentPoolT2List()),
		ResourcesContentPool:         []string{registry.GetResourcesContentPoolValues().StartZonePoor},
		GuardedContentValue:          120000,
		GuardedContentValuePerArea:   1000,
		UnguardedContentValue:        25000,
		UnguardedContentValuePerArea: 200,
		ResourcesValue:               30000,
		ResourcesValuePerArea:        240,
		PrimaryCityGuardValue:        4000,
		ExtraCityGuardValue:          2000,
		PrimaryBuildingsSid:          constructionValues.Poor,
		ExtraBuildingsSid:            constructionValues.Poor,
	}
}

func getMediumQualityProfile() neutral_zone.Profile {
	layoutValues := registry.GetLayoutValues()
	constructionValues := registry.GetBuildingsConstructionSidValues()
	return neutral_zone.Profile{
		Layout:                       layoutValues.TreasureZone,
		GuardReactionDistribution:    []int{0, 10, 10, 10, 10, 0},
		GuardMultiplier:              1.4,
		GuardedContentPool:           slices.Clone(registry.GetGuardedContentPoolT3List()),
		UnguardedContentPool:         slices.Clone(registry.GetUnguardedContentPoolT3List()),
		ResourcesContentPool:         []string{registry.GetResourcesContentPoolValues().StartZoneMedium},
		GuardedContentValue:          240000,
		GuardedContentValuePerArea:   2000,
		UnguardedContentValue:        38000,
		UnguardedContentValuePerArea: 300,
		ResourcesValue:               55000,
		ResourcesValuePerArea:        420,
		PrimaryCityGuardValue:        8000,
		ExtraCityGuardValue:          4000,
		PrimaryBuildingsSid:          constructionValues.Rich,
		ExtraBuildingsSid:            constructionValues.Poor,
	}
}

func getHighQualityProfile() neutral_zone.Profile {
	layoutValues := registry.GetLayoutValues()
	constructionValues := registry.GetBuildingsConstructionSidValues()
	return neutral_zone.Profile{
		Layout:                    layoutValues.TreasureZone,
		GuardReactionDistribution: []int{0, 10, 10, 20, 10, 0},
		GuardMultiplier:           1.8,
		GuardedContentPool: slices.Concat(
			registry.GetGuardedContentPoolT4List(),
			registry.GetGuardedContentPoolT5List(),
		),
		UnguardedContentPool: slices.Concat(
			registry.GetUnguardedContentPoolT4List(),
			registry.GetUnguardedContentPoolT5List(),
		),
		ResourcesContentPool:         []string{registry.GetResourcesContentPoolValues().StartZoneRich},
		GuardedContentValue:          480000,
		GuardedContentValuePerArea:   3000,
		UnguardedContentValue:        80000,
		UnguardedContentValuePerArea: 620,
		ResourcesValue:               90000,
		ResourcesValuePerArea:        580,
		PrimaryCityGuardValue:        16000,
		ExtraCityGuardValue:          8000,
		PrimaryBuildingsSid:          constructionValues.Rich,
		ExtraBuildingsSid:            constructionValues.Rich,
	}
}

func getHighestQualityProfile() neutral_zone.Profile {
	layoutValues := registry.GetLayoutValues()
	constructionValues := registry.GetBuildingsConstructionSidValues()
	return neutral_zone.Profile{
		Layout:                    layoutValues.Center,
		GuardReactionDistribution: []int{0, 10, 10, 20, 10, 0},
		GuardMultiplier:           2.3,
		GuardedContentPool: slices.Concat(
			registry.GetGuardedContentPoolT5List(),
			registry.GetGuardedContentPoolT5List(),
		),
		UnguardedContentPool: slices.Concat(
			registry.GetUnguardedContentPoolT5List(),
			registry.GetUnguardedContentPoolT5List(),
		),
		ResourcesContentPool:         []string{registry.GetResourcesContentPoolValues().TreasureZoneRich},
		GuardedContentValue:          960000,
		GuardedContentValuePerArea:   4500,
		UnguardedContentValue:        168000,
		UnguardedContentValuePerArea: 1280,
		ResourcesValue:               147000,
		ResourcesValuePerArea:        800,
		PrimaryCityGuardValue:        32000,
		ExtraCityGuardValue:          16000,
		PrimaryBuildingsSid:          constructionValues.UltraRich,
		ExtraBuildingsSid:            constructionValues.UltraRich,
	}
}
