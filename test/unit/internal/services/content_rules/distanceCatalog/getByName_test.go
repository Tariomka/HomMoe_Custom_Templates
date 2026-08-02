package distanceCatalog_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/content_rules"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenPresetNameIsKnown_ReturnsItsVariation(t *testing.T) {
	t.Parallel()
	catalog := content_rules.NewDistanceCatalog()
	testCases := []struct {
		name       string
		lookupName string
		expected   models.DistancePreset
	}{
		{
			"WhenNameIsNextTo_ReturnsNextToBounds",
			"Next To",
			models.DistancePreset{Name: "Next To", Min: 0.05, Max: 0.1},
		},
		{"WhenNameIsNear_ReturnsNearBounds", "Near", models.DistancePreset{Name: "Near", Min: 0.1, Max: 0.25}},
		{"WhenNameIsMedium_ReturnsMediumBounds", "Medium", models.DistancePreset{Name: "Medium", Min: 0.25, Max: 0.5}},
		{"WhenNameIsFar_ReturnsFarBounds", "Far", models.DistancePreset{Name: "Far", Min: 0.5, Max: 0.75}},
		{
			"WhenNameIsVeryFar_ReturnsVeryFarBounds",
			"Very Far",
			models.DistancePreset{Name: "Very Far", Min: 0.75, Max: 0.9},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			// Arrange

			// Act
			variation, ok := catalog.GetByName(testCase.lookupName)

			// Assert
			require.True(t, ok)
			assert.Equal(t, testCase.expected, variation)
		})
	}
}

func TestWhenNameDiffersOnlyByCase_ResolvesIt(t *testing.T) {
	t.Parallel()
	// Arrange
	catalog := content_rules.NewDistanceCatalog()

	// Act
	variation, ok := catalog.GetByName("mEdIuM")

	// Assert
	require.True(t, ok)
	assert.Equal(t, models.DistancePreset{Name: "Medium", Min: 0.25, Max: 0.5}, variation)
}

func TestWhenNameHasSurroundingWhitespace_ResolvesIt(t *testing.T) {
	t.Parallel()
	// Arrange
	catalog := content_rules.NewDistanceCatalog()

	// Act
	variation, ok := catalog.GetByName("  Far  ")

	// Assert
	require.True(t, ok)
	assert.Equal(t, models.DistancePreset{Name: "Far", Min: 0.5, Max: 0.75}, variation)
}

func TestWhenNameIsUnknown_ReturnsNotOk(t *testing.T) {
	t.Parallel()
	// Arrange
	catalog := content_rules.NewDistanceCatalog()

	// Act
	_, ok := catalog.GetByName("Whatever")

	// Assert
	assert.False(t, ok)
}
