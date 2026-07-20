package neutral_zone

import "github.com/Tariomka/hommoe_custom_templates/internal/registry"

type Profile struct {
	Layout                       string
	GuardReactionDistribution    []int
	GuardMultiplier              float64
	GuardedContentPool           []string
	UnguardedContentPool         []string
	ResourcesContentPool         []string
	GuardedContentValue          int
	GuardedContentValuePerArea   int
	UnguardedContentValue        int
	UnguardedContentValuePerArea int
	ResourcesValue               int
	ResourcesValuePerArea        int
	PrimaryCityGuardValue        int
	ExtraCityGuardValue          int
	PrimaryBuildingsSid          string
	ExtraBuildingsSid            string
}

func NewNeutralZoneProfile(quality Quality) Profile {
	switch quality {
	case QualityHighest:
		return newNeutralZoneProfileHighestQuality()
	case QualityHigh:
		return newNeutralZoneProfileHighQuality()
	case QualityMedium:
		return newNeutralZoneProfileMediumQuality()
	case QualityLow:
		return newNeutralZoneProfileLowQuality()
	case QualityLowest, QualityUnknown:
		fallthrough
	default:
		return newNeutralZoneProfileLowestQuality()
	}
}

func newNeutralZoneProfileLowestQuality() Profile {
	layoutValues := registry.GetLayoutValues()
	constructionValues := registry.GetBuildingsConstructionSidValues()
	return Profile{
		Layout:                       layoutValues.Sides,
		GuardReactionDistribution:    []int{0, 10, 10, 10, 10, 0},
		GuardMultiplier:              0.9,
		GuardedContentPool:           registry.GetGuardedContentPoolT1List(),
		UnguardedContentPool:         registry.GetUnguardedContentPoolT1List(),
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

func newNeutralZoneProfileLowQuality() Profile {
	layoutValues := registry.GetLayoutValues()
	constructionValues := registry.GetBuildingsConstructionSidValues()
	return Profile{
		Layout:                       layoutValues.Sides,
		GuardReactionDistribution:    []int{0, 10, 10, 10, 10, 0},
		GuardMultiplier:              1.1,
		GuardedContentPool:           registry.GetGuardedContentPoolT2List(),
		UnguardedContentPool:         registry.GetUnguardedContentPoolT2List(),
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

func newNeutralZoneProfileMediumQuality() Profile {
	layoutValues := registry.GetLayoutValues()
	constructionValues := registry.GetBuildingsConstructionSidValues()
	return Profile{
		Layout:                       layoutValues.TreasureZone,
		GuardReactionDistribution:    []int{0, 10, 10, 10, 10, 0},
		GuardMultiplier:              1.4,
		GuardedContentPool:           registry.GetGuardedContentPoolT3List(),
		UnguardedContentPool:         registry.GetUnguardedContentPoolT3List(),
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

func newNeutralZoneProfileHighQuality() Profile {
	layoutValues := registry.GetLayoutValues()
	constructionValues := registry.GetBuildingsConstructionSidValues()
	return Profile{
		Layout:                    layoutValues.TreasureZone,
		GuardReactionDistribution: []int{0, 10, 10, 20, 10, 0},
		GuardMultiplier:           1.8,
		GuardedContentPool: append(
			registry.GetGuardedContentPoolT4List(),
			registry.GetGuardedContentPoolT5List()...),
		UnguardedContentPool: append(
			registry.GetUnguardedContentPoolT4List(),
			registry.GetUnguardedContentPoolT5List()...),
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

func newNeutralZoneProfileHighestQuality() Profile {
	layoutValues := registry.GetLayoutValues()
	constructionValues := registry.GetBuildingsConstructionSidValues()
	return Profile{
		Layout:                    layoutValues.Center,
		GuardReactionDistribution: []int{0, 10, 10, 20, 10, 0},
		GuardMultiplier:           2.3,
		GuardedContentPool: append(
			registry.GetGuardedContentPoolT5List(),
			registry.GetGuardedContentPoolT5List()...),
		UnguardedContentPool: append(
			registry.GetUnguardedContentPoolT5List(),
			registry.GetUnguardedContentPoolT5List()...),
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
