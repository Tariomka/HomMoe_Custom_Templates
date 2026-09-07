package topologyBase_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/base"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
)

func TestWhenFewerThanTwoLabelsExist_NoBridgesAreCreated(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())

	// Act
	connections := topologyBase.CreateMissingConnections(
		[]string{"A"}, []string{"A"}, models.Positions{{X: 0.5, Y: 0.5}}, nil, nil, newUnitTuning(), nil)

	// Assert
	assert.Empty(t, connections)
}

func TestWhenAllZonesAreAlreadyConnected_NoBridgesAreCreated(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())
	positions := models.Positions{{X: 0.1, Y: 0.1}, {X: 0.9, Y: 0.9}}
	existingConnections := []template_model.Connection{
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
	topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())
	positions := models.Positions{{X: 0.1, Y: 0.1}, {X: 0.9, Y: 0.9}}
	expectedConnections := []template_model.Connection{
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
	topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())
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
	topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())
	positions := models.Positions{{X: 0.1, Y: 0.1}, {X: 0.5, Y: 0.5}, {X: 0.55, Y: 0.5}}
	existingConnections := []template_model.Connection{
		{Name: "Ring-A-B", From: "Spawn-A", To: "Spawn-B", ConnectionType: "Direct"},
	}
	expectedConnections := []template_model.Connection{
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
	topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())
	positions := models.Positions{{X: 0.1, Y: 0.1}, {X: 0.9, Y: 0.9}}
	neutralPlans := neutral_zone.Plans{
		{Label: "C", Quality: neutral_zone.QualityLow, CastleCount: 0},
		{Label: "D", Quality: neutral_zone.QualityHigh, CastleCount: 0},
	}
	expectedConnections := []template_model.Connection{
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
	topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())
	positions := models.Positions{{X: 0.1, Y: 0.1}, {X: 0.9, Y: 0.9}}
	expectedConnections := []template_model.Connection{
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
		firstZone template_model.Zone
	}{
		{
			name: "WhenZoneHasMainObjects_RoadAttachesToPrimaryMainObject",
			firstZone: template_model.Zone{Name: "Spawn-A", MainObjects: []template_model.MainObject{
				{Type: "Spawn"},
			}},
		},
		{
			name: "WhenZoneHasRoadFromExistingConnection_RoadChainsFromIt",
			firstZone: template_model.Zone{Name: "Spawn-A", Roads: []template_model.Road{
				{From: template_model.TypedRef{Type: "Connection", Args: []string{"Old-Conn"}}},
			}},
		},
		{
			name: "WhenZoneHasRoadToExistingConnection_RoadChainsFromIt",
			firstZone: template_model.Zone{Name: "Spawn-A", Roads: []template_model.Road{
				{To: template_model.TypedRef{Type: "Connection", Args: []string{"Old-Conn"}}},
			}},
		},
		{
			name: "WhenZoneRoadsHaveNoConnectionRefs_RoadFallsBackToBridgeName",
			firstZone: template_model.Zone{Name: "Spawn-A", Roads: []template_model.Road{
				{
					From: template_model.TypedRef{Type: "MainObject", Args: []string{"0"}},
					To:   template_model.TypedRef{Type: "MainObject", Args: []string{"1"}},
				},
			}},
		},
		{
			name:      "WhenZoneHasNeitherMainObjectsNorRoads_RoadSelfReferencesBridge",
			firstZone: template_model.Zone{Name: "Spawn-A"},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			// Arrange
			topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())
			zones := []template_model.Zone{testCase.firstZone, {Name: "Spawn-B"}}

			// Act
			connections := topologyBase.CreateMissingConnections(
				[]string{"A", "B"}, []string{"A", "B"}, positions, zones, nil, newUnitTuning(), nil)

			// Assert
			assert.Len(t, connections, 1)
		})
	}
}

func TestWhenBridgeIsCreated_FirstZoneInSliceGainsBridgeRoad(t *testing.T) {
	t.Parallel()
	// Arrange
	positions := models.Positions{{X: 0.1, Y: 0.1}, {X: 0.9, Y: 0.9}}
	testCases := []struct {
		name          string
		firstZone     template_model.Zone
		expectedRoads []template_model.Road
	}{
		{
			name: "WhenZoneHasMainObjects_AppendedRoadStartsAtPrimaryMainObject",
			firstZone: template_model.Zone{Name: "Spawn-A", MainObjects: []template_model.MainObject{
				{Type: "Spawn"},
			}},
			expectedRoads: []template_model.Road{
				{
					From: template_model.TypedRef{Type: "MainObject", Args: []string{"0"}},
					To:   template_model.TypedRef{Type: "Connection", Args: []string{"Bridge-A-B"}},
				},
			},
		},
		{
			name: "WhenZoneHasRoadReferencingExistingConnection_AppendedRoadChainsFromIt",
			firstZone: template_model.Zone{Name: "Spawn-A", Roads: []template_model.Road{
				{From: template_model.TypedRef{Type: "Connection", Args: []string{"Old-Conn"}}},
			}},
			expectedRoads: []template_model.Road{
				{From: template_model.TypedRef{Type: "Connection", Args: []string{"Old-Conn"}}},
				{
					From: template_model.TypedRef{Type: "Connection", Args: []string{"Old-Conn"}},
					To:   template_model.TypedRef{Type: "Connection", Args: []string{"Bridge-A-B"}},
				},
			},
		},
		{
			name: "WhenZoneRoadsHaveNoConnectionRefs_AppendedRoadFallsBackToBridgeName",
			firstZone: template_model.Zone{Name: "Spawn-A", Roads: []template_model.Road{
				{
					From: template_model.TypedRef{Type: "MainObject", Args: []string{"0"}},
					To:   template_model.TypedRef{Type: "MainObject", Args: []string{"1"}},
				},
			}},
			expectedRoads: []template_model.Road{
				{
					From: template_model.TypedRef{Type: "MainObject", Args: []string{"0"}},
					To:   template_model.TypedRef{Type: "MainObject", Args: []string{"1"}},
				},
				{
					From: template_model.TypedRef{Type: "Connection", Args: []string{"Bridge-A-B"}},
					To:   template_model.TypedRef{Type: "Connection", Args: []string{"Bridge-A-B"}},
				},
			},
		},
		{
			name:      "WhenZoneHasNeitherMainObjectsNorRoads_AppendedRoadSelfReferencesBridge",
			firstZone: template_model.Zone{Name: "Spawn-A"},
			expectedRoads: []template_model.Road{
				{
					From: template_model.TypedRef{Type: "Connection", Args: []string{"Bridge-A-B"}},
					To:   template_model.TypedRef{Type: "Connection", Args: []string{"Bridge-A-B"}},
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			// Arrange
			topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())
			zones := []template_model.Zone{testCase.firstZone, {Name: "Spawn-B"}}

			// Act
			topologyBase.CreateMissingConnections(
				[]string{"A", "B"}, []string{"A", "B"}, positions, zones, nil, newUnitTuning(), nil)

			// Assert
			assert.Equal(t, testCase.expectedRoads, zones[0].Roads)
		})
	}
}

func TestWhenBridgeIsCreated_SecondZoneInSliceAlsoGainsBridgeRoad(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())
	positions := models.Positions{{X: 0.1, Y: 0.1}, {X: 0.9, Y: 0.9}}
	zones := []template_model.Zone{{Name: "Spawn-A"}, {Name: "Spawn-B"}}
	expectedRoads := []template_model.Road{
		{
			From: template_model.TypedRef{Type: "Connection", Args: []string{"Bridge-A-B"}},
			To:   template_model.TypedRef{Type: "Connection", Args: []string{"Bridge-A-B"}},
		},
	}

	// Act
	topologyBase.CreateMissingConnections(
		[]string{"A", "B"}, []string{"A", "B"}, positions, zones, nil, newUnitTuning(), nil)

	// Assert
	assert.Equal(t, expectedRoads, zones[1].Roads)
}

func TestWhenBridgeNameAlreadyExistsOnUnmappedConnection_LoopTerminatesWithoutDuplicateBridge(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())
	positions := models.Positions{{X: 0.1, Y: 0.1}, {X: 0.9, Y: 0.9}}
	existingConnections := []template_model.Connection{
		{Name: "Bridge-A-B", From: "Unknown-X", To: "Unknown-Y", ConnectionType: "Direct"},
	}

	// Act
	connections := topologyBase.CreateMissingConnections(
		[]string{"A", "B"}, []string{"A", "B"}, positions, nil, existingConnections, newUnitTuning(), nil)

	// Assert
	assert.Empty(t, connections)
}
