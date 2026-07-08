package mainObjectBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenGuardChanceIsResolved_ClampsValueIntoUnitRange(t *testing.T) {
	inRangeChance := gofakeit.Float64Range(0.01, 0.99)
	testCases := []struct {
		name           string
		providedChance float64
		expectedChance float64
	}{
		{
			name:           "WhenChanceIsWithinUnitRange_KeepsProvidedChance",
			providedChance: inRangeChance,
			expectedChance: inRangeChance,
		},
		{
			name:           "WhenChanceIsBelowZero_ClampsChanceToZero",
			providedChance: gofakeit.Float64Range(-10, -0.01),
			expectedChance: 0,
		},
		{
			name:           "WhenChanceIsAboveOne_ClampsChanceToOne",
			providedChance: gofakeit.Float64Range(1.01, 10),
			expectedChance: 1,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// Arrange
			builder := variant_content.NewObjectBuilder()

			// Act
			mainObject := builder.WithGuardChance(testCase.providedChance).Build()

			// Assert
			assert.Equal(t, entities.MainObject{GuardChance: testCase.expectedChance}, mainObject)
		})
	}
}
