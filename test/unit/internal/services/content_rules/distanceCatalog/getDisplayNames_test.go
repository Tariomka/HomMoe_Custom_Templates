package distanceCatalog_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/content_rules"
	"github.com/stretchr/testify/assert"
)

func TestWhenNamesAreListed_ReturnsPresetsInDeclarationOrder(t *testing.T) {
	t.Parallel()
	// Arrange
	catalog := content_rules.NewDistanceCatalog()
	expected := []string{"Next To", "Near", "Medium", "Far", "Very Far"}

	// Act
	names := catalog.GetDisplayNames()

	// Assert
	assert.Equal(t, expected, names)
}

func TestWhenReturnedNamesAreMutated_NextResultRetainsCatalogValue(t *testing.T) {
	t.Parallel()
	// Arrange
	catalog := content_rules.NewDistanceCatalog()
	firstResult := catalog.GetDisplayNames()

	// Act
	firstResult[0] = "mutated"
	secondResult := catalog.GetDisplayNames()

	// Assert
	assert.Equal(t, "Next To", secondResult[0])
}
