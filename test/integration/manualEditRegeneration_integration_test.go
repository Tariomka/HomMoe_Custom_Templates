//go:build integration_test

package integration_test

import (
	"testing"
	"time"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/drivers"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// movedPosition is where newMovedSession puts the first zone.
func movedPosition() data.Vec2[float64] { return data.NewVec2(0.15, 0.85) }

// settingChangeCase is one §1.16 setting: the layout it is tried on, any
// precondition applied before the manual move, and the change itself.
type settingChangeCase struct {
	topology config.MapTopology
	prepare  func(settings *editor_state_model.EditorState)
	change   func(settings *editor_state_model.EditorState)
}

// Review §1.16: after a manual move, a later settings change must reach every
// value the user did not edit. The regenerated layout must therefore equal a
// fresh generation for the new settings, apart from the moved zone's position.
func TestWhenASettingChangesAfterAManualMove_UntouchedValuesFollowTheSettings(t *testing.T) {
	unprepared := func(*editor_state_model.EditorState) {}
	advanced := func(settings *editor_state_model.EditorState) { settings.AdvancedMode = true }
	settingChanges := map[string]settingChangeCase{
		"MapSizeChanged_FollowsIt": {config.TopologySquare, unprepared,
			func(settings *editor_state_model.EditorState) { settings.MapSize = 256 }},
		"ResourceDensityChanged_FollowsIt": {config.TopologySquare, unprepared,
			func(settings *editor_state_model.EditorState) { settings.ResourceDensityPercent = 180 }},
		"StructureDensityChanged_FollowsIt": {config.TopologySquare, unprepared,
			func(settings *editor_state_model.EditorState) { settings.StructureDensityPercent = 60 }},
		"PlayerZoneSizeChanged_FollowsIt": {config.TopologySquare, unprepared,
			func(settings *editor_state_model.EditorState) { settings.PlayerZoneSize = 1.5 }},
		"NeutralZoneSizeChanged_FollowsIt": {config.TopologySquare, unprepared,
			func(settings *editor_state_model.EditorState) { settings.NeutralZoneSize = 1.4 }},
		"HubZoneSizeChanged_FollowsIt": {config.TopologyGeometricHub, unprepared,
			func(settings *editor_state_model.EditorState) { settings.HubZoneSize = 1.5 }},
		"NeutralStackStrengthChanged_FollowsIt": {config.TopologySquare, unprepared,
			func(settings *editor_state_model.EditorState) { settings.NeutralStackStrengthPercent = 150 }},
		"BorderGuardStrengthChanged_FollowsIt": {config.TopologySquare, unprepared,
			func(settings *editor_state_model.EditorState) { settings.BorderGuardStrengthPercent = 150 }},
		"AbandonedOutpostsEnabled_FollowsIt": {config.TopologySquare, unprepared,
			func(settings *editor_state_model.EditorState) {
				settings.SpawnAbandonedOutposts = true
				settings.AbandonedOutpostCount = 2
			}},
		"AdvancedGuardRandomizationChanged_FollowsIt": {config.TopologySquare, advanced,
			func(settings *editor_state_model.EditorState) { settings.GuardRandomization = 0.2 }},
		"PlayerCastlesChanged_FollowsIt": {config.TopologySquare, unprepared,
			func(settings *editor_state_model.EditorState) { settings.PlayerZoneCastles = 2 }},
		"CityHoldEnabledOnAHub_FollowsIt": {config.TopologyGeometricHub, unprepared,
			func(settings *editor_state_model.EditorState) { settings.CityHold = true }},
	}
	for name, testCase := range settingChanges {
		t.Run(name, func(t *testing.T) {
			// Arrange
			state := newMovedSessionWith(t, testCase.topology, testCase.prepare)
			before := freshLayoutOf(t, state)
			state.UpdateState(testCase.change)
			expected := freshLayoutOf(t, state)
			require.NotEqual(t, layoutOf(before), layoutOf(expected), "the setting does not change the generation")

			// Act
			regenerateAfterDebounce(state, time.Now())

			// Assert
			assert.Equal(t, withTheManualMove(expected), layoutOf(state.GetLastTemplate().Variants[0]))
		})
	}
}

func TestWhenAnOverriddenFieldsSettingChanges_TheOverrideSurvives(t *testing.T) {
	// Arrange
	state := newMovedSession(t, config.TopologySquare)
	zones := append([]template_model.Zone(nil), state.GetLastTemplate().Variants[0].Zones...)
	edited := firstNeutralIndex(t, zones)
	zones[edited].Size = 2.75
	state.ApplyEditedZones(dtos.ZoneEditorZonesDto{
		Zones:       zones,
		Connections: state.GetLastTemplate().Variants[0].Connections,
	})
	state.UpdateState(func(settings *editor_state_model.EditorState) { settings.NeutralZoneSize = 1.4 })

	// Act
	regenerateAfterDebounce(state, time.Now())

	// Assert
	assert.Equal(t, 2.75, state.GetLastTemplate().Variants[0].Zones[edited].Size)
}

func TestWhenOnlyTheTemplateIsRenamedAfterAManualMove_NoZoneChanges(t *testing.T) {
	// Arrange
	state := newMovedSession(t, config.TopologySquare)
	before := append([]template_model.Zone(nil), state.GetLastTemplate().Variants[0].Zones...)
	state.UpdateState(func(settings *editor_state_model.EditorState) { settings.TemplateName = "Renamed" })

	// Act
	regenerateAfterDebounce(state, time.Now())

	// Assert
	assert.Equal(t, before, state.GetLastTemplate().Variants[0].Zones)
}

// newMovedSession generates a layout on the topology and applies a manual edit
// that only moves the first zone.
func newMovedSession(t *testing.T, topology config.MapTopology) *drivers.State {
	t.Helper()

	return newMovedSessionWith(t, topology, func(*editor_state_model.EditorState) {})
}

// newMovedSessionWith is newMovedSession with extra settings applied before
// the first generation.
func newMovedSessionWith(
	t *testing.T,
	topology config.MapTopology,
	prepare func(settings *editor_state_model.EditorState)) *drivers.State {
	t.Helper()
	state := newUIState()
	state.UpdateState(func(settings *editor_state_model.EditorState) {
		settings.Topology = topology
		settings.NeutralZoneCount = 4
		prepare(settings)
	})
	state.AutoRegenerate(time.Now())
	template := state.GetLastTemplate()
	require.NotNil(t, template)
	require.NotEmpty(t, template.Variants)

	zones := append([]template_model.Zone(nil), template.Variants[0].Zones...)
	zones[0].ManualPosition = new(movedPosition())
	state.ApplyEditedZones(dtos.ZoneEditorZonesDto{Zones: zones, Connections: template.Variants[0].Connections})
	stateData := state.GetStateData()
	require.True(t, stateData.HasManualEdits())

	return state
}

func layoutOf(variant template_model.Variant) template_model.Variant {
	return template_model.Variant{Zones: variant.Zones, Connections: variant.Connections}
}

// withTheManualMove is the fresh layout with the first zone moved the way
// newMovedSession moved it.
func withTheManualMove(variant template_model.Variant) template_model.Variant {
	zones := append([]template_model.Zone(nil), variant.Zones...)
	zones[0].ManualPosition = new(movedPosition())

	return template_model.Variant{Zones: zones, Connections: variant.Connections}
}

func firstNeutralIndex(t *testing.T, zones []template_model.Zone) int {
	t.Helper()
	for index, zone := range zones {
		if len(zone.Name) > len("Neutral-") && zone.Name[:len("Neutral-")] == "Neutral-" {
			return index
		}
	}
	t.Fatal("the layout has no neutral zone")

	return -1
}
