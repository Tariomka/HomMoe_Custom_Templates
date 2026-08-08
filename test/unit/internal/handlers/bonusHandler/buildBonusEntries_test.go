package bonusHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenBonusEntriesAreBuilt_ReturnsTheServiceResult(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newBonusHandlerFixture()
	request := dtos.BonusCompositionRequestDto{
		PresetType:     config.BonusStartingItem,
		ReceiverFilter: gofakeit.Word(),
		ItemText:       gofakeit.UUID(),
	}
	expected := dtos.BonusCompositionResultDto{
		Entries: []config.BonusEntry{{PresetType: config.BonusStartingItem, Param: gofakeit.UUID()}},
	}
	fixture.bonusService.On("BuildBonusEntries", request).Return(expected)

	// Act
	result := fixture.handler.BuildBonusEntries(request)

	// Assert
	assert.Equal(t, expected, result)
}
