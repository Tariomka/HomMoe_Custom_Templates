package zoneNameType_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/zone_helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/preview"
	"github.com/stretchr/testify/assert"
)

func TestWhenZoneNamesVary_ReturnsMatchingZoneType(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		subtestName string
		zoneName    string
		expected    preview.ZoneType
	}{
		{"WhenNameIsHub_ReturnsHubType", "Hub", preview.ZoneTypeHub},
		{"WhenNameIsHubLowercase_ReturnsHubType", "hub", preview.ZoneTypeHub},
		{"WhenNameHasHubPrefix_ReturnsHubType", "Hub-A", preview.ZoneTypeHub},
		{"WhenNameHasSpawnPrefix_ReturnsPlayerType", "Spawn-A", preview.ZoneTypePlayer},
		{"WhenNameHasNeutralPrefix_ReturnsNeutralType", "Neutral-C", preview.ZoneTypeNeutralZone},
		{"WhenNameHasNoKnownPrefix_ReturnsUnknownType", "Colosseum", preview.ZoneTypeUnknown},
		{"WhenNameIsEmpty_ReturnsUnknownType", "", preview.ZoneTypeUnknown},
	}
	for _, testCase := range testCases {
		t.Run(testCase.subtestName, func(t *testing.T) {
			t.Parallel()
			// Arrange

			// Act
			actual := zone_helpers.GetZoneTypeFromName(testCase.zoneName)

			// Assert
			assert.Equal(t, testCase.expected, actual)
		})
	}
}
