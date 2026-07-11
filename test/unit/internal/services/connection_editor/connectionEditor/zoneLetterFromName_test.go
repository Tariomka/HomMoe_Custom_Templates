package connectionEditor_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/connection_editor"
	"github.com/stretchr/testify/assert"
)

func TestWhenNameIsParsed_ReturnsPartAfterFirstDash(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name     string
		zoneName string
		expected string
	}{
		{"WhenNameIsSpawnA_ReturnsA", "Spawn-A", "A"},
		{"WhenNameIsNeutralC_ReturnsC", "Neutral-C", "C"},
		{"WhenLabelHasTwoLetters_ReturnsBothLetters", "Neutral-AA", "AA"},
		{"WhenNameHasNoDash_ReturnsWholeName", "Hub", "Hub"},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			// Arrange

			// Act
			letter := connection_editor.ZoneLetterFromName(testCase.zoneName)

			// Assert
			assert.Equal(t, testCase.expected, letter)
		})
	}
}
