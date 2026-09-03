package zoneEditorGeometryService_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestWhenTheGeometryIsBuilt_TheNodePositionsAreReportedVerbatim(t *testing.T) {
	t.Parallel()
	// Arrange
	service, _ := newGeometryFixture(chordPositions())

	// Act
	geometry := service.BuildGeometry(nil, nil, config.MapTopology(gofakeit.Word()), fixtureCanvasSide)

	// Assert
	assert.Equal(t, chordPositions(), geometry.Positions)
}

func TestWhenTheGeometryIsBuilt_TheZoneRadiusIsReportedVerbatim(t *testing.T) {
	t.Parallel()
	// Arrange
	service, _ := newGeometryFixture(chordPositions())

	// Act
	geometry := service.BuildGeometry(nil, nil, config.MapTopology(gofakeit.Word()), fixtureCanvasSide)

	// Assert
	assert.InDelta(t, fixtureZoneRadius, geometry.ZoneRadius, 1e-9)
}

func TestWhenTheGeometryIsBuilt_TheCanvasSideAndTopologyReachThePreviewLayout(t *testing.T) {
	t.Parallel()
	// Arrange
	service, previewLayout := newGeometryFixture(chordPositions())
	topology := config.MapTopology(gofakeit.Word())

	// Act
	service.BuildGeometry(nil, nil, topology, fixtureCanvasSide)

	// Assert
	previewLayout.AssertCalled(t, "BuildPreviewLayout", mock.Anything, topology, float64(fixtureCanvasSide))
}

func TestWhenTheGeometryIsBuilt_TheZonesReachThePreviewLayout(t *testing.T) {
	t.Parallel()
	// Arrange
	service, previewLayout := newGeometryFixture(chordPositions())
	zones := []template_model.Zone{{Name: "A"}, {Name: "B"}}

	// Act
	service.BuildGeometry(zones, nil, config.MapTopology(gofakeit.Word()), fixtureCanvasSide)

	// Assert
	template, _ := previewLayout.Calls[0].Arguments.Get(0).(*template_model.Template)
	require.NotNil(t, template)
	require.NotEmpty(t, template.Variants)
	assert.Equal(t, zones, template.Variants[0].Zones)
}

func TestWhenASingleConnectionSpansAClearChord_ItsCurveStaysOnTheChord(t *testing.T) {
	t.Parallel()
	// Arrange
	service, _ := newGeometryFixture(chordPositions())
	connections := []entities.Connection{newConnection("ab", "A", "B")}

	// Act
	geometry := service.BuildGeometry(nil, connections, "", fixtureCanvasSide)

	// Assert
	require.Len(t, geometry.Edges, 1)
	assert.Equal(t, data.NewVec2(350.0, 350.0), geometry.Edges[0].ControlPoint)
}

func TestWhenAnEdgeIsLaidOut_ItsLabelSitsOnTheCurveMidpoint(t *testing.T) {
	t.Parallel()
	// Arrange
	service, _ := newGeometryFixture(chordPositions())
	connections := []entities.Connection{
		newConnection("ab", "A", "B"),
		newConnection("ba", "B", "A"),
	}

	// Act
	geometry := service.BuildGeometry(nil, connections, "", fixtureCanvasSide)

	// Assert
	require.Len(t, geometry.Edges, 2)
	assert.Equal(t, data.NewVec2(350.0, 359.0), geometry.Edges[0].MidPoint)
}

func TestWhenTwoConnectionsSharePair_TheirCurvesSpreadSymmetrically(t *testing.T) {
	t.Parallel()
	// Arrange
	service, _ := newGeometryFixture(chordPositions())
	connections := []entities.Connection{
		newConnection("ab", "A", "B"),
		newConnection("ba", "B", "A"),
	}

	// Act
	geometry := service.BuildGeometry(nil, connections, "", fixtureCanvasSide)

	// Assert
	require.Len(t, geometry.Edges, 2)
	assert.Equal(t,
		[]data.Vec2[float64]{data.NewVec2(350.0, 368.0), data.NewVec2(350.0, 332.0)},
		[]data.Vec2[float64]{geometry.Edges[0].ControlPoint, geometry.Edges[1].ControlPoint})
}

func TestWhenConnectionsSharePairs_TheyAreGroupedInFirstSeenOrder(t *testing.T) {
	t.Parallel()
	// Arrange
	service, _ := newGeometryFixture(trianglePositions())
	connections := []entities.Connection{
		newConnection("ab", "A", "B"),
		newConnection("ac", "A", "C"),
		newConnection("ba", "B", "A"),
	}

	// Act
	geometry := service.BuildGeometry(nil, connections, "", fixtureCanvasSide)

	// Assert
	assert.Equal(t, []int{0, 2, 1}, connectionIndices(geometry.Edges))
}

func TestWhenAZoneSitsNearTheChord_TheCurveBulgesClearOfIt(t *testing.T) {
	t.Parallel()
	// Arrange
	service, _ := newGeometryFixture(obstaclePositions())
	connections := []entities.Connection{newConnection("ab", "A", "B")}

	// Act
	geometry := service.BuildGeometry(nil, connections, "", fixtureCanvasSide)

	// Assert
	require.Len(t, geometry.Edges, 1)
	assert.Equal(t, data.NewVec2(350.0, 274.0), geometry.Edges[0].ControlPoint)
}

func TestWhenAnEndpointHasNoPosition_TheConnectionIsSkipped(t *testing.T) {
	t.Parallel()
	// Arrange
	service, _ := newGeometryFixture(chordPositions())
	connections := []entities.Connection{
		newConnection("ab", "A", "B"),
		newConnection("az", "A", "Z"),
	}

	// Act
	geometry := service.BuildGeometry(nil, connections, "", fixtureCanvasSide)

	// Assert
	assert.Equal(t, []int{0}, connectionIndices(geometry.Edges))
}

func TestWhenBothEndpointsShareAPosition_TheCurveCollapsesOntoThatPoint(t *testing.T) {
	t.Parallel()
	// Arrange
	stacked := map[string]models.Position{
		"A": data.NewVec2(140.0, 350.0),
		"B": data.NewVec2(140.0, 350.0),
	}
	service, _ := newGeometryFixture(stacked)
	connections := []entities.Connection{newConnection("ab", "A", "B")}

	// Act
	geometry := service.BuildGeometry(nil, connections, "", fixtureCanvasSide)

	// Assert
	require.Len(t, geometry.Edges, 1)
	assert.Equal(t, data.NewVec2(140.0, 350.0), geometry.Edges[0].ControlPoint)
}
