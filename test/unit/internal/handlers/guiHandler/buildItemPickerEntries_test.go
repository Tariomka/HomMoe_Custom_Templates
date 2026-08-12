package guiHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/stretchr/testify/assert"
)

func TestWhenItemEntriesAreBuilt_TheCategoryBecomesTheGroup(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()

	// Act
	entries := handler.BuildItemPickerEntries([]dtos.PickerItemDto{{
		Sid:      "Sid.Sword",
		Name:     "Sword",
		Category: "Weapons",
	}})

	// Assert
	assert.Equal(t, "Weapons", entries[0].Group)
}
