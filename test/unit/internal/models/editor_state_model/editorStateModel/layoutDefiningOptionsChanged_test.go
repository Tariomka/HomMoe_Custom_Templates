package editorStateModel_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/stretchr/testify/assert"
)

func TestWhenLayoutDefiningOptionChanges_ReportsChanged(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		subtestName string
		mutate      func(state *editor_state_model.EditorStateModel)
	}{
		{
			"WhenPlayerCountChanges_ReportsChanged",
			func(state *editor_state_model.EditorStateModel) { state.PlayerCount++ },
		},
		{
			"WhenTopologyChanges_ReportsChanged",
			func(state *editor_state_model.EditorStateModel) { state.Topology = config.TopologyChain },
		},
		{
			"WhenGenerateRoadsFlips_ReportsChanged",
			func(state *editor_state_model.EditorStateModel) { state.GenerateRoads = !state.GenerateRoads },
		},
		{
			"WhenRandomPortalsFlips_ReportsChanged",
			func(state *editor_state_model.EditorStateModel) { state.RandomPortals = !state.RandomPortals },
		},
		{
			"WhenNoDirectPlayerConnFlips_ReportsChanged",
			func(state *editor_state_model.EditorStateModel) { state.NoDirectPlayerConn = !state.NoDirectPlayerConn },
		},
		{
			"WhenMaxPortalConnectionsChanges_ReportsChanged",
			func(state *editor_state_model.EditorStateModel) { state.MaxPortalConnections++ },
		},
		{
			"WhenAdvancedModeFlips_ReportsChanged",
			func(state *editor_state_model.EditorStateModel) { state.AdvancedMode = !state.AdvancedMode },
		},
		{
			"WhenNeutralZoneCountChanges_ReportsChanged",
			func(state *editor_state_model.EditorStateModel) { state.NeutralZoneCount++ },
		},
		{
			"WhenNeutralLowNoCastleCountChanges_ReportsChanged",
			func(state *editor_state_model.EditorStateModel) { state.NeutralLowNoCastleCount++ },
		},
		{
			"WhenNeutralLowCastleCountChanges_ReportsChanged",
			func(state *editor_state_model.EditorStateModel) { state.NeutralLowCastleCount++ },
		},
		{
			"WhenNeutralMediumNoCastleCountChanges_ReportsChanged",
			func(state *editor_state_model.EditorStateModel) { state.NeutralMediumNoCastleCount++ },
		},
		{
			"WhenNeutralMediumCastleCountChanges_ReportsChanged",
			func(state *editor_state_model.EditorStateModel) { state.NeutralMediumCastleCount++ },
		},
		{
			"WhenNeutralHighNoCastleCountChanges_ReportsChanged",
			func(state *editor_state_model.EditorStateModel) { state.NeutralHighNoCastleCount++ },
		},
		{
			"WhenNeutralHighCastleCountChanges_ReportsChanged",
			func(state *editor_state_model.EditorStateModel) { state.NeutralHighCastleCount++ },
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.subtestName, func(t *testing.T) {
			t.Parallel()
			// Arrange
			previous := editor_state_model.NewDefaultEditorStateModel()
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
	previous := editor_state_model.NewDefaultEditorStateModel()
	incoming := previous

	// Act
	changed := previous.LayoutDefiningOptionsChanged(&incoming)

	// Assert
	assert.False(t, changed)
}

func TestWhenOnlyNonLayoutOptionsChange_ReportsUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	previous := editor_state_model.NewDefaultEditorStateModel()
	incoming := previous
	incoming.TemplateName = "Renamed"
	incoming.NeutralZoneCastles = previous.NeutralZoneCastles + 2
	incoming.ResourceDensityPercent = previous.ResourceDensityPercent + 50

	// Act
	changed := previous.LayoutDefiningOptionsChanged(&incoming)

	// Assert
	assert.False(t, changed)
}
