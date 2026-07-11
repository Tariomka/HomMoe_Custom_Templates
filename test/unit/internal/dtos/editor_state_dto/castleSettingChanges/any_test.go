package castleSettingChanges_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/stretchr/testify/assert"
)

func TestWhenSingleFlagIsSet_ReportsAnyChange(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		subtestName string
		changes     editor_state_dto.CastleSettingChanges
	}{
		{"WhenPlayerCastlesFlagIsSet_ReportsAnyChange", editor_state_dto.CastleSettingChanges{PlayerCastles: true}},
		{"WhenNeutralSimpleFlagIsSet_ReportsAnyChange", editor_state_dto.CastleSettingChanges{NeutralSimple: true}},
		{"WhenNeutralLowFlagIsSet_ReportsAnyChange", editor_state_dto.CastleSettingChanges{NeutralLow: true}},
		{"WhenNeutralMediumFlagIsSet_ReportsAnyChange", editor_state_dto.CastleSettingChanges{NeutralMedium: true}},
		{"WhenNeutralHighFlagIsSet_ReportsAnyChange", editor_state_dto.CastleSettingChanges{NeutralHigh: true}},
		{"WhenHubFlagIsSet_ReportsAnyChange", editor_state_dto.CastleSettingChanges{Hub: true}},
	}
	for _, testCase := range testCases {
		t.Run(testCase.subtestName, func(t *testing.T) {
			t.Parallel()
			// Arrange
			changes := testCase.changes

			// Act
			anyChange := changes.Any()

			// Assert
			assert.True(t, anyChange)
		})
	}
}

func TestWhenNoFlagIsSet_ReportsNoChange(t *testing.T) {
	t.Parallel()
	// Arrange
	changes := editor_state_dto.CastleSettingChanges{}

	// Act
	anyChange := changes.Any()

	// Assert
	assert.False(t, anyChange)
}
