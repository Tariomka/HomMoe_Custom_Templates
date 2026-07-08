package models

import "github.com/Tariomka/hommoe_custom_templates/internal/registry"

type NeutralZoneProfile struct {
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
	PrimaryBuildingsCSid         string
	ExtraBuildingsCSid           string
}

func NewNeutralZoneProfile(quality NeutralZoneQuality) NeutralZoneProfile {
	switch quality {
	case QualityHigh:
		return newNeutralZoneProfileHighQuality()
	case QualityMedium:
		return newNeutralZoneProfileMediumQuality()
	default: // QualityLow
		return newNeutralZoneProfileLowQuality()
	}
}

func newNeutralZoneProfileLowQuality() NeutralZoneProfile {
	return NeutralZoneProfile{
		Layout:                       "zone_layout_sides",
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
		PrimaryBuildingsCSid:         "poor_buildings_construction",
		ExtraBuildingsCSid:           "poor_buildings_construction",
	}
}

func newNeutralZoneProfileMediumQuality() NeutralZoneProfile {
	return NeutralZoneProfile{
		Layout:                       "zone_layout_treasure_zone",
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
		PrimaryBuildingsCSid:         "rich_buildings_construction",
		ExtraBuildingsCSid:           "poor_buildings_construction",
	}
}

func newNeutralZoneProfileHighQuality() NeutralZoneProfile {
	return NeutralZoneProfile{
		Layout:                    "zone_layout_treasure_zone",
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
		PrimaryBuildingsCSid:         "rich_buildings_construction",
		ExtraBuildingsCSid:           "rich_buildings_construction",
	}
}
