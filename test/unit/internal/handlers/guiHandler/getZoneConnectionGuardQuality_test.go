package guiHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/stretchr/testify/assert"
)

func TestWhenConnectionIncludesHub_ReturnsHighestQuality(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := handlers.NewGuiHandler()
	zoneList := []entities.Zone{{Name: "Spawn-A"}, {Name: "Hub"}}

	// Act
	quality := handler.GetZoneConnectionGuardQuality(
		"Spawn-A",
		"Hub",
		zoneList,
		map[string]bool{"Spawn-A": true},
	)

	// Assert
	assert.Equal(t, neutral_zone.QualityHighest, quality)
}
