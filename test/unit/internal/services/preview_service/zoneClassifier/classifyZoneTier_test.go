package zoneClassifier_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/preview"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/preview_service"
	"github.com/stretchr/testify/assert"
)

func TestWhenZoneCharacteristicsVary_ClassifiesTierAccordingly(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		subtestName string
		zone        entities.Zone
		expected    preview.ZoneTier
	}{
		{
			"WhenZoneIsSpawn_ReturnsPlasticTier",
			entities.Zone{Name: "Spawn-A"},
			preview.TierPlastic,
		},
		{
			"WhenGuardedPoolCarriesTier5Marker_ReturnsGoldTier",
			entities.Zone{Name: "Neutral-B", GuardedContentPool: []string{"pool_t5_treasure"}},
			preview.TierGold,
		},
		{
			"WhenGuardedPoolCarriesTier4Marker_ReturnsGoldTier",
			entities.Zone{Name: "Neutral-B", GuardedContentPool: []string{"pool_t4_treasure"}},
			preview.TierGold,
		},
		{
			"WhenGuardedPoolCarriesTier3Marker_ReturnsSilverTier",
			entities.Zone{Name: "Neutral-B", GuardedContentPool: []string{"pool_t3_stuff"}},
			preview.TierSilver,
		},
		{
			"WhenGuardedPoolCarriesTier2Marker_ReturnsBronzeTier",
			entities.Zone{Name: "Neutral-B", GuardedContentPool: []string{"pool_t2_stuff"}},
			preview.TierBronze,
		},
		{
			"WhenGuardedPoolCarriesTier1Marker_ReturnsPlasticTier",
			entities.Zone{Name: "Neutral-B", GuardedContentPool: []string{"pool_t1_stuff"}},
			preview.TierPlastic,
		},
		{
			"WhenTier5PoolPairsWithRichTreasureResources_ReturnsPlatinumTier",
			entities.Zone{
				Name:                 "Hub",
				GuardedContentPool:   []string{"pool_t5_treasure"},
				ResourcesContentPool: []string{"content_pool_general_resources_treasure_zone_rich"},
			},
			preview.TierPlatinum,
		},
		{
			"WhenTier4PoolPairsWithRichTreasureResources_ReturnsGoldTier",
			entities.Zone{
				Name:                 "Neutral-B",
				GuardedContentPool:   []string{"pool_t4_treasure"},
				ResourcesContentPool: []string{"content_pool_general_resources_treasure_zone_rich"},
			},
			preview.TierGold,
		},
		{
			"WhenOnlyUnguardedPoolCarriesTierMarker_ReturnsThatTier",
			entities.Zone{Name: "Neutral-B", UnguardedContentPool: []string{"pool_t3_stuff"}},
			preview.TierSilver,
		},
		{
			"WhenLayoutContainsSides_ReturnsBronzeTier",
			entities.Zone{Name: "Neutral-B", Layout: "zone_sides"},
			preview.TierBronze,
		},
		{
			"WhenLayoutContainsTreasure_ReturnsSilverTier",
			entities.Zone{Name: "Neutral-B", Layout: "zone_treasure"},
			preview.TierSilver,
		},
		{
			"WhenLayoutContainsCenter_ReturnsGoldTier",
			entities.Zone{Name: "Neutral-B", Layout: "zone_center"},
			preview.TierGold,
		},
		{
			"WhenNameContainsLow_ReturnsBronzeTier",
			entities.Zone{Name: "low-neutral"},
			preview.TierBronze,
		},
		{
			"WhenNameContainsMed_ReturnsSilverTier",
			entities.Zone{Name: "med-neutral"},
			preview.TierSilver,
		},
		{
			"WhenNameContainsHigh_ReturnsGoldTier",
			entities.Zone{Name: "high-neutral"},
			preview.TierGold,
		},
		{
			"WhenNoHintIsPresent_DefaultsToBronzeTier",
			entities.Zone{Name: "Neutral-B"},
			preview.TierBronze,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.subtestName, func(t *testing.T) {
			t.Parallel()
			// Arrange

			// Act
			actual := preview_service.ClassifyZoneTier(testCase.zone)

			// Assert
			assert.Equal(t, testCase.expected, actual)
		})
	}
}
