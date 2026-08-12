package bonusEntryService_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/bonuses"
	"github.com/stretchr/testify/assert"
)

func TestWhenTheSpellCountIsLabelled_UsesThePluralAwareCaption(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name     string
		count    int
		expected string
	}{
		{"WhenNoSpellIsPicked_ReturnsThePlainCaption", 0, "Spells"},
		{"WhenOneSpellIsPicked_ReturnsTheSingularCaption", 1, "1 spell picked"},
		{"WhenSeveralSpellsArePicked_ReturnsThePluralCaption", 4, "4 spells picked"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			// Arrange
			service := bonuses.NewBonusEntryService()

			// Act
			label := service.GetSpellCountLabel(testCase.count)

			// Assert
			assert.Equal(t, testCase.expected, label)
		})
	}
}
