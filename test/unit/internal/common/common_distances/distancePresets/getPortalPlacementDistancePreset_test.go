package distancePresets_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_distances"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestWhenNearIsRequested_ReturnsPortalNearBounds(t *testing.T) {
	t.Parallel()
	// Arrange
	expected := models.DistancePreset{Name: "Near", Min: 0.075, Max: 0.35}

	// Act
	result, found := common_distances.GetPortalPlacementDistancePreset("Near")

	// Assert
	assert.Equal(t, struct {
		Preset models.DistancePreset
		Found  bool
	}{expected, true}, struct {
		Preset models.DistancePreset
		Found  bool
	}{result, found})
}
