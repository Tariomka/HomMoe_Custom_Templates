package zoneNameType_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/zone_helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
	"github.com/stretchr/testify/assert"
)

func TestWhenZoneNameIsEmpty_ReturnsLowQuality(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []entities.Zone{{Name: "Spawn-A"}}

	// Act
	quality := zone_helpers.GetZoneGuardQuality("", zones, []string{"Spawn-A"})

	// Assert
	assert.Equal(t, neutral_zone.QualityLow, quality)
}

func TestWhenZoneIsNotFound_ReturnsLowQuality(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []entities.Zone{{Name: "Spawn-A"}}

	// Act
	quality := zone_helpers.GetZoneGuardQuality("Neutral-Missing", zones, []string{"Spawn-A"})

	// Assert
	assert.Equal(t, neutral_zone.QualityLow, quality)
}

func TestWhenZoneIsPlayerOwned_ReturnsUnknownQuality(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []entities.Zone{{Name: "Spawn-A"}}

	// Act
	quality := zone_helpers.GetZoneGuardQuality("Spawn-A", zones, []string{"Spawn-A"})

	// Assert
	assert.Equal(t, neutral_zone.QualityUnknown, quality)
}

func TestWhenZoneIsHub_ReturnsHighestQuality(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []entities.Zone{{Name: "Hub"}}

	// Act
	quality := zone_helpers.GetZoneGuardQuality("Hub", zones, nil)

	// Assert
	assert.Equal(t, neutral_zone.QualityHighest, quality)
}

func TestWhenZoneHasPlayerPrefixButIsNotListed_ReturnsUnknownQuality(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []entities.Zone{{Name: "Spawn-B"}}

	// Act
	quality := zone_helpers.GetZoneGuardQuality("Spawn-B", zones, []string{"Spawn-A"})

	// Assert
	assert.Equal(t, neutral_zone.QualityUnknown, quality)
}

func TestWhenZoneIsNeutral_ReturnsQualityFromZoneContent(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []entities.Zone{
		{
			Name:               "Neutral-C",
			Layout:             registry.GetLayoutValues().TreasureZone,
			GuardedContentPool: []string{"pool_t3_x"},
		},
	}

	// Act
	quality := zone_helpers.GetZoneGuardQuality("Neutral-C", zones, nil)

	// Assert
	assert.Equal(t, neutral_zone.QualityMedium, quality)
}

func TestWhenZoneNameHasNoKnownPrefix_ReturnsUnknownQuality(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []entities.Zone{{Name: "Colosseum"}}

	// Act
	quality := zone_helpers.GetZoneGuardQuality("Colosseum", zones, nil)

	// Assert
	assert.Equal(t, neutral_zone.QualityUnknown, quality)
}
