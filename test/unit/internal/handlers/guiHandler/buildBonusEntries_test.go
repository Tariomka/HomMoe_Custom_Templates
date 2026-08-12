package guiHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenBonusEntriesRequested_ComposesThemFromTheRequest(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	receiver := gofakeit.Word()

	// Act
	result := handler.BuildBonusEntries(dtos.BonusCompositionRequestDto{
		PresetType:     config.BonusTownPortalFree,
		ReceiverFilter: receiver,
	})

	// Assert
	assert.Equal(t, []config.BonusEntry{{
		PresetType:     config.BonusTownPortalFree,
		ReceiverFilter: receiver,
	}}, result.Entries)
}
