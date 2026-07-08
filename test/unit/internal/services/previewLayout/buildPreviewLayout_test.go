package previewLayout_test

import (
	"image"
	"math"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/preview"
	"github.com/Tariomka/hommoe_custom_templates/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ── shared fixture builders ──────────────────────────────────────────

func position(x, y float64) *[2]float64 {
	point := [2]float64{x, y}
	return &point
}

func ringIndex(value int) *int { return &value }

func namedZone(name string) entities.Zone { return entities.Zone{Name: name} }

func positionedZone(name string, x, y float64) entities.Zone {
	zone := namedZone(name)
	zone.GeneratorPosition = position(x, y)
	return zone
}

func ringedZone(name string, ring int, x, y float64) entities.Zone {
	zone := positionedZone(name, x, y)
	zone.GeneratorRing = ringIndex(ring)
	return zone
}

func directConnection(from, to string) entities.Connection {
	return entities.Connection{From: from, To: to, ConnectionType: "Direct"}
}

func templateWith(zones []entities.Zone, connections []entities.Connection) *entities.RmgTemplate {
	return &entities.RmgTemplate{
		Variants: []entities.Variant{{Zones: zones, Connections: connections}},
	}
}

// ── edge cases ───────────────────────────────────────────────────────

func TestWhenTemplateIsNil_ReturnsEmptyLayout(t *testing.T) {
	// Arrange
	expected := services.PreviewLayout{Positions: map[string]image.Point{}}

	// Act
	actual := services.BuildPreviewLayout(nil, config.TopologyRing, 600)

	// Assert
	assert.Equal(t, expected, actual)
}

func TestWhenTemplateHasNoVariants_ReturnsEmptyLayout(t *testing.T) {
	// Arrange
	expected := services.PreviewLayout{Positions: map[string]image.Point{}}

	// Act
	actual := services.BuildPreviewLayout(&entities.RmgTemplate{}, config.TopologyRing, 600)

	// Assert
	assert.Equal(t, expected, actual)
}

func TestWhenVariantHasNoZones_ReturnsEmptyLayout(t *testing.T) {
	// Arrange
	expected := services.PreviewLayout{Positions: map[string]image.Point{}}

	// Act
	actual := services.BuildPreviewLayout(templateWith(nil, nil), config.TopologyRing, 600)

	// Assert
	assert.Equal(t, expected, actual)
}

// ── ring / default dispatch ──────────────────────────────────────────

func TestWhenRingTopologyProvided_PositionsEveryZone(t *testing.T) {
	// Arrange
	zones := []entities.Zone{namedZone("Spawn-A"), namedZone("Spawn-B"), namedZone("Neutral-C")}
	connections := []entities.Connection{
		directConnection("Spawn-A", "Spawn-B"),
		directConnection("Spawn-B", "Neutral-C"),
	}

	// Act
	layout := services.BuildPreviewLayout(templateWith(zones, connections), config.TopologyRing, 600)

	// Assert
	assert.Len(t, layout.Positions, 3)
}

func TestWhenRingTopologyProvided_ComputesPositiveZoneRadius(t *testing.T) {
	// Arrange
	zones := []entities.Zone{namedZone("Spawn-A"), namedZone("Spawn-B"), namedZone("Neutral-C")}

	// Act
	layout := services.BuildPreviewLayout(templateWith(zones, nil), config.TopologyRing, 600)

	// Assert
	assert.Positive(t, layout.ZoneRadius)
}

func TestWhenOnlyOneZoneExists_CentresItOnCanvas(t *testing.T) {
	// Arrange
	zones := []entities.Zone{namedZone("Spawn-A")}

	// Act
	layout := services.BuildPreviewLayout(templateWith(zones, nil), config.TopologyRing, 600)

	// Assert
	assert.Equal(t, image.Pt(300, 300), layout.Positions["Spawn-A"])
}

func TestWhenZoneIsNamedHub_PlacesItAtCanvasCentre(t *testing.T) {
	// Arrange
	zones := []entities.Zone{namedZone("Hub"), namedZone("Spawn-A"), namedZone("Spawn-B")}

	// Act
	layout := services.BuildPreviewLayout(templateWith(zones, nil), config.TopologyRing, 600)

	// Assert
	assert.Equal(t, image.Pt(300, 300), layout.Positions["Hub"])
}

// ── implicit hub rejection ───────────────────────────────────────────

func TestWhenNeutralTouchesEverySpawn_DoesNotCentreIt(t *testing.T) {
	// Arrange
	zones := []entities.Zone{
		namedZone("Spawn-A"), namedZone("Spawn-B"), namedZone("Spawn-C"), namedZone("Neutral-H"),
	}
	connections := []entities.Connection{
		directConnection("Neutral-H", "Spawn-A"),
		directConnection("Neutral-H", "Spawn-B"),
		directConnection("Neutral-H", "Spawn-C"),
	}

	// Act
	layout := services.BuildPreviewLayout(templateWith(zones, connections), config.TopologyRing, 600)

	// Assert
	assert.NotEqual(t, image.Pt(300, 300), layout.Positions["Neutral-H"])
}

func TestWhenNeutralTouchesEverySpawn_DoesNotFlagItAsHub(t *testing.T) {
	// Arrange
	zones := []entities.Zone{
		namedZone("Spawn-A"), namedZone("Spawn-B"), namedZone("Spawn-C"), namedZone("Neutral-H"),
	}
	connections := []entities.Connection{
		directConnection("Neutral-H", "Spawn-A"),
		directConnection("Neutral-H", "Spawn-B"),
		directConnection("Neutral-H", "Spawn-C"),
	}

	// Act
	layout := services.BuildPreviewLayout(templateWith(zones, connections), config.TopologyRing, 600)

	// Assert
	flaggedHubs := []string{}
	for _, previewZone := range layout.Zones {
		if previewZone.IsHub {
			flaggedHubs = append(flaggedHubs, previewZone.Name)
		}
	}
	assert.Empty(t, flaggedHubs)
}

func TestWhenNeutralOnlyConnectsTwoSpawns_FlagsNoHub(t *testing.T) {
	// Arrange
	zones := []entities.Zone{namedZone("Spawn-A"), namedZone("Spawn-B"), namedZone("Neutral-H")}
	connections := []entities.Connection{
		directConnection("Neutral-H", "Spawn-A"),
		directConnection("Neutral-H", "Spawn-B"),
	}

	// Act
	layout := services.BuildPreviewLayout(templateWith(zones, connections), config.TopologyRing, 600)

	// Assert
	flaggedHubs := []string{}
	for _, previewZone := range layout.Zones {
		if previewZone.IsHub {
			flaggedHubs = append(flaggedHubs, previewZone.Name)
		}
	}
	assert.Empty(t, flaggedHubs)
}

func TestWhenZoneIsExplicitlyNamedHub_FlagsOnlyThatZoneAsHub(t *testing.T) {
	// Arrange
	zones := []entities.Zone{namedZone("Hub"), namedZone("Spawn-A"), namedZone("Spawn-B")}
	connections := []entities.Connection{
		directConnection("Hub", "Spawn-A"),
		directConnection("Hub", "Spawn-B"),
	}

	// Act
	layout := services.BuildPreviewLayout(templateWith(zones, connections), config.TopologyRing, 600)

	// Assert
	flaggedHubs := []string{}
	for _, previewZone := range layout.Zones {
		if previewZone.IsHub {
			flaggedHubs = append(flaggedHubs, previewZone.Name)
		}
	}
	assert.Equal(t, []string{"Hub"}, flaggedHubs)
}

func TestWhenTwoHubZonesExist_PositionsEveryZone(t *testing.T) {
	// Arrange
	zones := []entities.Zone{
		namedZone("Hub-A"), namedZone("Hub-B"), namedZone("Spawn-A"), namedZone("Spawn-B"),
	}
	connections := []entities.Connection{
		directConnection("Hub-A", "Spawn-A"),
		directConnection("Hub-B", "Spawn-B"),
		directConnection("Hub-A", "Hub-B"),
	}

	// Act
	layout := services.BuildPreviewLayout(templateWith(zones, connections), config.TopologyRing, 600)

	// Assert
	assert.Len(t, layout.Positions, 4)
}

// ── scatter (Random) dispatch ────────────────────────────────────────

func TestWhenRandomTopologyZonesHavePositions_PositionsEveryZone(t *testing.T) {
	// Arrange
	zones := []entities.Zone{
		positionedZone("Spawn-A", 0.2, 0.2),
		positionedZone("Spawn-B", 0.8, 0.8),
		positionedZone("Neutral-C", 0.5, 0.5),
	}
	connections := []entities.Connection{
		directConnection("Spawn-A", "Neutral-C"),
		directConnection("Neutral-C", "Spawn-B"),
	}

	// Act
	layout := services.BuildPreviewLayout(templateWith(zones, connections), config.TopologyRandom, 600)

	// Assert
	assert.Len(t, layout.Positions, 3)
}

func TestWhenRandomTopologyHasNoConnections_PositionsEveryZone(t *testing.T) {
	// Arrange
	zones := []entities.Zone{
		positionedZone("Spawn-A", 0.1, 0.1),
		positionedZone("Spawn-B", 0.9, 0.9),
		positionedZone("Neutral-C", 0.5, 0.5),
	}

	// Act
	layout := services.BuildPreviewLayout(templateWith(zones, nil), config.TopologyRandom, 600)

	// Assert
	assert.Len(t, layout.Positions, 3)
}

func TestWhenRandomTopologyZonesLackPositions_FallsBackToRingLayout(t *testing.T) {
	// Arrange
	zones := []entities.Zone{namedZone("Spawn-A"), namedZone("Spawn-B"), namedZone("Neutral-C")}

	// Act
	layout := services.BuildPreviewLayout(templateWith(zones, nil), config.TopologyRandom, 600)

	// Assert
	assert.Len(t, layout.Positions, 3)
}

// ── circles ring dispatch ────────────────────────────────────────────

func TestWhenCirclesZonesSpanMultipleRings_PositionsEveryZone(t *testing.T) {
	// Arrange
	zones := []entities.Zone{
		ringedZone("Spawn-A", 0, 0.1, 0.1),
		ringedZone("Spawn-B", 0, 0.9, 0.1),
		ringedZone("Neutral-C", 1, 0.5, 0.5),
	}

	// Act
	layout := services.BuildPreviewLayout(templateWith(zones, nil), config.TopologyCircles, 600)

	// Assert
	assert.Len(t, layout.Positions, 3)
}

func TestWhenCirclesZonesShareOneRing_PositionsEveryZone(t *testing.T) {
	// Arrange
	zones := []entities.Zone{
		ringedZone("Spawn-A", 0, 0.2, 0.2),
		ringedZone("Spawn-B", 0, 0.8, 0.8),
		ringedZone("Spawn-C", 0, 0.5, 0.5),
	}

	// Act
	layout := services.BuildPreviewLayout(templateWith(zones, nil), config.TopologyCircles, 600)

	// Assert
	assert.Len(t, layout.Positions, 3)
}

func TestWhenCirclesTopologyHasOneZone_CentresItOnCanvas(t *testing.T) {
	// Arrange
	zones := []entities.Zone{ringedZone("Spawn-A", 0, 0.5, 0.5)}

	// Act
	layout := services.BuildPreviewLayout(templateWith(zones, nil), config.TopologyCircles, 600)

	// Assert
	assert.Equal(t, image.Pt(300, 300), layout.Positions["Spawn-A"])
}

// ── connection rendering ─────────────────────────────────────────────

func TestWhenDirectConnectionExists_CollectsIt(t *testing.T) {
	// Arrange
	zones := []entities.Zone{namedZone("Spawn-A"), namedZone("Spawn-B")}
	connections := []entities.Connection{directConnection("Spawn-A", "Spawn-B")}

	// Act
	layout := services.BuildPreviewLayout(templateWith(zones, connections), config.TopologyRing, 600)

	// Assert
	assert.Len(t, layout.Connections, 1)
}

func TestWhenDirectConnectionExists_DoesNotFlagItAsPortal(t *testing.T) {
	// Arrange
	zones := []entities.Zone{namedZone("Spawn-A"), namedZone("Spawn-B")}
	connections := []entities.Connection{directConnection("Spawn-A", "Spawn-B")}

	// Act
	layout := services.BuildPreviewLayout(templateWith(zones, connections), config.TopologyRing, 600)

	// Assert
	require.Len(t, layout.Connections, 1)
	assert.False(t, layout.Connections[0].Portal)
}

func TestWhenPortalConnectionExists_FlagsExactlyOnePortal(t *testing.T) {
	// Arrange
	zones := []entities.Zone{namedZone("Spawn-A"), namedZone("Spawn-B"), namedZone("Neutral-C")}
	connections := []entities.Connection{
		directConnection("Spawn-A", "Neutral-C"),
		directConnection("Neutral-C", "Spawn-B"),
		{From: "Spawn-A", To: "Spawn-B", ConnectionType: "Portal"},
	}

	// Act
	layout := services.BuildPreviewLayout(templateWith(zones, connections), config.TopologyRing, 600)

	// Assert
	portalCount := 0
	for _, previewConnection := range layout.Connections {
		if previewConnection.Portal {
			portalCount++
		}
	}
	assert.Equal(t, 1, portalCount)
}

func TestWhenConnectionReferencesUnknownZone_SkipsIt(t *testing.T) {
	// Arrange
	zones := []entities.Zone{namedZone("Spawn-A"), namedZone("Spawn-B")}
	connections := []entities.Connection{directConnection("Spawn-A", "Missing-X")}

	// Act
	layout := services.BuildPreviewLayout(templateWith(zones, connections), config.TopologyRing, 600)

	// Assert
	assert.Empty(t, layout.Connections)
}

// ── two-cluster (tournament) templates ───────────────────────────────

func TestWhenTemplateHasTwoClusters_PositionsEveryZone(t *testing.T) {
	// Arrange
	zones := []entities.Zone{
		namedZone("Spawn-A"), namedZone("Neutral-X"),
		namedZone("Spawn-B"), namedZone("Neutral-Y"),
	}
	connections := []entities.Connection{
		directConnection("Spawn-A", "Neutral-X"),
		directConnection("Spawn-B", "Neutral-Y"),
	}

	// Act
	layout := services.BuildPreviewLayout(templateWith(zones, connections), config.TopologyRing, 600)

	// Assert
	assert.Len(t, layout.Positions, 4)
}

// ── zone classification side-effects ─────────────────────────────────

func TestWhenZoneHasSpawnMainObject_ClassifiesItAsOwnedPlayerZone(t *testing.T) {
	// Arrange
	zones := []entities.Zone{
		{Name: "Spawn-A", MainObjects: []entities.MainObject{{Type: "Spawn", Spawn: "Player1"}}},
	}
	expected := preview.PreviewZone{
		Name:      "Spawn-A",
		Letter:    "A",
		Center:    image.Pt(300, 300),
		IsPlayer:  true,
		HasCastle: true,
		Castles:   1,
		Owner:     1,
	}

	// Act
	layout := services.BuildPreviewLayout(templateWith(zones, nil), config.TopologyRing, 600)

	// Assert
	require.Len(t, layout.Zones, 1)
	assert.Equal(t, expected, layout.Zones[0])
}

func TestWhenZoneHasTwoCityMainObjects_CountsTwoCastles(t *testing.T) {
	// Arrange
	zones := []entities.Zone{
		{Name: "Neutral-Z", MainObjects: []entities.MainObject{{Type: "City"}, {Type: "City"}}},
	}

	// Act
	layout := services.BuildPreviewLayout(templateWith(zones, nil), config.TopologyRing, 600)

	// Assert
	require.Len(t, layout.Zones, 1)
	assert.Equal(t, 2, layout.Zones[0].Castles)
}

// ── parallel connection fanning ──────────────────────────────────────

func TestWhenOnlyOneEdgeConnectsAPair_KeepsControlPointOnMidpoint(t *testing.T) {
	// Arrange
	zones := []entities.Zone{
		positionedZone("Spawn-A", 0.2, 0.5),
		positionedZone("Spawn-B", 0.8, 0.5),
	}
	connections := []entities.Connection{directConnection("Spawn-A", "Spawn-B")}

	// Act
	layout := services.BuildPreviewLayout(templateWith(zones, connections), config.TopologyRandom, 600)

	// Assert
	require.Len(t, layout.Connections, 1)
	edge := layout.Connections[0]
	midX := float64(edge.A.X+edge.B.X) / 2.0
	midY := float64(edge.A.Y+edge.B.Y) / 2.0
	distanceFromMidpoint := math.Hypot(float64(edge.Ctrl.X)-midX, float64(edge.Ctrl.Y)-midY)
	assert.LessOrEqual(t, distanceFromMidpoint, 1.5)
}

func TestWhenParallelEdgesConnectSamePair_GivesThemDistinctControlPoints(t *testing.T) {
	// Arrange
	zones := []entities.Zone{
		positionedZone("Spawn-A", 0.2, 0.5),
		positionedZone("Spawn-B", 0.8, 0.5),
	}
	connections := []entities.Connection{
		directConnection("Spawn-A", "Spawn-B"),
		directConnection("Spawn-A", "Spawn-B"),
	}

	// Act
	layout := services.BuildPreviewLayout(templateWith(zones, connections), config.TopologyRandom, 600)

	// Assert
	require.Len(t, layout.Connections, 2)
	assert.NotEqual(t, layout.Connections[0].Ctrl, layout.Connections[1].Ctrl)
}

func TestWhenParallelEdgesConnectSamePair_BulgesThemSymmetricallyAboutMidpoint(t *testing.T) {
	// Arrange
	zones := []entities.Zone{
		positionedZone("Spawn-A", 0.2, 0.5),
		positionedZone("Spawn-B", 0.8, 0.5),
	}
	connections := []entities.Connection{
		directConnection("Spawn-A", "Spawn-B"),
		directConnection("Spawn-A", "Spawn-B"),
	}

	// Act
	layout := services.BuildPreviewLayout(templateWith(zones, connections), config.TopologyRandom, 600)

	// Assert
	require.Len(t, layout.Connections, 2)
	first, second := layout.Connections[0], layout.Connections[1]
	midX := float64(first.A.X+first.B.X) / 2.0
	midY := float64(first.A.Y+first.B.Y) / 2.0
	averageX := float64(first.Ctrl.X+second.Ctrl.X) / 2.0
	averageY := float64(first.Ctrl.Y+second.Ctrl.Y) / 2.0
	averageDistanceFromMidpoint := math.Hypot(averageX-midX, averageY-midY)
	assert.LessOrEqual(t, averageDistanceFromMidpoint, 1.5)
}

// ── fixed-geometry (Square/Geometric/Cross/Fractal) dispatch ─────────

func TestWhenFixedGeometryTopologyHasOneZone_CentresItOnCanvas(t *testing.T) {
	// Arrange
	zones := []entities.Zone{positionedZone("Spawn-A", 0.3, 0.7)}

	// Act
	layout := services.BuildPreviewLayout(templateWith(zones, nil), config.TopologySquare, 600)

	// Assert
	assert.Equal(t, image.Pt(300, 300), layout.Positions["Spawn-A"])
}

func TestWhenFixedGeometryZonesShareOnePosition_PositionsBothAtSamePoint(t *testing.T) {
	// Arrange
	zones := []entities.Zone{
		positionedZone("Spawn-A", 0.5, 0.5),
		positionedZone("Spawn-B", 0.5, 0.5),
	}

	// Act
	layout := services.BuildPreviewLayout(templateWith(zones, nil), config.TopologySquare, 600)

	// Assert
	assert.Equal(t, layout.Positions["Spawn-A"], layout.Positions["Spawn-B"])
}

func TestWhenFixedGeometryZonesLackPositions_FallsBackToRingLayout(t *testing.T) {
	// Arrange
	zones := []entities.Zone{namedZone("Spawn-A"), namedZone("Spawn-B"), namedZone("Neutral-C")}

	// Act
	layout := services.BuildPreviewLayout(templateWith(zones, nil), config.TopologySquare, 600)

	// Assert
	assert.Len(t, layout.Positions, 3)
}

// ── circles ring edge cases ──────────────────────────────────────────

func TestWhenCirclesZoneLacksRingStamp_StillPositionsEveryZone(t *testing.T) {
	// Arrange
	missingRingZone := positionedZone("Neutral-C", 0.5, 0.5)
	zones := []entities.Zone{
		ringedZone("Spawn-A", 0, 0.1, 0.1),
		ringedZone("Spawn-B", 0, 0.9, 0.1),
		missingRingZone,
	}

	// Act
	layout := services.BuildPreviewLayout(templateWith(zones, nil), config.TopologyCircles, 600)

	// Assert
	assert.Len(t, layout.Positions, 3)
}

func TestWhenCirclesOuterRingHasSingleZone_PositionsEveryZone(t *testing.T) {
	// Arrange
	zones := []entities.Zone{
		ringedZone("Neutral-C", 1, 0.4, 0.4),
		ringedZone("Neutral-D", 1, 0.6, 0.4),
		ringedZone("Neutral-E", 1, 0.5, 0.6),
		ringedZone("Spawn-A", 0, 0.5, 0.9),
	}

	// Act
	layout := services.BuildPreviewLayout(templateWith(zones, nil), config.TopologyCircles, 600)

	// Assert
	assert.Len(t, layout.Positions, 4)
}

// ── scatter edge cases ───────────────────────────────────────────────

func TestWhenRandomTopologyHasOneZone_CentresItOnCanvas(t *testing.T) {
	// Arrange
	zones := []entities.Zone{positionedZone("Spawn-A", 0.2, 0.8)}

	// Act
	layout := services.BuildPreviewLayout(templateWith(zones, nil), config.TopologyRandom, 600)

	// Assert
	assert.Equal(t, image.Pt(300, 300), layout.Positions["Spawn-A"])
}

func TestWhenConnectedScatterZonesShareOnePosition_PositionsBothZones(t *testing.T) {
	// Arrange
	zones := []entities.Zone{
		positionedZone("Spawn-A", 0.5, 0.5),
		positionedZone("Spawn-B", 0.5, 0.5),
	}
	connections := []entities.Connection{directConnection("Spawn-A", "Spawn-B")}

	// Act
	layout := services.BuildPreviewLayout(templateWith(zones, connections), config.TopologyRandom, 600)

	// Assert
	assert.Len(t, layout.Positions, 2)
}

func TestWhenUnconnectedZoneLiesFarFromTightPair_KeepsEveryZoneInsideCanvas(t *testing.T) {
	// Arrange
	zones := []entities.Zone{
		positionedZone("Spawn-A", 0.0, 0.0),
		positionedZone("Spawn-B", 0.001, 0.0),
		positionedZone("Neutral-C", 1.0, 1.0),
	}
	connections := []entities.Connection{directConnection("Spawn-A", "Spawn-B")}

	// Act
	layout := services.BuildPreviewLayout(templateWith(zones, connections), config.TopologyRandom, 600)

	// Assert
	outOfBounds := []string{}
	for name, point := range layout.Positions {
		if point.X < 0 || point.Y < 0 || point.X > 600 || point.Y > 600 {
			outOfBounds = append(outOfBounds, name)
		}
	}
	assert.Empty(t, outOfBounds)
}

func TestWhenThirdZoneLiesOnConnectionLine_NudgesItOffTheLine(t *testing.T) {
	// Arrange
	zones := []entities.Zone{
		positionedZone("Spawn-A", 0.1, 0.5),
		positionedZone("Spawn-B", 0.9, 0.5),
		positionedZone("Neutral-C", 0.5, 0.5),
	}
	connections := []entities.Connection{directConnection("Spawn-A", "Spawn-B")}

	// Act
	layout := services.BuildPreviewLayout(templateWith(zones, connections), config.TopologyRandom, 600)

	// Assert
	assert.NotEqual(t, layout.Positions["Spawn-A"].Y, layout.Positions["Neutral-C"].Y)
}

// ── manual-position dispatch ─────────────────────────────────────────

func manualZone(name string, x, y float64) entities.Zone {
	zone := namedZone(name)
	zone.ManualPosition = position(x, y)
	return zone
}

func TestWhenAllZonesHaveManualPositions_PlacesThemAtScaledCoordinates(t *testing.T) {
	// Arrange
	zones := []entities.Zone{
		manualZone("Spawn-A", 0.25, 0.5),
		manualZone("Spawn-B", 0.75, 0.5),
	}

	// Act
	layout := services.BuildPreviewLayout(templateWith(zones, nil), config.TopologyRing, 600)

	// Assert
	assert.Equal(t, image.Pt(150, 300), layout.Positions["Spawn-A"])
}

func TestWhenManualZonesCoincide_KeepsControlPointOnSharedPoint(t *testing.T) {
	// Arrange
	zones := []entities.Zone{
		manualZone("Spawn-A", 0.5, 0.5),
		manualZone("Spawn-B", 0.5, 0.5),
	}
	connections := []entities.Connection{directConnection("Spawn-A", "Spawn-B")}

	// Act
	layout := services.BuildPreviewLayout(templateWith(zones, connections), config.TopologyRing, 600)

	// Assert
	require.Len(t, layout.Connections, 1)
	assert.Equal(t, layout.Connections[0].A, layout.Connections[0].Ctrl)
}

// ── further connection edge cases ────────────────────────────────────

func TestWhenConnectionSourceIsUnknownZone_SkipsIt(t *testing.T) {
	// Arrange
	zones := []entities.Zone{namedZone("Spawn-A"), namedZone("Spawn-B")}
	connections := []entities.Connection{directConnection("Missing-X", "Spawn-B")}

	// Act
	layout := services.BuildPreviewLayout(templateWith(zones, connections), config.TopologyRing, 600)

	// Assert
	assert.Empty(t, layout.Connections)
}

// ── multi-hub edge cases ─────────────────────────────────────────────

func TestWhenZoneConnectsToNoHub_PlacesItAtCanvasCentre(t *testing.T) {
	// Arrange
	zones := []entities.Zone{
		namedZone("Hub-A"), namedZone("Hub-B"),
		namedZone("Spawn-A"), namedZone("Neutral-X"),
	}
	connections := []entities.Connection{directConnection("Hub-A", "Spawn-A")}

	// Act
	layout := services.BuildPreviewLayout(templateWith(zones, connections), config.TopologyHubAndSpoke, 600)

	// Assert
	assert.Equal(t, image.Pt(300, 300), layout.Positions["Neutral-X"])
}

// ── zero-angle-zone rotation ─────────────────────────────────────────

func TestWhenZeroAngleZoneIsSet_RotatesThatZoneToFirstRingSlot(t *testing.T) {
	// Arrange
	zones := []entities.Zone{namedZone("Spawn-A"), namedZone("Spawn-B"), namedZone("Neutral-C")}
	rmgTemplate := &entities.RmgTemplate{
		Variants: []entities.Variant{{
			Zones:       zones,
			Orientation: entities.Orientation{ZeroAngleZone: "Spawn-B"},
		}},
	}

	// Act
	layout := services.BuildPreviewLayout(rmgTemplate, config.TopologyRing, 600)

	// Assert
	assert.Equal(t, image.Pt(300, 48), layout.Positions["Spawn-B"])
}
