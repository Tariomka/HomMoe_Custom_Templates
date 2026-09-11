//go:build integration_test && gui

package gui_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template_entity"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers/integration_common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The between-zone road setting decides what a pending edit is worth before
// Apply: a connection the user draws, and one whose type they change, has to
// already show the road state the export will carry. The editor never decides
// that itself - it asks the handler - so these tests drive the real window and
// read the pending record back out of the open dialog.

// portalConnectionTypeLabel is what the side panel's type dropdown calls a
// portal, spelled out rather than read back from the production table.
const portalConnectionTypeLabel = "Portal"

// openZoneEditorWithRoads drives the window to the zone editor over the
// deterministic Geometric Hub layout, with the between-zone road setting put
// where the test needs it first. Snapshots stay off: these tests assert what a
// connection carries, not what the frame looks like.
func openZoneEditorWithRoads(t *testing.T, generateRoads bool) (
	*integration_common.AppRunner, *integration_common.ZoneEditorHandler) {
	t.Helper()
	runner := integration_common.NewAppRunner(t)
	tab := integration_common.NewHandler(runner).
		WithFixtureDirectory().
		ClickLayoutAndZonesTab().
		SelectTopology(geometricHubLayout)
	if !generateRoads {
		tab = tab.ToggleGenerateRoads()
	}
	require.Equal(t, generateRoads, runner.CurrentState().GenerateRoads,
		"the road setting was not put where the test needs it")

	return runner, tab.OpenZoneEditor()
}

// addConnectionBetweenSpawns arms the add-connection mode and draws an edge
// between the two player spawns, which is a connection the generator never
// produces for this layout and so is unambiguously the pending one. The mode
// stays armed after a drag so several edges can be chained, so it is disarmed
// again here - while it runs, a canvas click cannot select anything.
func addConnectionBetweenSpawns(
	zoneEditor *integration_common.ZoneEditorHandler) *integration_common.ZoneEditorHandler {
	return zoneEditor.ClickAddConnection().
		DragFromZoneTo(spawnAZoneName, spawnBZoneName).
		ClickAddConnection()
}

// pendingConnection reads a connection as the open editor currently holds it.
// A user-added connection carries no name, so it is found by its endpoints.
func pendingConnection(
	t *testing.T,
	zoneEditor *integration_common.ZoneEditorHandler,
	from string,
	to string) template_model.Connection {
	t.Helper()
	for _, connection := range zoneEditor.Dialog().EditedConnectionRecords() {
		if (connection.From == from && connection.To == to) ||
			(connection.From == to && connection.To == from) {
			return connection
		}
	}
	t.Fatalf("the editor is not holding a connection between %q and %q", from, to)

	return template_model.Connection{}
}

// makeHubToSpawnARoadless selects the generated Hub-to-Spawn-A portal and turns
// it into a Direct connection, which with roads off is what leaves it roadless.
// A generated connection is the one addressed here because a user-added one
// carries an extra note row that moves the side panel's dropdowns down.
func makeHubToSpawnARoadless(
	t *testing.T,
	zoneEditor *integration_common.ZoneEditorHandler) *integration_common.ZoneEditorHandler {
	t.Helper()
	zoneEditor.ClickConnection(hubToSpawnAName).SelectConnectionType(directConnectionTypeLabel)
	require.Equal(t, new(false),
		pendingConnection(t, zoneEditor, hubZoneName, spawnAZoneName).Road,
		"the connection was not roadless before it was handed to the portal rules")

	return zoneEditor
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenRoadsAreOnAndAConnectionIsDrawn_ThePendingConnectionCarriesARoad(t *testing.T) {
	// Arrange
	_, zoneEditor := openZoneEditorWithRoads(t, true)

	// Act
	addConnectionBetweenSpawns(zoneEditor)

	// Assert
	assert.Equal(t, new(true), pendingConnection(t, zoneEditor, spawnAZoneName, spawnBZoneName).Road)
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenRoadsAreOffAndAConnectionIsDrawn_ThePendingConnectionCarriesNoRoad(t *testing.T) {
	// Arrange
	_, zoneEditor := openZoneEditorWithRoads(t, false)

	// Act
	addConnectionBetweenSpawns(zoneEditor)

	// Assert
	assert.Equal(t, new(false), pendingConnection(t, zoneEditor, spawnAZoneName, spawnBZoneName).Road)
}

// Switching a connection to Portal hands it to the portal rules, which keep
// whatever road flag it already carries.
//
//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenARoadlessConnectionBecomesAPortal_ItKeepsItsRoadlessFlag(t *testing.T) {
	// Arrange
	_, zoneEditor := openZoneEditorWithRoads(t, false)
	makeHubToSpawnARoadless(t, zoneEditor)

	// Act
	zoneEditor.SelectConnectionType(portalConnectionTypeLabel)

	// Assert
	assert.Equal(t, new(false), pendingConnection(t, zoneEditor, hubZoneName, spawnAZoneName).Road)
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenAConnectionBecomesAPortal_ThePendingConnectionRecordsThatType(t *testing.T) {
	// Arrange
	_, zoneEditor := openZoneEditorWithRoads(t, false)
	makeHubToSpawnARoadless(t, zoneEditor)

	// Act
	zoneEditor.SelectConnectionType(portalConnectionTypeLabel)

	// Assert
	assert.Equal(t,
		portalConnectionTypeLabel,
		pendingConnection(t, zoneEditor, hubZoneName, spawnAZoneName).ConnectionType)
}

// Leaving Portal hands the connection back to the setting, which with roads off
// means it loses the road the portal was carrying.
//
//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenAPortalBecomesDirectWithRoadsOff_ItTakesTheRoadlessFlag(t *testing.T) {
	// Arrange
	_, zoneEditor := openZoneEditorWithRoads(t, false)
	zoneEditor.ClickConnection(hubToSpawnAName)

	// Act
	zoneEditor.SelectConnectionType(directConnectionTypeLabel)

	// Assert
	assert.Equal(t, new(false), pendingConnection(t, zoneEditor, hubZoneName, spawnAZoneName).Road)
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenAPortalBecomesDirectWithRoadsOn_ItTakesTheRoadedFlag(t *testing.T) {
	// Arrange
	_, zoneEditor := openZoneEditorWithRoads(t, true)
	zoneEditor.ClickConnection(hubToSpawnAName)

	// Act
	zoneEditor.SelectConnectionType(directConnectionTypeLabel)

	// Assert
	assert.Equal(t, new(true), pendingConnection(t, zoneEditor, hubZoneName, spawnAZoneName).Road)
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenARoadlessPendingConnectionIsApplied_TheCommittedRecordCarriesTheFlag(t *testing.T) {
	// Arrange
	runner, zoneEditor := openZoneEditorWithRoads(t, false)
	addConnectionBetweenSpawns(zoneEditor)

	// Act
	zoneEditor.ClickApply()

	// Assert
	assert.Equal(t, new(false), committedConnection(t, runner, spawnAZoneName, spawnBZoneName).Road)
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenARoadedPendingConnectionIsApplied_TheCommittedRecordCarriesTheFlag(t *testing.T) {
	// Arrange
	runner, zoneEditor := openZoneEditorWithRoads(t, true)
	addConnectionBetweenSpawns(zoneEditor)

	// Act
	zoneEditor.ClickApply()

	// Assert
	assert.Equal(t, new(true), committedConnection(t, runner, spawnAZoneName, spawnBZoneName).Road)
}

// Cancel throws the working copy away, so a pending connection - and the road
// flag stamped on it - must never reach the editor state.
//
//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenAPendingConnectionIsCancelled_NothingIsCommitted(t *testing.T) {
	// Arrange
	runner, zoneEditor := openZoneEditorWithRoads(t, false)
	addConnectionBetweenSpawns(zoneEditor)
	require.Equal(t, new(false),
		pendingConnection(t, zoneEditor, spawnAZoneName, spawnBZoneName).Road,
		"the pending connection was not stamped before Cancel")

	// Act
	zoneEditor.ClickCancel()

	// Assert
	assert.Empty(t, runner.CurrentState().ManualConnections)
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenAPendingTypeChangeIsCancelled_NothingIsCommitted(t *testing.T) {
	// Arrange
	runner, zoneEditor := openZoneEditorWithRoads(t, false)
	zoneEditor.ClickConnection(hubToSpawnAName)
	zoneEditor.SelectConnectionType(directConnectionTypeLabel)
	require.Equal(t, directConnectionTypeLabel,
		pendingConnection(t, zoneEditor, hubZoneName, spawnAZoneName).ConnectionType,
		"the pending type change did not happen before Cancel")

	// Act
	zoneEditor.ClickCancel()

	// Assert
	assert.Empty(t, runner.CurrentState().ManualConnections)
}

// An empty ManualConnections list only says Cancel committed nothing new. This
// is the case that matters once the user already has committed work: a second
// session's pending edits must leave the retained types and road flags alone.
//
//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenAPendingChangeIsCancelled_TheRetainedConnectionsAreUntouched(t *testing.T) {
	// Arrange
	runner, zoneEditor := openZoneEditorWithRoads(t, false)
	makeHubToSpawnARoadless(t, zoneEditor)
	editor := zoneEditor.ClickApply()
	retained := committedConnectionSummaries(runner)
	require.NotEmpty(t, retained, "nothing was committed for Cancel to leave alone")

	reopened := editor.ClickLayoutAndZonesTab().OpenZoneEditor()
	reopened.ClickConnection(hubToSpawnAName).SelectConnectionType(portalConnectionTypeLabel)
	addConnectionBetweenSpawns(reopened)
	require.Equal(t, portalConnectionTypeLabel,
		pendingConnection(t, reopened, hubZoneName, spawnAZoneName).ConnectionType,
		"the second session made no pending change for Cancel to throw away")

	// Act
	reopened.ClickCancel()

	// Assert
	assert.Equal(t, retained, committedConnectionSummaries(runner))
}

// committedConnectionSummary is what Cancel must not disturb about a connection
// the editor state already retains, copied out so the before/after comparison
// cannot end up looking at the same values twice.
type committedConnectionSummary struct {
	From string
	To   string
	Type string
	Road *bool
}

func committedConnectionSummaries(
	runner *integration_common.AppRunner) []committedConnectionSummary {
	summaries := make([]committedConnectionSummary, 0)
	for _, save := range runner.CurrentState().ManualConnections {
		connection := save.Connection
		summary := committedConnectionSummary{
			From: connection.From,
			To:   connection.To,
			Type: connection.ConnectionType,
		}
		if connection.Road != nil {
			summary.Road = new(*connection.Road)
		}
		summaries = append(summaries, summary)
	}

	return summaries
}

// committedConnection reads a connection out of what Apply wrote to the editor
// state, found by its endpoints because a user-added connection has no name.
func committedConnection(
	t *testing.T,
	runner *integration_common.AppRunner,
	from string,
	to string) template_entity.Connection {
	t.Helper()
	for _, save := range runner.CurrentState().ManualConnections {
		connection := save.Connection
		if (connection.From == from && connection.To == to) ||
			(connection.From == to && connection.To == from) {
			return connection
		}
	}
	t.Fatalf("the editor state committed no connection between %q and %q", from, to)

	return template_entity.Connection{}
}
