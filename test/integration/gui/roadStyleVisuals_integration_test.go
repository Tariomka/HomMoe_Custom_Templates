//go:build integration_test && gui

package gui_test

import (
	"image"
	"image/color"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/dialogs"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/themes"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers/integration_common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Both canvases have to report a connection's road state in the colour they
// paint it with, so these tests read the rendered pixels rather than a golden:
// a golden proves the frame did not change, not that an edge is the right
// colour. Each test renders a layout whose road states it fully controls and
// asserts which of the four styles the canvas actually painted.

const (
	roadedStyle         = "road"
	roadlessStyle       = "no road"
	roadedPortalStyle   = "portal"
	roadlessPortalStyle = "portal without road"
)

// squareLayoutName is a topology whose connections are plain direct ones, which
// is what the road setting restyles. Geometric Hub's are explicit portals and
// keep their own colours whatever the setting says.
const squareLayoutName = "Square"

// editorLegendStripPx is how much of the bottom of the editor canvas the road
// key occupies, cut away when the assertion is about the edges.
const editorLegendStripPx = 60

// styleColorTolerance absorbs the one-step difference between the two paths a
// translucent colour composites through - a stroked curve and the legend's
// filled swatch land on 0x4D and 0x4E red.
const styleColorTolerance = 2

// roadedPortalOnCanvas is what the existing translucent portal blue composites
// to over the canvas backdrop. The other three styles are opaque and are read
// straight off the theme.
//
//nolint:gochecknoglobals // A colour literal, allocated once for every test here.
var roadedPortalOnCanvas = color.NRGBA{R: 0x4D, G: 0x92, B: 0xB4, A: 0xFF}

func roadStyleColors() map[string]color.NRGBA {
	return map[string]color.NRGBA{
		roadedStyle:         themes.ColorsPreview.DirectLine,
		roadlessStyle:       themes.ColorsPreview.NoRoadLine,
		roadedPortalStyle:   roadedPortalOnCanvas,
		roadlessPortalStyle: themes.ColorsPreview.PortalNoRoadLine,
	}
}

// styleCensus reports which of the four road styles appear inside area, so one
// assertion can pin down both what was painted and what was not.
func styleCensus(frame *image.RGBA, area image.Rectangle) map[string]bool {
	census := map[string]bool{}
	styles := roadStyleColors()
	for name := range styles {
		census[name] = false
	}
	for y := area.Min.Y; y < area.Max.Y; y++ {
		for x := area.Min.X; x < area.Max.X; x++ {
			pixel := frame.RGBAAt(x, y)
			for name, want := range styles {
				if matchesColor(pixel, want) {
					census[name] = true
				}
			}
		}
	}

	return census
}

func matchesColor(pixel color.RGBA, want color.NRGBA) bool {
	return withinTolerance(pixel.R, want.R) &&
		withinTolerance(pixel.G, want.G) &&
		withinTolerance(pixel.B, want.B)
}

func withinTolerance(actual uint8, want uint8) bool {
	if actual > want {
		return int(actual)-int(want) <= styleColorTolerance
	}

	return int(want)-int(actual) <= styleColorTolerance
}

// editorEdgeArea is the part of the editor canvas the connections are drawn in,
// with the road key at the bottom cut away.
func editorEdgeArea(zoneEditor *integration_common.ZoneEditorHandler) image.Rectangle {
	canvas := zoneEditor.CanvasRect()
	return image.Rect(canvas.Min.X, canvas.Min.Y, canvas.Max.X, canvas.Max.Y-editorLegendStripPx)
}

// canvasEdge finds the laid-out curve between two zones, whichever way round
// the connection records its endpoints.
func canvasEdge(
	t *testing.T,
	zoneEditor *integration_common.ZoneEditorHandler,
	from string,
	to string) dialogs.EdgeGeometry {
	t.Helper()
	for _, edge := range zoneEditor.Dialog().EdgeGeometries() {
		if (edge.From == from && edge.To == to) || (edge.From == to && edge.To == from) {
			return edge
		}
	}
	t.Fatalf("the canvas is not showing an edge between %q and %q", from, to)

	return dialogs.EdgeGeometry{}
}

// edgeSamplePoints walks the middle of an edge's curve, staying clear of the
// zone rings drawn over its ends and of the guard label above its midpoint.
func edgeSamplePoints(
	zoneEditor *integration_common.ZoneEditorHandler,
	edge dialogs.EdgeGeometry) []image.Point {
	points := make([]image.Point, 0, 41)
	for step := 30; step <= 70; step++ {
		position := float64(step) / 100.0
		remainder := 1 - position
		curvePoint := data.NewVec2(
			remainder*remainder*float64(edge.StartPoint.X)+
				2*remainder*position*float64(edge.ControlPoint.X)+
				position*position*float64(edge.EndPoint.X),
			remainder*remainder*float64(edge.StartPoint.Y)+
				2*remainder*position*float64(edge.ControlPoint.Y)+
				position*position*float64(edge.EndPoint.Y))
		window := zoneEditor.CanvasPoint(curvePoint)
		points = append(points, image.Pt(int(window.X), int(window.Y)))
	}

	return points
}

// edgeStyleCensus reports which road styles are painted along one edge, in a
// narrow band around its curve rather than across the whole canvas.
func edgeStyleCensus(
	frame *image.RGBA,
	zoneEditor *integration_common.ZoneEditorHandler,
	edge dialogs.EdgeGeometry) map[string]bool {
	census := map[string]bool{}
	for name := range roadStyleColors() {
		census[name] = false
	}
	for _, point := range edgeSamplePoints(zoneEditor, edge) {
		band := image.Rect(point.X-4, point.Y, point.X+5, point.Y+1)
		for name, present := range styleCensus(frame, band) {
			census[name] = census[name] || present
		}
	}

	return census
}

// edgeStrokeWidth measures how wide the canvas paints an edge, by counting the
// run of non-backdrop pixels across the curve at every sample point and taking
// the width that occurs most often.
func edgeStrokeWidth(
	frame *image.RGBA,
	zoneEditor *integration_common.ZoneEditorHandler,
	edge dialogs.EdgeGeometry) int {
	widths := map[int]int{}
	for _, point := range edgeSamplePoints(zoneEditor, edge) {
		painted := 0
		for offset := -6; offset <= 6; offset++ {
			if !matchesColor(frame.RGBAAt(point.X+offset, point.Y), themes.ColorsPreview.Background) {
				painted++
			}
		}
		widths[painted]++
	}

	widest := 0
	for width, count := range widths {
		if count > widths[widest] || (count == widths[widest] && width > widest) {
			widest = width
		}
	}

	return widest
}

// previewWithTopology renders the preview panel for a topology with the road
// setting where the test needs it, without opening the zone editor.
func previewWithTopology(
	t *testing.T,
	topology string,
	generateRoads bool) *integration_common.AppRunner {
	t.Helper()
	runner := integration_common.NewAppRunner(t)
	tab := integration_common.NewHandler(runner).
		WithFixtureDirectory().
		ClickLayoutAndZonesTab().
		SelectTopology(topology)
	if !generateRoads {
		tab.ToggleGenerateRoads()
	}
	require.Equal(t, generateRoads, runner.CurrentState().GenerateRoads,
		"the road setting was not put where the test needs it")

	return runner
}

// makeHubToSpawnBARoadlessPortal turns the generated Hub-to-Spawn-B portal into
// a Direct connection and back, which with roads off is the only way to reach a
// portal that carries no road: the checkbox never touches portal flags.
func makeHubToSpawnBARoadlessPortal(
	t *testing.T,
	zoneEditor *integration_common.ZoneEditorHandler) *integration_common.ZoneEditorHandler {
	t.Helper()
	zoneEditor.ClickConnection(hubToSpawnBName).
		SelectConnectionType(directConnectionTypeLabel).
		SelectConnectionType(portalConnectionTypeLabel)
	require.Equal(t, new(false),
		pendingConnection(t, zoneEditor, hubZoneName, spawnBZoneName).Road,
		"the portal was not left roadless")

	return zoneEditor
}

// editorWithRoadlessStyles leaves the canvas showing one roadless direct
// connection and one roadless portal, which is both roadless styles at once.
func editorWithRoadlessStyles(t *testing.T) (
	*integration_common.AppRunner, *integration_common.ZoneEditorHandler) {
	t.Helper()
	runner, zoneEditor := openZoneEditorWithRoads(t, false)
	makeHubToSpawnARoadless(t, zoneEditor)
	makeHubToSpawnBARoadlessPortal(t, zoneEditor)

	return runner, zoneEditor
}

// editorWithRoadedStyles leaves the canvas showing one roaded direct connection
// and the untouched roaded portal.
func editorWithRoadedStyles(t *testing.T) (
	*integration_common.AppRunner, *integration_common.ZoneEditorHandler) {
	t.Helper()
	runner, zoneEditor := openZoneEditorWithRoads(t, true)
	zoneEditor.ClickConnection(hubToSpawnAName).SelectConnectionType(directConnectionTypeLabel)
	require.Equal(t, new(true),
		pendingConnection(t, zoneEditor, hubZoneName, spawnAZoneName).Road,
		"the direct connection did not take a road")

	return runner, zoneEditor
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenTheEditorHoldsRoadlessConnections_TheCanvasPaintsBothRoadlessStyles(t *testing.T) {
	// Arrange
	runner, zoneEditor := editorWithRoadlessStyles(t)

	// Act
	frame := runner.CaptureFrame()

	// Assert
	assert.Equal(t,
		map[string]bool{
			roadedStyle:         false,
			roadlessStyle:       true,
			roadedPortalStyle:   false,
			roadlessPortalStyle: true,
		},
		styleCensus(frame, editorEdgeArea(zoneEditor)))
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenTheEditorHoldsRoadedConnections_TheCanvasKeepsBothRoadedStyles(t *testing.T) {
	// Arrange
	runner, zoneEditor := editorWithRoadedStyles(t)

	// Act
	frame := runner.CaptureFrame()

	// Assert
	assert.Equal(t,
		map[string]bool{
			roadedStyle:         true,
			roadlessStyle:       false,
			roadedPortalStyle:   true,
			roadlessPortalStyle: false,
		},
		styleCensus(frame, editorEdgeArea(zoneEditor)))
}

// A selected edge is the one place the canvas used to override the colour, so
// it has to keep reporting the road state it is in.
//
//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenARoadlessEdgeIsSelected_ItKeepsItsRoadlessColour(t *testing.T) {
	// Arrange
	runner, zoneEditor := openZoneEditorWithRoads(t, false)
	makeHubToSpawnARoadless(t, zoneEditor)
	edge := canvasEdge(t, zoneEditor, hubZoneName, spawnAZoneName)

	// Act
	frame := runner.CaptureFrame()

	// Assert
	assert.Equal(t,
		map[string]bool{
			roadedStyle:         false,
			roadlessStyle:       true,
			roadedPortalStyle:   false,
			roadlessPortalStyle: false,
		},
		edgeStyleCensus(frame, zoneEditor, edge))
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenARoadedEdgeIsSelected_ItKeepsItsRoadColour(t *testing.T) {
	// Arrange
	runner, zoneEditor := editorWithRoadedStyles(t)
	edge := canvasEdge(t, zoneEditor, hubZoneName, spawnAZoneName)

	// Act
	frame := runner.CaptureFrame()

	// Assert
	assert.Equal(t,
		map[string]bool{
			roadedStyle:         true,
			roadlessStyle:       false,
			roadedPortalStyle:   false,
			roadlessPortalStyle: false,
		},
		edgeStyleCensus(frame, zoneEditor, edge))
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenAnEdgeIsNotSelected_ItIsPaintedAtTheNormalWidth(t *testing.T) {
	// Arrange
	runner, zoneEditor := openZoneEditorWithRoads(t, false)
	edge := canvasEdge(t, zoneEditor, hubZoneName, spawnAZoneName)

	// Act
	frame := runner.CaptureFrame()

	// Assert
	assert.Equal(t, 2, edgeStrokeWidth(frame, zoneEditor, edge))
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenAnEdgeIsSelected_ItIsPaintedThicker(t *testing.T) {
	// Arrange
	runner, zoneEditor := openZoneEditorWithRoads(t, false)
	zoneEditor.ClickConnection(hubToSpawnAName)
	edge := canvasEdge(t, zoneEditor, hubZoneName, spawnAZoneName)

	// Act
	frame := runner.CaptureFrame()

	// Assert
	assert.Equal(t, 4, edgeStrokeWidth(frame, zoneEditor, edge))
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenTheSelectionLeavesAnEdge_ItGoesBackToTheNormalWidth(t *testing.T) {
	// Arrange
	runner, zoneEditor := openZoneEditorWithRoads(t, false)
	zoneEditor.ClickConnection(hubToSpawnAName)
	edge := canvasEdge(t, zoneEditor, hubZoneName, spawnAZoneName)
	require.Equal(t, 4, edgeStrokeWidth(runner.CaptureFrame(), zoneEditor, edge),
		"the edge was not selected before the selection was moved away")

	// Act
	zoneEditor.ClickZone(hubZoneName)

	// Assert
	assert.Equal(t, 2, edgeStrokeWidth(runner.CaptureFrame(), zoneEditor, edge))
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenRoadsAreOn_ThePreviewPaintsItsConnectionsRoaded(t *testing.T) {
	// Arrange
	runner := previewWithTopology(t, squareLayoutName, true)

	// Act
	frame := runner.CaptureFrame()

	// Assert
	assert.Equal(t,
		map[string]bool{
			roadedStyle:         true,
			roadlessStyle:       false,
			roadedPortalStyle:   false,
			roadlessPortalStyle: false,
		},
		styleCensus(frame, integration_common.PreviewCanvasRect()))
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenRoadsAreOff_ThePreviewPaintsItsConnectionsRoadless(t *testing.T) {
	// Arrange
	runner := previewWithTopology(t, squareLayoutName, false)

	// Act
	frame := runner.CaptureFrame()

	// Assert
	assert.Equal(t,
		map[string]bool{
			roadedStyle:         false,
			roadlessStyle:       true,
			roadedPortalStyle:   false,
			roadlessPortalStyle: false,
		},
		styleCensus(frame, integration_common.PreviewCanvasRect()))
}

// Turning the setting back on regenerates the template, and no edge may be left
// wearing the style it had while the setting was off.
//
//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenRoadsAreTurnedBackOn_ThePreviewKeepsNoRoadlessStyle(t *testing.T) {
	// Arrange
	runner := previewWithTopology(t, squareLayoutName, false)
	require.True(t,
		styleCensus(runner.CaptureFrame(), integration_common.PreviewCanvasRect())[roadlessStyle],
		"the preview was not roadless before the setting was turned back on")

	// Act
	integration_common.NewHandler(runner).ClickLayoutAndZonesTab().ToggleGenerateRoads()

	// Assert
	assert.Equal(t,
		map[string]bool{
			roadedStyle:         true,
			roadlessStyle:       false,
			roadedPortalStyle:   false,
			roadlessPortalStyle: false,
		},
		styleCensus(runner.CaptureFrame(), integration_common.PreviewCanvasRect()))
}

// Portals are exempt from the setting, so the preview keeps painting them the
// colour they always had even with roads off.
//
//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenRoadsAreOffAndTheConnectionsArePortals_ThePreviewKeepsThePortalColour(t *testing.T) {
	// Arrange
	runner := previewWithTopology(t, geometricHubLayout, false)

	// Act
	frame := runner.CaptureFrame()

	// Assert
	assert.Equal(t,
		map[string]bool{
			roadedStyle:         false,
			roadlessStyle:       false,
			roadedPortalStyle:   true,
			roadlessPortalStyle: false,
		},
		styleCensus(frame, integration_common.PreviewCanvasRect()))
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenARoadlessPortalIsApplied_ThePreviewPaintsItWithoutARoad(t *testing.T) {
	// Arrange
	runner, zoneEditor := editorWithRoadlessStyles(t)

	// Act
	zoneEditor.ClickApply()

	// Assert
	assert.Equal(t,
		map[string]bool{
			roadedStyle:         false,
			roadlessStyle:       true,
			roadedPortalStyle:   false,
			roadlessPortalStyle: true,
		},
		styleCensus(runner.CaptureFrame(), integration_common.PreviewCanvasRect()))
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenThePreviewPanelIsDrawn_ItShowsTheRoadKey(t *testing.T) {
	// Arrange
	runner := previewWithTopology(t, squareLayoutName, true)

	// Act
	frame := runner.CaptureFrame()

	// Assert
	assert.Equal(t,
		map[string]bool{
			roadedStyle:         true,
			roadlessStyle:       true,
			roadedPortalStyle:   true,
			roadlessPortalStyle: true,
		},
		styleCensus(frame, integration_common.PreviewLegendRect()))
}
