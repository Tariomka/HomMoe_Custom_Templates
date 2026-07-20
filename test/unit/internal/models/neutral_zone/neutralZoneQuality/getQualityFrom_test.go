package neutralZoneQuality_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
	"github.com/stretchr/testify/assert"
)

func TestWhenZoneCharacteristicsVary_DetectsQualityAccordingly(t *testing.T) {
	t.Parallel()
	layoutValues := registry.GetLayoutValues()
	resourcePools := registry.GetResourcesContentPoolValues()
	testCases := []struct {
		subtestName string
		zone        entities.Zone
		expected    neutral_zone.Quality
	}{
		{
			"WhenZoneIsPlayerSpawn_ReturnsUnknown",
			entities.Zone{Name: "Spawn-A", Layout: layoutValues.Center},
			neutral_zone.QualityUnknown,
		},
		{
			"WhenLayoutIsUnrecognized_ReturnsUnknown",
			entities.Zone{Name: "Neutral-B", Layout: "zone_layout_back"},
			neutral_zone.QualityUnknown,
		},
		{
			"WhenCenterLayoutHasRichTreasureResources_ReturnsHighest",
			entities.Zone{
				Name:                 "Hub",
				Layout:               layoutValues.Center,
				GuardedContentPool:   []string{"pool_t5_treasure"},
				ResourcesContentPool: []string{resourcePools.TreasureZoneRich},
			},
			neutral_zone.QualityHighest,
		},
		{
			"WhenCenterLayoutHasOnlyTier5GuardedPool_ReturnsHighest",
			entities.Zone{
				Name:               "Hub",
				Layout:             layoutValues.Center,
				GuardedContentPool: []string{"pool_t5_treasure", "pool_t5_item"},
			},
			neutral_zone.QualityHighest,
		},
		{
			"WhenCenterLayoutHasNoContentPools_ReturnsUnknown",
			entities.Zone{Name: "Hub", Layout: layoutValues.Center},
			neutral_zone.QualityUnknown,
		},
		{
			"WhenCenterLayoutPoolsCarryNoTier5Marker_ReturnsUnknown",
			entities.Zone{
				Name:               "Hub",
				Layout:             layoutValues.Center,
				GuardedContentPool: []string{"pool_t3_stuff"},
			},
			neutral_zone.QualityUnknown,
		},
		{
			"WhenTreasureLayoutHasRichStartResources_ReturnsHigh",
			entities.Zone{
				Name:                 "Neutral-B",
				Layout:               layoutValues.TreasureZone,
				GuardedContentPool:   []string{"pool_without_tier_marker"},
				ResourcesContentPool: []string{resourcePools.StartZoneRich},
			},
			neutral_zone.QualityHigh,
		},
		{
			"WhenTreasureLayoutHasMixedTier4AndTier5Pools_ReturnsHigh",
			entities.Zone{
				Name:               "Neutral-B",
				Layout:             layoutValues.TreasureZone,
				GuardedContentPool: []string{"pool_t4_stuff", "pool_t5_stuff"},
			},
			neutral_zone.QualityHigh,
		},
		{
			"WhenTreasureLayoutHasOnlyTier3Pools_ReturnsMedium",
			entities.Zone{
				Name:                 "Neutral-B",
				Layout:               layoutValues.TreasureZone,
				GuardedContentPool:   []string{"pool_t3_stuff"},
				UnguardedContentPool: []string{"pool_t3_other"},
			},
			neutral_zone.QualityMedium,
		},
		{
			"WhenTreasureLayoutHasMediumStartResources_ReturnsMedium",
			entities.Zone{
				Name:                 "Neutral-B",
				Layout:               layoutValues.TreasureZone,
				GuardedContentPool:   []string{"pool_without_tier_marker"},
				ResourcesContentPool: []string{resourcePools.StartZoneMedium},
			},
			neutral_zone.QualityMedium,
		},
		{
			"WhenTreasureLayoutHasNoContentPools_ReturnsUnknown",
			entities.Zone{Name: "Neutral-B", Layout: layoutValues.TreasureZone},
			neutral_zone.QualityUnknown,
		},
		{
			"WhenTreasureLayoutPoolsCarryNoKnownMarker_ReturnsUnknown",
			entities.Zone{
				Name:               "Neutral-B",
				Layout:             layoutValues.TreasureZone,
				GuardedContentPool: []string{"pool_t1_stuff"},
			},
			neutral_zone.QualityUnknown,
		},
		{
			"WhenSidesLayoutHasOnlyTier2Pools_ReturnsLow",
			entities.Zone{
				Name:               "Neutral-B",
				Layout:             layoutValues.Sides,
				GuardedContentPool: []string{"pool_t2_stuff"},
			},
			neutral_zone.QualityLow,
		},
		{
			"WhenSidesLayoutHasPoorStartResources_ReturnsLow",
			entities.Zone{
				Name:                 "Neutral-B",
				Layout:               layoutValues.Sides,
				GuardedContentPool:   []string{"pool_without_tier_marker"},
				ResourcesContentPool: []string{resourcePools.StartZonePoor},
			},
			neutral_zone.QualityLow,
		},
		{
			"WhenSidesLayoutHasOnlyTier1Pools_ReturnsLowest",
			entities.Zone{
				Name:               "Neutral-B",
				Layout:             layoutValues.Sides,
				GuardedContentPool: []string{"pool_t1_stuff"},
			},
			neutral_zone.QualityLowest,
		},
		{
			"WhenSidesLayoutHasVeryPoorStartResources_ReturnsLowest",
			entities.Zone{
				Name:                 "Neutral-B",
				Layout:               layoutValues.Sides,
				GuardedContentPool:   []string{"pool_without_tier_marker"},
				ResourcesContentPool: []string{resourcePools.StartZoneVeryPoor},
			},
			neutral_zone.QualityLowest,
		},
		{
			"WhenSidesLayoutHasNoContentPools_ReturnsUnknown",
			entities.Zone{Name: "Neutral-B", Layout: layoutValues.Sides},
			neutral_zone.QualityUnknown,
		},
		{
			"WhenSidesLayoutPoolsCarryNoKnownMarker_ReturnsUnknown",
			entities.Zone{
				Name:               "Neutral-B",
				Layout:             layoutValues.Sides,
				GuardedContentPool: []string{"pool_t3_stuff"},
			},
			neutral_zone.QualityUnknown,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.subtestName, func(t *testing.T) {
			t.Parallel()
			// Arrange

			// Act
			quality := neutral_zone.GetQualityFrom(testCase.zone)

			// Assert
			assert.Equal(t, testCase.expected, quality)
		})
	}
}

func TestWhenGeneratedProfileRoundTrips_EveryQualityIsDetectedBack(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		subtestName string
		quality     neutral_zone.Quality
	}{
		{"WhenProfileQualityIsLowest_DetectsSameQuality", neutral_zone.QualityLowest},
		{"WhenProfileQualityIsLow_DetectsSameQuality", neutral_zone.QualityLow},
		{"WhenProfileQualityIsMedium_DetectsSameQuality", neutral_zone.QualityMedium},
		{"WhenProfileQualityIsHigh_DetectsSameQuality", neutral_zone.QualityHigh},
		{"WhenProfileQualityIsHighest_DetectsSameQuality", neutral_zone.QualityHighest},
	}
	for _, testCase := range testCases {
		t.Run(testCase.subtestName, func(t *testing.T) {
			t.Parallel()
			// Arrange
			profile := neutral_zone.NewNeutralZoneProfile(testCase.quality)
			zone := entities.Zone{
				Name:                 "Neutral-Z",
				Layout:               profile.Layout,
				GuardedContentPool:   profile.GuardedContentPool,
				UnguardedContentPool: profile.UnguardedContentPool,
				ResourcesContentPool: profile.ResourcesContentPool,
			}

			// Act
			quality := neutral_zone.GetQualityFrom(zone)

			// Assert
			assert.Equal(t, testCase.quality, quality)
		})
	}
}
