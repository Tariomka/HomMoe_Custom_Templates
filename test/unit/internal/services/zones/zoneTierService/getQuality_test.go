package zoneTierService_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_zones"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
	"github.com/stretchr/testify/assert"
)

func TestWhenZoneCharacteristicsVary_DetectsQualityAccordingly(t *testing.T) {
	t.Parallel()
	service := zones.NewZoneTierService()
	layoutValues := registry.GetLayoutValues()
	resourcePools := registry.GetResourcesContentPoolValues()
	testCases := []struct {
		subtestName string
		zone        template_model.Zone
		expected    neutral_zone.Quality
	}{
		{
			"WhenZoneIsPlayerSpawn_ReturnsUnknown",
			template_model.Zone{Name: "Spawn-A", Layout: layoutValues.Center},
			neutral_zone.QualityUnknown,
		},
		{
			"WhenLayoutIsUnrecognized_ReturnsUnknown",
			template_model.Zone{Name: "Neutral-B", Layout: "zone_layout_back"},
			neutral_zone.QualityUnknown,
		},
		{
			"WhenCenterLayoutHasRichTreasureResources_ReturnsHighest",
			template_model.Zone{
				Name:                 "Hub",
				Layout:               layoutValues.Center,
				GuardedContentPool:   []string{"pool_t5_treasure"},
				ResourcesContentPool: []string{resourcePools.TreasureZoneRich},
			},
			neutral_zone.QualityHighest,
		},
		{
			"WhenCenterLayoutHasOnlyTier5GuardedPool_ReturnsHighest",
			template_model.Zone{
				Name:               "Hub",
				Layout:             layoutValues.Center,
				GuardedContentPool: []string{"pool_t5_treasure", "pool_t5_item"},
			},
			neutral_zone.QualityHighest,
		},
		{
			"WhenCenterLayoutHasNoContentPools_ReturnsUnknown",
			template_model.Zone{Name: "Hub", Layout: layoutValues.Center},
			neutral_zone.QualityUnknown,
		},
		{
			"WhenCenterLayoutPoolsCarryNoTier5Marker_ReturnsUnknown",
			template_model.Zone{
				Name:               "Hub",
				Layout:             layoutValues.Center,
				GuardedContentPool: []string{"pool_t3_stuff"},
			},
			neutral_zone.QualityUnknown,
		},
		{
			"WhenTreasureLayoutHasRichStartResources_ReturnsHigh",
			template_model.Zone{
				Name:                 "Neutral-B",
				Layout:               layoutValues.TreasureZone,
				GuardedContentPool:   []string{"pool_without_tier_marker"},
				ResourcesContentPool: []string{resourcePools.StartZoneRich},
			},
			neutral_zone.QualityHigh,
		},
		{
			"WhenTreasureLayoutHasMixedTier4AndTier5Pools_ReturnsHigh",
			template_model.Zone{
				Name:               "Neutral-B",
				Layout:             layoutValues.TreasureZone,
				GuardedContentPool: []string{"pool_t4_stuff", "pool_t5_stuff"},
			},
			neutral_zone.QualityHigh,
		},
		{
			"WhenTreasureLayoutHasOnlyTier3Pools_ReturnsMedium",
			template_model.Zone{
				Name:                 "Neutral-B",
				Layout:               layoutValues.TreasureZone,
				GuardedContentPool:   []string{"pool_t3_stuff"},
				UnguardedContentPool: []string{"pool_t3_other"},
			},
			neutral_zone.QualityMedium,
		},
		{
			"WhenTreasureLayoutHasMediumStartResources_ReturnsMedium",
			template_model.Zone{
				Name:                 "Neutral-B",
				Layout:               layoutValues.TreasureZone,
				GuardedContentPool:   []string{"pool_without_tier_marker"},
				ResourcesContentPool: []string{resourcePools.StartZoneMedium},
			},
			neutral_zone.QualityMedium,
		},
		{
			"WhenTreasureLayoutHasNoContentPools_ReturnsUnknown",
			template_model.Zone{Name: "Neutral-B", Layout: layoutValues.TreasureZone},
			neutral_zone.QualityUnknown,
		},
		{
			"WhenTreasureLayoutPoolsCarryNoKnownMarker_ReturnsUnknown",
			template_model.Zone{
				Name:               "Neutral-B",
				Layout:             layoutValues.TreasureZone,
				GuardedContentPool: []string{"pool_t1_stuff"},
			},
			neutral_zone.QualityUnknown,
		},
		{
			"WhenSidesLayoutHasOnlyTier2Pools_ReturnsLow",
			template_model.Zone{
				Name:               "Neutral-B",
				Layout:             layoutValues.Sides,
				GuardedContentPool: []string{"pool_t2_stuff"},
			},
			neutral_zone.QualityLow,
		},
		{
			"WhenSidesLayoutHasPoorStartResources_ReturnsLow",
			template_model.Zone{
				Name:                 "Neutral-B",
				Layout:               layoutValues.Sides,
				GuardedContentPool:   []string{"pool_without_tier_marker"},
				ResourcesContentPool: []string{resourcePools.StartZonePoor},
			},
			neutral_zone.QualityLow,
		},
		{
			"WhenSidesLayoutHasOnlyTier1Pools_ReturnsLowest",
			template_model.Zone{
				Name:               "Neutral-B",
				Layout:             layoutValues.Sides,
				GuardedContentPool: []string{"pool_t1_stuff"},
			},
			neutral_zone.QualityLowest,
		},
		{
			"WhenSidesLayoutHasVeryPoorStartResources_ReturnsLowest",
			template_model.Zone{
				Name:                 "Neutral-B",
				Layout:               layoutValues.Sides,
				GuardedContentPool:   []string{"pool_without_tier_marker"},
				ResourcesContentPool: []string{resourcePools.StartZoneVeryPoor},
			},
			neutral_zone.QualityLowest,
		},
		{
			"WhenSidesLayoutHasNoContentPools_ReturnsUnknown",
			template_model.Zone{Name: "Neutral-B", Layout: layoutValues.Sides},
			neutral_zone.QualityUnknown,
		},
		{
			"WhenSidesLayoutPoolsCarryNoKnownMarker_ReturnsUnknown",
			template_model.Zone{
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
			quality := service.GetQuality(testCase.zone)

			// Assert
			assert.Equal(t, testCase.expected, quality)
		})
	}
}

func TestWhenGeneratedProfileRoundTrips_EveryQualityIsDetectedBack(t *testing.T) {
	t.Parallel()
	service := zones.NewZoneTierService()
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
			profile := common_zones.GetNeutralZoneProfile(testCase.quality)
			zone := template_model.Zone{
				Name:                 "Neutral-Z",
				Layout:               profile.Layout,
				GuardedContentPool:   profile.GuardedContentPool,
				UnguardedContentPool: profile.UnguardedContentPool,
				ResourcesContentPool: profile.ResourcesContentPool,
			}

			// Act
			quality := service.GetQuality(zone)

			// Assert
			assert.Equal(t, testCase.quality, quality)
		})
	}
}
