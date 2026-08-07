//go:build integration_test && gui

package gui_test

import (
	"image"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/dialogs"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/themes"
	"github.com/Tariomka/hommoe_custom_templates/internal/composition"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
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
	dialog.SnapDraggedPosition(image.Pt(200, 355))
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

func TestWhenTheEditorIsResetToGenerated_TheDeletedConnectionComesBack(t *testing.T) {
	t.Parallel()
	// Arrange
	dialog := newTriangleFixture(t)
	require.True(t, dialog.SelectConnection("ab"))
	dialog.ClickDeleteSelected()
	frameZoneEditor(t, dialog)
	dialog.ClickReset()

	// Act
	frameZoneEditor(t, dialog)

	// Assert
	assert.Equal(t, []string{"ab", "ac", "ba"}, dialog.EditedConnectionNames())
}

func TestWhenTheEditorIsResetToGenerated_TheSelectionIsCleared(t *testing.T) {
	t.Parallel()
	// Arrange
	dialog := newTriangleFixture(t)
	dialog.SelectZone("C")
	dialog.ClickReset()

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

func TestWhenTheEditorIsResetToGenerated_TheAddModeTurnsOff(t *testing.T) {
	t.Parallel()
	// Arrange
	dialog := newTriangleFixture(t)
	dialog.ClickAddConnection()
	frameZoneEditor(t, dialog)
	dialog.ClickReset()

	// Act
	frameZoneEditor(t, dialog)

	// Assert
	assert.False(t, dialog.AddConnectionModeActive())
}
