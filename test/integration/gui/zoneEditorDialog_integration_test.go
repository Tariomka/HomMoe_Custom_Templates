//go:build integration_test && gui

package gui_test

import (
	"image"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/dialogs"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/themes"
	"github.com/Tariomka/hommoe_custom_templates/internal/composition"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Everything the running window can reach lives in
// zoneEditorActions_integration_test.go and zoneEditorGeometry_integration_test.go,
// driven through AppRunner. What is left here is what a real click cannot
// express: the Apply and revert-to-base callbacks, which the window wires to the
// state driver and never hands back, and the snap arithmetic, which needs an
// exact dragged position rather than a pointer gesture rounded to a pixel.

// geometryCanvasSide is the canvas side at which the preview metrics scale is
// exactly 1.0, so a fixture position of 0.2 lands on 140 and stays readable.
const geometryCanvasSide = 700

// newGeometryZone builds a zone pinned at a normalized position.
func newGeometryZone(name string, x, y float64) template_model.Zone {
	return template_model.Zone{Name: name, ManualPosition: &[2]float64{x, y}}
}

// newGeometryConnection builds a plain connection between two zones.
func newGeometryConnection(name, from, to string) entities.Connection {
	return entities.Connection{Name: name, From: from, To: to, ConnectionType: "Direct"}
}

// newGeometryDialog builds a zone editor over the given layout and lays its
// canvas out once, so the geometry is available to read back.
func newGeometryDialog(
	t *testing.T,
	zones []template_model.Zone,
	connections []entities.Connection,
) *dialogs.ZoneEditorDialog {
	t.Helper()
	handler := composition.InitializeGuiHandler()
	stateDto := editor_state_dto.EditorStateDto{EditorState: editor_state_model.NewDefaultEditorStateModel()}
	options := handler.GetZoneEditorOptions(stateDto, len(zones))
	dialog := dialogs.NewZoneEditorDialog(
		zones,
		connections,
		options.Topology,
		options.Tuning,
		options.GenerateRoads,
		handler,
		nil,
		nil,
	)
	dialog.RecomputeGeometry(geometryCanvasSide)

	return dialog
}

// newTriangleFixture builds three zones with two of them sharing a horizontal
// row, which is what the snap guides latch onto.
func newTriangleFixture(t *testing.T) *dialogs.ZoneEditorDialog {
	t.Helper()

	return newGeometryDialog(t,
		[]template_model.Zone{
			newGeometryZone("A", 0.2, 0.5),
			newGeometryZone("B", 0.8, 0.5),
			newGeometryZone("C", 0.5, 0.2),
		},
		[]entities.Connection{
			newGeometryConnection("ab", "A", "B"),
			newGeometryConnection("ac", "A", "C"),
			newGeometryConnection("ba", "B", "A"),
		})
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
	baseZones []template_model.Zone,
) (*dialogs.ZoneEditorDialog, *applyCapture) {
	t.Helper()
	handler := composition.InitializeGuiHandler()
	zones := []template_model.Zone{
		newGeometryZone("A", 0.2, 0.5),
		newGeometryZone("B", 0.8, 0.5),
	}
	connections := []entities.Connection{newGeometryConnection("ab", "A", "B")}
	stateDto := editor_state_dto.EditorStateDto{EditorState: editor_state_model.NewDefaultEditorStateModel()}
	options := handler.GetZoneEditorOptions(stateDto, len(zones))
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

func TestWhenZoneEditorDialogRenders_UsesHandlerProvidedOptions(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := composition.InitializeGuiHandler()
	state := editor_state_dto.EditorStateDto{EditorState: editor_state_model.NewDefaultEditorStateModel()}
	generated, err := handler.GenerateTemplate(state)
	require.NoError(t, err)
	require.NotNil(t, generated.Template)
	require.NotEmpty(t, generated.Template.Variants)
	variant := generated.Template.Variants[0]
	options := handler.GetZoneEditorOptions(state, len(variant.Zones))
	dialog := dialogs.NewZoneEditorDialog(
		variant.Zones,
		template_model.ToConnectionEntities(variant.Connections),
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
	base := []template_model.Zone{newGeometryZone("Fresh1", 0.3, 0.3), newGeometryZone("Fresh2", 0.7, 0.7)}
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

func TestWhenSnappingIsDisabled_TheDraggedPositionIsUntouched(t *testing.T) {
	t.Parallel()
	// Arrange
	dialog := newTriangleFixture(t)
	dialog.SetSnapEnabled(false)
	dialog.BeginZoneDrag("A")

	// Act
	snapped := dialog.SnapDraggedPosition(data.NewVec2(203.0, 351.0))

	// Assert
	assert.Equal(t, data.NewVec2(203.0, 351.0), snapped)
}

func TestWhenSnappingIsEnabled_TheDraggedZoneHoldsOntoNearbyGuides(t *testing.T) {
	t.Parallel()
	// Arrange
	dialog := newTriangleFixture(t)
	dialog.SetSnapEnabled(true)
	dialog.BeginZoneDrag("A")

	// Act
	snapped := dialog.SnapDraggedPosition(data.NewVec2(203.0, 351.0))

	// Assert
	assert.Equal(t, data.NewVec2(200+6.0/7.0, 350.0), snapped)
}

func TestWhenAZoneGuideIsHeld_OnlyThatAxisReportsAGuide(t *testing.T) {
	t.Parallel()
	// Arrange
	dialog := newTriangleFixture(t)
	dialog.SetSnapEnabled(true)
	dialog.BeginZoneDrag("A")
	dialog.SnapDraggedPosition(data.NewVec2(203.0, 351.0))

	// Act
	_, xActive, _, yActive := dialog.SnapGuides()

	// Assert
	assert.Equal(t, []bool{false, true}, []bool{xActive, yActive})
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
