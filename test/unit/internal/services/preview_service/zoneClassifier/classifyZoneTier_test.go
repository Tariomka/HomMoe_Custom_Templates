package zoneClassifier_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/preview_service"
	"github.com/stretchr/testify/assert"
)

func TestWhenZoneCharacteristicsVary_ClassifiesTierAccordingly(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		subtestName string
		zone        entities.Zone
		expected    int
	}{
		{
			"WhenZoneIsSpawn_ReturnsTierZero",
			entities.Zone{Name: "Spawn-A"},
			0,
		},
		{
			"WhenGuardedPoolCarriesTier5Marker_ReturnsGoldTier",
			entities.Zone{Name: "Neutral-B", GuardedContentPool: []string{"pool_t5_treasure"}},
			3,
		},
		{
			"WhenGuardedPoolCarriesTier4Marker_ReturnsGoldTier",
			entities.Zone{Name: "Neutral-B", GuardedContentPool: []string{"pool_t4_treasure"}},
			3,
		},
		{
			"WhenGuardedPoolCarriesTier3Marker_ReturnsSilverTier",
			entities.Zone{Name: "Neutral-B", GuardedContentPool: []string{"pool_t3_stuff"}},
			2,
		},
		{
			"WhenGuardedPoolCarriesTier2Marker_ReturnsBronzeTier",
			entities.Zone{Name: "Neutral-B", GuardedContentPool: []string{"pool_t2_stuff"}},
			1,
		},
		{
			"WhenOnlyUnguardedPoolCarriesTierMarker_ReturnsThatTier",
			entities.Zone{Name: "Neutral-B", UnguardedContentPool: []string{"pool_t3_stuff"}},
			2,
		},
		{
			"WhenLayoutContainsSides_ReturnsBronzeTier",
			entities.Zone{Name: "Neutral-B", Layout: "zone_sides"},
			1,
		},
		{
			"WhenLayoutContainsTreasure_ReturnsSilverTier",
			entities.Zone{Name: "Neutral-B", Layout: "zone_treasure"},
			2,
		},
		{
			"WhenLayoutContainsCenter_ReturnsGoldTier",
			entities.Zone{Name: "Neutral-B", Layout: "zone_center"},
			3,
		},
		{
			"WhenNameContainsLow_ReturnsBronzeTier",
			entities.Zone{Name: "low-neutral"},
			1,
		},
		{
			"WhenNameContainsMed_ReturnsSilverTier",
			entities.Zone{Name: "med-neutral"},
			2,
		},
		{
			"WhenNameContainsHigh_ReturnsGoldTier",
			entities.Zone{Name: "high-neutral"},
			3,
		},
		{
			"WhenNoHintIsPresent_DefaultsToBronzeTier",
			entities.Zone{Name: "Neutral-B"},
			1,
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
