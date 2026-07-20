package zoneNameType_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/zone_helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
	"github.com/stretchr/testify/assert"
)

func TestWhenBothZonesArePlayerOwned_ReturnsUnknownQuality(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []entities.Zone{{Name: "Spawn-A"}, {Name: "Spawn-B"}}

	// Act
	quality := zone_helpers.GetZoneConnectionGuardQuality(
		"Spawn-A", "Spawn-B", zones, []string{"Spawn-A", "Spawn-B"})

	// Assert
	assert.Equal(t, neutral_zone.QualityUnknown, quality)
}

func TestWhenOneZoneIsHub_ReturnsHighestQuality(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []entities.Zone{{Name: "Spawn-A"}, {Name: "Hub"}}

	// Act
	quality := zone_helpers.GetZoneConnectionGuardQuality(
		"Spawn-A", "Hub", zones, []string{"Spawn-A"})

	// Assert
	assert.Equal(t, neutral_zone.QualityHighest, quality)
}

func TestWhenBothZonesAreNeutral_HigherQualityWins(t *testing.T) {
	t.Parallel()
	// Arrange
	layoutValues := registry.GetLayoutValues()
	zones := []entities.Zone{
		{
			Name:               "Neutral-C",
			Layout:             layoutValues.Sides,
			GuardedContentPool: []string{"pool_t2_x"},
		},
		{
			Name:               "Neutral-D",
			Layout:             layoutValues.TreasureZone,
			GuardedContentPool: []string{"pool_t4_x"},
		},
	}

	// Act
	quality := zone_helpers.GetZoneConnectionGuardQuality(
		"Neutral-C", "Neutral-D", zones, nil)

	// Assert
	assert.Equal(t, neutral_zone.QualityHigh, quality)
}
