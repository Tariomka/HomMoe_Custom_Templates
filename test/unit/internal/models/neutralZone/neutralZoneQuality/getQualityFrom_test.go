package neutralZoneQuality_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutralZone"
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
		expected    neutralZone.Quality
	}{
		{
			"WhenZoneIsPlayerSpawn_ReturnsUnknown",
			entities.Zone{Name: "Spawn-A", Layout: layoutValues.Center},
			neutralZone.QualityUnknown,
		},
		{
			"WhenLayoutIsUnrecognized_ReturnsUnknown",
			entities.Zone{Name: "Neutral-B", Layout: "zone_layout_back"},
			neutralZone.QualityUnknown,
		},
		{
			"WhenCenterLayoutHasRichTreasureResources_ReturnsHighest",
			entities.Zone{
				Name:                 "Hub",
				Layout:               layoutValues.Center,
				GuardedContentPool:   []string{"pool_t5_treasure"},
				ResourcesContentPool: []string{resourcePools.TreasureZoneRich},
			},
			neutralZone.QualityHighest,
		},
		{
			"WhenCenterLayoutHasOnlyTier5GuardedPool_ReturnsHighest",
			entities.Zone{
				Name:               "Hub",
				Layout:             layoutValues.Center,
				GuardedContentPool: []string{"pool_t5_treasure", "pool_t5_item"},
			},
			neutralZone.QualityHighest,
		},
		{
			"WhenCenterLayoutHasNoContentPools_ReturnsUnknown",
			entities.Zone{Name: "Hub", Layout: layoutValues.Center},
			neutralZone.QualityUnknown,
		},
		{
			"WhenCenterLayoutPoolsCarryNoTier5Marker_ReturnsUnknown",
			entities.Zone{
				Name:               "Hub",
				Layout:             layoutValues.Center,
				GuardedContentPool: []string{"pool_t3_stuff"},
			},
			neutralZone.QualityUnknown,
		},
		{
			"WhenTreasureLayoutHasRichStartResources_ReturnsHigh",
			entities.Zone{
				Name:                 "Neutral-B",
				Layout:               layoutValues.TreasureZone,
				GuardedContentPool:   []string{"pool_without_tier_marker"},
				ResourcesContentPool: []string{resourcePools.StartZoneRich},
			},
			neutralZone.QualityHigh,
		},
		{
			"WhenTreasureLayoutHasMixedTier4AndTier5Pools_ReturnsHigh",
			entities.Zone{
				Name:               "Neutral-B",
				Layout:             layoutValues.TreasureZone,
				GuardedContentPool: []string{"pool_t4_stuff", "pool_t5_stuff"},
			},
			neutralZone.QualityHigh,
		},
		{
			"WhenTreasureLayoutHasOnlyTier3Pools_ReturnsMedium",
			entities.Zone{
				Name:                 "Neutral-B",
				Layout:               layoutValues.TreasureZone,
				GuardedContentPool:   []string{"pool_t3_stuff"},
				UnguardedContentPool: []string{"pool_t3_other"},
			},
			neutralZone.QualityMedium,
		},
		{
			"WhenTreasureLayoutHasMediumStartResources_ReturnsMedium",
			entities.Zone{
				Name:                 "Neutral-B",
				Layout:               layoutValues.TreasureZone,
				GuardedContentPool:   []string{"pool_without_tier_marker"},
				ResourcesContentPool: []string{resourcePools.StartZoneMedium},
			},
			neutralZone.QualityMedium,
		},
		{
			"WhenTreasureLayoutHasNoContentPools_ReturnsUnknown",
			entities.Zone{Name: "Neutral-B", Layout: layoutValues.TreasureZone},
			neutralZone.QualityUnknown,
		},
		{
			"WhenTreasureLayoutPoolsCarryNoKnownMarker_ReturnsUnknown",
			entities.Zone{
				Name:               "Neutral-B",
				Layout:             layoutValues.TreasureZone,
				GuardedContentPool: []string{"pool_t1_stuff"},
			},
			neutralZone.QualityUnknown,
		},
		{
			"WhenSidesLayoutHasOnlyTier2Pools_ReturnsLow",
			entities.Zone{
				Name:               "Neutral-B",
				Layout:             layoutValues.Sides,
				GuardedContentPool: []string{"pool_t2_stuff"},
			},
			neutralZone.QualityLow,
		},
		{
			"WhenSidesLayoutHasPoorStartResources_ReturnsLow",
			entities.Zone{
				Name:                 "Neutral-B",
				Layout:               layoutValues.Sides,
				GuardedContentPool:   []string{"pool_without_tier_marker"},
				ResourcesContentPool: []string{resourcePools.StartZonePoor},
			},
			neutralZone.QualityLow,
		},
		{
			"WhenSidesLayoutHasOnlyTier1Pools_ReturnsLowest",
			entities.Zone{
				Name:               "Neutral-B",
				Layout:             layoutValues.Sides,
				GuardedContentPool: []string{"pool_t1_stuff"},
			},
			neutralZone.QualityLowest,
		},
		{
			"WhenSidesLayoutHasVeryPoorStartResources_ReturnsLowest",
			entities.Zone{
				Name:                 "Neutral-B",
				Layout:               layoutValues.Sides,
				GuardedContentPool:   []string{"pool_without_tier_marker"},
				ResourcesContentPool: []string{resourcePools.StartZoneVeryPoor},
			},
			neutralZone.QualityLowest,
		},
		{
			"WhenSidesLayoutHasNoContentPools_ReturnsUnknown",
			entities.Zone{Name: "Neutral-B", Layout: layoutValues.Sides},
			neutralZone.QualityUnknown,
		},
		{
			"WhenSidesLayoutPoolsCarryNoKnownMarker_ReturnsUnknown",
			entities.Zone{
				Name:               "Neutral-B",
				Layout:             layoutValues.Sides,
				GuardedContentPool: []string{"pool_t3_stuff"},
			},
			neutralZone.QualityUnknown,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.subtestName, func(t *testing.T) {
			t.Parallel()
			// Arrange

			// Act
			quality := neutralZone.GetQualityFrom(testCase.zone)

			// Assert
			assert.Equal(t, testCase.expected, quality)
		})
	}
}

func TestWhenGeneratedProfileRoundTrips_EveryQualityIsDetectedBack(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		subtestName string
		quality     neutralZone.Quality
	}{
		{"WhenProfileQualityIsLowest_DetectsSameQuality", neutralZone.QualityLowest},
		{"WhenProfileQualityIsLow_DetectsSameQuality", neutralZone.QualityLow},
		{"WhenProfileQualityIsMedium_DetectsSameQuality", neutralZone.QualityMedium},
		{"WhenProfileQualityIsHigh_DetectsSameQuality", neutralZone.QualityHigh},
		{"WhenProfileQualityIsHighest_DetectsSameQuality", neutralZone.QualityHighest},
	}
	for _, testCase := range testCases {
		t.Run(testCase.subtestName, func(t *testing.T) {
			t.Parallel()
			// Arrange
			profile := neutralZone.NewNeutralZoneProfile(testCase.quality)
			zone := entities.Zone{
				Name:                 "Neutral-Z",
				Layout:               profile.Layout,
				GuardedContentPool:   profile.GuardedContentPool,
				UnguardedContentPool: profile.UnguardedContentPool,
				ResourcesContentPool: profile.ResourcesContentPool,
			}

			// Act
			quality := neutralZone.GetQualityFrom(zone)

			// Assert
			assert.Equal(t, testCase.quality, quality)
		})
	}
}
