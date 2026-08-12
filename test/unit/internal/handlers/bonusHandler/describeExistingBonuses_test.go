package bonusHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenExistingBonusesAreDescribed_ReturnsTheServiceSummary(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newBonusHandlerFixture()
	existing := []config.BonusEntry{{PresetType: config.BonusSpell, Param: gofakeit.UUID()}}
	expected := dtos.ExistingBonusesDto{
		Keys:     map[string]bool{gofakeit.UUID(): true},
		SpellIDs: []string{gofakeit.UUID()},
	}
	fixture.bonusService.On("DescribeExistingBonuses", existing).Return(expected)

	// Act
	summary := fixture.handler.DescribeExistingBonuses(existing)

	// Assert
	assert.Equal(t, expected, summary)
}
