package distancePresets_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_distances"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestWhenNearIsRequested_ReturnsContentNearBounds(t *testing.T) {
	t.Parallel()
	// Arrange
	expected := models.DistancePreset{Name: "Near", Min: 0.1, Max: 0.25}

	// Act
	result, found := common_distances.GetContentDistancePreset("Near")

	// Assert
	assert.Equal(t, struct {
		Preset models.DistancePreset
		Found  bool
	}{expected, true}, struct {
		Preset models.DistancePreset
		Found  bool
	}{result, found})
}

func TestWhenUnknownNameIsRequested_ReturnsNotFound(t *testing.T) {
	t.Parallel()
	// Arrange
	unknownName := "Unknown"

	// Act
	_, found := common_distances.GetContentDistancePreset(unknownName)

	// Assert
	assert.False(t, found)
}
