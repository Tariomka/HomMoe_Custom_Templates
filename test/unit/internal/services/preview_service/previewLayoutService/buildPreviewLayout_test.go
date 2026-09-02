package previewLayoutService_test

import (
	"math"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/preview"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/preview_service"
	zone_services "github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const layoutSide = 700.0

func TestWhenTemplateIsNil_ReturnsEmptyLayout(t *testing.T) {
	t.Parallel()
	// Arrange
	service := preview_service.NewPreviewLayoutService(zone_services.NewZoneTierService())
	expected := preview.Layout{Positions: map[string]data.Vec2[float64]{}}

	// Act
	actual := service.BuildPreviewLayout(nil, config.TopologyRing, layoutSide)

	// Assert
	assert.Equal(t, expected, actual)
}

func TestWhenTemplateHasNoVariants_ReturnsEmptyLayout(t *testing.T) {
	t.Parallel()
	// Arrange
	service := preview_service.NewPreviewLayoutService(zone_services.NewZoneTierService())
	expected := preview.Layout{Positions: map[string]data.Vec2[float64]{}}

	// Act
	actual := service.BuildPreviewLayout(&entities.RmgTemplate{}, config.TopologyRing, layoutSide)

	// Assert
	assert.Equal(t, expected, actual)
}

func TestWhenVariantHasNoZones_ReturnsEmptyLayout(t *testing.T) {
	t.Parallel()
	// Arrange
	service := preview_service.NewPreviewLayoutService(zone_services.NewZoneTierService())
	expected := preview.Layout{Positions: map[string]data.Vec2[float64]{}}

	// Act
	actual := service.BuildPreviewLayout(templateWith(nil, nil), config.TopologyRing, layoutSide)

	// Assert
	assert.Equal(t, expected, actual)
}

func TestWhenRingTopologyIsLaidOut_EveryZoneStaysInsideTheCanvas(t *testing.T) {
	t.Parallel()
	// Arrange
	service := preview_service.NewPreviewLayoutService(zone_services.NewZoneTierService())
	zones := []entities.Zone{
		namedZone("Spawn-A"), namedZone("Spawn-B"),
		namedZone("Neutral-C"), namedZone("Neutral-D"),
	}
	connections := []entities.Connection{
		directConnection("Spawn-A", "Neutral-C"),
		directConnection("Neutral-C", "Spawn-B"),
		directConnection("Spawn-B", "Neutral-D"),
		directConnection("Neutral-D", "Spawn-A"),
	}

	// Act
	layout := service.BuildPreviewLayout(templateWith(zones, connections), config.TopologyRing, layoutSide)

	// Assert
	require.Len(t, layout.Positions, len(zones))
	for zoneName, zonePosition := range layout.Positions {
		assert.True(t,
			zonePosition.X >= 0 && zonePosition.X <= layoutSide &&
				zonePosition.Y >= 0 && zonePosition.Y <= layoutSide,
			"zone %s at %v escapes the canvas", zoneName, zonePosition)
	}
}

func TestWhenAllZonesHaveManualPositions_PlacesThemVerbatim(t *testing.T) {
	t.Parallel()
	// Arrange
	service := preview_service.NewPreviewLayoutService(zone_services.NewZoneTierService())
	zones := []entities.Zone{
		manualZone("Spawn-A", 0.25, 0.5),
		manualZone("Neutral-B", 0.75, 0.5),
	}
	expected := map[string]data.Vec2[float64]{
		"Spawn-A":   data.NewVec2(175.0, 350.0),
		"Neutral-B": data.NewVec2(525.0, 350.0),
	}

	// Act
	layout := service.BuildPreviewLayout(templateWith(zones, nil), config.TopologyRing, layoutSide)

	// Assert
	assert.Equal(t, expected, layout.Positions)
}

func TestWhenTwoZonesAreLessThanAPixelApart_TheirCentresDiffer(t *testing.T) {
	t.Parallel()
	// Arrange - the two manual positions are 0.3px apart on the canvas, which
	// the old integer layout collapsed onto the same pixel.
	service := preview_service.NewPreviewLayoutService(zone_services.NewZoneTierService())
	zones := []entities.Zone{
		manualZone("Spawn-A", 0.5, 0.5),
		manualZone("Neutral-B", 0.5+0.3/layoutSide, 0.5),
	}

	// Act
	layout := service.BuildPreviewLayout(templateWith(zones, nil), config.TopologyRing, layoutSide)

	// Assert
	assert.NotEqual(t, layout.Positions["Spawn-A"], layout.Positions["Neutral-B"])
}

func TestWhenFixedGeometryTopologyIsLaidOut_PreservesRelativeGeometry(t *testing.T) {
	t.Parallel()
	// Arrange
	service := preview_service.NewPreviewLayoutService(zone_services.NewZoneTierService())
	zones := []entities.Zone{
		positionedZone("Spawn-A", 0.1, 0.5),
		positionedZone("Neutral-B", 0.5, 0.5),
		positionedZone("Spawn-C", 0.9, 0.5),
	}

	// Act
	layout := service.BuildPreviewLayout(templateWith(zones, nil), config.TopologySquare, layoutSide)

	// Assert - the middle zone must stay the exact midpoint of the outer two.
	require.Len(t, layout.Positions, len(zones))
	left := layout.Positions["Spawn-A"]
	middle := layout.Positions["Neutral-B"]
	right := layout.Positions["Spawn-C"]
	assert.Equal(t, left.Add(right).MultiplyScalar(0.5), middle)
}

func TestWhenGeometricHubTopologyIsLaidOut_FigureKeepsExtraBorderClearance(t *testing.T) {
	t.Parallel()
	// Arrange
	service := preview_service.NewPreviewLayoutService(zone_services.NewZoneTierService())
	zones := []entities.Zone{
		positionedZone("Spawn-A", 0.1, 0.5),
		positionedZone("Neutral-B", 0.5, 0.5),
		positionedZone("Spawn-C", 0.9, 0.5),
	}

	// Act
	hubLayout := service.BuildPreviewLayout(
		templateWith(zones, nil), config.TopologyGeometricHub, layoutSide)
	squareLayout := service.BuildPreviewLayout(
		templateWith(zones, nil), config.TopologySquare, layoutSide)

	// Assert - the same figure must span less width under Geometric Hub than
	// under the other fixed-geometry topologies (extra edge inset applied).
	hubWidth := hubLayout.Positions["Spawn-C"].X - hubLayout.Positions["Spawn-A"].X
	squareWidth := squareLayout.Positions["Spawn-C"].X - squareLayout.Positions["Spawn-A"].X
	assert.Less(t, hubWidth, squareWidth)
}

func TestWhenGeometricHubHasSixPlayers_FigureSitsCloserToBorder(t *testing.T) {
	t.Parallel()
	// Arrange
	service := preview_service.NewPreviewLayoutService(zone_services.NewZoneTierService())
	crowdedZones := []entities.Zone{
		positionedZone("Spawn-A", 0.1, 0.5),
		positionedZone("Spawn-B", 0.3, 0.2),
		positionedZone("Spawn-C", 0.5, 0.8),
		positionedZone("Spawn-D", 0.7, 0.2),
		positionedZone("Spawn-E", 0.9, 0.5),
		positionedZone("Spawn-F", 0.5, 0.5),
	}
	sparseZones := []entities.Zone{
		positionedZone("Spawn-A", 0.1, 0.5),
		positionedZone("Neutral-B", 0.3, 0.2),
		positionedZone("Neutral-C", 0.5, 0.8),
		positionedZone("Neutral-D", 0.7, 0.2),
		positionedZone("Spawn-E", 0.9, 0.5),
		positionedZone("Neutral-F", 0.5, 0.5),
	}

	// Act
	crowdedLayout := service.BuildPreviewLayout(
		templateWith(crowdedZones, nil), config.TopologyGeometricHub, layoutSide)
	sparseLayout := service.BuildPreviewLayout(
		templateWith(sparseZones, nil), config.TopologyGeometricHub, layoutSide)

	// Assert - six or more players shrink the edge inset, letting the same
	// figure scale further toward the border (players further from the hub).
	crowdedWidth := crowdedLayout.Positions["Spawn-E"].X - crowdedLayout.Positions["Spawn-A"].X
	sparseWidth := sparseLayout.Positions["Spawn-E"].X - sparseLayout.Positions["Spawn-A"].X
	assert.Greater(t, crowdedWidth, sparseWidth)
}

func TestWhenZoneNameStartsWithSpawn_MarksZoneAsPlayer(t *testing.T) {
	t.Parallel()
	// Arrange
	service := preview_service.NewPreviewLayoutService(zone_services.NewZoneTierService())
	zones := []entities.Zone{namedZone("Spawn-A"), namedZone("Neutral-B")}
	connections := []entities.Connection{directConnection("Spawn-A", "Neutral-B")}

	// Act
	layout := service.BuildPreviewLayout(templateWith(zones, connections), config.TopologyRing, layoutSide)

	// Assert
	playerFlags := map[string]bool{}
	for _, zone := range layout.Zones {
		playerFlags[zone.Name] = zone.Type == preview.ZoneTypePlayer
	}
	assert.Equal(t, map[string]bool{"Spawn-A": true, "Neutral-B": false}, playerFlags)
}

func TestWhenZoneIsNamedHub_MarksZoneAsHub(t *testing.T) {
	t.Parallel()
	// Arrange
	service := preview_service.NewPreviewLayoutService(zone_services.NewZoneTierService())
	zones := []entities.Zone{namedZone("Spawn-A"), namedZone("Hub")}
	connections := []entities.Connection{directConnection("Spawn-A", "Hub")}

	// Act
	layout := service.BuildPreviewLayout(
		templateWith(zones, connections), config.TopologyHubAndSpoke, layoutSide)

	// Assert
	hubFlags := map[string]bool{}
	for _, zone := range layout.Zones {
		hubFlags[zone.Name] = zone.Type == preview.ZoneTypeHub
	}
	assert.Equal(t, map[string]bool{"Spawn-A": false, "Hub": true}, hubFlags)
}

func TestWhenSpawnMainObjectNamesPlayer_ParsesOwnerNumber(t *testing.T) {
	t.Parallel()
	// Arrange
	service := preview_service.NewPreviewLayoutService(zone_services.NewZoneTierService())
	zone := namedZone("Spawn-A")
	zone.MainObjects = []entities.MainObject{{Type: "Spawn", Spawn: "Player3"}}
	zones := []entities.Zone{zone, namedZone("Neutral-B")}
	connections := []entities.Connection{directConnection("Spawn-A", "Neutral-B")}

	// Act
	layout := service.BuildPreviewLayout(templateWith(zones, connections), config.TopologyRing, layoutSide)

	// Assert
	owners := map[string]int{}
	for _, previewZone := range layout.Zones {
		owners[previewZone.Name] = previewZone.Owner
	}
	assert.Equal(t, 3, owners["Spawn-A"])
}

func TestWhenZoneHasCityMainObjects_CountsCastles(t *testing.T) {
	t.Parallel()
	// Arrange
	service := preview_service.NewPreviewLayoutService(zone_services.NewZoneTierService())
	zone := namedZone("Neutral-B")
	zone.MainObjects = []entities.MainObject{{Type: "City"}, {Type: "City"}}
	zones := []entities.Zone{namedZone("Spawn-A"), zone}
	connections := []entities.Connection{directConnection("Spawn-A", "Neutral-B")}

	// Act
	layout := service.BuildPreviewLayout(templateWith(zones, connections), config.TopologyRing, layoutSide)

	// Assert
	castles := map[string]int{}
	for _, previewZone := range layout.Zones {
		castles[previewZone.Name] = previewZone.Castles
	}
	assert.Equal(t, 2, castles["Neutral-B"])
}

func TestWhenConnectionTypeIsPortal_MarksPreviewConnectionAsPortal(t *testing.T) {
	t.Parallel()
	// Arrange
	service := preview_service.NewPreviewLayoutService(zone_services.NewZoneTierService())
	zones := []entities.Zone{namedZone("Spawn-A"), namedZone("Neutral-B")}
	connections := []entities.Connection{
		{From: "Spawn-A", To: "Neutral-B", ConnectionType: "Portal"},
	}

	// Act
	layout := service.BuildPreviewLayout(templateWith(zones, connections), config.TopologyRing, layoutSide)

	// Assert
	require.Len(t, layout.Connections, 1)
	assert.True(t, layout.Connections[0].IsPortal())
}

func TestWhenZoneHasGladiatorArenaMainObject_MarksZoneAsArena(t *testing.T) {
	t.Parallel()
	// Arrange
	service := preview_service.NewPreviewLayoutService(zone_services.NewZoneTierService())
	zone := namedZone("Neutral-B")
	zone.MainObjects = []entities.MainObject{{Type: "GladiatorArena"}}
	zones := []entities.Zone{namedZone("Spawn-A"), zone}
	connections := []entities.Connection{directConnection("Spawn-A", "Neutral-B")}

	// Act
	layout := service.BuildPreviewLayout(templateWith(zones, connections), config.TopologyRing, layoutSide)

	// Assert
	arenaFlags := map[string]bool{}
	for _, previewZone := range layout.Zones {
		arenaFlags[previewZone.Name] = previewZone.HasArena()
	}
	assert.Equal(t, map[string]bool{"Spawn-A": false, "Neutral-B": true}, arenaFlags)
}

func TestWhenConnectionTypeIsGladiatorArena_MarksPreviewConnectionAsArena(t *testing.T) {
	t.Parallel()
	// Arrange
	service := preview_service.NewPreviewLayoutService(zone_services.NewZoneTierService())
	zones := []entities.Zone{namedZone("Spawn-A"), namedZone("Neutral-B")}
	connections := []entities.Connection{
		{From: "Spawn-A", To: "Neutral-B", ConnectionType: "GladiatorArena"},
	}

	// Act
	layout := service.BuildPreviewLayout(templateWith(zones, connections), config.TopologyRing, layoutSide)

	// Assert
	require.Len(t, layout.Connections, 1)
	assert.True(t, layout.Connections[0].IsGladiatorArena())
}

func TestWhenConnectionTypeIsProximity_MarksPreviewConnectionAsProximity(t *testing.T) {
	t.Parallel()
	// Arrange
	service := preview_service.NewPreviewLayoutService(zone_services.NewZoneTierService())
	zones := []entities.Zone{namedZone("Spawn-A"), namedZone("Neutral-B")}
	connections := []entities.Connection{
		{From: "Spawn-A", To: "Neutral-B", ConnectionType: "Proximity"},
	}

	// Act
	layout := service.BuildPreviewLayout(templateWith(zones, connections), config.TopologyRing, layoutSide)

	// Assert
	require.Len(t, layout.Connections, 1)
	assert.Equal(t, preview.ConnectionTypeProximity, layout.Connections[0].Type)
}

func TestWhenConnectionEndpointHasNoPosition_SkipsThatConnection(t *testing.T) {
	t.Parallel()
	// Arrange
	service := preview_service.NewPreviewLayoutService(zone_services.NewZoneTierService())
	zones := []entities.Zone{namedZone("Spawn-A"), namedZone("Neutral-B")}
	connections := []entities.Connection{
		directConnection("Spawn-A", "Neutral-B"),
		directConnection("Spawn-A", "Ghost-Zone"),
	}

	// Act
	layout := service.BuildPreviewLayout(templateWith(zones, connections), config.TopologyRing, layoutSide)

	// Assert
	assert.Len(t, layout.Connections, 1)
}

func TestWhenTwoConnectionsShareTheSameZonePair_FansOutTheirControlPoints(t *testing.T) {
	t.Parallel()
	// Arrange
	service := preview_service.NewPreviewLayoutService(zone_services.NewZoneTierService())
	zones := []entities.Zone{namedZone("Spawn-A"), namedZone("Neutral-B")}
	connections := []entities.Connection{
		directConnection("Spawn-A", "Neutral-B"),
		directConnection("Spawn-A", "Neutral-B"),
	}

	// Act
	layout := service.BuildPreviewLayout(templateWith(zones, connections), config.TopologyRing, layoutSide)

	// Assert
	require.Len(t, layout.Connections, 2)
	assert.NotEqual(t, layout.Connections[0].Ctrl, layout.Connections[1].Ctrl)
}

func TestWhenZeroAngleZoneIsSet_RotatesTheRingToStartAtThatZone(t *testing.T) {
	t.Parallel()
	// Arrange
	service := preview_service.NewPreviewLayoutService(zone_services.NewZoneTierService())
	zones := []entities.Zone{
		namedZone("Spawn-A"), namedZone("Neutral-B"),
		namedZone("Spawn-C"), namedZone("Neutral-D"),
	}
	connections := []entities.Connection{
		directConnection("Spawn-A", "Neutral-B"),
		directConnection("Neutral-B", "Spawn-C"),
		directConnection("Spawn-C", "Neutral-D"),
		directConnection("Neutral-D", "Spawn-A"),
	}
	defaultLayout := service.BuildPreviewLayout(
		templateWith(zones, connections), config.TopologyRing, layoutSide)
	pivotedTemplate := templateWith(zones, connections)
	pivotedTemplate.Variants[0].Orientation = entities.Orientation{ZeroAngleZone: "Spawn-C"}

	// Act
	pivotedLayout := service.BuildPreviewLayout(pivotedTemplate, config.TopologyRing, layoutSide)

	// Assert - Spawn-C takes the ring slot Spawn-A had without the pivot.
	assert.Equal(t, defaultLayout.Positions["Spawn-A"], pivotedLayout.Positions["Spawn-C"])
}

// ── ring / default dispatch ──────────────────────────────────────────

func TestWhenRingTopologyProvided_PositionsEveryZone(t *testing.T) {
	t.Parallel()
	// Arrange
	service := preview_service.NewPreviewLayoutService(zone_services.NewZoneTierService())
	zones := []entities.Zone{namedZone("Spawn-A"), namedZone("Spawn-B"), namedZone("Neutral-C")}
	connections := []entities.Connection{
		directConnection("Spawn-A", "Spawn-B"),
		directConnection("Spawn-B", "Neutral-C"),
	}

	// Act
	layout := service.BuildPreviewLayout(templateWith(zones, connections), config.TopologyRing, 600)

	// Assert
	assert.Len(t, layout.Positions, 3)
}

func TestWhenTopologyIsUnknown_UsesRingLayout(t *testing.T) {
	t.Parallel()
	// Arrange
	service := preview_service.NewPreviewLayoutService(zone_services.NewZoneTierService())
	zones := []entities.Zone{namedZone("Spawn-A"), namedZone("Neutral-B"), namedZone("Neutral-C")}
	connections := []entities.Connection{
		directConnection("Spawn-A", "Neutral-B"),
		directConnection("Neutral-B", "Neutral-C"),
	}
	template := templateWith(zones, connections)
	expected := service.BuildPreviewLayout(template, config.TopologyRing, layoutSide)

	// Act
	actual := service.BuildPreviewLayout(template, config.MapTopology("Unknown"), layoutSide)

	// Assert
	assert.Equal(t, expected.Positions, actual.Positions)
}

func TestWhenRingTopologyProvided_ComputesPositiveZoneRadius(t *testing.T) {
	t.Parallel()
	// Arrange
	service := preview_service.NewPreviewLayoutService(zone_services.NewZoneTierService())
	zones := []entities.Zone{namedZone("Spawn-A"), namedZone("Spawn-B"), namedZone("Neutral-C")}

	// Act
	layout := service.BuildPreviewLayout(templateWith(zones, nil), config.TopologyRing, 600)

	// Assert
	assert.Positive(t, layout.ZoneRadius)
}

func TestWhenOnlyOneZoneExists_CentersItOnCanvas(t *testing.T) {
	t.Parallel()
	// Arrange
	service := preview_service.NewPreviewLayoutService(zone_services.NewZoneTierService())
	zones := []entities.Zone{namedZone("Spawn-A")}

	// Act
	layout := service.BuildPreviewLayout(templateWith(zones, nil), config.TopologyRing, 600)

	// Assert
	assert.Equal(t, data.NewVec2(300.0, 300.0), layout.Positions["Spawn-A"])
}

func TestWhenZoneIsNamedHub_PlacesItAtCanvasCenter(t *testing.T) {
	t.Parallel()
	// Arrange
	service := preview_service.NewPreviewLayoutService(zone_services.NewZoneTierService())
	zones := []entities.Zone{namedZone("Hub"), namedZone("Spawn-A"), namedZone("Spawn-B")}

	// Act
	layout := service.BuildPreviewLayout(templateWith(zones, nil), config.TopologyRing, 600)

	// Assert
	assert.Equal(t, data.NewVec2(300.0, 300.0), layout.Positions["Hub"])
}

// ── implicit hub rejection ───────────────────────────────────────────

func TestWhenNeutralTouchesEverySpawn_DoesNotCenterIt(t *testing.T) {
	t.Parallel()
	// Arrange
	service := preview_service.NewPreviewLayoutService(zone_services.NewZoneTierService())
	zones := []entities.Zone{
		namedZone("Spawn-A"), namedZone("Spawn-B"), namedZone("Spawn-C"), namedZone("Neutral-H"),
	}
	connections := []entities.Connection{
		directConnection("Neutral-H", "Spawn-A"),
		directConnection("Neutral-H", "Spawn-B"),
		directConnection("Neutral-H", "Spawn-C"),
	}

	// Act
	layout := service.BuildPreviewLayout(templateWith(zones, connections), config.TopologyRing, 600)

	// Assert
	assert.NotEqual(t, data.NewVec2(300.0, 300.0), layout.Positions["Neutral-H"])
}

func TestWhenNeutralTouchesEverySpawn_DoesNotFlagItAsHub(t *testing.T) {
	t.Parallel()
	// Arrange
	service := preview_service.NewPreviewLayoutService(zone_services.NewZoneTierService())
	zones := []entities.Zone{
		namedZone("Spawn-A"), namedZone("Spawn-B"), namedZone("Spawn-C"), namedZone("Neutral-H"),
	}
	connections := []entities.Connection{
		directConnection("Neutral-H", "Spawn-A"),
		directConnection("Neutral-H", "Spawn-B"),
		directConnection("Neutral-H", "Spawn-C"),
	}

	// Act
	layout := service.BuildPreviewLayout(templateWith(zones, connections), config.TopologyRing, 600)

	// Assert
	flaggedHubs := []string{}
	for _, previewZone := range layout.Zones {
		if previewZone.Type == preview.ZoneTypeHub {
			flaggedHubs = append(flaggedHubs, previewZone.Name)
		}
	}
	assert.Empty(t, flaggedHubs)
}

func TestWhenNeutralOnlyConnectsTwoSpawns_FlagsNoHub(t *testing.T) {
	t.Parallel()
	// Arrange
	service := preview_service.NewPreviewLayoutService(zone_services.NewZoneTierService())
	zones := []entities.Zone{namedZone("Spawn-A"), namedZone("Spawn-B"), namedZone("Neutral-H")}
	connections := []entities.Connection{
		directConnection("Neutral-H", "Spawn-A"),
		directConnection("Neutral-H", "Spawn-B"),
	}

	// Act
	layout := service.BuildPreviewLayout(templateWith(zones, connections), config.TopologyRing, 600)

	// Assert
	flaggedHubs := []string{}
	for _, previewZone := range layout.Zones {
		if previewZone.Type == preview.ZoneTypeHub {
			flaggedHubs = append(flaggedHubs, previewZone.Name)
		}
	}
	assert.Empty(t, flaggedHubs)
}

func TestWhenZoneIsExplicitlyNamedHub_FlagsOnlyThatZoneAsHub(t *testing.T) {
	t.Parallel()
	// Arrange
	service := preview_service.NewPreviewLayoutService(zone_services.NewZoneTierService())
	zones := []entities.Zone{namedZone("Hub"), namedZone("Spawn-A"), namedZone("Spawn-B")}
	connections := []entities.Connection{
		directConnection("Hub", "Spawn-A"),
		directConnection("Hub", "Spawn-B"),
	}

	// Act
	layout := service.BuildPreviewLayout(templateWith(zones, connections), config.TopologyRing, 600)

	// Assert
	flaggedHubs := []string{}
	for _, previewZone := range layout.Zones {
		if previewZone.Type == preview.ZoneTypeHub {
			flaggedHubs = append(flaggedHubs, previewZone.Name)
		}
	}
	assert.Equal(t, []string{"Hub"}, flaggedHubs)
}

func TestWhenTwoHubZonesExist_PositionsEveryZone(t *testing.T) {
	t.Parallel()
	// Arrange
	service := preview_service.NewPreviewLayoutService(zone_services.NewZoneTierService())
	zones := []entities.Zone{
		namedZone("Hub-A"), namedZone("Hub-B"), namedZone("Spawn-A"), namedZone("Spawn-B"),
	}
	connections := []entities.Connection{
		directConnection("Hub-A", "Spawn-A"),
		directConnection("Hub-B", "Spawn-B"),
		directConnection("Hub-A", "Hub-B"),
	}

	// Act
	layout := service.BuildPreviewLayout(templateWith(zones, connections), config.TopologyRing, 600)

	// Assert
	assert.Len(t, layout.Positions, 4)
}

// ── scatter (Random) dispatch ────────────────────────────────────────

func TestWhenRandomTopologyZonesHavePositions_PositionsEveryZone(t *testing.T) {
	t.Parallel()
	// Arrange
	service := preview_service.NewPreviewLayoutService(zone_services.NewZoneTierService())
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
	layout := service.BuildPreviewLayout(templateWith(zones, connections), config.TopologyRandom, 600)

	// Assert
	assert.Len(t, layout.Positions, 3)
}

func TestWhenRandomTopologyHasNoConnections_PositionsEveryZone(t *testing.T) {
	t.Parallel()
	// Arrange
	service := preview_service.NewPreviewLayoutService(zone_services.NewZoneTierService())
	zones := []entities.Zone{
		positionedZone("Spawn-A", 0.1, 0.1),
		positionedZone("Spawn-B", 0.9, 0.9),
		positionedZone("Neutral-C", 0.5, 0.5),
	}

	// Act
	layout := service.BuildPreviewLayout(templateWith(zones, nil), config.TopologyRandom, 600)

	// Assert
	assert.Len(t, layout.Positions, 3)
}

func TestWhenRandomTopologyZonesLackPositions_FallsBackToRingLayout(t *testing.T) {
	t.Parallel()
	// Arrange
	service := preview_service.NewPreviewLayoutService(zone_services.NewZoneTierService())
	zones := []entities.Zone{namedZone("Spawn-A"), namedZone("Spawn-B"), namedZone("Neutral-C")}

	// Act
	layout := service.BuildPreviewLayout(templateWith(zones, nil), config.TopologyRandom, 600)

	// Assert
	assert.Len(t, layout.Positions, 3)
}

// ── circles ring dispatch ────────────────────────────────────────────

func TestWhenCirclesZonesSpanMultipleRings_PositionsEveryZone(t *testing.T) {
	t.Parallel()
	// Arrange
	service := preview_service.NewPreviewLayoutService(zone_services.NewZoneTierService())
	zones := []entities.Zone{
		ringedZone("Spawn-A", 0, 0.1, 0.1),
		ringedZone("Spawn-B", 0, 0.9, 0.1),
		ringedZone("Neutral-C", 1, 0.5, 0.5),
	}

	// Act
	layout := service.BuildPreviewLayout(templateWith(zones, nil), config.TopologyCircles, 600)

	// Assert
	assert.Len(t, layout.Positions, 3)
}

func TestWhenCirclesZonesShareOneRing_PositionsEveryZone(t *testing.T) {
	t.Parallel()
	// Arrange
	service := preview_service.NewPreviewLayoutService(zone_services.NewZoneTierService())
	zones := []entities.Zone{
		ringedZone("Spawn-A", 0, 0.2, 0.2),
		ringedZone("Spawn-B", 0, 0.8, 0.8),
		ringedZone("Spawn-C", 0, 0.5, 0.5),
	}

	// Act
	layout := service.BuildPreviewLayout(templateWith(zones, nil), config.TopologyCircles, 600)

	// Assert
	assert.Len(t, layout.Positions, 3)
}

func TestWhenCirclesTopologyHasOneZone_CentersItOnCanvas(t *testing.T) {
	t.Parallel()
	// Arrange
	service := preview_service.NewPreviewLayoutService(zone_services.NewZoneTierService())
	zones := []entities.Zone{ringedZone("Spawn-A", 0, 0.5, 0.5)}

	// Act
	layout := service.BuildPreviewLayout(templateWith(zones, nil), config.TopologyCircles, 600)

	// Assert
	assert.Equal(t, data.NewVec2(300.0, 300.0), layout.Positions["Spawn-A"])
}

// ── connection rendering ─────────────────────────────────────────────

func TestWhenDirectConnectionExists_CollectsIt(t *testing.T) {
	t.Parallel()
	// Arrange
	service := preview_service.NewPreviewLayoutService(zone_services.NewZoneTierService())
	zones := []entities.Zone{namedZone("Spawn-A"), namedZone("Spawn-B")}
	connections := []entities.Connection{directConnection("Spawn-A", "Spawn-B")}

	// Act
	layout := service.BuildPreviewLayout(templateWith(zones, connections), config.TopologyRing, 600)

	// Assert
	assert.Len(t, layout.Connections, 1)
}

func TestWhenDirectConnectionExists_DoesNotFlagItAsPortal(t *testing.T) {
	t.Parallel()
	// Arrange
	service := preview_service.NewPreviewLayoutService(zone_services.NewZoneTierService())
	zones := []entities.Zone{namedZone("Spawn-A"), namedZone("Spawn-B")}
	connections := []entities.Connection{directConnection("Spawn-A", "Spawn-B")}

	// Act
	layout := service.BuildPreviewLayout(templateWith(zones, connections), config.TopologyRing, 600)

	// Assert
	require.Len(t, layout.Connections, 1)
	assert.False(t, layout.Connections[0].IsPortal())
}

func TestWhenPortalConnectionExists_FlagsExactlyOnePortal(t *testing.T) {
	t.Parallel()
	// Arrange
	service := preview_service.NewPreviewLayoutService(zone_services.NewZoneTierService())
	zones := []entities.Zone{namedZone("Spawn-A"), namedZone("Spawn-B"), namedZone("Neutral-C")}
	connections := []entities.Connection{
		directConnection("Spawn-A", "Neutral-C"),
		directConnection("Neutral-C", "Spawn-B"),
		{From: "Spawn-A", To: "Spawn-B", ConnectionType: "Portal"},
	}

	// Act
	layout := service.BuildPreviewLayout(templateWith(zones, connections), config.TopologyRing, 600)

	// Assert
	portalCount := 0
	for _, previewConnection := range layout.Connections {
		if previewConnection.IsPortal() {
			portalCount++
		}
	}
	assert.Equal(t, 1, portalCount)
}

func TestWhenConnectionReferencesUnknownZone_SkipsIt(t *testing.T) {
	t.Parallel()
	// Arrange
	service := preview_service.NewPreviewLayoutService(zone_services.NewZoneTierService())
	zones := []entities.Zone{namedZone("Spawn-A"), namedZone("Spawn-B")}
	connections := []entities.Connection{directConnection("Spawn-A", "Missing-X")}

	// Act
	layout := service.BuildPreviewLayout(templateWith(zones, connections), config.TopologyRing, 600)

	// Assert
	assert.Empty(t, layout.Connections)
}

func TestWhenConnectionSourceIsUnknownZone_SkipsIt(t *testing.T) {
	t.Parallel()
	// Arrange
	service := preview_service.NewPreviewLayoutService(zone_services.NewZoneTierService())
	zones := []entities.Zone{namedZone("Spawn-A"), namedZone("Spawn-B")}
	connections := []entities.Connection{directConnection("Missing-X", "Spawn-B")}

	// Act
	layout := service.BuildPreviewLayout(templateWith(zones, connections), config.TopologyRing, 600)

	// Assert
	assert.Empty(t, layout.Connections)
}

// ── two-cluster (tournament) templates ───────────────────────────────

func TestWhenTemplateHasTwoClusters_PositionsEveryZone(t *testing.T) {
	t.Parallel()
	// Arrange
	service := preview_service.NewPreviewLayoutService(zone_services.NewZoneTierService())
	zones := []entities.Zone{
		namedZone("Spawn-A"), namedZone("Neutral-X"),
		namedZone("Spawn-B"), namedZone("Neutral-Y"),
	}
	connections := []entities.Connection{
		directConnection("Spawn-A", "Neutral-X"),
		directConnection("Spawn-B", "Neutral-Y"),
	}

	// Act
	layout := service.BuildPreviewLayout(templateWith(zones, connections), config.TopologyRing, 600)

	// Assert
	assert.Len(t, layout.Positions, 4)
}

// ── zone classification side-effects ─────────────────────────────────

func TestWhenZoneHasSpawnMainObject_ClassifiesItAsOwnedPlayerZone(t *testing.T) {
	t.Parallel()
	// Arrange
	service := preview_service.NewPreviewLayoutService(zone_services.NewZoneTierService())
	zones := []entities.Zone{
		{Name: "Spawn-A", MainObjects: []entities.MainObject{{Type: "Spawn", Spawn: "Player1"}}},
	}
	expected := preview.Zone{
		Name:    "Spawn-A",
		Label:   "A",
		Center:  data.NewVec2(300.0, 300.0),
		Type:    preview.ZoneTypePlayer,
		Quality: neutral_zone.QualityUnknown,
		Castles: 1,
		Owner:   1,
	}

	// Act
	layout := service.BuildPreviewLayout(templateWith(zones, nil), config.TopologyRing, 600)

	// Assert
	require.Len(t, layout.Zones, 1)
	assert.Equal(t, expected, layout.Zones[0])
}

func TestWhenZoneHasTwoCityMainObjects_CountsTwoCastles(t *testing.T) {
	t.Parallel()
	// Arrange
	service := preview_service.NewPreviewLayoutService(zone_services.NewZoneTierService())
	zones := []entities.Zone{
		{Name: "Neutral-Z", MainObjects: []entities.MainObject{{Type: "City"}, {Type: "City"}}},
	}

	// Act
	layout := service.BuildPreviewLayout(templateWith(zones, nil), config.TopologyRing, 600)

	// Assert
	require.Len(t, layout.Zones, 1)
	assert.Equal(t, 2, layout.Zones[0].Castles)
}

// ── parallel connection fanning ──────────────────────────────────────

func TestWhenOnlyOneEdgeConnectsAPair_KeepsControlPointOnMidpoint(t *testing.T) {
	t.Parallel()
	// Arrange
	service := preview_service.NewPreviewLayoutService(zone_services.NewZoneTierService())
	zones := []entities.Zone{
		positionedZone("Spawn-A", 0.2, 0.5),
		positionedZone("Spawn-B", 0.8, 0.5),
	}
	connections := []entities.Connection{directConnection("Spawn-A", "Spawn-B")}

	// Act
	layout := service.BuildPreviewLayout(templateWith(zones, connections), config.TopologyRandom, 600)

	// Assert
	require.Len(t, layout.Connections, 1)
	edge := layout.Connections[0]
	midX := float64(edge.Start.X+edge.End.X) / 2.0
	midY := float64(edge.Start.Y+edge.End.Y) / 2.0
	distanceFromMidpoint := math.Hypot(float64(edge.Ctrl.X)-midX, float64(edge.Ctrl.Y)-midY)
	assert.LessOrEqual(t, distanceFromMidpoint, 1.5)
}

func TestWhenParallelEdgesConnectSamePair_GivesThemDistinctControlPoints(t *testing.T) {
	t.Parallel()
	// Arrange
	service := preview_service.NewPreviewLayoutService(zone_services.NewZoneTierService())
	zones := []entities.Zone{
		positionedZone("Spawn-A", 0.2, 0.5),
		positionedZone("Spawn-B", 0.8, 0.5),
	}
	connections := []entities.Connection{
		directConnection("Spawn-A", "Spawn-B"),
		directConnection("Spawn-A", "Spawn-B"),
	}

	// Act
	layout := service.BuildPreviewLayout(templateWith(zones, connections), config.TopologyRandom, 600)

	// Assert
	require.Len(t, layout.Connections, 2)
	assert.NotEqual(t, layout.Connections[0].Ctrl, layout.Connections[1].Ctrl)
}

func TestWhenParallelEdgesConnectSamePair_BulgesThemSymmetricallyAboutMidpoint(t *testing.T) {
	t.Parallel()
	// Arrange
	service := preview_service.NewPreviewLayoutService(zone_services.NewZoneTierService())
	zones := []entities.Zone{
		positionedZone("Spawn-A", 0.2, 0.5),
		positionedZone("Spawn-B", 0.8, 0.5),
	}
	connections := []entities.Connection{
		directConnection("Spawn-A", "Spawn-B"),
		directConnection("Spawn-A", "Spawn-B"),
	}

	// Act
	layout := service.BuildPreviewLayout(templateWith(zones, connections), config.TopologyRandom, 600)

	// Assert
	require.Len(t, layout.Connections, 2)
	first, second := layout.Connections[0], layout.Connections[1]
	midX := float64(first.Start.X+first.End.X) / 2.0
	midY := float64(first.Start.Y+first.End.Y) / 2.0
	averageX := float64(first.Ctrl.X+second.Ctrl.X) / 2.0
	averageY := float64(first.Ctrl.Y+second.Ctrl.Y) / 2.0
	averageDistanceFromMidpoint := math.Hypot(averageX-midX, averageY-midY)
	assert.LessOrEqual(t, averageDistanceFromMidpoint, 1.5)
}

// ── fixed-geometry (Square/Geometric/Cross/Fractal) dispatch ─────────

func TestWhenFixedGeometryTopologyHasOneZone_CentersItOnCanvas(t *testing.T) {
	t.Parallel()
	// Arrange
	service := preview_service.NewPreviewLayoutService(zone_services.NewZoneTierService())
	zones := []entities.Zone{positionedZone("Spawn-A", 0.3, 0.7)}

	// Act
	layout := service.BuildPreviewLayout(templateWith(zones, nil), config.TopologySquare, 600)

	// Assert
	assert.Equal(t, data.NewVec2(300.0, 300.0), layout.Positions["Spawn-A"])
}

func TestWhenFixedGeometryZonesShareOnePosition_PositionsBothAtSamePoint(t *testing.T) {
	t.Parallel()
	// Arrange
	service := preview_service.NewPreviewLayoutService(zone_services.NewZoneTierService())
	zones := []entities.Zone{
		positionedZone("Spawn-A", 0.5, 0.5),
		positionedZone("Spawn-B", 0.5, 0.5),
	}

	// Act
	layout := service.BuildPreviewLayout(templateWith(zones, nil), config.TopologySquare, 600)

	// Assert
	assert.Equal(t, layout.Positions["Spawn-A"], layout.Positions["Spawn-B"])
}

func TestWhenFixedGeometryZonesLackPositions_FallsBackToRingLayout(t *testing.T) {
	t.Parallel()
	// Arrange
	service := preview_service.NewPreviewLayoutService(zone_services.NewZoneTierService())
	zones := []entities.Zone{namedZone("Spawn-A"), namedZone("Spawn-B"), namedZone("Neutral-C")}

	// Act
	layout := service.BuildPreviewLayout(templateWith(zones, nil), config.TopologySquare, 600)

	// Assert
	assert.Len(t, layout.Positions, 3)
}

// ── circles ring edge cases ──────────────────────────────────────────

func TestWhenCirclesZoneLacksRingStamp_StillPositionsEveryZone(t *testing.T) {
	t.Parallel()
	// Arrange
	service := preview_service.NewPreviewLayoutService(zone_services.NewZoneTierService())
	missingRingZone := positionedZone("Neutral-C", 0.5, 0.5)
	zones := []entities.Zone{
		ringedZone("Spawn-A", 0, 0.1, 0.1),
		ringedZone("Spawn-B", 0, 0.9, 0.1),
		missingRingZone,
	}

	// Act
	layout := service.BuildPreviewLayout(templateWith(zones, nil), config.TopologyCircles, 600)

	// Assert
	assert.Len(t, layout.Positions, 3)
}

func TestWhenCirclesOuterRingHasSingleZone_PositionsEveryZone(t *testing.T) {
	t.Parallel()
	// Arrange
	service := preview_service.NewPreviewLayoutService(zone_services.NewZoneTierService())
	zones := []entities.Zone{
		ringedZone("Neutral-C", 1, 0.4, 0.4),
		ringedZone("Neutral-D", 1, 0.6, 0.4),
		ringedZone("Neutral-E", 1, 0.5, 0.6),
		ringedZone("Spawn-A", 0, 0.5, 0.9),
	}

	// Act
	layout := service.BuildPreviewLayout(templateWith(zones, nil), config.TopologyCircles, 600)

	// Assert
	assert.Len(t, layout.Positions, 4)
}

func TestWhenCirclesOuterRingIsOvercrowded_ShrinksZoneRadiusBelowMaximum(t *testing.T) {
	t.Parallel()
	// Arrange - 23 zones on one ring force the ring circumference past the
	// draw radius at the maximum zone size, so the binary search must shrink.
	service := preview_service.NewPreviewLayoutService(zone_services.NewZoneTierService())
	zones := []entities.Zone{ringedZone("Neutral-Center", 1, 0.5, 0.5)}
	for index := range 23 {
		angle := 2.0 * math.Pi * float64(index) / 23.0
		zones = append(zones, ringedZone(
			"Neutral-"+string(rune('A'+index)), 0,
			0.5+0.4*math.Cos(angle), 0.5+0.4*math.Sin(angle)))
	}

	// Act
	layout := service.BuildPreviewLayout(templateWith(zones, nil), config.TopologyCircles, 600)

	// Assert
	assert.Less(t, layout.ZoneRadius, 33.0)
}

// ── scatter edge cases ───────────────────────────────────────────────

func TestWhenRandomTopologyHasOneZone_CentersItOnCanvas(t *testing.T) {
	t.Parallel()
	// Arrange
	service := preview_service.NewPreviewLayoutService(zone_services.NewZoneTierService())
	zones := []entities.Zone{positionedZone("Spawn-A", 0.2, 0.8)}

	// Act
	layout := service.BuildPreviewLayout(templateWith(zones, nil), config.TopologyRandom, 600)

	// Assert
	assert.Equal(t, data.NewVec2(300.0, 300.0), layout.Positions["Spawn-A"])
}

func TestWhenConnectedScatterZonesShareOnePosition_PositionsBothZones(t *testing.T) {
	t.Parallel()
	// Arrange
	service := preview_service.NewPreviewLayoutService(zone_services.NewZoneTierService())
	zones := []entities.Zone{
		positionedZone("Spawn-A", 0.5, 0.5),
		positionedZone("Spawn-B", 0.5, 0.5),
	}
	connections := []entities.Connection{directConnection("Spawn-A", "Spawn-B")}

	// Act
	layout := service.BuildPreviewLayout(templateWith(zones, connections), config.TopologyRandom, 600)

	// Assert
	assert.Len(t, layout.Positions, 2)
}

func TestWhenUnconnectedZoneLiesFarFromTightPair_KeepsEveryZoneInsideCanvas(t *testing.T) {
	t.Parallel()
	// Arrange
	service := preview_service.NewPreviewLayoutService(zone_services.NewZoneTierService())
	zones := []entities.Zone{
		positionedZone("Spawn-A", 0.0, 0.0),
		positionedZone("Spawn-B", 0.001, 0.0),
		positionedZone("Neutral-C", 1.0, 1.0),
	}
	connections := []entities.Connection{directConnection("Spawn-A", "Spawn-B")}

	// Act
	layout := service.BuildPreviewLayout(templateWith(zones, connections), config.TopologyRandom, 600)

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
	t.Parallel()
	// Arrange
	service := preview_service.NewPreviewLayoutService(zone_services.NewZoneTierService())
	zones := []entities.Zone{
		positionedZone("Spawn-A", 0.1, 0.5),
		positionedZone("Spawn-B", 0.9, 0.5),
		positionedZone("Neutral-C", 0.5, 0.5),
	}
	connections := []entities.Connection{directConnection("Spawn-A", "Spawn-B")}

	// Act
	layout := service.BuildPreviewLayout(templateWith(zones, connections), config.TopologyRandom, 600)

	// Assert
	assert.NotEqual(t, layout.Positions["Spawn-A"].Y, layout.Positions["Neutral-C"].Y)
}

// ── scatter adjacency filtering ──────────────────────────────────────

func TestWhenScatterConnectionIsPortal_IgnoresItForAdjacency(t *testing.T) {
	t.Parallel()
	// Arrange
	service := preview_service.NewPreviewLayoutService(zone_services.NewZoneTierService())
	zones := []entities.Zone{
		positionedZone("Spawn-A", 0.2, 0.2),
		positionedZone("Spawn-B", 0.8, 0.8),
		positionedZone("Neutral-C", 0.5, 0.5),
	}
	portalOnly := []entities.Connection{{From: "Spawn-A", To: "Spawn-B", ConnectionType: "Portal"}}
	unconnectedLayout := service.BuildPreviewLayout(templateWith(zones, nil), config.TopologyRandom, 600)

	// Act
	layout := service.BuildPreviewLayout(templateWith(zones, portalOnly), config.TopologyRandom, 600)

	// Assert
	assert.Equal(t, unconnectedLayout.Positions, layout.Positions)
}

func TestWhenScatterConnectionReferencesUnknownZone_IgnoresItForAdjacency(t *testing.T) {
	t.Parallel()
	// Arrange
	service := preview_service.NewPreviewLayoutService(zone_services.NewZoneTierService())
	zones := []entities.Zone{
		positionedZone("Spawn-A", 0.2, 0.2),
		positionedZone("Spawn-B", 0.8, 0.8),
	}
	danglingOnly := []entities.Connection{directConnection("Spawn-A", "Missing-X")}
	unconnectedLayout := service.BuildPreviewLayout(templateWith(zones, nil), config.TopologyRandom, 600)

	// Act
	layout := service.BuildPreviewLayout(templateWith(zones, danglingOnly), config.TopologyRandom, 600)

	// Assert
	assert.Equal(t, unconnectedLayout.Positions, layout.Positions)
}

func TestWhenScatterConnectionIsSelfLoop_IgnoresItForAdjacency(t *testing.T) {
	t.Parallel()
	// Arrange
	service := preview_service.NewPreviewLayoutService(zone_services.NewZoneTierService())
	zones := []entities.Zone{
		positionedZone("Spawn-A", 0.2, 0.2),
		positionedZone("Spawn-B", 0.8, 0.8),
	}
	selfLoopOnly := []entities.Connection{directConnection("Spawn-A", "Spawn-A")}
	unconnectedLayout := service.BuildPreviewLayout(templateWith(zones, nil), config.TopologyRandom, 600)

	// Act
	layout := service.BuildPreviewLayout(templateWith(zones, selfLoopOnly), config.TopologyRandom, 600)

	// Assert
	assert.Equal(t, unconnectedLayout.Positions, layout.Positions)
}

// ── manual-position dispatch ─────────────────────────────────────────

func TestWhenAllZonesHaveManualPositions_PlacesThemAtScaledCoordinates(t *testing.T) {
	t.Parallel()
	// Arrange
	service := preview_service.NewPreviewLayoutService(zone_services.NewZoneTierService())
	zones := []entities.Zone{
		manualZone("Spawn-A", 0.25, 0.5),
		manualZone("Spawn-B", 0.75, 0.5),
		manualZone("Spawn-C", 0.5, 0.25),
	}
	expected := map[string]data.Vec2[float64]{
		"Spawn-A": data.NewVec2(150.0, 300.0),
		"Spawn-B": data.NewVec2(450.0, 300.0),
		"Spawn-C": data.NewVec2(300.0, 150.0),
	}

	// Act
	layout := service.BuildPreviewLayout(templateWith(zones, nil), config.TopologyRing, 600)

	// Assert
	assert.Equal(t, expected, layout.Positions)
}

func TestWhenManualZonesCoincide_KeepsControlPointOnSharedPoint(t *testing.T) {
	t.Parallel()
	// Arrange
	service := preview_service.NewPreviewLayoutService(zone_services.NewZoneTierService())
	zones := []entities.Zone{
		manualZone("Spawn-A", 0.5, 0.5),
		manualZone("Spawn-B", 0.5, 0.5),
	}
	connections := []entities.Connection{directConnection("Spawn-A", "Spawn-B")}

	// Act
	layout := service.BuildPreviewLayout(templateWith(zones, connections), config.TopologyRing, 600)

	// Assert
	require.Len(t, layout.Connections, 1)
	assert.Equal(t, layout.Connections[0].Start, layout.Connections[0].Ctrl)
}

// ── multi-hub edge cases ─────────────────────────────────────────────

func TestWhenZoneConnectsToNoHub_PlacesItAtCanvasCenter(t *testing.T) {
	t.Parallel()
	// Arrange
	service := preview_service.NewPreviewLayoutService(zone_services.NewZoneTierService())
	zones := []entities.Zone{
		namedZone("Hub-A"), namedZone("Hub-B"),
		namedZone("Spawn-A"), namedZone("Neutral-X"),
	}
	connections := []entities.Connection{directConnection("Hub-A", "Spawn-A")}

	// Act
	layout := service.BuildPreviewLayout(templateWith(zones, connections), config.TopologyHubAndSpoke, 600)

	// Assert
	assert.Equal(t, data.NewVec2(300.0, 300.0), layout.Positions["Neutral-X"])
}

func TestWhenZoneOnlyPortalsToAHub_PlacesItAtCanvasCenter(t *testing.T) {
	t.Parallel()
	// Arrange - portal connections never count as spokes, so the zone
	// collapses to the canvas center as a straggler.
	service := preview_service.NewPreviewLayoutService(zone_services.NewZoneTierService())
	zones := []entities.Zone{
		namedZone("Hub-A"), namedZone("Hub-B"),
		namedZone("Spawn-A"), namedZone("Neutral-X"),
	}
	connections := []entities.Connection{
		directConnection("Hub-A", "Spawn-A"),
		{From: "Hub-B", To: "Neutral-X", ConnectionType: "Portal"},
	}

	// Act
	layout := service.BuildPreviewLayout(templateWith(zones, connections), config.TopologyHubAndSpoke, 600)

	// Assert
	assert.Equal(t, data.NewVec2(300.0, 300.0), layout.Positions["Neutral-X"])
}

func TestWhenHubSpokeConnectionIsDuplicated_PlacesTheSpokeOnce(t *testing.T) {
	t.Parallel()
	// Arrange
	service := preview_service.NewPreviewLayoutService(zone_services.NewZoneTierService())
	zones := []entities.Zone{
		namedZone("Hub-A"), namedZone("Hub-B"),
		namedZone("Spawn-A"), namedZone("Spawn-B"),
	}
	singleConnections := []entities.Connection{
		directConnection("Hub-A", "Spawn-A"),
		directConnection("Hub-B", "Spawn-B"),
	}
	duplicatedConnections := append(
		[]entities.Connection{directConnection("Hub-A", "Spawn-A")}, singleConnections...)
	singleLayout := service.BuildPreviewLayout(
		templateWith(zones, singleConnections), config.TopologyHubAndSpoke, 600)

	// Act
	layout := service.BuildPreviewLayout(
		templateWith(zones, duplicatedConnections), config.TopologyHubAndSpoke, 600)

	// Assert
	assert.Equal(t, singleLayout.Positions, layout.Positions)
}

// ── zero-angle-zone rotation ─────────────────────────────────────────

func TestWhenZeroAngleZoneIsSet_RotatesThatZoneToFirstRingSlot(t *testing.T) {
	t.Parallel()
	// Arrange
	service := preview_service.NewPreviewLayoutService(zone_services.NewZoneTierService())
	zones := []entities.Zone{namedZone("Spawn-A"), namedZone("Spawn-B"), namedZone("Neutral-C")}
	rmgTemplate := &entities.RmgTemplate{
		Variants: []entities.Variant{{
			Zones:       zones,
			Orientation: entities.Orientation{ZeroAngleZone: "Spawn-B"},
		}},
	}

	// Act
	layout := service.BuildPreviewLayout(rmgTemplate, config.TopologyRing, 600)

	// Assert
	assert.InDeltaSlice(t,
		[]float64{300.0, 48.0},
		[]float64{layout.Positions["Spawn-B"].X, layout.Positions["Spawn-B"].Y},
		1e-9)
}
