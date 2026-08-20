//go:build integration_test && gui

package gui_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers/integration_common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The two player spawns of the Geometric Hub layout, plus the zone the editor
// appends when one is placed by hand.
const (
	spawnAZoneName = "Spawn-A"
	spawnBZoneName = "Spawn-B"
	placedZoneName = "Neutral-C"
)

// The Geometric Hub layout stacks its three zones on one vertical line at
// x = 290, so a spot well to the left of that line is empty canvas whatever the
// row, and a target three pixels off the line is inside the zone-alignment snap
// threshold. Both facts are what the pointer tests below are built on.
var (
	emptyCanvasSpot  = data.NewVec2(120.0, 120.0)
	draggedZoneSpot  = data.NewVec2(200.0, 200.0)
	nearZoneLineSpot = data.NewVec2(296.0, 250.0)
)

// manualZoneSave finds the committed manual record of a zone, which is what
// Apply writes and what a reload would restore the layout from.
func manualZoneSave(
	t *testing.T,
	runner *integration_common.AppRunner,
	name string) editor_state_dto.ManualZoneSave {
	t.Helper()
	for _, save := range runner.CurrentState().ManualZones {
		if save.Zone.Name == name {
			return save
		}
	}
	t.Fatalf("the editor state committed no manual zone called %q", name)

	return editor_state_dto.ManualZoneSave{}
}

// A drag has to survive the round trip through the canvas' normalized manual
// position, so the assertion is on what Apply committed rather than on where
// the canvas drew the zone mid-drag.
//
//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenAZoneIsDraggedToANewPosition_TheAppliedLayoutRecordsIt(t *testing.T) {
	// Arrange
	runner, zoneEditor := openZoneEditor(t, geometricHubLayout, true)
	zoneEditor.DragZone(hubZoneName, draggedZoneSpot)
	require.Equal(t, draggedZoneSpot, zoneEditor.ZonePosition(hubZoneName),
		"the drag must have moved the zone for Apply to have something to commit")

	// Act
	zoneEditor.ClickApply()

	// Assert
	assert.Equal(t,
		&[2]float64{0.3448275862068966, 0.3448275862068966},
		manualZoneSave(t, runner, hubZoneName).ManualPosition)
}

// With snapping on, the drop point is pulled onto the other zones' centre line
// horizontally and onto the nearest grid intersection vertically. Neither
// correction is rounded to a whole pixel, which is what the fractional y pins.
//
//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenAZoneIsDraggedNearAGuide_ItSnapsToTheGuide(t *testing.T) {
	// Arrange
	_, zoneEditor := openZoneEditor(t, geometricHubLayout, true)
	zoneEditor.ToggleSnap()
	require.True(t, zoneEditor.Dialog().SnapEnabled(), "the drag below only snaps while snapping is on")

	// Act
	zoneEditor.DragZone(hubZoneName, nearZoneLineSpot)

	// Assert
	assert.Equal(t, data.NewVec2(290.0, 251.88571428571436), zoneEditor.ZonePosition(hubZoneName))
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenADragStartsOnAZoneInAddConnectionMode_AConnectionIsCreated(t *testing.T) {
	// Arrange
	_, zoneEditor := openZoneEditor(t, geometricHubLayout, true)
	zoneEditor.ClickAddConnection()

	// Act
	zoneEditor.DragFromZoneTo(spawnAZoneName, spawnBZoneName)

	// Assert
	assert.Equal(t,
		[][2]string{{hubZoneName, spawnAZoneName}, {hubZoneName, spawnBZoneName}, {spawnAZoneName, spawnBZoneName}},
		edgePairs(zoneEditor.Dialog().EdgeGeometries()))
}

// The same gesture released over empty canvas has no target zone, so the
// rubber band is simply abandoned.
//
//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenADragEndsOnEmptyCanvas_NoConnectionIsCreated(t *testing.T) {
	// Arrange
	_, zoneEditor := openZoneEditor(t, geometricHubLayout, true)
	zoneEditor.ClickAddConnection()

	// Act
	zoneEditor.DragZone(spawnAZoneName, emptyCanvasSpot)

	// Assert
	assert.Equal(t,
		[][2]string{{hubZoneName, spawnAZoneName}, {hubZoneName, spawnBZoneName}},
		edgePairs(zoneEditor.Dialog().EdgeGeometries()))
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenACurveIsRightClicked_ThatConnectionIsDeleted(t *testing.T) {
	// Arrange
	_, zoneEditor := openZoneEditor(t, geometricHubLayout, true)

	// Act
	zoneEditor.RightClickEdge(hubToSpawnAName)

	// Assert
	assert.Equal(t, []string{hubToSpawnBName}, zoneEditor.Dialog().EditedConnectionNames())
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenAddZoneModeIsArmedAndEmptyCanvasIsClicked_AZoneIsPlaced(t *testing.T) {
	// Arrange
	_, zoneEditor := openZoneEditor(t, geometricHubLayout, true)
	zoneEditor.ClickAddZone()

	// Act
	zoneEditor.ClickCanvasAt(emptyCanvasSpot)

	// Assert
	assert.Equal(t,
		[]string{hubZoneName, spawnAZoneName, spawnBZoneName, placedZoneName},
		editedZoneNames(zoneEditor))
}

// The placed zone lands where it was clicked rather than at some free slot the
// editor picked, so a second click places a second zone somewhere else.
//
//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenAZoneIsPlacedOnEmptyCanvas_ItSitsWhereItWasClicked(t *testing.T) {
	// Arrange
	_, zoneEditor := openZoneEditor(t, geometricHubLayout, true)
	zoneEditor.ClickAddZone()

	// Act
	zoneEditor.ClickCanvasAt(emptyCanvasSpot)

	// Assert
	assert.Equal(t, emptyCanvasSpot, zoneEditor.ZonePosition(placedZoneName))
}

// A press that never leaves the dead zone is a selection, not a move, so the
// zone stays exactly where the generator put it.
//
//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenAZoneIsDraggedInsideTheDeadZone_ItDoesNotMove(t *testing.T) {
	// Arrange
	_, zoneEditor := openZoneEditor(t, geometricHubLayout, true)
	start := zoneEditor.ZonePosition(hubZoneName)

	// Act
	zoneEditor.DragZone(hubZoneName, shiftedBy(start, 4.0))

	// Assert
	assert.Equal(t, start, zoneEditor.ZonePosition(hubZoneName))
}

// shiftedBy offsets a canvas position diagonally, for the drags that are about
// how far the pointer travelled rather than where it ended up.
func shiftedBy(position models.Position, distance float64) models.Position {
	return data.NewVec2(position.X+distance, position.Y+distance)
}
