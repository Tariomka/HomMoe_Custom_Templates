package zoneTierService_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
	"github.com/stretchr/testify/assert"
)

func TestWhenBothZonesArePlayerOwned_ReturnsUnknownQuality(t *testing.T) {
	t.Parallel()
	// Arrange
	service := zones.NewZoneTierService()
	zoneList := []entities.Zone{{Name: "Spawn-A"}, {Name: "Spawn-B"}}

	// Act
	quality := service.GetConnectionGuardQuality(
		"Spawn-A", "Spawn-B", zoneList, []string{"Spawn-A", "Spawn-B"})

	// Assert
	assert.Equal(t, neutral_zone.QualityUnknown, quality)
}

func TestWhenOneZoneIsHub_ReturnsHighestQuality(t *testing.T) {
	t.Parallel()
	// Arrange
	service := zones.NewZoneTierService()
	zoneList := []entities.Zone{{Name: "Spawn-A"}, {Name: "Hub"}}

	// Act
	quality := service.GetConnectionGuardQuality("Spawn-A", "Hub", zoneList, []string{"Spawn-A"})

	// Assert
	assert.Equal(t, neutral_zone.QualityHighest, quality)
}

func TestWhenBothZonesAreNeutral_HigherQualityWins(t *testing.T) {
	t.Parallel()
	// Arrange
	service := zones.NewZoneTierService()
	layoutValues := registry.GetLayoutValues()
	zoneList := []entities.Zone{
		{Name: "Neutral-C", Layout: layoutValues.Sides, GuardedContentPool: []string{"pool_t2_x"}},
		{Name: "Neutral-D", Layout: layoutValues.TreasureZone, GuardedContentPool: []string{"pool_t4_x"}},
	}

	// Act
	quality := service.GetConnectionGuardQuality("Neutral-C", "Neutral-D", zoneList, nil)

	// Assert
	assert.Equal(t, neutral_zone.QualityHigh, quality)
}
