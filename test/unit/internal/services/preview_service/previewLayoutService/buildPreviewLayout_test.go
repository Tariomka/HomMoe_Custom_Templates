package previewLayoutService_test

import (
	"image"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/preview"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/preview_service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const layoutSide = 700.0

func TestWhenTemplateIsNil_ReturnsEmptyLayout(t *testing.T) {
	// Arrange
	expected := preview.PreviewLayout{Positions: map[string]image.Point{}}

	// Act
	actual := preview_service.BuildPreviewLayout(nil, config.TopologyRing, layoutSide)

	// Assert
	assert.Equal(t, expected, actual)
}

func TestWhenTemplateHasNoVariants_ReturnsEmptyLayout(t *testing.T) {
	// Arrange
	expected := preview.PreviewLayout{Positions: map[string]image.Point{}}

	// Act
	actual := preview_service.BuildPreviewLayout(&entities.RmgTemplate{}, config.TopologyRing, layoutSide)

	// Assert
	assert.Equal(t, expected, actual)
}

func TestWhenVariantHasNoZones_ReturnsEmptyLayout(t *testing.T) {
	// Arrange
	expected := preview.PreviewLayout{Positions: map[string]image.Point{}}

	// Act
	actual := preview_service.BuildPreviewLayout(templateWith(nil, nil), config.TopologyRing, layoutSide)

	// Assert
	assert.Equal(t, expected, actual)
}

func TestWhenRingTopologyIsLaidOut_EveryZoneStaysInsideTheCanvas(t *testing.T) {
	// Arrange
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
	layout := preview_service.BuildPreviewLayout(templateWith(zones, connections), config.TopologyRing, layoutSide)

	// Assert
	require.Len(t, layout.Positions, len(zones))
	canvasBounds := image.Rect(0, 0, layoutSide, layoutSide)
	for zoneName, zonePosition := range layout.Positions {
		assert.True(t, zonePosition.In(canvasBounds), "zone %s at %v escapes the canvas", zoneName, zonePosition)
	}
}

func TestWhenAllZonesHaveManualPositions_PlacesThemVerbatim(t *testing.T) {
	// Arrange
	zones := []entities.Zone{
		manualZone("Spawn-A", 0.25, 0.5),
		manualZone("Neutral-B", 0.75, 0.5),
	}
	expected := map[string]image.Point{
		"Spawn-A":   image.Pt(175, 350),
		"Neutral-B": image.Pt(525, 350),
	}

	// Act
	layout := preview_service.BuildPreviewLayout(templateWith(zones, nil), config.TopologyRing, layoutSide)

	// Assert
	assert.Equal(t, expected, layout.Positions)
}

func TestWhenFixedGeometryTopologyIsLaidOut_PreservesRelativeGeometry(t *testing.T) {
	// Arrange
	zones := []entities.Zone{
		positionedZone("Spawn-A", 0.1, 0.5),
		positionedZone("Neutral-B", 0.5, 0.5),
		positionedZone("Spawn-C", 0.9, 0.5),
	}

	// Act
	layout := preview_service.BuildPreviewLayout(templateWith(zones, nil), config.TopologySquare, layoutSide)

	// Assert - the middle zone must stay the exact midpoint of the outer two.
	require.Len(t, layout.Positions, len(zones))
	left := layout.Positions["Spawn-A"]
	middle := layout.Positions["Neutral-B"]
	right := layout.Positions["Spawn-C"]
	assert.Equal(t, image.Pt((left.X+right.X)/2, (left.Y+right.Y)/2), middle)
}

func TestWhenZoneNameStartsWithSpawn_MarksZoneAsPlayer(t *testing.T) {
	// Arrange
	zones := []entities.Zone{namedZone("Spawn-A"), namedZone("Neutral-B")}
	connections := []entities.Connection{directConnection("Spawn-A", "Neutral-B")}

	// Act
	layout := preview_service.BuildPreviewLayout(templateWith(zones, connections), config.TopologyRing, layoutSide)

	// Assert
	playerFlags := map[string]bool{}
	for _, zone := range layout.Zones {
		playerFlags[zone.Name] = zone.IsPlayer
	}
	assert.Equal(t, map[string]bool{"Spawn-A": true, "Neutral-B": false}, playerFlags)
}

func TestWhenZoneIsNamedHub_MarksZoneAsHub(t *testing.T) {
	// Arrange
	zones := []entities.Zone{namedZone("Spawn-A"), namedZone("Hub")}
	connections := []entities.Connection{directConnection("Spawn-A", "Hub")}

	// Act
	layout := preview_service.BuildPreviewLayout(
		templateWith(zones, connections), config.TopologyHubAndSpoke, layoutSide)

	// Assert
	hubFlags := map[string]bool{}
	for _, zone := range layout.Zones {
		hubFlags[zone.Name] = zone.IsHub
	}
	assert.Equal(t, map[string]bool{"Spawn-A": false, "Hub": true}, hubFlags)
}

func TestWhenSpawnMainObjectNamesPlayer_ParsesOwnerNumber(t *testing.T) {
	// Arrange
	zone := namedZone("Spawn-A")
	zone.MainObjects = []entities.MainObject{{Type: "Spawn", Spawn: "Player3"}}
	zones := []entities.Zone{zone, namedZone("Neutral-B")}
	connections := []entities.Connection{directConnection("Spawn-A", "Neutral-B")}

	// Act
	layout := preview_service.BuildPreviewLayout(templateWith(zones, connections), config.TopologyRing, layoutSide)

	// Assert
	owners := map[string]int{}
	for _, previewZone := range layout.Zones {
		owners[previewZone.Name] = previewZone.Owner
	}
	assert.Equal(t, 3, owners["Spawn-A"])
}

func TestWhenZoneHasCityMainObjects_CountsCastles(t *testing.T) {
	// Arrange
	zone := namedZone("Neutral-B")
	zone.MainObjects = []entities.MainObject{{Type: "City"}, {Type: "City"}}
	zones := []entities.Zone{namedZone("Spawn-A"), zone}
	connections := []entities.Connection{directConnection("Spawn-A", "Neutral-B")}

	// Act
	layout := preview_service.BuildPreviewLayout(templateWith(zones, connections), config.TopologyRing, layoutSide)

	// Assert
	castles := map[string]int{}
	for _, previewZone := range layout.Zones {
		castles[previewZone.Name] = previewZone.Castles
	}
	assert.Equal(t, 2, castles["Neutral-B"])
}

func TestWhenConnectionTypeIsPortal_MarksPreviewConnectionAsPortal(t *testing.T) {
	// Arrange
	zones := []entities.Zone{namedZone("Spawn-A"), namedZone("Neutral-B")}
	connections := []entities.Connection{
		{From: "Spawn-A", To: "Neutral-B", ConnectionType: "Portal"},
	}

	// Act
	layout := preview_service.BuildPreviewLayout(templateWith(zones, connections), config.TopologyRing, layoutSide)

	// Assert
	require.Len(t, layout.Connections, 1)
	assert.True(t, layout.Connections[0].Portal)
}

func TestWhenConnectionEndpointHasNoPosition_SkipsThatConnection(t *testing.T) {
	// Arrange
	zones := []entities.Zone{namedZone("Spawn-A"), namedZone("Neutral-B")}
	connections := []entities.Connection{
		directConnection("Spawn-A", "Neutral-B"),
		directConnection("Spawn-A", "Ghost-Zone"),
	}

	// Act
	layout := preview_service.BuildPreviewLayout(templateWith(zones, connections), config.TopologyRing, layoutSide)

	// Assert
	assert.Len(t, layout.Connections, 1)
}

func TestWhenTwoConnectionsShareTheSameZonePair_FansOutTheirControlPoints(t *testing.T) {
	// Arrange
	zones := []entities.Zone{namedZone("Spawn-A"), namedZone("Neutral-B")}
	connections := []entities.Connection{
		directConnection("Spawn-A", "Neutral-B"),
		directConnection("Spawn-A", "Neutral-B"),
	}

	// Act
	layout := preview_service.BuildPreviewLayout(templateWith(zones, connections), config.TopologyRing, layoutSide)

	// Assert
	require.Len(t, layout.Connections, 2)
	assert.NotEqual(t, layout.Connections[0].Ctrl, layout.Connections[1].Ctrl)
}

func TestWhenZeroAngleZoneIsSet_RotatesTheRingToStartAtThatZone(t *testing.T) {
	// Arrange
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
	defaultLayout := preview_service.BuildPreviewLayout(
		templateWith(zones, connections), config.TopologyRing, layoutSide)
	pivotedTemplate := templateWith(zones, connections)
	pivotedTemplate.Variants[0].Orientation = entities.Orientation{ZeroAngleZone: "Spawn-C"}

	// Act
	pivotedLayout := preview_service.BuildPreviewLayout(pivotedTemplate, config.TopologyRing, layoutSide)

	// Assert - Spawn-C takes the ring slot Spawn-A had without the pivot.
	assert.Equal(t, defaultLayout.Positions["Spawn-A"], pivotedLayout.Positions["Spawn-C"])
}
