package topologyBase_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/base"
	"github.com/stretchr/testify/assert"
)

func TestWhenFewerThanTwoLabelsExist_NoBridgesAreCreated(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase()

	// Act
	connections := topologyBase.CreateMissingConnections(
		[]string{"A"}, []string{"A"}, models.Positions{{X: 0.5, Y: 0.5}}, nil, nil, newUnitTuning(), nil)

	// Assert
	assert.Empty(t, connections)
}

func TestWhenAllZonesAreAlreadyConnected_NoBridgesAreCreated(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase()
	positions := models.Positions{{X: 0.1, Y: 0.1}, {X: 0.9, Y: 0.9}}
	existingConnections := []entities.Connection{
		{Name: "Ring-A-B", From: "Spawn-A", To: "Spawn-B", ConnectionType: "Direct"},
	}

	// Act
	connections := topologyBase.CreateMissingConnections(
		[]string{"A", "B"}, []string{"A", "B"}, positions, nil, existingConnections, newUnitTuning(), nil)

	// Assert
	assert.Empty(t, connections)
}

func TestWhenTwoPlayerZonesAreDisconnected_BridgeLinksThemWithPlayerBorderGuard(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase()
	positions := models.Positions{{X: 0.1, Y: 0.1}, {X: 0.9, Y: 0.9}}
	expectedConnections := []entities.Connection{
		{
			Name: "Bridge-A-B", From: "Spawn-A", To: "Spawn-B",
			ConnectionType: "Direct", GuardZone: "Spawn-A", SimTurnSquad: true,
			GuardValue: 30000, GuardWeeklyIncrement: 0.15,
			GuardMatchGroup: "bridge_guard_A-B",
		},
	}

	// Act
	connections := topologyBase.CreateMissingConnections(
		[]string{"A", "B"}, []string{"A", "B"}, positions, nil, nil, newUnitTuning(), nil)

	// Assert
	assert.Equal(t, expectedConnections, connections)
}

func TestWhenThreeZonesAreAllDisconnected_BridgesAreAddedUntilFullyConnected(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase()
	positions := models.Positions{{X: 0.1, Y: 0.1}, {X: 0.5, Y: 0.5}, {X: 0.9, Y: 0.9}}

	// Act
	connections := topologyBase.CreateMissingConnections(
		[]string{"A", "B", "C"}, []string{"A", "B", "C"}, positions, nil, nil, newUnitTuning(), nil)

	// Assert
	assert.Len(t, connections, 2)
}

func TestWhenIsolatedZoneSitsClosestToSecondZone_BridgeAttachesToClosestPair(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase()
	positions := models.Positions{{X: 0.1, Y: 0.1}, {X: 0.5, Y: 0.5}, {X: 0.55, Y: 0.5}}
	existingConnections := []entities.Connection{
		{Name: "Ring-A-B", From: "Spawn-A", To: "Spawn-B", ConnectionType: "Direct"},
	}
	expectedConnections := []entities.Connection{
		{
			Name: "Bridge-B-C", From: "Spawn-B", To: "Spawn-C",
			ConnectionType: "Direct", GuardZone: "Spawn-B", SimTurnSquad: true,
			GuardValue: 30000, GuardWeeklyIncrement: 0.15,
			GuardMatchGroup: "bridge_guard_B-C",
		},
	}

	// Act
	connections := topologyBase.CreateMissingConnections(
		[]string{"A", "B", "C"}, []string{"A", "B", "C"}, positions, nil,
		existingConnections, newUnitTuning(), nil)

	// Assert
	assert.Equal(t, expectedConnections, connections)
}

func TestWhenDisconnectedZonesAreNeutral_BridgeGuardUsesHigherNeutralQuality(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase()
	positions := models.Positions{{X: 0.1, Y: 0.1}, {X: 0.9, Y: 0.9}}
	neutralPlans := models.NeutralZonePlans{
		{Label: "C", Quality: models.QualityLow, CastleCount: 0},
		{Label: "D", Quality: models.QualityHigh, CastleCount: 0},
	}
	expectedConnections := []entities.Connection{
		{
			Name: "Bridge-C-D", From: "Neutral-C", To: "Neutral-D",
			ConnectionType: "Direct", GuardZone: "Neutral-C", SimTurnSquad: true,
			GuardValue: 25000, GuardWeeklyIncrement: 0.15,
			GuardMatchGroup: "bridge_guard_C-D",
		},
	}

	// Act
	connections := topologyBase.CreateMissingConnections(
		nil, []string{"C", "D"}, positions, nil, nil, newUnitTuning(), neutralPlans)

	// Assert
	assert.Equal(t, expectedConnections, connections)
}

func TestWhenLabelOrderIsReversed_BridgeNameStillSortsLabelsAlphabetically(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase()
	positions := models.Positions{{X: 0.1, Y: 0.1}, {X: 0.9, Y: 0.9}}
	expectedConnections := []entities.Connection{
		{
			Name: "Bridge-A-B", From: "Spawn-B", To: "Spawn-A",
			ConnectionType: "Direct", GuardZone: "Spawn-B", SimTurnSquad: true,
			GuardValue: 30000, GuardWeeklyIncrement: 0.15,
			GuardMatchGroup: "bridge_guard_A-B",
		},
	}

	// Act
	connections := topologyBase.CreateMissingConnections(
		[]string{"A", "B"}, []string{"B", "A"}, positions, nil, nil, newUnitTuning(), nil)

	// Assert
	assert.Equal(t, expectedConnections, connections)
}

func TestWhenBridgedZonesHaveVariousRoadShapes_BridgeIsStillReturned(t *testing.T) {
	t.Parallel()
	// Arrange
	positions := models.Positions{{X: 0.1, Y: 0.1}, {X: 0.9, Y: 0.9}}
	testCases := []struct {
		name      string
		firstZone entities.Zone
	}{
		{
			name: "WhenZoneHasMainObjects_RoadAttachesToPrimaryMainObject",
			firstZone: entities.Zone{Name: "Spawn-A", MainObjects: []entities.MainObject{
				{Type: "Spawn"},
			}},
		},
		{
			name: "WhenZoneHasRoadFromExistingConnection_RoadChainsFromIt",
			firstZone: entities.Zone{Name: "Spawn-A", Roads: []entities.Road{
				{From: entities.TypedRef{Type: "Connection", Args: []string{"Old-Conn"}}},
			}},
		},
		{
			name: "WhenZoneHasRoadToExistingConnection_RoadChainsFromIt",
			firstZone: entities.Zone{Name: "Spawn-A", Roads: []entities.Road{
				{To: entities.TypedRef{Type: "Connection", Args: []string{"Old-Conn"}}},
			}},
		},
		{
			name: "WhenZoneRoadsHaveNoConnectionRefs_RoadFallsBackToBridgeName",
			firstZone: entities.Zone{Name: "Spawn-A", Roads: []entities.Road{
				{
					From: entities.TypedRef{Type: "MainObject", Args: []string{"0"}},
					To:   entities.TypedRef{Type: "MainObject", Args: []string{"1"}},
				},
			}},
		},
		{
			name:      "WhenZoneHasNeitherMainObjectsNorRoads_RoadSelfReferencesBridge",
			firstZone: entities.Zone{Name: "Spawn-A"},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			// Arrange
			topologyBase := base.NewTopologyBase()
			zones := []entities.Zone{testCase.firstZone, {Name: "Spawn-B"}}

			// Act
			connections := topologyBase.CreateMissingConnections(
				[]string{"A", "B"}, []string{"A", "B"}, positions, zones, nil, newUnitTuning(), nil)

			// Assert
			assert.Len(t, connections, 1)
		})
	}
}
