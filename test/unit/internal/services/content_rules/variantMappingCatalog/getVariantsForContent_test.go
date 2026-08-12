package variantMappingCatalog_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/content_rules"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// SIDs are duplicated here on purpose: internal-layer tests must not import the GUI catalogue.
const (
	dragonUtopiaSid = "dragon_utopia"
	pandoraBoxSid   = "pandora_box"
	montyHallSid    = "monty_hall"
	watchtowerSid   = "watchtower"
)

func TestWhenContentHasVariants_ReturnsOneMappingPerVariant(t *testing.T) {
	t.Parallel()
	catalog := content_rules.NewVariantMappingCatalog()
	testCases := []struct {
		name          string
		content       models.SidMapping
		expectedCount int
	}{
		{"WhenContentIsDragonUtopia_ReturnsFourMappings", models.SidMapping{Sid: dragonUtopiaSid}, 4},
		{"WhenContentIsPandoraBox_ReturnsTwentyEightMappings", models.SidMapping{Sid: pandoraBoxSid}, 28},
		{"WhenContentIsMontyHall_ReturnsFourMappings", models.SidMapping{Sid: montyHallSid}, 4},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			// Arrange

			// Act
			mappings := catalog.GetVariantsForContent(testCase.content)

			// Assert
			assert.Len(t, mappings, testCase.expectedCount)
		})
	}
}

func TestWhenContentHasNoVariants_ReturnsEmptySlice(t *testing.T) {
	t.Parallel()
	// Arrange
	catalog := content_rules.NewVariantMappingCatalog()

	// Act
	mappings := catalog.GetVariantsForContent(models.SidMapping{Sid: watchtowerSid})

	// Assert
	assert.Empty(t, mappings)
}

func TestWhenVariantsAreReturned_OrdersThemByVariantId(t *testing.T) {
	t.Parallel()
	// Arrange
	catalog := content_rules.NewVariantMappingCatalog()

	// Act
	mappings := catalog.GetVariantsForContent(models.SidMapping{Sid: dragonUtopiaSid})

	// Assert
	variantIDs := make([]int, 0, len(mappings))
	for _, mapping := range mappings {
		require.Len(t, mapping.Variants, 1)
		for _, variant := range mapping.Variants {
			variantIDs = append(variantIDs, variant.Key)
		}
	}
	assert.Equal(t, []int{0, 1, 2, 3}, variantIDs)
}

func TestWhenVariantsAreReturned_BindsRequestedContent(t *testing.T) {
	t.Parallel()
	// Arrange
	catalog := content_rules.NewVariantMappingCatalog()
	montyHall := models.SidMapping{Sid: montyHallSid, Name: "The Monty Hall"}

	// Act
	mappings := catalog.GetVariantsForContent(montyHall)

	// Assert
	require.NotEmpty(t, mappings)
	assert.Equal(t, montyHall, mappings[0].Content)
}

func TestWhenReturnedVariantIsMutated_NextResultRetainsCatalogValue(t *testing.T) {
	t.Parallel()
	// Arrange
	catalog := content_rules.NewVariantMappingCatalog()
	dragonUtopia := models.SidMapping{Sid: dragonUtopiaSid}
	firstResult := catalog.GetVariantsForContent(dragonUtopia)
	expected := firstResult[0].Variants[0].Value

	// Act
	firstResult[0].Variants[0].Value = "mutated"
	secondResult := catalog.GetVariantsForContent(dragonUtopia)

	// Assert
	assert.Equal(t, expected, secondResult[0].Variants[0].Value)
}
