package topologyConnectionService_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/base"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
	"github.com/stretchr/testify/assert"
)

func TestWhenExistingConnectionHasValidPlayerEndpoints_NoFallbackIsCreated(t *testing.T) {
	t.Parallel()
	// Arrange
	connectionService := base.NewTopologyConnectionService(zones.NewZoneLabelProvider())
	zoneList := []template_model.Zone{{Name: "Spawn-A"}, {Name: "Spawn-B"}}
	existingConnections := []template_model.Connection{{Name: "Custom-A-B", From: "Spawn-A", To: "Spawn-B"}}

	// Act
	connections := connectionService.CreateMissingPlayerConnections(
		[]string{"A", "B"}, zoneList, existingConnections, newUnitTuning(), false)

	// Assert
	assert.Empty(t, connections)
}

func TestWhenExistingConnectionHasMissingEndpoint_FallbackIsCreated(t *testing.T) {
	t.Parallel()
	// Arrange
	connectionService := base.NewTopologyConnectionService(zones.NewZoneLabelProvider())
	zoneList := []template_model.Zone{{Name: "Spawn-A"}, {Name: "Spawn-B"}}
	existingConnections := []template_model.Connection{{Name: "Ghost-A-X", From: "Spawn-A", To: "Unknown-X"}}

	// Act
	connections := connectionService.CreateMissingPlayerConnections(
		[]string{"A", "B"}, zoneList, existingConnections, newUnitTuning(), false)

	// Assert
	assert.Equal(t, []string{"Fallback-A-B"}, connectionNames(connections))
}

func TestWhenExistingConnectionIsSelfReferential_FallbackIsCreated(t *testing.T) {
	t.Parallel()
	// Arrange
	connectionService := base.NewTopologyConnectionService(zones.NewZoneLabelProvider())
	zoneList := []template_model.Zone{{Name: "Spawn-A"}, {Name: "Spawn-B"}}
	existingConnections := []template_model.Connection{{Name: "Loop-A", From: "Spawn-A", To: "Spawn-A"}}

	// Act
	connections := connectionService.CreateMissingPlayerConnections(
		[]string{"A", "B"}, zoneList, existingConnections, newUnitTuning(), false)

	// Assert
	assert.Equal(t, []string{"Fallback-A-B"}, connectionNames(connections))
}

func TestWhenOnePlayerZoneIsMissing_ExistingPlayerZoneDoesNotCreateFallback(t *testing.T) {
	t.Parallel()
	// Arrange
	connectionService := base.NewTopologyConnectionService(zones.NewZoneLabelProvider())
	zoneList := []template_model.Zone{{Name: "Spawn-B"}}

	// Act
	connections := connectionService.CreateMissingPlayerConnections(
		[]string{"A", "B"}, zoneList, nil, newUnitTuning(), false)

	// Assert
	assert.Empty(t, connections)
}

func TestWhenFallbackRoadGenerationIsEnabled_FallbackRoadsAreMaterialized(t *testing.T) {
	t.Parallel()
	// Arrange
	connectionService := base.NewTopologyConnectionService(zones.NewZoneLabelProvider())
	zoneList := []template_model.Zone{{Name: "Spawn-A"}, {Name: "Spawn-B"}}

	// Act
	connectionService.CreateMissingPlayerConnections([]string{"A", "B"}, zoneList, nil, newUnitTuning(), true)

	// Assert
	assert.Len(t, zoneList[0].Roads, 1)
}

func TestWhenFallbackRoadGenerationIsDisabled_FallbackRoadsAreNotMaterialized(t *testing.T) {
	t.Parallel()
	// Arrange
	connectionService := base.NewTopologyConnectionService(zones.NewZoneLabelProvider())
	zoneList := []template_model.Zone{{Name: "Spawn-A"}, {Name: "Spawn-B"}}

	// Act
	connectionService.CreateMissingPlayerConnections([]string{"A", "B"}, zoneList, nil, newUnitTuning(), false)

	// Assert
	assert.Empty(t, zoneList[0].Roads)
}

func TestWhenFallbackNamesAreOccupied_UsesFirstFreeSuffix(t *testing.T) {
	t.Parallel()
	// Arrange
	connectionService := base.NewTopologyConnectionService(zones.NewZoneLabelProvider())
	zoneList := []template_model.Zone{{Name: "Spawn-A"}, {Name: "Spawn-B"}}
	existingConnections := []template_model.Connection{
		{Name: "Fallback-A-B", From: "Unknown-X", To: "Unknown-Y"},
		{Name: "Fallback-A-B-2", From: "Unknown-Y", To: "Unknown-Z"},
	}

	// Act
	connections := connectionService.CreateMissingPlayerConnections(
		[]string{"A", "B"}, zoneList, existingConnections, newUnitTuning(), false)

	// Assert
	assert.Equal(t, []string{"Fallback-A-B-3"}, connectionNames(connections))
}

func TestWhenFallbackNamesAreOccupied_RepeatedRepairIsEmpty(t *testing.T) {
	t.Parallel()
	// Arrange
	connectionService := base.NewTopologyConnectionService(zones.NewZoneLabelProvider())
	zoneList := []template_model.Zone{{Name: "Spawn-A"}, {Name: "Spawn-B"}}
	existingConnections := []template_model.Connection{
		{Name: "Fallback-A-B", From: "Unknown-X", To: "Unknown-Y"},
		{Name: "Fallback-A-B-2", From: "Unknown-Y", To: "Unknown-Z"},
	}
	firstRepair := connectionService.CreateMissingPlayerConnections(
		[]string{"A", "B"}, zoneList, existingConnections, newUnitTuning(), false)

	// Act
	secondRepair := connectionService.CreateMissingPlayerConnections(
		[]string{"A", "B"}, zoneList, append(existingConnections, firstRepair...),
		newUnitTuning(), false)

	// Assert
	assert.Empty(t, secondRepair)
}

func TestWhenFallbackNamesAreOccupied_RoadTargetsMatchRoadGeneration(t *testing.T) {
	t.Parallel()
	// Arrange
	testCases := []struct {
		name            string
		generateRoads   bool
		expectedTargets []string
	}{
		{
			name:            "WhenRoadsAreEnabled_UsesResolvedFallbackName",
			generateRoads:   true,
			expectedTargets: []string{"Fallback-A-B-3", "Fallback-A-B-3"},
		},
		{name: "WhenRoadsAreDisabled_HasNoRoadTargets", expectedTargets: nil},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			// Arrange
			connectionService := base.NewTopologyConnectionService(zones.NewZoneLabelProvider())
			zoneList := []template_model.Zone{{Name: "Spawn-A"}, {Name: "Spawn-B"}}
			existingConnections := []template_model.Connection{
				{Name: "Fallback-A-B", From: "Unknown-X", To: "Unknown-Y"},
				{Name: "Fallback-A-B-2", From: "Unknown-Y", To: "Unknown-Z"},
			}

			// Act
			connectionService.CreateMissingPlayerConnections(
				[]string{"A", "B"}, zoneList, existingConnections, newUnitTuning(),
				testCase.generateRoads)

			// Assert
			assert.Equal(t, testCase.expectedTargets, connectionRoadTargets(zoneList))
		})
	}
}

func TestWhenEachPlayerHasAValidIncidentNeutralEdge_NoFallbackIsCreated(t *testing.T) {
	t.Parallel()
	// Arrange
	connectionService := base.NewTopologyConnectionService(zones.NewZoneLabelProvider())
	zoneList := []template_model.Zone{
		{Name: "Spawn-A"}, {Name: "Spawn-B"}, {Name: "Neutral-N1"}, {Name: "Neutral-N2"},
	}
	existingConnections := []template_model.Connection{
		{Name: "Rnd-A-N1", From: "Spawn-A", To: "Neutral-N1"},
		{Name: "Rnd-B-N2", From: "Spawn-B", To: "Neutral-N2"},
	}

	// Act
	connections := connectionService.CreateMissingPlayerConnections(
		[]string{"A", "B"}, zoneList, existingConnections, newUnitTuning(), false)

	// Assert
	assert.Empty(t, connections)
}

func connectionNames(connections []template_model.Connection) []string {
	names := make([]string, len(connections))
	for index, connection := range connections {
		names[index] = connection.Name
	}
	return names
}
