package topologyConnectionService_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/base"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
	"github.com/stretchr/testify/assert"
)

func TestWhenExistingConnectionHasValidEndpoints_NoBridgeIsCreated(t *testing.T) {
	t.Parallel()
	// Arrange
	connectionService := base.NewTopologyConnectionService(zones.NewZoneLabelProvider())
	positions := models.Positions{{X: 0.1, Y: 0.1}, {X: 0.9, Y: 0.9}}
	existingConnections := []template_model.Connection{
		{Name: "Custom-A-B", From: "Spawn-A", To: "Spawn-B", ConnectionType: "Direct"},
	}

	// Act
	connections := connectionService.CreateMissingConnections(
		[]string{"A", "B"}, []string{"A", "B"}, positions, nil, existingConnections, newUnitTuning(), nil, false)

	// Assert
	assert.Empty(t, connections)
}

func TestWhenExistingConnectionHasMissingEndpoint_BridgeIsCreated(t *testing.T) {
	t.Parallel()
	// Arrange
	connectionService := base.NewTopologyConnectionService(zones.NewZoneLabelProvider())
	positions := models.Positions{{X: 0.1, Y: 0.1}, {X: 0.9, Y: 0.9}}
	existingConnections := []template_model.Connection{
		{Name: "Ghost-A-X", From: "Spawn-A", To: "Unknown-X", ConnectionType: "Direct"},
	}

	// Act
	connections := connectionService.CreateMissingConnections(
		[]string{"A", "B"}, []string{"A", "B"}, positions, nil, existingConnections, newUnitTuning(), nil, false)

	// Assert
	assert.Equal(t, []string{"Bridge-A-B"}, connectionNames(connections))
}

func TestWhenExistingConnectionIsSelfReferential_BridgeIsCreated(t *testing.T) {
	t.Parallel()
	// Arrange
	connectionService := base.NewTopologyConnectionService(zones.NewZoneLabelProvider())
	positions := models.Positions{{X: 0.1, Y: 0.1}, {X: 0.9, Y: 0.9}}
	existingConnections := []template_model.Connection{
		{Name: "Loop-A", From: "Spawn-A", To: "Spawn-A", ConnectionType: "Direct"},
	}

	// Act
	connections := connectionService.CreateMissingConnections(
		[]string{"A", "B"}, []string{"A", "B"}, positions, nil, existingConnections, newUnitTuning(), nil, false)

	// Assert
	assert.Equal(t, []string{"Bridge-A-B"}, connectionNames(connections))
}

func TestWhenBridgeRoadGenerationIsEnabled_BridgeRoadsAreMaterialized(t *testing.T) {
	t.Parallel()
	// Arrange
	connectionService := base.NewTopologyConnectionService(zones.NewZoneLabelProvider())
	positions := models.Positions{{X: 0.1, Y: 0.1}, {X: 0.9, Y: 0.9}}
	zoneList := []template_model.Zone{{Name: "Spawn-A"}, {Name: "Spawn-B"}}

	// Act
	connectionService.CreateMissingConnections(
		[]string{"A", "B"}, []string{"A", "B"}, positions, zoneList, nil, newUnitTuning(), nil, true)

	// Assert
	assert.Len(t, zoneList[0].Roads, 1)
}

func TestWhenBridgeRoadGenerationIsDisabled_BridgeRoadsAreNotMaterialized(t *testing.T) {
	t.Parallel()
	// Arrange
	connectionService := base.NewTopologyConnectionService(zones.NewZoneLabelProvider())
	positions := models.Positions{{X: 0.1, Y: 0.1}, {X: 0.9, Y: 0.9}}
	zoneList := []template_model.Zone{{Name: "Spawn-A"}, {Name: "Spawn-B"}}

	// Act
	connectionService.CreateMissingConnections(
		[]string{"A", "B"}, []string{"A", "B"}, positions, zoneList, nil, newUnitTuning(), nil, false)

	// Assert
	assert.Empty(t, zoneList[0].Roads)
}

func TestWhenBridgeNamesAreOccupied_UsesFirstFreeSuffix(t *testing.T) {
	t.Parallel()
	// Arrange
	connectionService := base.NewTopologyConnectionService(zones.NewZoneLabelProvider())
	positions := models.Positions{{X: 0.1, Y: 0.1}, {X: 0.9, Y: 0.9}}
	existingConnections := []template_model.Connection{
		{Name: "Bridge-A-B", From: "Unknown-X", To: "Unknown-Y", ConnectionType: "Direct"},
		{Name: "Bridge-A-B-2", From: "Unknown-Y", To: "Unknown-Z", ConnectionType: "Direct"},
	}

	// Act
	connections := connectionService.CreateMissingConnections(
		[]string{"A", "B"}, []string{"A", "B"}, positions, nil,
		existingConnections, newUnitTuning(), nil, false)

	// Assert
	assert.Equal(t, []string{"Bridge-A-B-3"}, connectionNames(connections))
}

func TestWhenBridgeNamesAreOccupied_RepeatedRepairIsEmpty(t *testing.T) {
	t.Parallel()
	// Arrange
	connectionService := base.NewTopologyConnectionService(zones.NewZoneLabelProvider())
	positions := models.Positions{{X: 0.1, Y: 0.1}, {X: 0.9, Y: 0.9}}
	existingConnections := []template_model.Connection{
		{Name: "Bridge-A-B", From: "Unknown-X", To: "Unknown-Y", ConnectionType: "Direct"},
		{Name: "Bridge-A-B-2", From: "Unknown-Y", To: "Unknown-Z", ConnectionType: "Direct"},
	}
	firstRepair := connectionService.CreateMissingConnections(
		[]string{"A", "B"}, []string{"A", "B"}, positions, nil,
		existingConnections, newUnitTuning(), nil, false)

	// Act
	secondRepair := connectionService.CreateMissingConnections(
		[]string{"A", "B"}, []string{"A", "B"}, positions, nil,
		append(existingConnections, firstRepair...), newUnitTuning(), nil, false)

	// Assert
	assert.Empty(t, secondRepair)
}

func TestWhenBridgeNamesAreOccupied_RoadTargetsMatchRoadGeneration(t *testing.T) {
	t.Parallel()
	// Arrange
	testCases := []struct {
		name            string
		generateRoads   bool
		expectedTargets []string
	}{
		{
			name:          "WhenRoadsAreEnabled_UsesResolvedBridgeName",
			generateRoads: true,
			expectedTargets: []string{
				"Bridge-A-B-3", "Bridge-A-B-3", "Bridge-A-B-3", "Bridge-A-B-3",
			},
		},
		{name: "WhenRoadsAreDisabled_HasNoRoadTargets", expectedTargets: nil},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			// Arrange
			connectionService := base.NewTopologyConnectionService(zones.NewZoneLabelProvider())
			positions := models.Positions{{X: 0.1, Y: 0.1}, {X: 0.9, Y: 0.9}}
			zoneList := []template_model.Zone{{Name: "Spawn-A"}, {Name: "Spawn-B"}}
			existingConnections := []template_model.Connection{
				{Name: "Bridge-A-B", From: "Unknown-X", To: "Unknown-Y", ConnectionType: "Direct"},
				{Name: "Bridge-A-B-2", From: "Unknown-Y", To: "Unknown-Z", ConnectionType: "Direct"},
			}

			// Act
			connectionService.CreateMissingConnections(
				[]string{"A", "B"}, []string{"A", "B"}, positions, zoneList,
				existingConnections, newUnitTuning(), nil, testCase.generateRoads)

			// Assert
			assert.Equal(t, testCase.expectedTargets, connectionRoadTargets(zoneList))
		})
	}
}

func TestWhenThreeDisconnectedZonesAreRepaired_TwoBridgesAreCreated(t *testing.T) {
	t.Parallel()
	// Arrange
	connectionService := base.NewTopologyConnectionService(zones.NewZoneLabelProvider())
	labels := []string{"A", "B", "C"}
	positions := models.Positions{{X: 0.1, Y: 0.1}, {X: 0.5, Y: 0.5}, {X: 0.9, Y: 0.9}}

	// Act
	connections := connectionService.CreateMissingConnections(
		labels, labels, positions, nil, nil, newUnitTuning(), nil, false)

	// Assert
	assert.Len(t, connections, 2)
}

func TestWhenAccumulatedBridgeEndpointsConnectAllZones_NoFurtherBridgesAreCreated(t *testing.T) {
	t.Parallel()
	// Arrange
	connectionService := base.NewTopologyConnectionService(zones.NewZoneLabelProvider())
	labels := []string{"A", "B", "C"}
	positions := models.Positions{{X: 0.1, Y: 0.1}, {X: 0.5, Y: 0.5}, {X: 0.9, Y: 0.9}}
	firstRepair := connectionService.CreateMissingConnections(
		labels, labels, positions, nil, nil, newUnitTuning(), nil, false)

	// Act
	secondRepair := connectionService.CreateMissingConnections(
		labels, labels, positions, nil, firstRepair, newUnitTuning(), nil, false)

	// Assert
	assert.Empty(t, secondRepair)
}
