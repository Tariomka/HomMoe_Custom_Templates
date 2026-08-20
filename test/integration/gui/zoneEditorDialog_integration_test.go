//go:build integration_test && gui

package gui_test

import (
	"image"
	"testing"

	"gioui.org/io/semantic"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/dialogs"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/themes"
	"github.com/Tariomka/hommoe_custom_templates/internal/composition"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenZoneEditorDialogRenders_UsesHandlerProvidedOptions(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := composition.InitializeGuiHandler()
	state := dtos.NewDefaultEditorStateDto()
	generated, err := handler.GenerateTemplate(state)
	require.NoError(t, err)
	require.NotNil(t, generated.Template)
	require.NotEmpty(t, generated.Template.Variants)
	variant := generated.Template.Variants[0]
	options := handler.GetZoneEditorOptions(state, len(variant.Zones))
	dialog := dialogs.NewZoneEditorDialog(
		variant.Zones,
		variant.Connections,
		options.Topology,
		options.Tuning,
		options.GenerateRoads,
		handler,
		nil,
		nil,
	)
	gtx, frameRouter := newDialogContext(image.Pt(1000, 720))

	// Act
	dimensions, closed := dialog.Body(gtx, themes.NewTheme())
	frameRouter.Frame(gtx.Ops)

	// Assert
	assert.Equal(t, image.Pt(1000, 720), dimensions.Size)
	assert.False(t, closed)
}

// frameZoneEditor lays the dialog out once at its preferred size and reports
// whether it asked to close.
func frameZoneEditor(t *testing.T, dialog *dialogs.ZoneEditorDialog) bool {
	t.Helper()
	gtx, frameRouter := newDialogContext(image.Pt(1000, 720))
	_, closed := dialog.Body(gtx, themes.NewTheme())
	frameRouter.Frame(gtx.Ops)

	return closed
}

// zoneEditorButtonLabels lays the dialog out and reports the label of every
// button the frame published semantics for.
func zoneEditorButtonLabels(t *testing.T, dialog *dialogs.ZoneEditorDialog) []string {
	t.Helper()
	gtx, frameRouter := newDialogContext(image.Pt(1000, 720))
	dialog.Body(gtx, themes.NewTheme())
	frameRouter.Frame(gtx.Ops)
	labels := make([]string, 0)
	for _, node := range frameRouter.AppendSemantics(nil) {
		if node.Desc.Class == semantic.Button && node.Desc.Label != "" {
			labels = append(labels, node.Desc.Label)
		}
	}
	require.NotEmpty(t, labels, "the frame must publish button semantics for the assertions to mean anything")

	return labels
}

// applyCapture records what the zone editor hands to its Apply callback.
type applyCapture struct {
	request dtos.ZoneEditorZonesDto
	fired   bool
}

// newApplyCaptureFixture builds a two-zone editor whose Apply result is
// captured instead of committed. baseZones, when non-empty, is what a revert
// to base hands back; an empty slice makes the revert report failure.
func newApplyCaptureFixture(
	t *testing.T,
	baseZones []entities.Zone,
) (*dialogs.ZoneEditorDialog, *applyCapture) {
	t.Helper()
	handler := composition.InitializeGuiHandler()
	zones := []entities.Zone{
		newGeometryZone("A", 0.2, 0.5),
		newGeometryZone("B", 0.8, 0.5),
	}
	connections := []entities.Connection{newGeometryConnection("ab", "A", "B")}
	options := handler.GetZoneEditorOptions(dtos.NewDefaultEditorStateDto(), len(zones))
	capture := &applyCapture{}
	dialog := dialogs.NewZoneEditorDialog(
		zones,
		connections,
		options.Topology,
		options.Tuning,
		options.GenerateRoads,
		handler,
		func(request dtos.ZoneEditorZonesDto) {
			capture.request = request
			capture.fired = true
		},
		func() (dtos.ZoneEditorZonesDto, bool) {
			if len(baseZones) == 0 {
				return dtos.ZoneEditorZonesDto{}, false
			}

			return dtos.ZoneEditorZonesDto{Zones: baseZones}, true
		},
	)

	return dialog, capture
}

// zoneNames reports the names of the zones the editor currently holds.
func zoneNames(dialog *dialogs.ZoneEditorDialog) []string {
	names := make([]string, 0)
	for _, zone := range dialog.EditedZones() {
		names = append(names, zone.Name)
	}

	return names
}

// The button restores the layout the dialog was opened with, which is not the
// generated one once manual edits have been applied before - so it must not
// claim to reset to generated.
func TestWhenTheZoneEditorRenders_TheUndoButtonIsLabelledForWhatItDoes(t *testing.T) {
	t.Parallel()
	// Arrange
	dialog := newTriangleFixture(t)

	// Act
	labels := zoneEditorButtonLabels(t, dialog)

	// Assert
	assert.Contains(t, labels, "Undo")
}

func TestWhenTheZoneEditorRenders_TheRevertToBaseButtonIsOffered(t *testing.T) {
	t.Parallel()
	// Arrange
	dialog := newTriangleFixture(t)

	// Act
	labels := zoneEditorButtonLabels(t, dialog)

	// Assert
	assert.Contains(t, labels, "Revert to Base")
}

func TestWhenTheZoneEditorRenders_NoButtonClaimsToResetToGenerated(t *testing.T) {
	t.Parallel()
	// Arrange
	dialog := newTriangleFixture(t)

	// Act
	labels := zoneEditorButtonLabels(t, dialog)

	// Assert
	assert.NotContains(t, labels, "Reset to generated")
}

// Undo is session-scoped: it must not tell the driver to drop anything, so
// Apply after an Undo simply commits the zones the editor started from.
func TestWhenApplyFollowsAnUndo_TheStartingZonesAreReported(t *testing.T) {
	t.Parallel()
	// Arrange
	dialog, capture := newApplyCaptureFixture(t, nil)
	dialog.SelectConnection("ab")
	dialog.ClickDeleteSelected()
	frameZoneEditor(t, dialog)
	dialog.ClickUndo()
	frameZoneEditor(t, dialog)
	dialog.ClickApply()

	// Act
	frameZoneEditor(t, dialog)

	// Assert
	require.True(t, capture.fired, "Apply must have reached the callback")
	assert.Len(t, capture.request.Connections, 1)
}

func TestWhenRevertToBaseIsPressed_TheEditorShowsTheRegeneratedZones(t *testing.T) {
	t.Parallel()
	// Arrange
	base := []entities.Zone{newGeometryZone("Fresh1", 0.3, 0.3), newGeometryZone("Fresh2", 0.7, 0.7)}
	dialog, _ := newApplyCaptureFixture(t, base)
	dialog.ClickRevertToBase()

	// Act
	frameZoneEditor(t, dialog)

	// Assert
	assert.Equal(t, []string{"Fresh1", "Fresh2"}, zoneNames(dialog))
}

// After a revert the base becomes the new baseline, so Undo returns to it
// rather than to the edits that were discarded.
func TestWhenUndoFollowsARevertToBase_TheBaseZonesComeBack(t *testing.T) {
	t.Parallel()
	// Arrange
	base := []entities.Zone{newGeometryZone("Fresh1", 0.3, 0.3), newGeometryZone("Fresh2", 0.7, 0.7)}
	dialog, _ := newApplyCaptureFixture(t, base)
	dialog.ClickRevertToBase()
	frameZoneEditor(t, dialog)
	dialog.SelectZone("Fresh1")
	dialog.ClickDeleteSelected()
	frameZoneEditor(t, dialog)
	dialog.ClickUndo()

	// Act
	frameZoneEditor(t, dialog)

	// Assert
	assert.Equal(t, []string{"Fresh1", "Fresh2"}, zoneNames(dialog))
}

func TestWhenRevertToBaseCannotRegenerate_TheEditorKeepsItsZones(t *testing.T) {
	t.Parallel()
	// Arrange
	dialog, _ := newApplyCaptureFixture(t, nil)
	dialog.ClickRevertToBase()

	// Act
	frameZoneEditor(t, dialog)

	// Assert
	assert.Equal(t, []string{"A", "B"}, zoneNames(dialog))
}

func TestWhenRevertToBaseCannotRegenerate_TheEditorSaysSo(t *testing.T) {
	t.Parallel()
	// Arrange
	dialog, _ := newApplyCaptureFixture(t, nil)
	dialog.ClickRevertToBase()

	// Act
	frameZoneEditor(t, dialog)

	// Assert
	assert.Contains(t, dialog.StatusHint(), "Could not regenerate")
}

// The revert only takes effect when the user applies, so the editor has to
// tell the driver that this apply came from a revert.
func TestWhenApplyFollowsARevertToBase_TheRevertIsReported(t *testing.T) {
	t.Parallel()
	// Arrange
	base := []entities.Zone{newGeometryZone("Fresh1", 0.3, 0.3), newGeometryZone("Fresh2", 0.7, 0.7)}
	dialog, capture := newApplyCaptureFixture(t, base)
	dialog.ClickRevertToBase()
	frameZoneEditor(t, dialog)
	dialog.ClickApply()

	// Act
	frameZoneEditor(t, dialog)

	// Assert
	require.True(t, capture.fired, "Apply must have reached the callback")
	assert.True(t, capture.request.RevertToBase)
}

func TestWhenApplyDoesNotFollowARevertToBase_NoRevertIsReported(t *testing.T) {
	t.Parallel()
	// Arrange
	dialog, capture := newApplyCaptureFixture(t, nil)
	dialog.ClickApply()

	// Act
	frameZoneEditor(t, dialog)

	// Assert
	require.True(t, capture.fired, "Apply must have reached the callback")
	assert.False(t, capture.request.RevertToBase)
}

// A failed reroll changes nothing, so it must not make the next apply drop the
// user's manual edits.
func TestWhenRevertToBaseFailed_TheApplyReportsNoRevert(t *testing.T) {
	t.Parallel()
	// Arrange
	dialog, capture := newApplyCaptureFixture(t, nil)
	dialog.ClickRevertToBase()
	frameZoneEditor(t, dialog)
	dialog.ClickApply()

	// Act
	frameZoneEditor(t, dialog)

	// Assert
	require.True(t, capture.fired, "Apply must have reached the callback")
	assert.False(t, capture.request.RevertToBase)
}

func TestWhenAConnectionIsSelected_TheEditorRendersItsPropertyPanel(t *testing.T) {
	t.Parallel()
	// Arrange
	dialog := newTriangleFixture(t)
	require.True(t, dialog.SelectConnection("ab"))

	// Act
	closed := frameZoneEditor(t, dialog)

	// Assert
	assert.Equal(t, []any{false, "ab"}, []any{closed, dialog.SelectedConnection()})
}

func TestWhenAZoneIsSelected_TheEditorRendersItsPropertyPanel(t *testing.T) {
	t.Parallel()
	// Arrange
	dialog := newTriangleFixture(t)
	dialog.SelectZone("A")

	// Act
	closed := frameZoneEditor(t, dialog)

	// Assert
	assert.Equal(t, []any{false, "A", ""}, []any{closed, dialog.SelectedZone(), dialog.SelectedConnection()})
}

func TestWhenSnapIsOnAndAZoneIsDragged_TheEditorRendersTheGuideOverlay(t *testing.T) {
	t.Parallel()
	// Arrange
	dialog := newTriangleFixture(t)
	dialog.SetSnapEnabled(true)
	dialog.BeginZoneDrag("A")
	dialog.SnapDraggedPosition(data.NewVec2(200.0, 355.0))
	_, _, _, yActive := dialog.SnapGuides()
	require.True(t, yActive, "the fixture must hold a horizontal guide for the overlay to render")

	// Act
	closed := frameZoneEditor(t, dialog)

	// Assert
	assert.False(t, closed)
}

func TestWhenTheSelectedConnectionIsDeleted_ItLeavesTheWorkingSet(t *testing.T) {
	t.Parallel()
	// Arrange
	dialog := newTriangleFixture(t)
	require.True(t, dialog.SelectConnection("ab"))
	dialog.ClickDeleteSelected()

	// Act
	frameZoneEditor(t, dialog)

	// Assert
	assert.Equal(t, []string{"ac", "ba"}, dialog.EditedConnectionNames())
}

func TestWhenTheSelectedZoneIsDeleted_ItsConnectionsGoWithIt(t *testing.T) {
	t.Parallel()
	// Arrange
	dialog := newTriangleFixture(t)
	dialog.SelectZone("C")
	dialog.ClickDeleteSelected()

	// Act
	frameZoneEditor(t, dialog)

	// Assert
	assert.Equal(t, []string{"ab", "ba"}, dialog.EditedConnectionNames())
}

func TestWhenTheSessionEditsAreUndone_TheDeletedConnectionComesBack(t *testing.T) {
	t.Parallel()
	// Arrange
	dialog := newTriangleFixture(t)
	require.True(t, dialog.SelectConnection("ab"))
	dialog.ClickDeleteSelected()
	frameZoneEditor(t, dialog)
	dialog.ClickUndo()

	// Act
	frameZoneEditor(t, dialog)

	// Assert
	assert.Equal(t, []string{"ab", "ac", "ba"}, dialog.EditedConnectionNames())
}

func TestWhenTheSessionEditsAreUndone_TheSelectionIsCleared(t *testing.T) {
	t.Parallel()
	// Arrange
	dialog := newTriangleFixture(t)
	dialog.SelectZone("C")
	dialog.ClickUndo()

	// Act
	frameZoneEditor(t, dialog)

	// Assert
	assert.Empty(t, dialog.SelectedZone())
}

func TestWhenAddConnectionIsClicked_TheEditorEntersAddConnectionMode(t *testing.T) {
	t.Parallel()
	// Arrange
	dialog := newTriangleFixture(t)
	dialog.ClickAddConnection()

	// Act
	frameZoneEditor(t, dialog)

	// Assert
	assert.True(t, dialog.AddConnectionModeActive())
}

func TestWhenAddConnectionIsClickedTwice_TheEditorLeavesAddConnectionMode(t *testing.T) {
	t.Parallel()
	// Arrange
	dialog := newTriangleFixture(t)
	dialog.ClickAddConnection()
	frameZoneEditor(t, dialog)
	dialog.ClickAddConnection()

	// Act
	frameZoneEditor(t, dialog)

	// Assert
	assert.False(t, dialog.AddConnectionModeActive())
}

func TestWhenAddZoneIsClickedWhileAddingAConnection_TheAddConnectionModeTurnsOff(t *testing.T) {
	t.Parallel()
	// Arrange
	dialog := newTriangleFixture(t)
	dialog.ClickAddConnection()
	frameZoneEditor(t, dialog)
	dialog.ClickAddZone()

	// Act
	frameZoneEditor(t, dialog)

	// Assert
	assert.False(t, dialog.AddConnectionModeActive())
}

func TestWhenAddZoneIsClickedWhileAddingAConnection_TheAddZoneModeTurnsOn(t *testing.T) {
	t.Parallel()
	// Arrange
	dialog := newTriangleFixture(t)
	dialog.ClickAddConnection()
	frameZoneEditor(t, dialog)
	dialog.ClickAddZone()

	// Act
	frameZoneEditor(t, dialog)

	// Assert
	assert.True(t, dialog.AddZoneModeActive())
}

func TestWhenTheSessionEditsAreUndone_TheAddModeTurnsOff(t *testing.T) {
	t.Parallel()
	// Arrange
	dialog := newTriangleFixture(t)
	dialog.ClickAddConnection()
	frameZoneEditor(t, dialog)
	dialog.ClickUndo()

	// Act
	frameZoneEditor(t, dialog)

	// Assert
	assert.False(t, dialog.AddConnectionModeActive())
}
