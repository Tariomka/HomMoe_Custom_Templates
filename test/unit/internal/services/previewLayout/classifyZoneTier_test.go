package previewLayout_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services"
	"github.com/stretchr/testify/assert"
)

func TestWhenZoneIsASpawnZone_ReturnsTierZero(t *testing.T) {
	// Arrange
	zone := entities.Zone{Name: "Spawn-A"}

	// Act
	tier := services.ClassifyZoneTier(zone)

	// Assert
	assert.Equal(t, 0, tier)
}

func TestWhenGuardedPoolContainsTierFiveMarker_ReturnsGoldTier(t *testing.T) {
	// Arrange
	zone := entities.Zone{Name: "Neutral-A", GuardedContentPool: []string{"classic_template_pool_random_t5_item"}}

	// Act
	tier := services.ClassifyZoneTier(zone)

	// Assert
	assert.Equal(t, 3, tier)
}

func TestWhenGuardedPoolContainsTierFourMarker_ReturnsGoldTier(t *testing.T) {
	// Arrange
	zone := entities.Zone{Name: "Neutral-A", GuardedContentPool: []string{"classic_template_pool_random_t4_item"}}

	// Act
	tier := services.ClassifyZoneTier(zone)

	// Assert
	assert.Equal(t, 3, tier)
}

func TestWhenGuardedPoolContainsTierThreeMarker_ReturnsSilverTier(t *testing.T) {
	// Arrange
	zone := entities.Zone{Name: "Neutral-A", GuardedContentPool: []string{"classic_template_pool_random_t3_item"}}

	// Act
	tier := services.ClassifyZoneTier(zone)

	// Assert
	assert.Equal(t, 2, tier)
}

func TestWhenGuardedPoolContainsTierTwoMarker_ReturnsBronzeTier(t *testing.T) {
	// Arrange
	zone := entities.Zone{Name: "Neutral-A", GuardedContentPool: []string{"classic_template_pool_random_t2_item"}}

	// Act
	tier := services.ClassifyZoneTier(zone)

	// Assert
	assert.Equal(t, 1, tier)
}

func TestWhenOnlyUnguardedPoolHasTierMarker_FallsBackToUnguardedPool(t *testing.T) {
	// Arrange
	zone := entities.Zone{
		Name:                 "Neutral-A",
		UnguardedContentPool: []string{"classic_template_pool_random_unguarded_t3_item"},
	}

	// Act
	tier := services.ClassifyZoneTier(zone)

	// Assert
	assert.Equal(t, 2, tier)
}

func TestWhenLayoutNameContainsSides_ReturnsBronzeTier(t *testing.T) {
	// Arrange
	zone := entities.Zone{Name: "Neutral-A", Layout: "zone_layout_sides"}

	// Act
	tier := services.ClassifyZoneTier(zone)

	// Assert
	assert.Equal(t, 1, tier)
}

func TestWhenLayoutNameContainsTreasure_ReturnsSilverTier(t *testing.T) {
	// Arrange
	zone := entities.Zone{Name: "Neutral-A", Layout: "zone_layout_treasure_zone"}

	// Act
	tier := services.ClassifyZoneTier(zone)

	// Assert
	assert.Equal(t, 2, tier)
}

func TestWhenLayoutNameContainsCenter_ReturnsGoldTier(t *testing.T) {
	// Arrange
	zone := entities.Zone{Name: "Neutral-A", Layout: "zone_layout_center"}

	// Act
	tier := services.ClassifyZoneTier(zone)

	// Assert
	assert.Equal(t, 3, tier)
}

func TestWhenZoneNameContainsLow_ReturnsBronzeTier(t *testing.T) {
	// Arrange
	zone := entities.Zone{Name: "Neutral-low"}

	// Act
	tier := services.ClassifyZoneTier(zone)

	// Assert
	assert.Equal(t, 1, tier)
}

func TestWhenZoneNameContainsMed_ReturnsSilverTier(t *testing.T) {
	// Arrange
	zone := entities.Zone{Name: "Neutral-med"}

	// Act
	tier := services.ClassifyZoneTier(zone)

	// Assert
	assert.Equal(t, 2, tier)
}

func TestWhenZoneNameContainsHigh_ReturnsGoldTier(t *testing.T) {
	// Arrange
	zone := entities.Zone{Name: "Neutral-high"}

	// Act
	tier := services.ClassifyZoneTier(zone)

	// Assert
	assert.Equal(t, 3, tier)
}

func TestWhenZoneNameContainsCore_ReturnsGoldTier(t *testing.T) {
	// Arrange
	zone := entities.Zone{Name: "Neutral-core"}

	// Act
	tier := services.ClassifyZoneTier(zone)

	// Assert
	assert.Equal(t, 3, tier)
}

func TestWhenNoTierHintIsPresent_DefaultsToBronzeTier(t *testing.T) {
	// Arrange
	zone := entities.Zone{Name: "Neutral-Z"}

	// Act
	tier := services.ClassifyZoneTier(zone)

	// Assert
	assert.Equal(t, 1, tier)
}
