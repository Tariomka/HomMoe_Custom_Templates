package topologyBase_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/base"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
)

func TestWhenFewerThanTwoPlayerLabelsExist_NoFallbackConnectionsAreCreated(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())
	zones := []template_model.Zone{{Name: "Spawn-A"}}

	// Act
	connections := topologyBase.CreateMissingPlayerConnections(
		[]string{"A"}, zones, nil, newUnitTuning(), true)

	// Assert
	assert.Empty(t, connections)
}

func TestWhenBothSpawnZonesLackConnections_SingleSharedFallbackLinksThem(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())
	zones := []template_model.Zone{{Name: "Spawn-A"}, {Name: "Spawn-B"}}
	expectedConnections := []template_model.Connection{
		{
			Name: "Fallback-A-B", From: "Spawn-A", To: "Spawn-B",
			ConnectionType: "Direct", GuardZone: "Spawn-A", SimTurnSquad: true,
			GuardValue: 30000, GuardWeeklyIncrement: 0.15,
			GuardMatchGroup: "fallback_guard_Fallback-A-B",
		},
	}

	// Act
	connections := topologyBase.CreateMissingPlayerConnections(
		[]string{"A", "B"}, zones, nil, newUnitTuning(), true)

	// Assert
	assert.Equal(t, expectedConnections, connections)
}

func TestWhenSpawnZonesHaveKnownConnectionEndpoints_NoFallbackIsCreated(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())
	zones := []template_model.Zone{
		{Name: "Spawn-A", Roads: []template_model.Road{
			{To: template_model.TypedRef{Type: "Connection", Args: []string{"Ring-A-B"}}},
		}},
		{Name: "Spawn-B", Roads: []template_model.Road{
			{To: template_model.TypedRef{Type: "Connection", Args: []string{"Ring-A-B"}}},
		}},
	}
	existingConnections := []template_model.Connection{{Name: "Ring-A-B", From: "Spawn-A", To: "Spawn-B"}}

	// Act
	connections := topologyBase.CreateMissingPlayerConnections(
		[]string{"A", "B"}, zones, existingConnections, newUnitTuning(), true)

	// Assert
	assert.Empty(t, connections)
}

func TestWhenPlayersAreConnectedThroughNeutralEndpoints_NoFallbackIsCreated(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())
	zones := []template_model.Zone{{Name: "Spawn-A"}, {Name: "Neutral-N"}, {Name: "Spawn-B"}}
	existingConnections := []template_model.Connection{
		{Name: "Rnd-A-N", From: "Spawn-A", To: "Neutral-N"},
		{Name: "Rnd-N-B", From: "Neutral-N", To: "Spawn-B"},
	}

	// Act
	connections := topologyBase.CreateMissingPlayerConnections(
		[]string{"A", "B"}, zones, existingConnections, newUnitTuning(), true)

	// Assert
	assert.Empty(t, connections)
}

func TestWhenPlayersAreConnectedThroughReversedNeutralEndpoints_NoFallbackIsCreated(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())
	zones := []template_model.Zone{{Name: "Spawn-A"}, {Name: "Neutral-N"}, {Name: "Spawn-B"}}
	existingConnections := []template_model.Connection{
		{Name: "Rnd-N-A", From: "Neutral-N", To: "Spawn-A"},
		{Name: "Rnd-B-N", From: "Spawn-B", To: "Neutral-N"},
	}

	// Act
	connections := topologyBase.CreateMissingPlayerConnections(
		[]string{"A", "B"}, zones, existingConnections, newUnitTuning(), true)

	// Assert
	assert.Empty(t, connections)
}

func TestWhenUnnamedConnectionsLinkPlayersThroughNeutral_NoFallbackIsCreated(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())
	zones := []template_model.Zone{{Name: "Spawn-A"}, {Name: "Neutral-N"}, {Name: "Spawn-B"}}
	existingConnections := []template_model.Connection{
		{From: "Spawn-A", To: "Neutral-N"},
		{From: "Neutral-N", To: "Spawn-B"},
	}

	// Act
	connections := topologyBase.CreateMissingPlayerConnections(
		[]string{"A", "B"}, zones, existingConnections, newUnitTuning(), true)

	// Assert
	assert.Empty(t, connections)
}

func TestWhenRoadReferencesUnrelatedKnownConnection_OnlyIsolatedPlayerIsRepaired(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())
	zones := []template_model.Zone{
		{Name: "Spawn-A", Roads: []template_model.Road{{
			To: template_model.TypedRef{Type: "Connection", Args: []string{"Rnd-B-C"}},
		}}},
		{Name: "Spawn-B"},
		{Name: "Spawn-C"},
	}
	existingConnections := []template_model.Connection{{Name: "Rnd-B-C", From: "Spawn-B", To: "Spawn-C"}}

	// Act
	connections := topologyBase.CreateMissingPlayerConnections(
		[]string{"A", "B", "C"}, zones, existingConnections, newUnitTuning(), true)

	// Assert
	assert.Equal(t, []string{"Fallback-A-B"}, connectionNames(connections))
}

func TestWhenEarlierFallbackRepairsBothIsolatedPlayers_NoRedundantFallbackIsCreated(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())
	zones := []template_model.Zone{
		{Name: "Spawn-A"},
		{Name: "Spawn-B"},
		{Name: "Spawn-C"},
		{Name: "Neutral-N"},
	}
	existingConnections := []template_model.Connection{{Name: "Rnd-C-N", From: "Spawn-C", To: "Neutral-N"}}

	// Act
	connections := topologyBase.CreateMissingPlayerConnections(
		[]string{"A", "B", "C"}, zones, existingConnections, newUnitTuning(), true)

	// Assert
	assert.Equal(t, []string{"Fallback-A-B"}, connectionNames(connections))
}

func TestWhenRepairOrderIsNotAlphabetical_UsesAccumulatedFallbackEndpoints(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())
	zones := []template_model.Zone{
		{Name: "Spawn-B"},
		{Name: "Spawn-A"},
		{Name: "Spawn-C"},
		{Name: "Neutral-N"},
	}
	existingConnections := []template_model.Connection{{Name: "Rnd-A-N", From: "Spawn-A", To: "Neutral-N"}}

	// Act
	connections := topologyBase.CreateMissingPlayerConnections(
		[]string{"B", "A", "C"}, zones, existingConnections, newUnitTuning(), true)

	// Assert
	assert.Equal(t, []string{"Spawn-B->Spawn-A", "Spawn-C->Spawn-B"}, connectionEndpoints(connections))
}

func TestWhenZoneRoadReferencesUnknownConnection_FallbackIsStillCreated(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())
	zones := []template_model.Zone{
		{Name: "Spawn-A", Roads: []template_model.Road{
			{To: template_model.TypedRef{Type: "Connection", Args: []string{"Ghost-Conn"}}},
		}},
		{Name: "Spawn-B"},
	}

	// Act
	connections := topologyBase.CreateMissingPlayerConnections(
		[]string{"A", "B"}, zones, nil, newUnitTuning(), true)

	// Assert
	assert.Len(t, connections, 1)
}

func TestWhenSpawnZonesAreMissingFromZoneList_NoFallbacksAreCreated(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())

	// Act
	connections := topologyBase.CreateMissingPlayerConnections(
		[]string{"A", "B"}, nil, nil, newUnitTuning(), true)

	// Assert
	assert.Empty(t, connections)
}

func TestWhenPlayerLabelsAreReversed_FallbackNameStillSortsLabelsAlphabetically(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())
	zones := []template_model.Zone{{Name: "Spawn-B"}, {Name: "Spawn-A"}}
	expectedConnections := []template_model.Connection{
		{
			Name: "Fallback-A-B", From: "Spawn-B", To: "Spawn-A",
			ConnectionType: "Direct", GuardZone: "Spawn-B", SimTurnSquad: true,
			GuardValue: 30000, GuardWeeklyIncrement: 0.15,
			GuardMatchGroup: "fallback_guard_Fallback-A-B",
		},
	}

	// Act
	connections := topologyBase.CreateMissingPlayerConnections(
		[]string{"B", "A"}, zones, nil, newUnitTuning(), true)

	// Assert
	assert.Equal(t, expectedConnections, connections)
}

func TestWhenBorderGuardMultiplierIsDoubled_FallbackGuardValueIsScaled(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())
	tuning := newUnitTuning()
	tuning.BorderGuardStrengthMultiplier = 2.0
	zones := []template_model.Zone{{Name: "Spawn-A"}, {Name: "Spawn-B"}}

	// Act
	connections := topologyBase.CreateMissingPlayerConnections(
		[]string{"A", "B"}, zones, nil, tuning, true)

	// Assert
	assert.Equal(t, 60000, connections[0].GuardValue)
}

func TestWhenFallbackConnectionIsCreated_BothSpawnZonesInSliceGainFallbackRoad(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())
	zones := []template_model.Zone{{Name: "Spawn-A"}, {Name: "Spawn-B"}}
	expectedRoad := template_model.Road{
		From: template_model.TypedRef{Type: "MainObject", Args: []string{"0"}},
		To:   template_model.TypedRef{Type: "Connection", Args: []string{"Fallback-A-B"}},
	}
	expectedZones := []template_model.Zone{
		{Name: "Spawn-A", Roads: []template_model.Road{expectedRoad}},
		{Name: "Spawn-B", Roads: []template_model.Road{expectedRoad}},
	}

	// Act
	topologyBase.CreateMissingPlayerConnections([]string{"A", "B"}, zones, nil, newUnitTuning(), true)

	// Assert
	assert.Equal(t, expectedZones, zones)
}

func TestWhenFallbackRoadGenerationIsDisabled_ExistingConnectionsAreRetained(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())
	zones := []template_model.Zone{{Name: "Spawn-A"}, {Name: "Spawn-B"}}
	existingConnections := []template_model.Connection{{Name: "Custom-A-N", From: "Spawn-A", To: "Neutral-N"}}

	// Act
	connections := topologyBase.CreateMissingPlayerConnections(
		[]string{"A", "B"}, zones, existingConnections, newUnitTuning(), false)

	// Assert
	assert.Equal(t, []string{"Fallback-A-B"}, connectionNames(connections))
}

func TestWhenFallbackRoadGenerationIsDisabled_WholeZonesIncludingCustomRoadsArePreserved(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())
	expectedZones := []template_model.Zone{
		{Name: "Spawn-A", Roads: []template_model.Road{{
			From: template_model.TypedRef{Type: "MainObject", Args: []string{"0"}},
			To:   template_model.TypedRef{Type: "MainObject", Args: []string{"1"}},
		}}},
		{Name: "Spawn-B", Roads: []template_model.Road{{
			From: template_model.TypedRef{Type: "MainObject", Args: []string{"0"}},
			To:   template_model.TypedRef{Type: "MainObject", Args: []string{"2"}},
		}}},
	}
	zones := append([]template_model.Zone(nil), expectedZones...)

	// Act
	topologyBase.CreateMissingPlayerConnections([]string{"A", "B"}, zones, nil, newUnitTuning(), false)

	// Assert
	assert.Equal(t, expectedZones, zones)
}

func connectionNames(connections []template_model.Connection) []string {
	names := make([]string, len(connections))
	for index, connection := range connections {
		names[index] = connection.Name
	}
	return names
}

func connectionEndpoints(connections []template_model.Connection) []string {
	endpoints := make([]string, len(connections))
	for index, connection := range connections {
		endpoints[index] = connection.From + "->" + connection.To
	}
	return endpoints
}
