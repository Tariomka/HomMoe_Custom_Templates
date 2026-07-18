package string_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
	"github.com/stretchr/testify/assert"
)

func TestWhenZoneNamesVary_ExtractsTrailingLabelAccordingly(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		subtestName string
		zoneName    string
		expected    string
	}{
		{"WhenNameHasSpawnPrefix_ReturnsLabelAfterPrefix", "Spawn-A", "A"},
		{"WhenNameHasNeutralPrefix_ReturnsLabelAfterPrefix", "Neutral-C", "C"},
		{"WhenNameHasNoKnownPrefix_ReturnsNameUnchanged", "Hub", "Hub"},
		{"WhenNameIsEmpty_ReturnsEmptyString", "", ""},
	}
	for _, testCase := range testCases {
		t.Run(testCase.subtestName, func(t *testing.T) {
			t.Parallel()
			// Arrange

			// Act
			actual := helpers.ExtractZoneLabel(testCase.zoneName)

			// Assert
			assert.Equal(t, testCase.expected, actual)
		})
	}
}
