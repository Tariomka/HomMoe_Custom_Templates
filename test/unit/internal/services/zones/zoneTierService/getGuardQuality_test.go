package zoneTierService_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
	"github.com/stretchr/testify/assert"
)

func TestWhenZoneNameIsEmpty_ReturnsLowQuality(t *testing.T) {
	t.Parallel()
	// Arrange
	service := zones.NewZoneTierService()
	zoneList := []template_model.Zone{{Name: "Spawn-A"}}

	// Act
	quality := service.GetGuardQuality("", zoneList, []string{"Spawn-A"})

	// Assert
	assert.Equal(t, neutral_zone.QualityLow, quality)
}

func TestWhenZoneIsNotFound_ReturnsLowQuality(t *testing.T) {
	t.Parallel()
	// Arrange
	service := zones.NewZoneTierService()
	zoneList := []template_model.Zone{{Name: "Spawn-A"}}

	// Act
	quality := service.GetGuardQuality("Neutral-Missing", zoneList, []string{"Spawn-A"})

	// Assert
	assert.Equal(t, neutral_zone.QualityLow, quality)
}

func TestWhenZoneIsPlayerOwned_ReturnsUnknownQuality(t *testing.T) {
	t.Parallel()
	// Arrange
	service := zones.NewZoneTierService()
	zoneList := []template_model.Zone{{Name: "Spawn-A"}}

	// Act
	quality := service.GetGuardQuality("Spawn-A", zoneList, []string{"Spawn-A"})

	// Assert
	assert.Equal(t, neutral_zone.QualityUnknown, quality)
}

func TestWhenZoneIsHub_ReturnsHighestQuality(t *testing.T) {
	t.Parallel()
	// Arrange
	service := zones.NewZoneTierService()
	zoneList := []template_model.Zone{{Name: "Hub"}}

	// Act
	quality := service.GetGuardQuality("Hub", zoneList, nil)

	// Assert
	assert.Equal(t, neutral_zone.QualityHighest, quality)
}

func TestWhenPlayerPrefixIsNotListed_ReturnsUnknownQuality(t *testing.T) {
	t.Parallel()
	// Arrange
	service := zones.NewZoneTierService()
	zoneList := []template_model.Zone{{Name: "Spawn-B"}}

	// Act
	quality := service.GetGuardQuality("Spawn-B", zoneList, []string{"Spawn-A"})

	// Assert
	assert.Equal(t, neutral_zone.QualityUnknown, quality)
}

func TestWhenZoneIsNeutral_ReturnsQualityFromContent(t *testing.T) {
	t.Parallel()
	// Arrange
	service := zones.NewZoneTierService()
	zoneList := []template_model.Zone{{
		Name:               "Neutral-C",
		Layout:             registry.GetLayoutValues().TreasureZone,
		GuardedContentPool: []string{"pool_t3_x"},
	}}

	// Act
	quality := service.GetGuardQuality("Neutral-C", zoneList, nil)

	// Assert
	assert.Equal(t, neutral_zone.QualityMedium, quality)
}

func TestWhenZoneNameHasNoKnownPrefix_ReturnsUnknownQuality(t *testing.T) {
	t.Parallel()
	// Arrange
	service := zones.NewZoneTierService()
	zoneList := []template_model.Zone{{Name: "Colosseum"}}

	// Act
	quality := service.GetGuardQuality("Colosseum", zoneList, nil)

	// Assert
	assert.Equal(t, neutral_zone.QualityUnknown, quality)
}
