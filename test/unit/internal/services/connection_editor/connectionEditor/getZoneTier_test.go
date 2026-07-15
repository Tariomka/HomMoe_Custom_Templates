package connectionEditor_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/connection_editor"
	"github.com/stretchr/testify/assert"
)

func TestWhenZoneIsClassified_ReturnsExpectedTier(t *testing.T) {
	t.Parallel()
	zones := []entities.Zone{
		{Name: "Spawn-A"},
		{Name: "Hub"},
		{Name: "Hub-1"},
		{Name: "Neutral-Gold4", GuardedContentPool: []string{"pool_t4_x"}},
		{Name: "Neutral-Gold5", GuardedContentPool: []string{"pool_t5_x"}},
		{Name: "Neutral-Bronze1", GuardedContentPool: []string{"pool_t1_x"}},
		{Name: "Neutral-Bronze2", GuardedContentPool: []string{"pool_t2_x"}},
		{Name: "Neutral-Silver", GuardedContentPool: []string{"pool_t3_x"}},
		{Name: "Neutral-NoPool"},
	}
	playerZoneNames := map[string]bool{"Spawn-A": true}

	testCases := []struct {
		name     string
		zoneName string
		expected connection_editor.ZoneTier
	}{
		{"WhenZoneNameIsEmpty_ReturnsBronze", "", connection_editor.ZoneTierBronze},
		{"WhenZoneIsUnknown_ReturnsBronze", "Nope", connection_editor.ZoneTierBronze},
		{"WhenZoneIsPlayerZone_ReturnsBronze", "Spawn-A", connection_editor.ZoneTierBronze},
		{"WhenZoneIsHub_ReturnsGold", "Hub", connection_editor.ZoneTierGold},
		{"WhenZoneHasHubPrefix_ReturnsGold", "Hub-1", connection_editor.ZoneTierGold},
		{"WhenNeutralPoolIsTier4_ReturnsGold", "Neutral-Gold4", connection_editor.ZoneTierGold},
		{"WhenNeutralPoolIsTier5_ReturnsGold", "Neutral-Gold5", connection_editor.ZoneTierGold},
		{"WhenNeutralPoolIsTier1_ReturnsBronze", "Neutral-Bronze1", connection_editor.ZoneTierBronze},
		{"WhenNeutralPoolIsTier2_ReturnsBronze", "Neutral-Bronze2", connection_editor.ZoneTierBronze},
		{"WhenNeutralPoolIsTier3_ReturnsSilver", "Neutral-Silver", connection_editor.ZoneTierSilver},
		{"WhenNeutralHasNoPool_ReturnsSilver", "Neutral-NoPool", connection_editor.ZoneTierSilver},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			// Arrange

			// Act
			tier := connection_editor.GetZoneTier(testCase.zoneName, zones, playerZoneNames)

			// Assert
			assert.Equal(t, testCase.expected, tier)
		})
	}
}
