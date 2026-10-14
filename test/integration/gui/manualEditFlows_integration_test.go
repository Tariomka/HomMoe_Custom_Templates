//go:build integration_test && gui

package gui_test

import (
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/composition"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers/integration_common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// savedManualEdits saves the window's state and reads the manualEdits section of
// the written file back as generic entries, which is what a later session would
// load.
func savedManualEdits(
	t *testing.T, runner *integration_common.AppRunner) (zones []map[string]any, connections []map[string]any) {
	t.Helper()
	runner.SaveStateToFile(filepath.Join(t.TempDir(), "saved.gen.json"))
	status, isError := runner.Status()
	require.False(t, isError, status)
	section, ok := readStateDocument(t, runner.CurrentPath())["manualEdits"].(map[string]any)
	require.True(t, ok, "the saved file holds no manualEdits section")

	return asEntries(section["zones"]), asEntries(section["connections"])
}

func asEntries(value any) []map[string]any {
	items, _ := value.([]any)
	entries := make([]map[string]any, 0, len(items))
	for _, item := range items {
		entry, _ := item.(map[string]any)
		entries = append(entries, entry)
	}

	return entries
}

func entryNames(entries []map[string]any) []string {
	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		name, _ := entry["name"].(string)
		names = append(names, name)
	}

	return names
}

func zoneEntryNamed(t *testing.T, entries []map[string]any, name string) map[string]any {
	t.Helper()
	for _, entry := range entries {
		if entry["name"] == name {
			return entry
		}
	}
	t.Fatalf("the saved file records no zone called %q; it records %v", name, entryNames(entries))

	return nil
}

func connectionEntryBetween(t *testing.T, entries []map[string]any, from string, to string) map[string]any {
	t.Helper()
	for _, entry := range entries {
		if (entry["from"] == from && entry["to"] == to) || (entry["from"] == to && entry["to"] == from) {
			return entry
		}
	}
	t.Fatalf("the saved file records no connection between %q and %q", from, to)

	return nil
}

// presentKeys reports which of the given keys the entry carries, in the order
// asked for.
func presentKeys(entry map[string]any, keys ...string) []string {
	present := make([]string, 0, len(keys))
	for _, key := range keys {
		if _, found := entry[key]; found {
			present = append(present, key)
		}
	}

	return present
}

// fieldOverride reads one overridden field of a saved entry, or nil when the
// entry overrides nothing.
func fieldOverride(entry map[string]any, field string) any {
	overrides, _ := entry["fieldOverrides"].(map[string]any)
	return overrides[field]
}

// guardValueOf reads the guard of the named connection.
func guardValueOf(t *testing.T, connections []template_model.Connection, name string) int {
	t.Helper()
	for _, connection := range connections {
		if connection.Name == name {
			return connection.GuardValue
		}
	}
	t.Fatalf("there is no connection called %q", name)

	return 0
}

// reopenZoneEditor opens the editor again on the tab the previous session left.
func reopenZoneEditor(base *integration_common.BaseHandler) *integration_common.ZoneEditorHandler {
	return base.ClickLayoutAndZonesTab().OpenZoneEditor()
}

// openEditorOverAppliedUserConnection draws a connection between the two spawns,
// applies it and reopens the editor over the result.
func openEditorOverAppliedUserConnection(t *testing.T) (
	*integration_common.AppRunner, *integration_common.ZoneEditorHandler) {
	t.Helper()
	runner, zoneEditor := openZoneEditor(t, geometricHubLayout, false)
	base := zoneEditor.
		ClickAddConnection().
		DragFromZoneTo(spawnAZoneName, spawnBZoneName).
		ClickAddConnection().
		ClickApply()

	return runner, reopenZoneEditor(base)
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenAUserConnectionIsAppliedAndTheEditorReopened_TheConnectionIsUserAdded(t *testing.T) {
	// Arrange
	_, zoneEditor := openEditorOverAppliedUserConnection(t)

	// Act
	zoneEditor.ClickConnectionBetween(spawnAZoneName, spawnBZoneName)

	// Assert
	assert.True(t, zoneEditor.Dialog().SelectedConnectionIsUserAdded())
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenTheSessionEditsAreUndoneAndApplied_TheReopenedConnectionIsStillUserAdded(t *testing.T) {
	// Arrange
	_, zoneEditor := openEditorOverAppliedUserConnection(t)
	zoneEditor.ClickUndo()
	reopened := reopenZoneEditor(zoneEditor.ClickApply())

	// Act
	reopened.ClickConnectionBetween(spawnAZoneName, spawnBZoneName)

	// Assert
	assert.True(t, reopened.Dialog().SelectedConnectionIsUserAdded())
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenTheSessionEditsAreUndoneAndApplied_TheSavedConnectionIsUserAdded(t *testing.T) {
	// Arrange
	runner, zoneEditor := openEditorOverAppliedUserConnection(t)
	zoneEditor.ClickUndo().ClickApply()

	// Act
	_, connections := savedManualEdits(t, runner)

	// Assert
	assert.Equal(t, true, connectionEntryBetween(t, connections, spawnAZoneName, spawnBZoneName)["isUserAdded"])
}

// A moved zone is stored as a position on top of what the generator makes, not
// as a rewritten copy of the zone.
//
//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenAZoneIsMovedAndApplied_TheSavedZoneHoldsOnlyItsManualPosition(t *testing.T) {
	// Arrange
	runner, zoneEditor := openZoneEditor(t, geometricHubLayout, false)
	zoneEditor.DragZone(hubZoneName, draggedZoneSpot).ClickApply()

	// Act
	zones, _ := savedManualEdits(t, runner)

	// Assert
	assert.Equal(t,
		[]string{"manualPosition"},
		presentKeys(zoneEntryNamed(t, zones, hubZoneName), "manualPosition", "fieldOverrides"))
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenAZoneIsAddedAndApplied_TheSavedZoneIsMarkedAdded(t *testing.T) {
	// Arrange
	runner, zoneEditor := openZoneEditor(t, geometricHubLayout, false)
	zoneEditor.ClickAddZone().ClickCanvasAt(emptyCanvasSpot).ClickApply()

	// Act
	zones, _ := savedManualEdits(t, runner)

	// Assert
	assert.Equal(t, "added", zoneEntryNamed(t, zones, placedZoneName)["origin"])
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenAZoneIsDeletedAndApplied_TheSavedEditsDoNotListIt(t *testing.T) {
	// Arrange
	runner, zoneEditor := openZoneEditor(t, geometricHubLayout, false)
	zoneEditor.ClickZone(hubZoneName).ClickDeleteSelected().ClickApply()

	// Act
	zones, _ := savedManualEdits(t, runner)

	// Assert
	assert.NotContains(t, entryNames(zones), hubZoneName)
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenAZoneQualityIsChangedAndApplied_TheSavedZoneRecordsItsCastles(t *testing.T) {
	// Arrange
	runner, zoneEditor := openZoneEditor(t, geometricHubLayout, false)
	selectPlacedNeutralZone(zoneEditor)
	zoneEditor.SelectZoneQuality(goldZoneQualityLabel).ClickApply()

	// Act
	zones, _ := savedManualEdits(t, runner)

	// Assert
	assert.Contains(t, zoneEntryNamed(t, zones, placedZoneName), "castles")
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenAConnectionGuardIsChangedAndApplied_TheSavedConnectionOverridesTheGuardValue(t *testing.T) {
	// Arrange
	runner, zoneEditor := openZoneEditor(t, geometricHubLayout, false)
	zoneEditor.ClickConnection(hubToSpawnAName).TypeConnectionGuardValue("1").ClickApply()

	// Act
	_, connections := savedManualEdits(t, runner)

	// Assert
	assert.Equal(t,
		json.Number("135000"),
		fieldOverride(connectionEntryBetween(t, connections, hubZoneName, spawnAZoneName), "guardValue"))
}

// A setting the user never overrode keeps following the settings after a manual
// edit: the untouched connection must carry what a fresh generation gives it.
//
//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenBorderGuardStrengthChangesAfterAMove_AnUneditedConnectionFollowsAFreshGeneration(t *testing.T) {
	// Arrange
	runner, zoneEditor := openZoneEditor(t, geometricHubLayout, false)
	layoutTab := zoneEditor.DragZone(hubZoneName, draggedZoneSpot).ClickApply().ClickLayoutAndZonesTab()

	// Act
	layoutTab.DragBorderGuardStrengthToMaximum().WaitForRegeneration()

	// Assert
	state := runner.CurrentState()
	require.Greater(t, state.BorderGuardStrengthPercent, 100,
		"the slider must have moved for the comparison to mean anything")
	fresh, err := composition.InitializeGuiHandler().GenerateTemplate(editor_state_dto.EditorStateDto{
		EditorState: state.CloneWithoutManualEdits(),
	})
	require.NoError(t, err)
	assert.Equal(t,
		guardValueOf(t, fresh.Template.Variants[0].Connections, hubToSpawnBName),
		guardValueOf(t, runner.AppliedConnections(), hubToSpawnBName))
}
