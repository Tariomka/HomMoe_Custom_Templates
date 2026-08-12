package guiHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
	"github.com/stretchr/testify/assert"
)

func TestWhenNeutralZoneHasTierThreeContent_ReturnsMediumQuality(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	zone := entities.Zone{
		Name:               "Neutral-C",
		Layout:             registry.GetLayoutValues().TreasureZone,
		GuardedContentPool: []string{"pool_t3_x"},
	}

	// Act
	quality := handler.GetZoneQuality(zone)

	// Assert
	assert.Equal(t, neutral_zone.QualityMedium, quality)
}
