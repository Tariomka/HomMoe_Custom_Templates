package previewGeneratorService_test

import (
	"image"
	"image/color"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/preview_service"
	zone_services "github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
