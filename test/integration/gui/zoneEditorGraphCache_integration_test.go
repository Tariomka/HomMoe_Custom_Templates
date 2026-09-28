//go:build integration_test && gui

package gui_test

import (
	"image"
	"testing"
	"time"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/dialogs"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/themes"
	"github.com/Tariomka/hommoe_custom_templates/internal/composition"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers/integration_common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The zone editor caches its status-line graph diagnostics, so these tests
// count how often the window asks the handler for them. Every helper already
// ends on a frame; one more frame after it is where a stale cache would
// rebuild, so a delta taken across both proves the rebuild happened once and
// only once.

// redrawSettleTimeout bounds how long a test waits for animations started by
// earlier clicks (button ink runs on wall-clock time) to stop asking for frames.
const redrawSettleTimeout = 3 * time.Second

// openCountedZoneEditor opens the zone editor over the Geometric Hub layout in
// a window whose GUI handler counts graph descriptions, and lays out one more
// frame so the count has settled.
func openCountedZoneEditor(t *testing.T) (
	*integration_common.AppRunner,
	*integration_common.GraphDescriptionCounter,
	*integration_common.ZoneEditorHandler) {
	t.Helper()
	counter := integration_common.NewGraphDescriptionCounter(composition.InitializeGuiHandler())
	runner := integration_common.NewAppRunnerWithGuiHandler(t, counter)
	zoneEditor := integration_common.NewHandler(runner).
		ClickLayoutAndZonesTab().
		SelectTopology(geometricHubLayout).
		OpenZoneEditor()
	runner.NextFrame()

	return runner, counter, zoneEditor
}

// describedAcross reports how many graph descriptions an action and the frame
// after it asked for.
func describedAcross(
	runner *integration_common.AppRunner,
	counter *integration_common.GraphDescriptionCounter,
	action func()) int {
	before := counter.Calls()
	action()
	runner.NextFrame()

	return counter.Calls() - before
}

// settleRedraws lays out frames until the window stops asking for an immediate
// redraw, so a test observes only the request its own action causes.
func settleRedraws(t *testing.T, runner *integration_common.AppRunner) {
	t.Helper()
	deadline := time.Now().Add(redrawSettleTimeout)
	for time.Now().Before(deadline) {
		runner.NextFrame()
		if !runner.ImmediateRedrawRequested() {
			return
		}
	}
	t.Fatal("the window never stopped asking for redraws")
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenTheZoneEditorOpens_TheGraphIsDescribedOnce(t *testing.T) {
	// Arrange
	_, counter, _ := openCountedZoneEditor(t)

	// Act
	calls := counter.Calls()

	// Assert
	assert.Equal(t, 1, calls)
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenTheZoneEditorIsIdle_TheGraphIsNotDescribedAgain(t *testing.T) {
	// Arrange
	runner, counter, _ := openCountedZoneEditor(t)

	// Act
	calls := describedAcross(runner, counter, func() {
		for range idleFrameCount {
			runner.NextFrame()
		}
	})

	// Assert
	assert.Equal(t, 0, calls)
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenAConnectionIsDrawn_TheGraphIsDescribedOnce(t *testing.T) {
	// Arrange
	runner, counter, zoneEditor := openCountedZoneEditor(t)
	zoneEditor.ClickAddConnection()

	// Act
	calls := describedAcross(runner, counter, func() {
		zoneEditor.DragFromZoneTo(spawnAZoneName, spawnBZoneName)
	})

	// Assert
	assert.Equal(t, 1, calls)
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenAConnectionIsRightClicked_TheGraphIsDescribedOnce(t *testing.T) {
	// Arrange
	runner, counter, zoneEditor := openCountedZoneEditor(t)

	// Act
	calls := describedAcross(runner, counter, func() {
		zoneEditor.RightClickEdge(hubToSpawnAName)
	})

	// Assert
	assert.Equal(t, 1, calls)
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenASelectedConnectionIsDeleted_TheGraphIsDescribedOnce(t *testing.T) {
	// Arrange
	runner, counter, zoneEditor := openCountedZoneEditor(t)
	zoneEditor.ClickConnection(hubToSpawnAName)

	// Act
	calls := describedAcross(runner, counter, func() {
		zoneEditor.ClickDeleteSelected()
	})

	// Assert
	assert.Equal(t, 1, calls)
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenAZoneIsPlaced_TheGraphIsDescribedOnce(t *testing.T) {
	// Arrange
	runner, counter, zoneEditor := openCountedZoneEditor(t)
	zoneEditor.ClickAddZone()

	// Act
	calls := describedAcross(runner, counter, func() {
		zoneEditor.ClickCanvasAt(emptyCanvasSpot)
	})

	// Assert
	assert.Equal(t, 1, calls)
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenANeutralZoneIsDeleted_TheGraphIsDescribedOnce(t *testing.T) {
	// Arrange
	runner, counter, zoneEditor := openCountedZoneEditor(t)
	zoneEditor.ClickZone(hubZoneName)

	// Act
	calls := describedAcross(runner, counter, func() {
		zoneEditor.ClickDeleteSelected()
	})

	// Assert
	assert.Equal(t, 1, calls)
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenSessionEditsAreUndone_TheGraphIsDescribedOnce(t *testing.T) {
	// Arrange
	runner, counter, zoneEditor := openCountedZoneEditor(t)
	zoneEditor.RightClickEdge(hubToSpawnAName)

	// Act
	calls := describedAcross(runner, counter, func() {
		zoneEditor.ClickUndo()
	})

	// Assert
	assert.Equal(t, 1, calls)
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenTheLayoutIsRevertedToBase_TheGraphIsDescribedOnce(t *testing.T) {
	// Arrange
	runner, counter, zoneEditor := openCountedZoneEditor(t)

	// Act
	calls := describedAcross(runner, counter, func() {
		zoneEditor.ClickRevertToBase()
	})

	// Assert
	assert.Equal(t, 1, calls)
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenAZonesQualityChanges_TheGraphIsDescribedOnce(t *testing.T) {
	// Arrange
	runner, counter, zoneEditor := openCountedZoneEditor(t)
	selectPlacedNeutralZone(zoneEditor)

	// Act
	calls := describedAcross(runner, counter, func() {
		zoneEditor.SelectZoneQuality(goldZoneQualityLabel)
	})

	// Assert
	assert.Equal(t, 1, calls)
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenAZoneIsDragged_TheGraphIsNotDescribedAgain(t *testing.T) {
	// Arrange
	runner, counter, zoneEditor := openCountedZoneEditor(t)

	// Act
	calls := describedAcross(runner, counter, func() {
		zoneEditor.DragZone(hubZoneName, draggedZoneSpot)
	})

	// Assert
	assert.Equal(t, 0, calls)
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenAGuardValueIsTyped_TheGraphIsNotDescribedAgain(t *testing.T) {
	// Arrange
	runner, counter, zoneEditor := openCountedZoneEditor(t)
	zoneEditor.ClickConnection(hubToSpawnAName)

	// Act
	calls := describedAcross(runner, counter, func() {
		zoneEditor.TypeConnectionGuardValue("1")
	})

	// Assert
	assert.Equal(t, 0, calls)
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenAConnectionTypeChanges_TheGraphIsNotDescribedAgain(t *testing.T) {
	// Arrange
	runner, counter, zoneEditor := openCountedZoneEditor(t)
	zoneEditor.ClickConnection(hubToSpawnAName)

	// Act
	calls := describedAcross(runner, counter, func() {
		zoneEditor.SelectConnectionType(directConnectionTypeLabel)
	})

	// Assert
	assert.Equal(t, 0, calls)
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenASpawnZoneDeleteIsRefused_TheGraphIsNotDescribedAgain(t *testing.T) {
	// Arrange
	runner, counter, zoneEditor := openCountedZoneEditor(t)
	zoneEditor.ClickZone(spawnAZoneName)

	// Act
	calls := describedAcross(runner, counter, func() {
		zoneEditor.ClickDeleteSelected()
	})

	// Assert
	assert.Equal(t, 0, calls)
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenASpawnLosesItsOnlyConnection_TheGraphReportsItIsolated(t *testing.T) {
	// Arrange
	runner, counter, zoneEditor := openCountedZoneEditor(t)

	// Act
	zoneEditor.RightClickEdge(hubToSpawnAName)
	runner.NextFrame()

	// Assert
	assert.Equal(t, 1, counter.Last().IsolatedZoneCount)
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenAnIsolatedSpawnIsReconnected_TheGraphReportsNoIsolatedZones(t *testing.T) {
	// Arrange
	runner, counter, zoneEditor := openCountedZoneEditor(t)
	zoneEditor.RightClickEdge(hubToSpawnAName)
	runner.NextFrame()
	require.Equal(t, 1, counter.Last().IsolatedZoneCount, "the spawn must start isolated for the reconnect to mean anything")
	zoneEditor.ClickAddConnection()

	// Act
	zoneEditor.DragFromZoneTo(spawnAZoneName, hubZoneName)
	runner.NextFrame()

	// Assert
	assert.Equal(t, 0, counter.Last().IsolatedZoneCount)
}

// A right click is handled after the toolbar drew the status line, so the
// editor has to ask for the frame that shows the new counts.
//
//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenACanvasEditChangesTheStatus_AnImmediateRedrawIsRequested(t *testing.T) {
	// Arrange
	runner, _, zoneEditor := openCountedZoneEditor(t)
	settleRedraws(t, runner)

	// Act
	zoneEditor.RightClickEdge(hubToSpawnAName)

	// Assert
	assert.True(t, runner.ImmediateRedrawRequested())
}

// The same gesture on empty canvas changes nothing the status line shows.
//
//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenACanvasClickChangesNoStatus_NoImmediateRedrawIsRequested(t *testing.T) {
	// Arrange
	runner, _, zoneEditor := openCountedZoneEditor(t)
	settleRedraws(t, runner)

	// Act
	runner.RightClickAt(zoneEditor.CanvasPoint(emptyCanvasSpot))
	runner.NextFrame()

	// Assert
	assert.False(t, runner.ImmediateRedrawRequested())
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenTheRedrawAfterACanvasEditHasRun_NoFurtherRedrawIsRequested(t *testing.T) {
	// Arrange
	runner, _, zoneEditor := openCountedZoneEditor(t)
	settleRedraws(t, runner)
	zoneEditor.RightClickEdge(hubToSpawnAName)
	require.True(t, runner.ImmediateRedrawRequested(), "the edit must have asked for a redraw for its end to mean anything")

	// Act
	runner.NextFrame()

	// Assert
	assert.False(t, runner.ImmediateRedrawRequested())
}

// In the error state the status line draws neither the hint nor the counts, so
// a comparison against only what it drew would never settle.
func TestWhenTheStatusShowsAnError_IdleFramesRequestNoRedraw(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := integration_common.NewGraphDescriptionCounter(composition.InitializeGuiHandler())
	zones := []template_model.Zone{newGeometryZone("A", 0.2, 0.5), newGeometryZone("B", 0.8, 0.5)}
	stateDto := editor_state_dto.EditorStateDto{EditorState: editor_state_model.NewDefaultEditorStateModel()}
	options := handler.GetZoneEditorOptions(stateDto, len(zones))
	dialog := dialogs.NewZoneEditorDialog(
		zones,
		[]template_model.Connection{newGeometryConnection("ab", "A", "B"), newGeometryConnection("am", "A", "Missing")},
		options.Topology,
		options.Tuning,
		options.GenerateRoads,
		handler,
		nil,
		func() (dtos.ZoneEditorZonesDto, bool) { return dtos.ZoneEditorZonesDto{}, false })
	dialog.ClickRevertToBase()
	gtx, frameRouter := newDialogContext(image.Pt(1000, 720))
	theme := themes.NewTheme()
	frame := func() {
		gtx.Ops.Reset()
		dialog.Body(gtx, theme)
		frameRouter.Frame(gtx.Ops)
	}
	frame()
	require.True(t, handler.Last().HasErrors, "the dialog must be in the error state for this test to mean anything")
	require.NotEmpty(t, dialog.StatusHint(), "the failed revert must have left a hint for the error state to hide")
	frameRouter.WakeupTime()

	// Act
	frame()

	// Assert
	_, requested := frameRouter.WakeupTime()
	assert.False(t, requested)
}
