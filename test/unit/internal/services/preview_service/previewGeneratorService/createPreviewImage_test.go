package previewGeneratorService_test

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/preview"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/asset_provider"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/preview_service"
	zone_services "github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	// roadlessStrokeAlpha is the eight-bit approximation of the half opacity a
	// roadless connector is composited with, once per edge.
	roadlessStrokeAlpha = 128

	// fixtureZoneRadius keeps every fixed layout at asset scale 1, so the generator's
	// trimming radius and marker scale are the plain asset values.
	fixtureZoneRadius = 21.0
	arenaMarkerScale  = 0.75
	zoneAssetScale    = 1.0

	// previewAnchorZone keeps a fixed layout non-empty without placing any zone
	// artwork on it, which leaves the background and the connectors alone on the canvas.
	previewAnchorZone = "Anchor"
)

// recordedFootprint pins one fixture's connector pixels to a fingerprint captured
// before the roadless stroke work began.
type recordedFootprint struct{ fixture, fingerprint string }

func TestWhenTemplateIsRendered_ReturnsFullSizeCanvas(t *testing.T) {
	t.Parallel()
	// Arrange
	generator := mustNewGenerator(t)

	// Act
	canvas := generator.CreatePreviewImage(ringTemplate(), config.TopologyRing)

	// Assert
	assert.Equal(t, image.Rect(0, 0, 700, 700), canvas.Bounds())
}

func TestWhenTemplateIsNil_ReturnsBackgroundOnlyCanvas(t *testing.T) {
	t.Parallel()
	// Arrange
	generator := mustNewGenerator(t)
	backgroundOnly := generator.CreatePreviewImage(&template_model.Template{}, config.TopologyRing)

	// Act
	canvas := generator.CreatePreviewImage(nil, config.TopologyRing)

	// Assert
	assert.Equal(t, backgroundOnly.Pix, canvas.Pix)
}

func TestWhenTemplateHasZones_DrawsThemOverTheBackground(t *testing.T) {
	t.Parallel()
	// Arrange
	generator := mustNewGenerator(t)
	backgroundOnly := generator.CreatePreviewImage(&template_model.Template{}, config.TopologyRing)

	// Act
	canvas := generator.CreatePreviewImage(ringTemplate(), config.TopologyRing)

	// Assert
	assert.NotEqual(t, backgroundOnly.Pix, canvas.Pix)
}

func TestWhenTemplateHasConnections_DrawsLinesBetweenZones(t *testing.T) {
	t.Parallel()
	// Arrange
	generator := mustNewGenerator(t)
	disconnectedTemplate := ringTemplate()
	disconnectedTemplate.Variants[0].Connections = nil
	withoutConnections := generator.CreatePreviewImage(disconnectedTemplate, config.TopologyRing)

	// Act
	canvas := generator.CreatePreviewImage(ringTemplate(), config.TopologyRing)

	// Assert
	assert.NotEqual(t, withoutConnections.Pix, canvas.Pix)
}

func TestWhenConnectionIsPortal_DrawsDashedLineDifferentFromSolid(t *testing.T) {
	t.Parallel()
	// Arrange
	generator := mustNewGenerator(t)
	solidTemplate := ringTemplate()
	solidRender := generator.CreatePreviewImage(solidTemplate, config.TopologyRing)
	portalTemplate := ringTemplate()
	for index := range portalTemplate.Variants[0].Connections {
		portalTemplate.Variants[0].Connections[index].ConnectionType = "Portal"
	}

	// Act
	canvas := generator.CreatePreviewImage(portalTemplate, config.TopologyRing)

	// Assert
	assert.NotEqual(t, solidRender.Pix, canvas.Pix)
}

func TestWhenConnectionIsGladiatorArena_DrawsArenaMarkerOverTheSolidLine(t *testing.T) {
	t.Parallel()
	// Arrange
	generator := mustNewGenerator(t)
	solidRender := generator.CreatePreviewImage(ringTemplate(), config.TopologyRing)
	arenaTemplate := ringTemplate()
	arenaTemplate.Variants[0].Connections[0].ConnectionType = "GladiatorArena"

	// Act
	canvas := generator.CreatePreviewImage(arenaTemplate, config.TopologyRing)

	// Assert
	assert.NotEqual(t, solidRender.Pix, canvas.Pix)
}

func TestWhenZoneHostsTheArena_DrawsArenaBubbleInsteadOfThePlainOne(t *testing.T) {
	t.Parallel()
	// Arrange
	generator := mustNewGenerator(t)
	plainRender := generator.CreatePreviewImage(ringTemplate(), config.TopologyRing)
	arenaTemplate := ringTemplate()
	arenaTemplate.Variants[0].Zones[1].MainObjects = append(
		arenaTemplate.Variants[0].Zones[1].MainObjects,
		template_model.MainObject{Type: "GladiatorArena"})

	// Act
	canvas := generator.CreatePreviewImage(arenaTemplate, config.TopologyRing)

	// Assert
	assert.NotEqual(t, plainRender.Pix, canvas.Pix)
}

func TestWhenSameTemplateIsRenderedTwice_ProducesIdenticalImages(t *testing.T) {
	t.Parallel()
	// Arrange
	generator := mustNewGenerator(t)
	firstRender := generator.CreatePreviewImage(ringTemplate(), config.TopologyRing)

	// Act
	secondRender := generator.CreatePreviewImage(ringTemplate(), config.TopologyRing)

	// Assert
	assert.Equal(t, firstRender.Pix, secondRender.Pix)
}

// Every case renders a fixed layout twice, once with and once without its single
// connection, so the zones land on identical pixels and the only possible
// difference is the connector itself. A pixel that still carries the connector
// color in the final canvas is a pixel no zone artwork painted over.
func TestWhenConnectorIsRendered_PaintsVisibleConnectorPixels(t *testing.T) {
	t.Parallel()
	testCases := []struct{ name, fixture string }{
		{"WhenPortalIsShortAndHorizontal_PaintsVisibleConnectorPixels",
			test_helpers.PreviewFixturePortalHorizontalShort},
		{"WhenDirectIsShortAndHorizontal_PaintsVisibleConnectorPixels",
			test_helpers.PreviewFixtureDirectHorizontalShort},
		{"WhenPortalChordIsJustBelowTheSubdivisionThreshold_PaintsVisibleConnectorPixels",
			test_helpers.PreviewFixturePortalHorizontalBelowThreshold},
		{"WhenPortalChordIsJustAboveTheSubdivisionThreshold_PaintsVisibleConnectorPixels",
			test_helpers.PreviewFixturePortalHorizontalAboveThreshold},
		{"WhenDirectChordIsJustBelowTheSubdivisionThreshold_PaintsVisibleConnectorPixels",
			test_helpers.PreviewFixtureDirectHorizontalBelowThreshold},
		{"WhenDirectChordIsJustAboveTheSubdivisionThreshold_PaintsVisibleConnectorPixels",
			test_helpers.PreviewFixtureDirectHorizontalAboveThreshold},
		{"WhenPortalIsShortAndVertical_PaintsVisibleConnectorPixels",
			test_helpers.PreviewFixturePortalVerticalShort},
		{"WhenPortalIsLongAndVertical_PaintsVisibleConnectorPixels",
			test_helpers.PreviewFixturePortalVerticalLong},
		{"WhenDirectIsVerticalBelowTheSubdivisionThreshold_PaintsVisibleConnectorPixels",
			test_helpers.PreviewFixtureDirectVerticalBelowThreshold},
		{"WhenDirectIsVerticalAboveTheSubdivisionThreshold_PaintsVisibleConnectorPixels",
			test_helpers.PreviewFixtureDirectVerticalAboveThreshold},
		{"WhenPortalIsShortAndDiagonal_PaintsVisibleConnectorPixels",
			test_helpers.PreviewFixturePortalDiagonalShort},
		{"WhenPortalIsLongAndDiagonal_PaintsVisibleConnectorPixels",
			test_helpers.PreviewFixturePortalDiagonalLong},
		{"WhenDirectIsDiagonalBelowTheSubdivisionThreshold_PaintsVisibleConnectorPixels",
			test_helpers.PreviewFixtureDirectDiagonalBelowThreshold},
		{"WhenDirectIsDiagonalAboveTheSubdivisionThreshold_PaintsVisibleConnectorPixels",
			test_helpers.PreviewFixtureDirectDiagonalAboveThreshold},
		{"WhenShortPortalRunsInTheReverseDirection_PaintsVisibleConnectorPixels",
			test_helpers.PreviewFixturePortalReversedShort},
		{"WhenShortPortalIsCurved_PaintsVisibleConnectorPixels",
			test_helpers.PreviewFixturePortalCurvedShort},
		{"WhenShortDirectIsCurved_PaintsVisibleConnectorPixels",
			test_helpers.PreviewFixtureDirectCurvedShort},
		{"WhenPortalSitsNextToTheBorderAndIsFitted_PaintsVisibleConnectorPixels",
			test_helpers.PreviewFixturePortalNearBorder},
		{"WhenDirectSitsNextToTheBorderAndIsFitted_PaintsVisibleConnectorPixels",
			test_helpers.PreviewFixtureDirectNearBorder},
		{"WhenDirectReachesPastTheCanvasEdge_PaintsVisibleConnectorPixels",
			test_helpers.PreviewFixtureDirectOffCanvas},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			// Arrange
			fixture := mustFixture(t, testCase.fixture)
			withoutConnector := renderFixture(t, fixture.WithoutConnections())

			// Act
			withConnector := renderFixture(t, fixture)

			// Assert
			assert.Positive(t, countVisibleConnectorPixels(withConnector, withoutConnector))
		})
	}
}

func TestWhenBothConnectorEndsCollapseOntoTheControlPoint_LeavesTheCanvasUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := mustFixture(t, test_helpers.PreviewFixtureDirectCoincidentEnds)
	withoutConnector := renderFixture(t, fixture.WithoutConnections())

	// Act
	withConnector := renderFixture(t, fixture)

	// Assert
	assert.Equal(t, withoutConnector.Pix, withConnector.Pix)
}

func TestWhenTheConnectorStartSitsOnTheControlPoint_LeavesTheCanvasUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := mustFixture(t, test_helpers.PreviewFixtureDirectStartOnControlPoint)
	withoutConnector := renderFixture(t, fixture.WithoutConnections())

	// Act
	withConnector := renderFixture(t, fixture)

	// Assert
	assert.Equal(t, withoutConnector.Pix, withConnector.Pix)
}

func TestWhenPortalIsLongEnoughToDash_KeepsTheGapsBetweenTheDashes(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := mustFixture(t, test_helpers.PreviewFixturePortalHorizontalShort)
	withoutConnector := renderFixture(t, fixture.WithoutConnections())

	// Act
	withConnector := renderFixture(t, fixture)
	painted := visibleConnectorBounds(withConnector, withoutConnector)

	// Assert
	require.False(t, painted.Empty(), "the short portal painted no visible pixel")
	assert.GreaterOrEqual(t, countDashRunsOnRow(withConnector, withoutConnector, 350, painted), 2)
}

func TestWhenPortalIsAboveTheSubdivisionThreshold_KeepsTheGapsBetweenTheDashes(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := mustFixture(t, test_helpers.PreviewFixturePortalHorizontalAboveThreshold)
	withoutConnector := renderFixture(t, fixture.WithoutConnections())
	withConnector := renderFixture(t, fixture)
	span := visibleConnectorBounds(withConnector, withoutConnector)
	require.False(t, span.Empty(), "the dashed connector painted no visible pixel")

	// Act
	runs := countDashRunsOnRow(withConnector, withoutConnector, 350, span)

	// Assert
	assert.GreaterOrEqual(t, runs, 2)
}

// Unfitted, the right-hand zone sits at x=660 and hides the connector only up to
// x=632, so anything painted further right proves the fitter did not pull the
// connector inward with its zones.
func TestWhenTheLayoutIsFittedAwayFromTheBorder_PaintsTheConnectorInsideTheFittedSpan(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := mustFixture(t, test_helpers.PreviewFixtureDirectNearBorder)
	withoutConnector := renderFixture(t, fixture.WithoutConnections())

	// Act
	painted := visibleConnectorBounds(renderFixture(t, fixture), withoutConnector)

	// Assert
	require.False(t, painted.Empty(), "the fitted connector painted no visible pixel")
	assert.Less(t, painted.Max.X, 632)
}

// The connector starts off canvas and runs along y=350, so a correctly clipped
// brush paints the leftmost column of the four-pixel-wide stroke and nothing before it.
func TestWhenConnectorReachesPastTheCanvasEdge_ClipsItToTheCanvas(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := mustFixture(t, test_helpers.PreviewFixtureDirectOffCanvas)
	withoutConnector := renderFixture(t, fixture.WithoutConnections())

	// Act
	painted := visibleConnectorBounds(renderFixture(t, fixture), withoutConnector)

	// Assert
	assert.Equal(t, image.Pt(0, 348), painted.Min)
}

func TestWhenTheSameFixtureIsRenderedTwice_ProducesIdenticalImages(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := mustFixture(t, test_helpers.PreviewFixturePortalHorizontalShort)
	firstRender := renderFixture(t, fixture)

	// Act
	secondRender := renderFixture(t, fixture)

	// Assert
	assert.Equal(t, firstRender.Pix, secondRender.Pix)
}

// The recorded fingerprints cover every pixel a roaded connector owns - the stamp
// positions, the dash gaps, the clipping and the painted color - and were captured
// from this renderer before any roadless work started. Geometry that silently moves,
// a dropped dash gap or a changed brush color all fail here.
func TestWhenRoadedFixtureIsRendered_KeepsItsRecordedConnectorFootprint(t *testing.T) {
	t.Parallel()
	for _, testCase := range recordedConnectorFootprints() {
		t.Run(testCase.fixture+"_KeepsItsRecordedConnectorFootprint", func(t *testing.T) {
			t.Parallel()
			// Arrange
			fixture := mustFixture(t, testCase.fixture)
			withoutConnector := renderFixture(t, fixture.WithoutConnections())

			// Act
			withConnector := renderFixture(t, fixture)

			// Assert
			assert.Equal(t, testCase.fingerprint, connectorFootprintFingerprint(withConnector, withoutConnector))
		})
	}
}

// Unlike the footprint fingerprints above, these cover the whole canvas, so the
// background and the zone artwork are pinned as well. They assume the renderer's
// float math is reproducible, which holds for every platform this project builds on.
func TestWhenRoadedFixtureIsRendered_KeepsItsRecordedCanvas(t *testing.T) {
	t.Parallel()
	testCases := []struct{ fixture, canvasHash string }{
		{test_helpers.PreviewFixtureDirectHorizontalShort,
			"2ffa0f7f52729c071827ad4ab38aed914ec62d8eaf2063200cea7cf0502511ff"},
		{test_helpers.PreviewFixturePortalCurvedShort,
			"9cf81d268821fd61afbdf505f23a7b850549b754729e1c037f965a41bda375f0"},
	}
	for _, testCase := range testCases {
		t.Run(testCase.fixture+"_KeepsItsRecordedCanvas", func(t *testing.T) {
			t.Parallel()
			// Arrange
			fixture := mustFixture(t, testCase.fixture)

			// Act
			canvas := renderFixture(t, fixture)

			// Assert
			assert.Equal(t, testCase.canvasHash, hex.EncodeToString(sliceHash(canvas.Pix)))
		})
	}
}

// A roadless edge reaches the canvas through exactly one half-opacity composite of
// the stamps its roaded twin paints opaquely. Every case compares the whole canvas
// against a composite built from the real background and that edge's own opaque
// footprint, so a dropped stroke, a moved stamp, a lost dash gap and a pixel
// composited twice by overlapping stamps of one edge all fail.
func TestWhenEdgeHasNoRoad_CompositesItsStrokeOnceAtHalfOpacity(t *testing.T) {
	t.Parallel()
	point := data.NewVec2[float64]
	testCases := []struct {
		name       string
		connection preview.Connection
	}{
		{"WhenDirectEdgeIsShort_CompositesItsStrokeOnceAtHalfOpacity",
			roadlessEdge(point(250, 350), point(350, 350), point(450, 350), preview.ConnectionTypeDirect)},
		{"WhenDirectEdgeSpansTheCanvas_CompositesItsStrokeOnceAtHalfOpacity",
			roadlessEdge(point(70, 350), point(350, 350), point(630, 350), preview.ConnectionTypeDirect)},
		{"WhenDirectEdgeIsVertical_CompositesItsStrokeOnceAtHalfOpacity",
			roadlessEdge(point(350, 250), point(350, 350), point(350, 450), preview.ConnectionTypeDirect)},
		{"WhenDirectEdgeIsDiagonal_CompositesItsStrokeOnceAtHalfOpacity",
			roadlessEdge(point(200, 200), point(350, 350), point(500, 500), preview.ConnectionTypeDirect)},
		{"WhenDirectEdgeIsCurved_CompositesItsStrokeOnceAtHalfOpacity",
			roadlessEdge(point(250, 420), point(350, 240), point(450, 420), preview.ConnectionTypeDirect)},
		{"WhenPortalEdgeIsDashed_CompositesItsStrokeOnceAtHalfOpacity",
			explicitPortal(roadlessEdge(
				point(200, 300), point(350, 300), point(500, 300), preview.ConnectionTypePortal))},
		{"WhenPortalEdgeIsCurved_CompositesItsStrokeOnceAtHalfOpacity",
			explicitPortal(roadlessEdge(
				point(250, 430), point(350, 250), point(450, 430), preview.ConnectionTypePortal))},
		// A Direct connection with portal placement rules keeps the portal shape while
		// its road state stays non-portal, so the dashes must be composited too.
		{"WhenEdgeIsPortalShapedWithoutBeingAnExplicitPortal_CompositesItsStrokeOnceAtHalfOpacity",
			roadlessEdge(point(200, 300), point(350, 300), point(500, 300), preview.ConnectionTypePortal)},
		{"WhenEdgeReachesPastTheCanvasEdge_CompositesItsStrokeOnceAtHalfOpacity",
			roadlessEdge(point(-60, 200), point(120, 200), point(300, 200), preview.ConnectionTypeDirect)},
		{"WhenEdgeRunsAlongTheCanvasEdge_CompositesItsStrokeOnceAtHalfOpacity",
			roadlessEdge(point(200, 1), point(350, 1), point(500, 1), preview.ConnectionTypeDirect)},
		{"WhenPortalEdgeReachesPastTheCanvasEdge_CompositesItsStrokeOnceAtHalfOpacity",
			roadlessEdge(point(-60, 200), point(120, 200), point(300, 200), preview.ConnectionTypePortal)},
		{"WhenLongPortalEdgeIsCurved_CompositesItsStrokeOnceAtHalfOpacity",
			roadlessEdge(point(100, 500), point(350, 80), point(600, 500), preview.ConnectionTypePortal)},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			// Arrange
			require.NotEmpty(t, strokeFootprint(t, testCase.connection), "the edge stamped no pixel at all")
			expected := expectedConnectorRender(t, testCase.connection)

			// Act
			canvas := renderConnections(t, testCase.connection)

			// Assert
			assert.Empty(t, describeCanvasMismatch(expected, canvas))
		})
	}
}

func TestWhenTwoRoadlessEdgesCross_CompositesTheSharedPixelsTwice(t *testing.T) {
	t.Parallel()
	point := data.NewVec2[float64]
	horizontal := roadlessEdge(point(200, 350), point(350, 350), point(500, 350), preview.ConnectionTypeDirect)
	vertical := roadlessEdge(point(350, 200), point(350, 350), point(350, 500), preview.ConnectionTypeDirect)
	// Arrange
	require.NotEmpty(t, sharedFootprint(t, horizontal, vertical), "the two edges never cross")
	expected := expectedConnectorRender(t, horizontal, vertical)

	// Act
	canvas := renderConnections(t, horizontal, vertical)

	// Assert
	assert.Empty(t, describeCanvasMismatch(expected, canvas))
}

func TestWhenTheSameRoadlessEdgeIsDrawnTwice_CompositesItsPixelsTwice(t *testing.T) {
	t.Parallel()
	point := data.NewVec2[float64]
	retraced := roadlessEdge(point(200, 350), point(350, 350), point(500, 350), preview.ConnectionTypeDirect)
	// Arrange
	expected := expectedConnectorRender(t, retraced, retraced)

	// Act
	canvas := renderConnections(t, retraced, retraced)

	// Assert
	assert.Empty(t, describeCanvasMismatch(expected, canvas))
}

func TestWhenRoadlessEdgesAreFarApart_CompositesEachStrokeOnce(t *testing.T) {
	t.Parallel()
	point := data.NewVec2[float64]
	topLeft := roadlessEdge(point(80, 120), point(180, 120), point(280, 120), preview.ConnectionTypeDirect)
	bottomRight := roadlessEdge(point(420, 600), point(520, 600), point(620, 600), preview.ConnectionTypeDirect)
	// Arrange
	expected := expectedConnectorRender(t, topLeft, bottomRight)

	// Act
	canvas := renderConnections(t, topLeft, bottomRight)

	// Assert
	assert.Empty(t, describeCanvasMismatch(expected, canvas))
}

// Sequential edges with different spans must each composite only their own pixels.
func TestWhenAShortRoadlessEdgeFollowsALongOne_CompositesEachStrokeOnce(t *testing.T) {
	t.Parallel()
	point := data.NewVec2[float64]
	long := roadlessEdge(point(70, 200), point(350, 200), point(630, 200), preview.ConnectionTypeDirect)
	short := roadlessEdge(point(300, 500), point(330, 500), point(360, 500), preview.ConnectionTypeDirect)
	// Arrange
	expected := expectedConnectorRender(t, long, short)

	// Act
	canvas := renderConnections(t, long, short)

	// Assert
	assert.Empty(t, describeCanvasMismatch(expected, canvas))
}

// The two diagonals share bounds but no stamped pixels. Uncleared mask coverage
// from the first edge would be composited again inside the second edge's bounds.
func TestWhenRoadlessEdgeBoundsOverlapWithoutSharedPixels_CompositesEachStrokeOnce(t *testing.T) {
	t.Parallel()
	point := data.NewVec2[float64]
	lower := roadlessEdge(point(200, 200), point(300, 300), point(400, 400), preview.ConnectionTypeDirect)
	upper := roadlessEdge(point(250, 180), point(350, 280), point(450, 380), preview.ConnectionTypeDirect)
	// Arrange
	require.Empty(t, sharedFootprint(t, lower, upper), "the two edges were expected to stay apart")
	expected := expectedConnectorRender(t, lower, upper)

	// Act
	canvas := renderConnections(t, lower, upper)

	// Assert
	assert.Empty(t, describeCanvasMismatch(expected, canvas))
}

func TestWhenARoadedEdgeCrossesARoadlessOneDrawnBefore_PaintsTheCrossingOpaque(t *testing.T) {
	t.Parallel()
	point := data.NewVec2[float64]
	roadless := roadlessEdge(point(200, 350), point(350, 350), point(500, 350), preview.ConnectionTypeDirect)
	roaded := roadedEdge(point(350, 200), point(350, 350), point(350, 500), preview.ConnectionTypeDirect)
	// Arrange
	require.NotEmpty(t, sharedFootprint(t, roadless, roaded), "the two edges never cross")
	expected := expectedConnectorRender(t, roadless, roaded)

	// Act
	canvas := renderConnections(t, roadless, roaded)

	// Assert
	assert.Empty(t, describeCanvasMismatch(expected, canvas))
}

func TestWhenARoadlessEdgeCrossesARoadedOneDrawnBefore_CompositesOverTheOpaqueStroke(t *testing.T) {
	t.Parallel()
	point := data.NewVec2[float64]
	roaded := roadedEdge(point(200, 350), point(350, 350), point(500, 350), preview.ConnectionTypeDirect)
	roadless := roadlessEdge(point(350, 200), point(350, 350), point(350, 500), preview.ConnectionTypeDirect)
	// Arrange
	require.NotEmpty(t, sharedFootprint(t, roaded, roadless), "the two edges never cross")
	expected := expectedConnectorRender(t, roaded, roadless)

	// Act
	canvas := renderConnections(t, roaded, roadless)

	// Assert
	assert.Empty(t, describeCanvasMismatch(expected, canvas))
}

func TestWhenARoadlessConnectorCollapsesOntoTheControlPoint_LeavesTheCanvasUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := mustFixture(t, test_helpers.PreviewFixtureDirectCoincidentEnds)
	withoutConnector := renderFixture(t, fixture.WithoutConnections())

	// Act
	withConnector := renderFixture(t, fixture.WithoutRoads())

	// Assert
	assert.Empty(t, describeCanvasMismatch(withoutConnector, withConnector))
}

func TestWhenARoadlessConnectorStartsOnTheControlPoint_LeavesTheCanvasUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := mustFixture(t, test_helpers.PreviewFixtureDirectStartOnControlPoint)
	withoutConnector := renderFixture(t, fixture.WithoutConnections())

	// Act
	withConnector := renderFixture(t, fixture.WithoutRoads())

	// Assert
	assert.Empty(t, describeCanvasMismatch(withoutConnector, withConnector))
}

// The arena sprite is stamped by the shared asset provider, so the expected canvas
// draws it with that very provider over the half-opacity stroke. A marker that
// inherits the stroke's transparency cannot match it.
func TestWhenARoadlessEdgeHostsTheArena_StampsTheMarkerThroughTheAssetProvider(t *testing.T) {
	t.Parallel()
	point := data.NewVec2[float64]
	start, controlPoint, end := point(250, 350), point(350, 350), point(450, 350)
	arena := roadlessEdge(start, controlPoint, end, preview.ConnectionTypeGladiatorArena)
	solidTwin := roadlessEdge(start, controlPoint, end, preview.ConnectionTypeDirect)
	// Arrange
	expected := expectedConnectorRender(t, solidTwin)
	drawArenaMarker(t, expected, arena)

	// Act
	canvas := renderConnections(t, arena)

	// Assert
	assert.Empty(t, describeCanvasMismatch(expected, canvas))
}

// Both bubbles sit on the stroke, so they are composited over its half-opacity
// pixels by the shared asset provider exactly as the expected canvas does it.
func TestWhenAnEdgeHasNoRoad_StampsZoneArtworkThroughTheAssetProvider(t *testing.T) {
	t.Parallel()
	point := data.NewVec2[float64]
	connection := roadlessEdge(point(250, 350), point(350, 350), point(450, 350), preview.ConnectionTypeDirect)
	neutralZone := preview.Zone{Name: "Neutral", Label: "Neutral",
		Center: point(300, 350), Type: preview.ZoneTypeNeutral}
	playerZone := preview.Zone{Name: "Player", Label: "Player",
		Center: point(400, 350), Type: preview.ZoneTypePlayer, Owner: 1}
	// Arrange
	expected := expectedConnectorRender(t, connection)
	drawZoneArtwork(t, expected, neutralZone, playerZone)

	// Act
	canvas := renderFixture(t, newZoneFixture(connection, neutralZone, playerZone))

	// Assert
	assert.Empty(t, describeCanvasMismatch(expected, canvas))
}

// mustFixture fails the test immediately when the named raster fixture is missing.
func mustFixture(t *testing.T, name string) test_helpers.PreviewRasterFixture {
	t.Helper()
	fixture, found := test_helpers.FindPreviewRasterFixture(name)
	require.True(t, found, "unknown preview raster fixture %q", name)
	return fixture
}

// renderFixture draws the fixture's fixed layout through the public preview API.
func renderFixture(t *testing.T, fixture test_helpers.PreviewRasterFixture) *image.RGBA {
	t.Helper()
	generator, err := preview_service.NewPreviewGenerator(fixture)
	require.NoError(t, err)
	return generator.CreatePreviewImage(nil, config.TopologyRing)
}

// countVisibleConnectorPixels counts the pixels the connector changed and still
// owns, which excludes everything a zone asset painted or blended over.
func countVisibleConnectorPixels(withConnector, withoutConnector *image.RGBA) int {
	count := 0
	bounds := withConnector.Bounds()
	for pixelY := bounds.Min.Y; pixelY < bounds.Max.Y; pixelY++ {
		for pixelX := bounds.Min.X; pixelX < bounds.Max.X; pixelX++ {
			if isVisibleConnectorPixel(withConnector, withoutConnector, pixelX, pixelY) {
				count++
			}
		}
	}

	return count
}

// visibleConnectorBounds is the smallest rectangle holding every visible connector
// pixel, empty when the connector painted nothing at all.
func visibleConnectorBounds(withConnector, withoutConnector *image.RGBA) image.Rectangle {
	painted := image.Rectangle{}
	bounds := withConnector.Bounds()
	for pixelY := bounds.Min.Y; pixelY < bounds.Max.Y; pixelY++ {
		for pixelX := bounds.Min.X; pixelX < bounds.Max.X; pixelX++ {
			if isVisibleConnectorPixel(withConnector, withoutConnector, pixelX, pixelY) {
				painted = painted.Union(image.Rect(pixelX, pixelY, pixelX+1, pixelY+1))
			}
		}
	}

	return painted
}

// countDashRunsOnRow counts the stretches of visible connector pixels on one row of
// the connector's own span. Only a pixel the connector left exactly as the
// connector-free render painted it closes a stretch, so two stretches can only be
// counted when an untouched gap really separates them.
func countDashRunsOnRow(withConnector, withoutConnector *image.RGBA, row int, span image.Rectangle) int {
	runs, inRun := 0, false
	for pixelX := span.Min.X; pixelX < span.Max.X; pixelX++ {
		switch {
		case isVisibleConnectorPixel(withConnector, withoutConnector, pixelX, row):
			if !inRun {
				runs++
			}
			inRun = true
		case withConnector.RGBAAt(pixelX, row) == withoutConnector.RGBAAt(pixelX, row):
			inRun = false
		}
	}

	return runs
}

func isVisibleConnectorPixel(withConnector, withoutConnector *image.RGBA, pixelX, pixelY int) bool {
	return withConnector.RGBAAt(pixelX, pixelY) == connectorColor() &&
		withConnector.RGBAAt(pixelX, pixelY) != withoutConnector.RGBAAt(pixelX, pixelY)
}

// connectorColor is the opaque brush color the generator paints connectors with.
func connectorColor() color.RGBA {
	return color.RGBA{R: 0x33, G: 0x18, B: 0x18, A: 0xFF}
}

// mustNewGenerator fails the test immediately when the embedded assets cannot load.
func mustNewGenerator(t *testing.T) preview_service.IPreviewGeneratorService {
	t.Helper()
	generator, err := preview_service.NewPreviewGenerator(
		preview_service.NewPreviewLayoutService(zone_services.NewZoneTierService()),
	)
	require.NoError(t, err)
	return generator
}

// ringTemplate builds a small two-player ring template with plain connections.
func ringTemplate() *template_model.Template {
	return &template_model.Template{
		Variants: []template_model.Variant{{
			Zones: []template_model.Zone{
				{Name: "Spawn-A"}, {Name: "Neutral-B"},
				{Name: "Spawn-C"}, {Name: "Neutral-D"},
			},
			Connections: []template_model.Connection{
				{From: "Spawn-A", To: "Neutral-B", ConnectionType: "Direct"},
				{From: "Neutral-B", To: "Spawn-C", ConnectionType: "Direct"},
				{From: "Spawn-C", To: "Neutral-D", ConnectionType: "Direct"},
				{From: "Neutral-D", To: "Spawn-A", ConnectionType: "Direct"},
			},
		}},
	}
}

func roadedEdge(start, ctrl, end data.Vec2[float64], connectionType preview.ConnectionType) preview.Connection {
	return preview.Connection{Start: start, Ctrl: ctrl, End: end, Type: connectionType, HasRoad: true}
}

func roadlessEdge(start, ctrl, end data.Vec2[float64], connectionType preview.ConnectionType) preview.Connection {
	return preview.Connection{Start: start, Ctrl: ctrl, End: end, Type: connectionType}
}

func explicitPortal(connection preview.Connection) preview.Connection {
	connection.ExplicitPortal = true
	return connection
}

// renderConnections draws the given connections on a zone-free canvas, so the only
// thing between the background and the assertion is the connectors themselves.
func renderConnections(t *testing.T, connections ...preview.Connection) *image.RGBA {
	t.Helper()
	return renderFixture(t, test_helpers.PreviewRasterFixture{Layout: preview.Layout{
		Positions:   map[string]data.Vec2[float64]{previewAnchorZone: data.NewVec2(350.0, 350.0)},
		Connections: connections,
		ZoneRadius:  fixtureZoneRadius,
	}})
}

// newZoneFixture places real zone artwork on the canvas next to one connection.
func newZoneFixture(connection preview.Connection, zones ...preview.Zone) test_helpers.PreviewRasterFixture {
	positions := make(map[string]data.Vec2[float64], len(zones))
	for _, zone := range zones {
		positions[zone.Name] = zone.Center
	}

	return test_helpers.PreviewRasterFixture{Layout: preview.Layout{
		Positions:   positions,
		Zones:       zones,
		Connections: []preview.Connection{connection},
		ZoneRadius:  fixtureZoneRadius,
	}}
}

// expectedConnectorRender composes the canvas the generator owes for the given
// connections: each one paints its own stamps once, opaque when it carries a road
// and composited at half opacity when it does not. The stamps come from an opaque
// render of that single connection, so the geometry is never restated here.
func expectedConnectorRender(t *testing.T, connections ...preview.Connection) *image.RGBA {
	t.Helper()
	expected := cloneCanvas(renderConnections(t))
	for _, connection := range connections {
		for _, pixel := range strokeFootprint(t, connection) {
			if connection.HasRoad {
				expected.SetRGBA(pixel.X, pixel.Y, connectorColor())
				continue
			}

			expected.SetRGBA(pixel.X, pixel.Y, halfOpacityOver(expected.RGBAAt(pixel.X, pixel.Y)))
		}
	}

	return expected
}

// strokeFootprint is every canvas pixel one connection stamps, read from an opaque
// render of that connection on its own.
func strokeFootprint(t *testing.T, connection preview.Connection) []image.Point {
	t.Helper()
	connection.HasRoad = true
	return changedPixels(renderConnections(t, connection), renderConnections(t))
}

// sharedFootprint is the pixels two connections both stamp.
func sharedFootprint(t *testing.T, first, second preview.Connection) []image.Point {
	t.Helper()
	secondPixels := make(map[image.Point]bool)
	for _, pixel := range strokeFootprint(t, second) {
		secondPixels[pixel] = true
	}

	shared := []image.Point{}
	for _, pixel := range strokeFootprint(t, first) {
		if secondPixels[pixel] {
			shared = append(shared, pixel)
		}
	}

	return shared
}

// halfOpacityOver composites the connector brush onto one background pixel exactly
// once, through the half-opacity mask a roadless edge is drawn with.
func halfOpacityOver(background color.RGBA) color.RGBA {
	pixel := image.NewRGBA(image.Rect(0, 0, 1, 1))
	pixel.SetRGBA(0, 0, background)
	draw.DrawMask(pixel, pixel.Bounds(),
		image.NewUniform(connectorColor()), image.Point{},
		image.NewUniform(color.Alpha{A: roadlessStrokeAlpha}), image.Point{}, draw.Over)
	return pixel.RGBAAt(0, 0)
}

// drawArenaMarker stamps the marker the generator owes a Gladiator Arena edge, on
// the trimmed curve's midpoint and through the very asset provider the generator uses.
func drawArenaMarker(t *testing.T, canvas *image.RGBA, connection preview.Connection) {
	t.Helper()
	startPoint, startFits := helpers.CalculatePointTowards(connection.Start, connection.Ctrl, fixtureZoneRadius)
	endPoint, endFits := helpers.CalculatePointTowards(connection.End, connection.Ctrl, fixtureZoneRadius)
	require.True(t, startFits && endFits, "the arena connector collapsed before its marker was placed")

	midPoint := helpers.GetVectorOnQuadraticBezierCurve(startPoint, connection.Ctrl, endPoint, 0.5)
	mustAssetProvider(t).DrawArenaMarker(canvas, midPoint, arenaMarkerScale)
}

// drawZoneArtwork stamps both bubbles in the order the generator draws them.
func drawZoneArtwork(t *testing.T, canvas *image.RGBA, neutralZone, playerZone preview.Zone) {
	t.Helper()
	assetProvider := mustAssetProvider(t)
	assetProvider.DrawNeutralZone(canvas, neutralZone, neutralZone.Center, zoneAssetScale)
	assetProvider.DrawPlayerZone(canvas, playerZone, playerZone.Center, zoneAssetScale)
}

func mustAssetProvider(t *testing.T) *asset_provider.AssetProvider {
	t.Helper()
	assetProvider, err := asset_provider.NewAssetProvider()
	require.NoError(t, err)
	return assetProvider
}

func cloneCanvas(canvas *image.RGBA) *image.RGBA {
	clone := image.NewRGBA(canvas.Bounds())
	copy(clone.Pix, canvas.Pix)
	return clone
}

func changedPixels(canvas, reference *image.RGBA) []image.Point {
	changed := []image.Point{}
	bounds := canvas.Bounds()
	for pixelY := bounds.Min.Y; pixelY < bounds.Max.Y; pixelY++ {
		for pixelX := bounds.Min.X; pixelX < bounds.Max.X; pixelX++ {
			if canvas.RGBAAt(pixelX, pixelY) != reference.RGBAAt(pixelX, pixelY) {
				changed = append(changed, image.Pt(pixelX, pixelY))
			}
		}
	}

	return changed
}

// describeCanvasMismatch reports the first pixel two canvases disagree on and how
// many disagree in total. Comparing the pixel slices themselves would dump two
// megabytes of bytes into the failure output.
func describeCanvasMismatch(expected, actual *image.RGBA) string {
	if expected.Bounds() != actual.Bounds() {
		return fmt.Sprintf("canvas bounds are %v, expected %v", actual.Bounds(), expected.Bounds())
	}

	firstMismatch, differingCount := "", 0
	bounds := expected.Bounds()
	for pixelY := bounds.Min.Y; pixelY < bounds.Max.Y; pixelY++ {
		for pixelX := bounds.Min.X; pixelX < bounds.Max.X; pixelX++ {
			expectedPixel, actualPixel := expected.RGBAAt(pixelX, pixelY), actual.RGBAAt(pixelX, pixelY)
			if expectedPixel == actualPixel {
				continue
			}

			differingCount++
			if firstMismatch == "" {
				firstMismatch = fmt.Sprintf("pixel (%d,%d) is %v, expected %v",
					pixelX, pixelY, actualPixel, expectedPixel)
			}
		}
	}

	if differingCount == 0 {
		return ""
	}

	return fmt.Sprintf("%s; %d pixels differ", firstMismatch, differingCount)
}

// connectorFootprintFingerprint hashes the position and the color of every pixel the
// connector owns, leaving the untouched background out of the recorded value.
func connectorFootprintFingerprint(withConnector, withoutConnector *image.RGBA) string {
	hasher := sha256.New()
	position := make([]byte, 8)
	bounds := withConnector.Bounds()
	for pixelY := bounds.Min.Y; pixelY < bounds.Max.Y; pixelY++ {
		for pixelX := bounds.Min.X; pixelX < bounds.Max.X; pixelX++ {
			painted := withConnector.RGBAAt(pixelX, pixelY)
			if painted == withoutConnector.RGBAAt(pixelX, pixelY) {
				continue
			}

			binary.BigEndian.PutUint32(position[0:4], uint32(pixelX))
			binary.BigEndian.PutUint32(position[4:8], uint32(pixelY))
			hasher.Write(position)
			hasher.Write([]byte{painted.R, painted.G, painted.B, painted.A})
		}
	}

	return hex.EncodeToString(hasher.Sum(nil))
}

func sliceHash(pixels []byte) []byte {
	sum := sha256.Sum256(pixels)
	return sum[:]
}

// recordedConnectorFootprints holds the fingerprints captured from this renderer
// before the roadless stroke work started. The two degenerate fixtures are left out
// on purpose: they paint nothing, and their own tests assert exactly that.
func recordedConnectorFootprints() []recordedFootprint {
	return []recordedFootprint{
		{test_helpers.PreviewFixturePortalHorizontalShort,
			"eedddb78e337b5d75f9c44e704f201e65770906380f51b30781762b04e6f0a17"},
		{test_helpers.PreviewFixtureDirectHorizontalShort,
			"4b1d844031c2904e1665a850c635b373fc82307d189a0340d1a2e2b56002c7e0"},
		{test_helpers.PreviewFixturePortalHorizontalBelowThreshold,
			"983b9823ca1ddadce2539be064168fb34684e37f5ef1edbcde86e5dc3b9d3340"},
		{test_helpers.PreviewFixturePortalHorizontalAboveThreshold,
			"7f36b9929822f8178ed10403601e1264c88c1f40d6b2f7af9d371da02795975a"},
		{test_helpers.PreviewFixtureDirectHorizontalBelowThreshold,
			"b71cb3b6e18aa4c6472c0c0ec54f0d58334f9e0691b31d96c64ad1a4f490ab63"},
		{test_helpers.PreviewFixtureDirectHorizontalAboveThreshold,
			"cc1ef9d37c3450d7ba50a88824313ab25ad5c9607cf3701683150b251f31c601"},
		{test_helpers.PreviewFixturePortalVerticalShort,
			"8faa6ef2dc2bf8cafc34beb70122e4cac2775a761e33f7d7f65d03f00a4ce009"},
		{test_helpers.PreviewFixturePortalVerticalLong,
			"ff2a7b4f29012ac3a4dc3b1948a701a24986cf8251ee9c2d492ab7596ff571f3"},
		{test_helpers.PreviewFixtureDirectVerticalBelowThreshold,
			"6aa2587ffd50c55e9da072cf2af22a02a58c78c56fd7515401ca3681144dc75d"},
		{test_helpers.PreviewFixtureDirectVerticalAboveThreshold,
			"3a2c2ac0fb38cd18ec838de9eb0901541835bdda06886fcd4a799e3394802925"},
		{test_helpers.PreviewFixturePortalDiagonalShort,
			"a5144e05a4eb4433a1860cb4e57550a2bbea8ee2edb498c5da4ba8e3464b6dd6"},
		{test_helpers.PreviewFixturePortalDiagonalLong,
			"413611881ef390c376352aa1f42036ba06a712931f97d466b0207853488b968d"},
		{test_helpers.PreviewFixtureDirectDiagonalBelowThreshold,
			"5475dacd5509fb896ce6dae96a46f5c5a2b16eced4f8131d9c0db9705e96c710"},
		{test_helpers.PreviewFixtureDirectDiagonalAboveThreshold,
			"8c61ecf5ab2d95f181b9143f75baf9fa1f009a7aa4f7bce1a9e3fcfa8bfbbe98"},
		{test_helpers.PreviewFixturePortalReversedShort,
			"3dd012ac276f8ca35732601bf87408c609a69e327a3741dac36315d8adcb2694"},
		{test_helpers.PreviewFixturePortalCurvedShort,
			"b675021449d5e6507f4c1f0f51b6aafbea8e1047e5cd62797080bc3f36c6ecb7"},
		{test_helpers.PreviewFixtureDirectCurvedShort,
			"55db9c47777f0d318df79ef988bdc16d1637dfa180e45c8d1d1e6b55e1fd3ca7"},
		{test_helpers.PreviewFixturePortalNearBorder,
			"0506fa47086c51fd7ee33fb63631bc5739bdad18a13d0ef82e778a4252d93be1"},
		{test_helpers.PreviewFixtureDirectNearBorder,
			"4b899612bd096f1bad25d1a8b1366499c76cc52604a2c460af0b2ccce1605ea5"},
		{test_helpers.PreviewFixtureDirectOffCanvas,
			"4a40ecf6c0be0d52dc30217d9bbd0859500aacc64e65607f15104a9ba03d7f5b"},
	}
}
