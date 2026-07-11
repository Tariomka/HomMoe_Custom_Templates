package editorStateDto_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/stretchr/testify/assert"
)

func TestWhenLayoutDefiningOptionChanges_ReportsChanged(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		subtestName string
		mutate      func(state *dtos.EditorStateDto)
	}{
		{"WhenPlayerCountChanges_ReportsChanged", func(state *dtos.EditorStateDto) { state.PlayerCount++ }},
		{
			"WhenTopologyChanges_ReportsChanged",
			func(state *dtos.EditorStateDto) { state.Topology = config.TopologyChain },
		},
		{
			"WhenGenerateRoadsFlips_ReportsChanged",
			func(state *dtos.EditorStateDto) { state.GenerateRoads = !state.GenerateRoads },
		},
		{
			"WhenRandomPortalsFlips_ReportsChanged",
			func(state *dtos.EditorStateDto) { state.RandomPortals = !state.RandomPortals },
		},
		{
			"WhenNoDirectPlayerConnFlips_ReportsChanged",
			func(state *dtos.EditorStateDto) { state.NoDirectPlayerConn = !state.NoDirectPlayerConn },
		},
		{
			"WhenMaxPortalConnectionsChanges_ReportsChanged",
			func(state *dtos.EditorStateDto) { state.MaxPortalConnections++ },
		},
		{
			"WhenMinNeutralZonesBetweenPlayersChanges_ReportsChanged",
			func(state *dtos.EditorStateDto) { state.MinNeutralZonesBetweenPlayers++ },
		},
		{
			"WhenAdvancedModeFlips_ReportsChanged",
			func(state *dtos.EditorStateDto) { state.AdvancedMode = !state.AdvancedMode },
		},
		{"WhenNeutralZoneCountChanges_ReportsChanged", func(state *dtos.EditorStateDto) { state.NeutralZoneCount++ }},
		{
			"WhenNeutralLowNoCastleCountChanges_ReportsChanged",
			func(state *dtos.EditorStateDto) { state.NeutralLowNoCastleCount++ },
		},
		{
			"WhenNeutralLowCastleCountChanges_ReportsChanged",
			func(state *dtos.EditorStateDto) { state.NeutralLowCastleCount++ },
		},
		{
			"WhenNeutralMediumNoCastleCountChanges_ReportsChanged",
			func(state *dtos.EditorStateDto) { state.NeutralMediumNoCastleCount++ },
		},
		{
			"WhenNeutralMediumCastleCountChanges_ReportsChanged",
			func(state *dtos.EditorStateDto) { state.NeutralMediumCastleCount++ },
		},
		{
			"WhenNeutralHighNoCastleCountChanges_ReportsChanged",
			func(state *dtos.EditorStateDto) { state.NeutralHighNoCastleCount++ },
		},
		{
			"WhenNeutralHighCastleCountChanges_ReportsChanged",
			func(state *dtos.EditorStateDto) { state.NeutralHighCastleCount++ },
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.subtestName, func(t *testing.T) {
			t.Parallel()
			// Arrange
			previous := dtos.NewDefaultEditorStateDto()
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
	previous := dtos.NewDefaultEditorStateDto()
	incoming := previous

	// Act
	changed := previous.LayoutDefiningOptionsChanged(&incoming)

	// Assert
	assert.False(t, changed)
}

func TestWhenOnlyNonLayoutOptionsChange_ReportsUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	previous := dtos.NewDefaultEditorStateDto()
	incoming := previous
	incoming.TemplateName = "Renamed"
	incoming.NeutralZoneCastles = previous.NeutralZoneCastles + 2
	incoming.ResourceDensityPercent = previous.ResourceDensityPercent + 50

	// Act
	changed := previous.LayoutDefiningOptionsChanged(&incoming)

	// Assert
	assert.False(t, changed)
}
