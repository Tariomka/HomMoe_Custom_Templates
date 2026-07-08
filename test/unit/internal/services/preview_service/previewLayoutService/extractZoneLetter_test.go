package previewLayoutService_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/preview_service"
	"github.com/stretchr/testify/assert"
)

func TestWhenZoneNamesVary_ExtractsTrailingLetterAccordingly(t *testing.T) {
	testCases := []struct {
		subtestName string
		zoneName    string
		expected    string
	}{
		{"WhenNameHasSpawnPrefix_ReturnsLetterAfterPrefix", "Spawn-A", "A"},
		{"WhenNameHasNeutralPrefix_ReturnsLetterAfterPrefix", "Neutral-C", "C"},
		{"WhenNameHasNoKnownPrefix_ReturnsNameUnchanged", "Hub", "Hub"},
		{"WhenNameIsEmpty_ReturnsEmptyString", "", ""},
	}
	for _, testCase := range testCases {
		t.Run(testCase.subtestName, func(t *testing.T) {
			// Arrange

			// Act
			actual := preview_service.ExtractZoneLetter(testCase.zoneName)

			// Assert
			assert.Equal(t, testCase.expected, actual)
		})
	}
}
