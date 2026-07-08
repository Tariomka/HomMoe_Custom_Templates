package distancePresets_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/content_rules"
	"github.com/stretchr/testify/assert"
)

func TestWhenNamesAreListed_ReturnsPresetsInDeclarationOrder(t *testing.T) {
	// Arrange
	expected := []string{"Next To", "Near", "Medium", "Far", "Very Far"}

	// Act
	names := content_rules.GetDistanceDisplayNames()

	// Assert
	assert.Equal(t, expected, names)
}
