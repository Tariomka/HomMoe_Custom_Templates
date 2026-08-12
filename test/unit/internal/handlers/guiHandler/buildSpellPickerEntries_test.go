package guiHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/stretchr/testify/assert"
)

func TestWhenSpellEntriesAreBuilt_TheTierBecomesTheBadge(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()

	// Act
	entries := handler.BuildSpellPickerEntries([]dtos.PickerSpellDto{{
		Sid:  "Sid.Bless",
		Name: "Bless",
		Tier: 2,
	}})

	// Assert
	assert.Equal(t, "[T2]", entries[0].Badge)
}
