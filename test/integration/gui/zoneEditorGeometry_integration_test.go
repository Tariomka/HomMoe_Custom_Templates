//go:build integration_test && gui

package gui_test

import (
	"testing"

	"gioui.org/f32"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/dialogs"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Each topology below is picked for the one geometry property it is the only
// reachable layout to exhibit: Square for a plain two-zone placement, Ring for a
// pair of connections sharing the same two zones, Hub for a zone sitting on the
// chord between two others, and Geometric Hub for a curve that runs straight.
const (
	squareLayout    = "Square"
	ringLayout      = "Ring"
	hubLayout       = "Hub"
	hubToSpawnAEdge = "Portal-Hub-A"
)

// edgePairs reports the zone pair of every edge in layout order, which is what
// the grouping assertions are about. The Hub topology emits more than one
// connection per pair, so names would not tell the grouping apart.
func edgePairs(edges []dialogs.EdgeGeometry) [][2]string {
	pairs := make([][2]string, 0, len(edges))
	for _, edge := range edges {
		pairs = append(pairs, [2]string{edge.From, edge.To})
	}

	return pairs
}

// findEdge returns the first edge with the given name.
func findEdge(t *testing.T, edges []dialogs.EdgeGeometry, name string) dialogs.EdgeGeometry {
	t.Helper()
	for _, edge := range edges {
		if edge.Name == name {
			return edge
		}
	}
	t.Fatalf("the canvas laid out no edge called %q", name)

	return dialogs.EdgeGeometry{}
}

// The canvas rounds once, at draw time. Every coordinate it hands back is the
// unrounded layout value, so pinning them exactly is what stops a rounding step
// from creeping back into the pipeline.
//
//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenTheCanvasPlacesAGeneratedLayout_TheZonePositionsAreUnrounded(t *testing.T) {
	// Arrange
	_, zoneEditor := openZoneEditor(t, squareLayout, false)

	// Act
	positions := zoneEditor.Dialog().ZonePositions()

	// Assert
	assert.Equal(t, map[string]models.Position{
		"Spawn-A": data.NewVec2(46.39999999999995, 46.39999999999995),
		"Spawn-B": data.NewVec2(533.6, 533.6),
	}, positions)
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenTheDialogOpens_TheCanvasIsSquaredToTheBoxItWasGiven(t *testing.T) {
	// Arrange
	_, zoneEditor := openZoneEditor(t, squareLayout, false)

	// Act
	side := zoneEditor.Dialog().CanvasSquareSide()

	// Assert
	assert.Equal(t, 580, side)
}

// Two zones are one pair whichever way round a connection names them, so the
// reversed Spawn-B to Spawn-A edge groups with its forward twin rather than
// starting a group of its own.
//
//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenConnectionsSharePairs_TheyAreGroupedInFirstSeenOrder(t *testing.T) {
	// Arrange
	_, zoneEditor := openZoneEditor(t, hubLayout, false)

	// Act
	pairs := edgePairs(zoneEditor.Dialog().EdgeGeometries())

	// Assert
	assert.Equal(t, [][2]string{
		{"Hub", "Spawn-A"},
		{"Hub", "Spawn-A"},
		{"Spawn-A", "Spawn-B"},
		{"Spawn-B", "Spawn-A"},
		{"Hub", "Spawn-B"},
		{"Hub", "Spawn-B"},
	}, pairs)
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenTwoConnectionsSharePair_TheirCurvesSpreadSymmetrically(t *testing.T) {
	// Arrange
	_, zoneEditor := openZoneEditor(t, ringLayout, false)
	edges := zoneEditor.Dialog().EdgeGeometries()
	require.Len(t, edges, 2, "the ring layout must pair the two spawns twice for the spread to show")

	// Act
	controlPoints := []f32.Point{edges[0].ControlPoint, edges[1].ControlPoint}

	// Assert
	assert.Equal(t, []f32.Point{f32.Pt(272, 290), f32.Pt(308, 290)}, controlPoints)
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenAnEdgeIsLaidOut_ItsLabelSitsOnTheCurveMidpoint(t *testing.T) {
	// Arrange
	_, zoneEditor := openZoneEditor(t, ringLayout, false)

	// Act
	edge := findEdge(t, zoneEditor.Dialog().EdgeGeometries(), "Ring-A-B")

	// Assert
	assert.Equal(t, data.NewVec2(281.0, 290.0), edge.MidPoint)
}

// The hub sits exactly on the chord between the two spawns, so the pseudo edge
// joining them has to bow far enough sideways to clear it.
//
//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenAZoneSitsNearTheChord_TheCurveBulgesClearOfIt(t *testing.T) {
	// Arrange
	_, zoneEditor := openZoneEditor(t, hubLayout, false)

	// Act
	edge := findEdge(t, zoneEditor.Dialog().EdgeGeometries(), "Pseudo-A-B")

	// Assert
	assert.Equal(t, f32.Pt(181.02856, 290), edge.ControlPoint)
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenAPointIsInsideAZone_TheHitTestNamesThatZone(t *testing.T) {
	// Arrange
	_, zoneEditor := openZoneEditor(t, squareLayout, false)

	// Act
	hit := zoneEditor.Dialog().HitTestCanvasNode(zoneEditor.ZonePosition("Spawn-A"))

	// Assert
	assert.Equal(t, "Spawn-A", hit)
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenAPointIsJustOutsideEveryZone_TheHitTestNamesNone(t *testing.T) {
	// Arrange
	_, zoneEditor := openZoneEditor(t, squareLayout, false)
	center := zoneEditor.ZonePosition("Spawn-A")
	justOutside := data.NewVec2(center.X+zoneEditor.Dialog().CanvasZoneRadius()+1, center.Y)

	// Act
	hit := zoneEditor.Dialog().HitTestCanvasNode(justOutside)

	// Assert
	assert.Empty(t, hit)
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenAPointSitsOnACurve_TheEdgeHitTestNamesThatConnection(t *testing.T) {
	// Arrange
	_, zoneEditor := openZoneEditor(t, geometricHubLayout, false)
	edge := findEdge(t, zoneEditor.Dialog().EdgeGeometries(), hubToSpawnAEdge)

	// Act
	hit := zoneEditor.Dialog().HitTestCanvasEdge(edge.MidPoint)

	// Assert
	assert.Equal(t, hubToSpawnAEdge, hit)
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenAPointIsFarFromEveryCurve_TheEdgeHitTestNamesNone(t *testing.T) {
	// Arrange
	_, zoneEditor := openZoneEditor(t, geometricHubLayout, false)

	// Act
	hit := zoneEditor.Dialog().HitTestCanvasEdge(data.NewVec2(10.0, 10.0))

	// Assert
	assert.Empty(t, hit)
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenTheZoneRadiusIsKnown_TheSnapGridHoldsSevenCellsPerDiameter(t *testing.T) {
	// Arrange
	_, zoneEditor := openZoneEditor(t, squareLayout, false)
	radius := zoneEditor.Dialog().CanvasZoneRadius()

	// Act
	step := zoneEditor.Dialog().CanvasGridStep()

	// Assert
	assert.InDelta(t, radius*2/7, step, 1e-9)
}
