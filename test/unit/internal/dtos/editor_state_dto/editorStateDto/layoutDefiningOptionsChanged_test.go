package editorStateDto_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/stretchr/testify/assert"
)

func TestWhenLayoutDefiningOptionChanges_ReportsChanged(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		subtestName string
		mutate      func(state *editor_state_dto.EditorStateDto)
	}{
		{"WhenPlayerCountChanges_ReportsChanged", func(state *editor_state_dto.EditorStateDto) { state.PlayerCount++ }},
		{
			"WhenTopologyChanges_ReportsChanged",
			func(state *editor_state_dto.EditorStateDto) { state.Topology = config.TopologyChain },
		},
		{
			"WhenGenerateRoadsFlips_ReportsChanged",
			func(state *editor_state_dto.EditorStateDto) { state.GenerateRoads = !state.GenerateRoads },
		},
		{
			"WhenRandomPortalsFlips_ReportsChanged",
			func(state *editor_state_dto.EditorStateDto) { state.RandomPortals = !state.RandomPortals },
		},
		{
			"WhenNoDirectPlayerConnFlips_ReportsChanged",
			func(state *editor_state_dto.EditorStateDto) { state.NoDirectPlayerConn = !state.NoDirectPlayerConn },
		},
		{
			"WhenMaxPortalConnectionsChanges_ReportsChanged",
			func(state *editor_state_dto.EditorStateDto) { state.MaxPortalConnections++ },
		},
		{
			"WhenAdvancedModeFlips_ReportsChanged",
			func(state *editor_state_dto.EditorStateDto) { state.AdvancedMode = !state.AdvancedMode },
		},
		{
			"WhenNeutralZoneCountChanges_ReportsChanged",
			func(state *editor_state_dto.EditorStateDto) { state.NeutralZoneCount++ },
		},
		{
			"WhenNeutralLowNoCastleCountChanges_ReportsChanged",
			func(state *editor_state_dto.EditorStateDto) { state.NeutralLowNoCastleCount++ },
		},
		{
			"WhenNeutralLowCastleCountChanges_ReportsChanged",
			func(state *editor_state_dto.EditorStateDto) { state.NeutralLowCastleCount++ },
		},
		{
			"WhenNeutralMediumNoCastleCountChanges_ReportsChanged",
			func(state *editor_state_dto.EditorStateDto) { state.NeutralMediumNoCastleCount++ },
		},
		{
			"WhenNeutralMediumCastleCountChanges_ReportsChanged",
			func(state *editor_state_dto.EditorStateDto) { state.NeutralMediumCastleCount++ },
		},
		{
			"WhenNeutralHighNoCastleCountChanges_ReportsChanged",
			func(state *editor_state_dto.EditorStateDto) { state.NeutralHighNoCastleCount++ },
		},
		{
			"WhenNeutralHighCastleCountChanges_ReportsChanged",
			func(state *editor_state_dto.EditorStateDto) { state.NeutralHighCastleCount++ },
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.subtestName, func(t *testing.T) {
			t.Parallel()
			// Arrange
			previous := editor_state_dto.NewDefaultEditorStateDto()
			incoming := previous
			testCase.mutate(&incoming)

			// Act
			changed := previous.LayoutDefiningOptionsChanged(&incoming)

			// Assert
			assert.True(t, changed)
		})
	}
}

func TestWhenStatesAreIdentical_ReportsUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	previous := editor_state_dto.NewDefaultEditorStateDto()
	incoming := previous

	// Act
	changed := previous.LayoutDefiningOptionsChanged(&incoming)

	// Assert
	assert.False(t, changed)
}

func TestWhenOnlyNonLayoutOptionsChange_ReportsUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	previous := editor_state_dto.NewDefaultEditorStateDto()
	incoming := previous
	incoming.TemplateName = "Renamed"
	incoming.NeutralZoneCastles = previous.NeutralZoneCastles + 2
	incoming.ResourceDensityPercent = previous.ResourceDensityPercent + 50

	// Act
	changed := previous.LayoutDefiningOptionsChanged(&incoming)

	// Assert
	assert.False(t, changed)
}
