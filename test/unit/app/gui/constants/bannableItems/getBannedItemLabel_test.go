package bannableItems_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/constants"
	"github.com/stretchr/testify/assert"
)

func TestWhenSidIsCatalogued_ReturnsItsNameAndCategory(t *testing.T) {
	t.Parallel()
	// Arrange
	expected := constants.GetBannableItemsWithExclusions(nil)[0]

	// Act
	name, category := constants.GetBannedItemLabel(expected.Sid)

	// Assert
	assert.Equal(t, expected.Name, name)
	assert.Equal(t, expected.Category, category)
}

func TestWhenSidIsUncatalogued_FallsBackToTheDerivedDisplayName(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	name, _ := constants.GetBannedItemLabel("modded_relic_artifact")

	// Assert
	assert.Equal(t, "Modded relic", name)
}

func TestWhenSidIsUncatalogued_FallsBackToTheMiscCategory(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	_, category := constants.GetBannedItemLabel("modded_relic_artifact")

	// Assert
	assert.Equal(t, "Misc", category)
}
