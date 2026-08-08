//go:build integration_test && gui

package gui_test

import (
	"image"
	"testing"

	"gioui.org/f32"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/dialogs"
	"github.com/Tariomka/hommoe_custom_templates/internal/composition"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// geometryCanvasSide is the canvas side these tests lay geometry out on. At
// exactly 700px the preview metrics scale is 1.0, so every expected coordinate
// below is an exact integer rather than a scaled approximation.
const geometryCanvasSide = 700

// geometryZoneRadius is the zone radius the preview layout settles on for the
// fixtures below: the zones are far enough apart that the radius is not
// shrunk, so it stays at the unscaled maximum.
const geometryZoneRadius = 38

// newGeometryZone builds a zone pinned to a normalized manual position. When
// every zone carries one, the preview layout honours them verbatim, which makes
// the resulting canvas geometry fully deterministic.
func newGeometryZone(name string, x, y float64) entities.Zone {
	position := [2]float64{x, y}

	return entities.Zone{Name: name, ManualPosition: &position}
}

func newGeometryConnection(name, from, to string) entities.Connection {
	return entities.Connection{Name: name, From: from, To: to, ConnectionType: "Direct"}
}

// newGeometryDialog builds a zone editor over the given deterministic fixture
// and lays its geometry out once.
func newGeometryDialog(
	t *testing.T,
	zones []entities.Zone,
	connections []entities.Connection,
) *dialogs.ZoneEditorDialog {
	t.Helper()
	handler := composition.InitializeGuiHandler()
	options := handler.GetZoneEditorOptions(dtos.NewDefaultEditorStateDto(), len(zones))
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
	require.Equal(t, geometryZoneRadius, dialog.CanvasZoneRadius(),
		"the fixture must keep the unscaled zone radius, otherwise every expectation below shifts")

	return dialog
}

// newTriangleFixture places three zones at (140,350), (560,350) and (350,140)
// and links A-B twice plus A-C once, so the parallel-edge spread is exercised.
func newTriangleFixture(t *testing.T) *dialogs.ZoneEditorDialog {
	t.Helper()

	return newGeometryDialog(t,
		[]entities.Zone{
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

// newObstacleFixture places a third zone 14px off the A-B chord, close enough
// to push the curve clear of it.
func newObstacleFixture(t *testing.T) *dialogs.ZoneEditorDialog {
	t.Helper()

	return newGeometryDialog(t,
		[]entities.Zone{
			newGeometryZone("A", 0.2, 0.5),
			newGeometryZone("B", 0.8, 0.5),
			newGeometryZone("D", 0.5, 0.52),
		},
		[]entities.Connection{newGeometryConnection("ab", "A", "B")})
}

func TestWhenZonesCarryManualPositions_TheCanvasPlacesThemVerbatim(t *testing.T) {
	t.Parallel()
	// Arrange
	dialog := newTriangleFixture(t)

	// Act
	positions := dialog.ZonePositions()

	// Assert
	assert.Equal(t, map[string]image.Point{
		"A": image.Pt(140, 350),
		"B": image.Pt(560, 350),
		"C": image.Pt(350, 140),
	}, positions)
}

func TestWhenConnectionsSharePairs_TheyAreGroupedInFirstSeenOrder(t *testing.T) {
	t.Parallel()
	// Arrange
	dialog := newTriangleFixture(t)

	// Act
	edges := dialog.EdgeGeometries()

	// Assert
	assert.Equal(t, []string{"ab", "ba", "ac"}, edgeNames(edges))
}

func TestWhenTwoConnectionsSharePair_TheirCurvesSpreadSymmetrically(t *testing.T) {
	t.Parallel()
	// Arrange
	dialog := newTriangleFixture(t)

	// Act
	edges := dialog.EdgeGeometries()

	// Assert
	assert.Equal(t,
		[]f32.Point{f32.Pt(350, 368), f32.Pt(350, 332)},
		[]f32.Point{edges[0].ControlPoint, edges[1].ControlPoint})
}

func TestWhenAnEdgeIsLaidOut_ItsLabelSitsOnTheCurveMidpoint(t *testing.T) {
	t.Parallel()
	// Arrange
	dialog := newTriangleFixture(t)

	// Act
	edges := dialog.EdgeGeometries()

	// Assert
	assert.Equal(t, image.Pt(350, 359), edges[0].MidPoint)
}

func TestWhenAZoneSitsNearTheChord_TheCurveBulgesClearOfIt(t *testing.T) {
	t.Parallel()
	// Arrange
	dialog := newObstacleFixture(t)

	// Act
	edges := dialog.EdgeGeometries()

	// Assert
	assert.Equal(t, f32.Pt(350, 274), edges[0].ControlPoint)
}

func TestWhenAPointIsInsideAZone_TheHitTestNamesThatZone(t *testing.T) {
	t.Parallel()
	// Arrange
	dialog := newTriangleFixture(t)

	// Act
	hit := dialog.HitTestCanvasNode(image.Pt(140, 350))

	// Assert
	assert.Equal(t, "A", hit)
}

func TestWhenAPointIsJustOutsideEveryZone_TheHitTestNamesNone(t *testing.T) {
	t.Parallel()
	// Arrange
	dialog := newTriangleFixture(t)

	// Act
	hit := dialog.HitTestCanvasNode(image.Pt(140+geometryZoneRadius+1, 350))

	// Assert
	assert.Empty(t, hit)
}

func TestWhenAPointSitsOnACurve_TheEdgeHitTestNamesThatConnection(t *testing.T) {
	t.Parallel()
	// Arrange
	dialog := newObstacleFixture(t)

	// Act
	hit := dialog.HitTestCanvasEdge(image.Pt(350, 312))

	// Assert
	assert.Equal(t, "ab", hit)
}

func TestWhenAPointIsFarFromEveryCurve_TheEdgeHitTestNamesNone(t *testing.T) {
	t.Parallel()
	// Arrange
	dialog := newObstacleFixture(t)

	// Act
	hit := dialog.HitTestCanvasEdge(image.Pt(350, 350))

	// Assert
	assert.Empty(t, hit)
}

func TestWhenTheZoneRadiusIsKnown_TheSnapGridHoldsSevenCellsPerDiameter(t *testing.T) {
	t.Parallel()
	// Arrange
	dialog := newTriangleFixture(t)

	// Act
	step := dialog.CanvasGridStep()

	// Assert
	assert.InDelta(t, float64(geometryZoneRadius)*2.0/7.0, step, 1e-9)
}

func TestWhenSnappingIsDisabled_TheDraggedPositionIsUntouched(t *testing.T) {
	t.Parallel()
	// Arrange
	dialog := newTriangleFixture(t)
	dialog.SetSnapEnabled(false)
	dialog.BeginZoneDrag("A")

	// Act
	snapped := dialog.SnapDraggedPosition(image.Pt(200, 355))

	// Assert
	assert.Equal(t, image.Pt(200, 355), snapped)
}

func TestWhenSnappingIsEnabled_TheDraggedZoneHoldsOntoNearbyGuides(t *testing.T) {
	t.Parallel()
	// Arrange
	dialog := newTriangleFixture(t)
	dialog.SetSnapEnabled(true)
	dialog.BeginZoneDrag("A")

	// Act
	snapped := dialog.SnapDraggedPosition(image.Pt(200, 355))

	// Assert
	assert.Equal(t, image.Pt(201, 350), snapped)
}

func TestWhenAZoneGuideIsHeld_OnlyThatAxisReportsAGuide(t *testing.T) {
	t.Parallel()
	// Arrange
	dialog := newTriangleFixture(t)
	dialog.SetSnapEnabled(true)
	dialog.BeginZoneDrag("A")

	// Act
	dialog.SnapDraggedPosition(image.Pt(200, 355))
	_, xActive, _, yActive := dialog.SnapGuides()

	// Assert
	// The guide coordinate itself is not pinned: the three candidate guides at
	// 312/350/388 are all exactly 5px away, so which one wins depends on the
	// position map's iteration order. Only the axis flags are deterministic.
	assert.Equal(t, []bool{false, true}, []bool{xActive, yActive})
}

func edgeNames(edges []dialogs.EdgeGeometry) []string {
	names := make([]string, 0, len(edges))
	for _, edge := range edges {
		names = append(names, edge.Name)
	}

	return names
}
