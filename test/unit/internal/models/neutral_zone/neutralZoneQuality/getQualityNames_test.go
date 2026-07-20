package neutralZoneQuality_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/stretchr/testify/assert"
)

func TestWhenCalled_ReturnsSelectableQualityNamesWithoutPlatinum(t *testing.T) {
	t.Parallel()
	// Arrange
	expected := []string{"Plastic", "Bronze", "Silver", "Gold"}

	// Act
	qualityNames := neutral_zone.GetQualityNames()

	// Assert
	assert.Equal(t, expected, qualityNames)
}
