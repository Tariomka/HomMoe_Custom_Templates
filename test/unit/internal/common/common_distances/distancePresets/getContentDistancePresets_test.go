package distancePresets_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_distances"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestWhenRequested_ReturnsContentBoundsInDisplayOrder(t *testing.T) {
	t.Parallel()
	// Arrange
	expected := []models.DistancePreset{
		{Name: "Next To", Min: 0.05, Max: 0.1},
		{Name: "Near", Min: 0.1, Max: 0.25},
		{Name: "Medium", Min: 0.25, Max: 0.5},
		{Name: "Far", Min: 0.5, Max: 0.75},
		{Name: "Very Far", Min: 0.75, Max: 0.9},
	}

	// Act
	result := common_distances.GetContentDistancePresets()

	// Assert
	assert.Equal(t, expected, result)
}

func TestWhenReturnedContentPresetsAreMutated_KeepsCatalogUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	expected := common_distances.GetContentDistancePresets()
	result := common_distances.GetContentDistancePresets()

	// Act
	result[0].Name = "mutated"

	// Assert
	assert.Equal(t, expected, common_distances.GetContentDistancePresets())
}
