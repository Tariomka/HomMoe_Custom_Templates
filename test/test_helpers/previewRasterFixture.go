package test_helpers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/preview"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
)

// Names of the raster fixtures returned by NewPreviewRasterFixtures.
//
// The fitted chord of a straight connector is the center distance minus one zone
// radius per end (2 * 21 px). A solid connector is subdivided into 24 samples and
// a portal one into 96, so a fitted dominant-axis chord below 24 (solid) or 96
// (portal) leaves every sub-segment shorter than one pixel. The "below"/"above"
// fixtures straddle those two boundaries by a single pixel-scale step, staying
// clear of floating point noise at the exact boundary.
const (
	PreviewFixturePortalHorizontalShort          = "portal_horizontal_short"
	PreviewFixtureDirectHorizontalShort          = "direct_horizontal_short"
	PreviewFixturePortalHorizontalBelowThreshold = "portal_horizontal_below_threshold"
	PreviewFixturePortalHorizontalAboveThreshold = "portal_horizontal_above_threshold"
	PreviewFixtureDirectHorizontalBelowThreshold = "direct_horizontal_below_threshold"
	PreviewFixtureDirectHorizontalAboveThreshold = "direct_horizontal_above_threshold"
	PreviewFixturePortalVerticalShort            = "portal_vertical_short"
	PreviewFixturePortalVerticalLong             = "portal_vertical_long"
	PreviewFixtureDirectVerticalBelowThreshold   = "direct_vertical_below_threshold"
	PreviewFixtureDirectVerticalAboveThreshold   = "direct_vertical_above_threshold"
	PreviewFixturePortalDiagonalShort            = "portal_diagonal_short"
	PreviewFixturePortalDiagonalLong             = "portal_diagonal_long"
	PreviewFixtureDirectDiagonalBelowThreshold   = "direct_diagonal_below_threshold"
	PreviewFixtureDirectDiagonalAboveThreshold   = "direct_diagonal_above_threshold"
	PreviewFixturePortalReversedShort            = "portal_reversed_short"
	PreviewFixturePortalCurvedShort              = "portal_curved_short"
	PreviewFixtureDirectCurvedShort              = "direct_curved_short"
	PreviewFixtureDirectCoincidentEnds           = "direct_coincident_ends"
	PreviewFixtureDirectStartOnControlPoint      = "direct_start_on_control_point"
	PreviewFixturePortalNearBorder               = "portal_near_border"
	PreviewFixtureDirectNearBorder               = "direct_near_border"
	PreviewFixtureDirectOffCanvas                = "direct_off_canvas"
)

const (
	// previewFixtureZoneRadius matches the generator's asset radius, so every fixture
	// renders at asset scale 1 and its geometry stays predictable in canvas pixels.
	previewFixtureZoneRadius = 21.0

	previewFixtureStartZone = "Start"
	previewFixtureEndZone   = "End"
)

// PreviewRasterFixture is one fixed preview geometry. It implements
// preview_service.IPreviewLayoutService, so raster tests can push an exact layout
// through the public CreatePreviewImage API without computing a real layout.
type PreviewRasterFixture struct {
	Name   string
	Layout preview.Layout
}

// BuildPreviewLayout returns the fixed layout, ignoring the template and topology.
func (this PreviewRasterFixture) BuildPreviewLayout(
	_ *template_model.Template, _ config.MapTopology, _ float64) preview.Layout {
	return this.Layout
}

// WithoutConnections returns the same fixture with its connections removed, so a
// render can be compared against an otherwise identical edge-free canvas.
func (this PreviewRasterFixture) WithoutConnections() PreviewRasterFixture {
	this.Layout.Connections = nil
	return this
}

// NewPreviewRasterFixtures returns every connector raster fixture shared by the
// unit tests and the PNG review capture.
func NewPreviewRasterFixtures() []PreviewRasterFixture {
	point := data.NewVec2[float64]
	portal, direct := preview.ConnectionTypePortal, preview.ConnectionTypeDirect

	return []PreviewRasterFixture{
		// Reported reproducer: a 58 px fitted chord, invisible for a portal connector.
		newStraightFixture(PreviewFixturePortalHorizontalShort, point(300, 350), point(400, 350), portal),
		newStraightFixture(PreviewFixtureDirectHorizontalShort, point(300, 350), point(400, 350), direct),
		newStraightFixture(PreviewFixturePortalHorizontalBelowThreshold, point(300, 350), point(437, 350), portal),
		newStraightFixture(PreviewFixturePortalHorizontalAboveThreshold, point(300, 350), point(439, 350), portal),
		newStraightFixture(PreviewFixtureDirectHorizontalBelowThreshold, point(300, 350), point(365, 350), direct),
		newStraightFixture(PreviewFixtureDirectHorizontalAboveThreshold, point(300, 350), point(367, 350), direct),
		newStraightFixture(PreviewFixturePortalVerticalShort, point(350, 300), point(350, 400), portal),
		newStraightFixture(PreviewFixturePortalVerticalLong, point(350, 275), point(350, 425), portal),
		newStraightFixture(PreviewFixtureDirectVerticalBelowThreshold, point(350, 300), point(350, 365), direct),
		newStraightFixture(PreviewFixtureDirectVerticalAboveThreshold, point(350, 300), point(350, 367), direct),
		newStraightFixture(PreviewFixturePortalDiagonalShort, point(250, 250), point(350, 350), portal),
		newStraightFixture(PreviewFixturePortalDiagonalLong, point(230, 230), point(360, 360), portal),
		newStraightFixture(PreviewFixtureDirectDiagonalBelowThreshold, point(330, 330), point(380, 380), direct),
		newStraightFixture(PreviewFixtureDirectDiagonalAboveThreshold, point(330, 330), point(385, 385), direct),
		newStraightFixture(PreviewFixturePortalReversedShort, point(400, 350), point(300, 350), portal),
		newCurvedFixture(PreviewFixturePortalCurvedShort, point(300, 350), point(350, 300), point(400, 350), portal),
		newCurvedFixture(PreviewFixtureDirectCurvedShort, point(300, 350), point(350, 300), point(400, 350), direct),
		// Both ends sit exactly one zone radius from the control point, so trimming
		// collapses the whole connector onto that single point.
		newCurvedFixture(PreviewFixtureDirectCoincidentEnds, point(300, 350), point(321, 350), point(342, 350), direct),
		// The start sits on the control point, so it has no direction to be trimmed along.
		newCurvedFixture(PreviewFixtureDirectStartOnControlPoint,
			point(300, 350), point(300, 350), point(400, 350), direct),
		// Near the border the asset fitter is no longer the identity and pulls the
		// zones, the control point and therefore the connector back toward the center.
		newStraightFixture(PreviewFixturePortalNearBorder, point(660, 350), point(560, 350), portal),
		newStraightFixture(PreviewFixtureDirectNearBorder, point(660, 350), point(560, 350), direct),
		newConnectorOnlyFixture(PreviewFixtureDirectOffCanvas, point(-40, 350), point(300, 350), direct),
	}
}

// FindPreviewRasterFixture returns the fixture registered under the given name.
func FindPreviewRasterFixture(name string) (PreviewRasterFixture, bool) {
	for _, fixture := range NewPreviewRasterFixtures() {
		if fixture.Name == name {
			return fixture, true
		}
	}

	return PreviewRasterFixture{}, false
}

// newStraightFixture places the control point on the midpoint, which renders the
// quadratic curve as a straight line.
func newStraightFixture(
	name string, start, end data.Vec2[float64], connectionType preview.ConnectionType) PreviewRasterFixture {
	midpoint := data.NewVec2((start.X+end.X)/2, (start.Y+end.Y)/2)
	return newCurvedFixture(name, start, midpoint, end, connectionType)
}

func newCurvedFixture(
	name string, start, controlPoint, end data.Vec2[float64],
	connectionType preview.ConnectionType) PreviewRasterFixture {
	return PreviewRasterFixture{
		Name: name,
		Layout: preview.Layout{
			Positions: map[string]data.Vec2[float64]{
				previewFixtureStartZone: start,
				previewFixtureEndZone:   end,
			},
			Zones: []preview.Zone{
				{Name: previewFixtureStartZone, Label: previewFixtureStartZone,
					Center: start, Type: preview.ZoneTypeNeutral},
				{Name: previewFixtureEndZone, Label: previewFixtureEndZone,
					Center: end, Type: preview.ZoneTypeNeutral},
			},
			Connections: []preview.Connection{{Start: start, Ctrl: controlPoint, End: end, Type: connectionType}},
			ZoneRadius:  previewFixtureZoneRadius,
		},
	}
}

// newConnectorOnlyFixture keeps the zone list empty, which leaves the asset fitter
// at the identity. That is the only way a layout can hold a connector reaching past
// the canvas edge, because the fitter always pulls zones back inside it.
func newConnectorOnlyFixture(
	name string, start, end data.Vec2[float64], connectionType preview.ConnectionType) PreviewRasterFixture {
	midpoint := data.NewVec2((start.X+end.X)/2, (start.Y+end.Y)/2)
	return PreviewRasterFixture{
		Name: name,
		Layout: preview.Layout{
			Positions:   map[string]data.Vec2[float64]{previewFixtureStartZone: start},
			Connections: []preview.Connection{{Start: start, Ctrl: midpoint, End: end, Type: connectionType}},
			ZoneRadius:  previewFixtureZoneRadius,
		},
	}
}
