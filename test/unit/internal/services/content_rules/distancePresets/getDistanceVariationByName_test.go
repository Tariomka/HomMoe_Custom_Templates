package distancePresets_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/content_rules"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The expected Min/Max bounds pin the C# DistancePresets values exactly.
func TestWhenPresetNameIsKnown_ReturnsItsVariation(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name       string
		lookupName string
		expected   content_rules.DistanceVariation
	}{
		{
			"WhenNameIsNextTo_ReturnsNextToBounds",
			"Next To",
			content_rules.DistanceVariation{Name: "Next To", Min: 0.05, Max: 0.1},
		},
		{
			"WhenNameIsNear_ReturnsNearBounds",
			"Near",
			content_rules.DistanceVariation{Name: "Near", Min: 0.1, Max: 0.25},
		},
		{
			"WhenNameIsMedium_ReturnsMediumBounds",
			"Medium",
			content_rules.DistanceVariation{Name: "Medium", Min: 0.25, Max: 0.5},
		},
		{"WhenNameIsFar_ReturnsFarBounds", "Far", content_rules.DistanceVariation{Name: "Far", Min: 0.5, Max: 0.75}},
		{
			"WhenNameIsVeryFar_ReturnsVeryFarBounds",
			"Very Far",
			content_rules.DistanceVariation{Name: "Very Far", Min: 0.75, Max: 0.9},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			// Arrange

			// Act
			variation, ok := content_rules.GetDistanceVariationByName(testCase.lookupName)

			// Assert
			require.True(t, ok)
			assert.Equal(t, testCase.expected, variation)
		})
	}
}

func TestWhenNameDiffersOnlyByCase_ResolvesIt(t *testing.T) {
	t.Parallel()
	// Arrange

	// Act
	variation, ok := content_rules.GetDistanceVariationByName("mEdIuM")

	// Assert
	require.True(t, ok)
	assert.Equal(t, content_rules.DistanceMedium, variation)
}

func TestWhenNameHasSurroundingWhitespace_ResolvesIt(t *testing.T) {
	t.Parallel()
	// Arrange

	// Act
	variation, ok := content_rules.GetDistanceVariationByName("  Far  ")

	// Assert
	require.True(t, ok)
	assert.Equal(t, content_rules.DistanceFar, variation)
}

func TestWhenNameIsUnknown_ReturnsNotOk(t *testing.T) {
	t.Parallel()
	// Arrange

	// Act
	_, ok := content_rules.GetDistanceVariationByName("Whatever")

	// Assert
	assert.False(t, ok)
}
