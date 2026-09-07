//go:build integration_test && gui

package gui_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers/integration_common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The Geometric Hub layout for the default two-player template is what these
// tests edit: a neutral Hub the generator lets you delete, two player spawns it
// does not, and one named portal per spawn. Everything below is that layout, so
// the names are constants rather than lookups.
const (
	hubZoneName        = "Hub"
	hubToSpawnAName    = "Portal-Hub-A"
	hubToSpawnBName    = "Portal-Hub-B"
	geometricHubLayout = "Geometric Hub"
)

// openZoneEditor drives the window to the zone editor over a deterministic
// layout. A fixture directory is seeded even though no file is opened: the
// toolbar reports the current file's path, which is per-machine, and the mask
// that hides it comes with the fixture directory.
func openZoneEditor(t *testing.T, topology string, withSnapshots bool) (
	*integration_common.AppRunner, *integration_common.ZoneEditorHandler) {
	t.Helper()
	runner := integration_common.NewAppRunner(t)
	handler := integration_common.NewHandler(runner).WithFixtureDirectory()
	if withSnapshots {
		handler = handler.WithSnapshots()
	}

	return runner, handler.
		ClickLayoutAndZonesTab().
		SelectTopology(topology).
		OpenZoneEditor()
}

// editedZoneNames reports the zones the editor currently holds.
func editedZoneNames(zoneEditor *integration_common.ZoneEditorHandler) []string {
	names := make([]string, 0)
	for _, zone := range zoneEditor.Dialog().EditedZones() {
		names = append(names, zone.Name)
	}

	return names
}

// committedConnectionNames reports the connections the editor state records as
// manually edited, which is what Apply writes.
func committedConnectionNames(runner *integration_common.AppRunner) []string {
	names := make([]string, 0)
	for _, save := range runner.CurrentState().ManualConnections {
		names = append(names, save.Connection.Name)
	}

	return names
}

// The button restores the layout the dialog was opened with, which is not the
// generated one once manual edits have been applied before - so it must not
// claim to reset to generated.
//
//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenTheZoneEditorRenders_TheUndoButtonIsLabelledForWhatItDoes(t *testing.T) {
	// Arrange
	_, zoneEditor := openZoneEditor(t, geometricHubLayout, false)

	// Act
	labels := zoneEditor.ButtonLabels()

	// Assert
	assert.Contains(t, labels, "Undo")
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenTheZoneEditorRenders_TheRevertToBaseButtonIsOffered(t *testing.T) {
	// Arrange
	_, zoneEditor := openZoneEditor(t, geometricHubLayout, false)

	// Act
	labels := zoneEditor.ButtonLabels()

	// Assert
	assert.Contains(t, labels, "Revert to Base")
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenTheZoneEditorRenders_NoButtonClaimsToResetToGenerated(t *testing.T) {
	// Arrange
	_, zoneEditor := openZoneEditor(t, geometricHubLayout, false)

	// Act
	labels := zoneEditor.ButtonLabels()

	// Assert
	assert.NotContains(t, labels, "Reset to generated")
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenAZoneIsSelected_TheEditorRendersItsPropertyPanel(t *testing.T) {
	// Arrange
	_, zoneEditor := openZoneEditor(t, geometricHubLayout, true)

	// Act
	zoneEditor.ClickZone(hubZoneName)

	// Assert
	assert.Equal(t,
		[]string{hubZoneName, ""},
		[]string{zoneEditor.Dialog().SelectedZone(), zoneEditor.Dialog().SelectedConnection()})
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenAConnectionIsSelected_TheEditorRendersItsPropertyPanel(t *testing.T) {
	// Arrange
	_, zoneEditor := openZoneEditor(t, geometricHubLayout, true)

	// Act
	zoneEditor.ClickConnection(hubToSpawnAName)

	// Assert
	assert.Equal(t,
		[]string{"", hubToSpawnAName},
		[]string{zoneEditor.Dialog().SelectedZone(), zoneEditor.Dialog().SelectedConnection()})
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenTheSelectedConnectionIsDeleted_ItLeavesTheWorkingSet(t *testing.T) {
	// Arrange
	_, zoneEditor := openZoneEditor(t, geometricHubLayout, true)
	zoneEditor.ClickConnection(hubToSpawnAName)

	// Act
	zoneEditor.ClickDeleteSelected()

	// Assert
	assert.Equal(t, []string{hubToSpawnBName}, zoneEditor.Dialog().EditedConnectionNames())
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenTheSelectedZoneIsDeleted_ItsConnectionsGoWithIt(t *testing.T) {
	// Arrange
	_, zoneEditor := openZoneEditor(t, geometricHubLayout, true)
	zoneEditor.ClickZone(hubZoneName)

	// Act
	zoneEditor.ClickDeleteSelected()

	// Assert
	assert.Empty(t, zoneEditor.Dialog().EditedConnectionNames())
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenTheSessionEditsAreUndone_TheDeletedConnectionComesBack(t *testing.T) {
	// Arrange
	_, zoneEditor := openZoneEditor(t, geometricHubLayout, true)
	zoneEditor.ClickConnection(hubToSpawnAName).ClickDeleteSelected()

	// Act
	zoneEditor.ClickUndo()

	// Assert
	assert.Equal(t,
		[]string{hubToSpawnAName, hubToSpawnBName},
		zoneEditor.Dialog().EditedConnectionNames())
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenTheSessionEditsAreUndone_TheSelectionIsCleared(t *testing.T) {
	// Arrange
	_, zoneEditor := openZoneEditor(t, geometricHubLayout, true)
	zoneEditor.ClickZone(hubZoneName)

	// Act
	zoneEditor.ClickUndo()

	// Assert
	assert.Empty(t, zoneEditor.Dialog().SelectedZone())
}

// Undo is session-scoped: it must not tell the driver to drop anything, so
// Apply after an Undo simply commits the layout the editor started from.
//
//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenApplyFollowsAnUndo_TheStartingConnectionsAreCommitted(t *testing.T) {
	// Arrange
	runner, zoneEditor := openZoneEditor(t, geometricHubLayout, true)
	zoneEditor.ClickConnection(hubToSpawnAName).ClickDeleteSelected().ClickUndo()

	// Act
	zoneEditor.ClickApply()

	// Assert
	assert.Equal(t, []string{hubToSpawnAName, hubToSpawnBName}, committedConnectionNames(runner))
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenRevertToBaseIsPressed_TheEditorShowsTheRegeneratedZones(t *testing.T) {
	// Arrange
	_, zoneEditor := openZoneEditor(t, geometricHubLayout, true)
	zoneEditor.ClickZone(hubZoneName).ClickDeleteSelected()
	require.NotContains(t, editedZoneNames(zoneEditor), hubZoneName,
		"the delete must have taken for the revert to have something to restore")

	// Act
	zoneEditor.ClickRevertToBase()

	// Assert
	assert.Contains(t, editedZoneNames(zoneEditor), hubZoneName)
}

// After a revert the base becomes the new baseline, so Undo returns to it
// rather than to the edits that were discarded.
//
//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenUndoFollowsARevertToBase_TheBaseZonesComeBack(t *testing.T) {
	// Arrange
	_, zoneEditor := openZoneEditor(t, geometricHubLayout, true)
	zoneEditor.ClickRevertToBase().ClickZone(hubZoneName).ClickDeleteSelected()

	// Act
	zoneEditor.ClickUndo()

	// Assert
	assert.Contains(t, editedZoneNames(zoneEditor), hubZoneName)
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenAddConnectionIsClicked_TheEditorEntersAddConnectionMode(t *testing.T) {
	// Arrange
	_, zoneEditor := openZoneEditor(t, geometricHubLayout, true)

	// Act
	zoneEditor.ClickAddConnection()

	// Assert
	assert.True(t, zoneEditor.Dialog().AddConnectionModeActive())
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenAddConnectionIsClickedTwice_TheEditorLeavesAddConnectionMode(t *testing.T) {
	// Arrange
	_, zoneEditor := openZoneEditor(t, geometricHubLayout, true)
	zoneEditor.ClickAddConnection()

	// Act
	zoneEditor.ClickAddConnection()

	// Assert
	assert.False(t, zoneEditor.Dialog().AddConnectionModeActive())
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenAddZoneIsClickedWhileAddingAConnection_TheAddConnectionModeTurnsOff(t *testing.T) {
	// Arrange
	_, zoneEditor := openZoneEditor(t, geometricHubLayout, true)
	zoneEditor.ClickAddConnection()

	// Act
	zoneEditor.ClickAddZone()

	// Assert
	assert.False(t, zoneEditor.Dialog().AddConnectionModeActive())
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenAddZoneIsClickedWhileAddingAConnection_TheAddZoneModeTurnsOn(t *testing.T) {
	// Arrange
	_, zoneEditor := openZoneEditor(t, geometricHubLayout, true)
	zoneEditor.ClickAddConnection()

	// Act
	zoneEditor.ClickAddZone()

	// Assert
	assert.True(t, zoneEditor.Dialog().AddZoneModeActive())
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenTheSessionEditsAreUndone_TheAddModeTurnsOff(t *testing.T) {
	// Arrange
	_, zoneEditor := openZoneEditor(t, geometricHubLayout, true)
	zoneEditor.ClickAddConnection()

	// Act
	zoneEditor.ClickUndo()

	// Assert
	assert.False(t, zoneEditor.Dialog().AddConnectionModeActive())
}
