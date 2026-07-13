package variantMappingManager_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/constants"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/content_rules"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenContentHasVariants_ReturnsOneMappingPerVariant(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name          string
		content       models.SidMapping
		expectedCount int
	}{
		{"WhenContentIsDragonUtopia_ReturnsFourMappings", constants.ContentIDs.DragonUtopia, 4},
		{"WhenContentIsPandoraBox_ReturnsTwentyEightMappings", constants.ContentIDs.PandoraBox, 28},
		{"WhenContentIsMontyHall_ReturnsFourMappings", constants.ContentIDs.MontyHall, 4},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			// Arrange

			// Act
			mappings := content_rules.GetVariantsForContent(testCase.content)

			// Assert
			assert.Len(t, mappings, testCase.expectedCount)
		})
	}
}

func TestWhenContentHasNoVariants_ReturnsEmptySlice(t *testing.T) {
	t.Parallel()
	// Arrange

	// Act
	mappings := content_rules.GetVariantsForContent(constants.ContentIDs.Watchtower)

	// Assert
	assert.Empty(t, mappings)
}

func TestWhenVariantsAreReturned_OrdersThemByVariantId(t *testing.T) {
	t.Parallel()
	// Arrange

	// Act
	mappings := content_rules.GetVariantsForContent(constants.ContentIDs.DragonUtopia)

	// Assert
	variantIDs := make([]int, 0, len(mappings))
	for _, mapping := range mappings {
		require.Len(t, mapping.Variants, 1)

		for _, variantID := range mapping.Variants {
			variantIDs = append(variantIDs, variantID.Key)
		}
	}
	assert.Equal(t, []int{0, 1, 2, 3}, variantIDs)
}

func TestWhenVariantsAreReturned_BindsRequestedContent(t *testing.T) {
	t.Parallel()
	// Arrange

	// Act
	mappings := content_rules.GetVariantsForContent(constants.ContentIDs.MontyHall)

	// Assert
	require.NotEmpty(t, mappings)
	assert.Equal(t, constants.ContentIDs.MontyHall, mappings[0].Content)
}
