package topologyBase_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/base"
	"github.com/stretchr/testify/assert"
)

func TestWhenFewerThanTwoLabelsExist_NoPortalsAreCreated(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase()

	// Act
	connections := topologyBase.CreateRandomPortalConnections(
		[]string{"A"}, []string{"A"}, newUnitTuning(), 10)

	// Assert
	assert.Empty(t, connections)
}

func TestWhenMaxCountIsBelowLabelCount_PortalCountEqualsMaxCount(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase()
	orderedLabels := []string{"A", "B", "C", "D", "E"}

	// Act
	connections := topologyBase.CreateRandomPortalConnections(
		[]string{"A", "B"}, orderedLabels, newUnitTuning(), 2)

	// Assert
	assert.Len(t, connections, 2)
}

func TestWhenMaxCountExceedsLabelCount_EveryLabelGetsOnePortal(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase()
	orderedLabels := []string{"A", "B", "C", "D"}

	// Act
	connections := topologyBase.CreateRandomPortalConnections(
		[]string{"A", "B"}, orderedLabels, newUnitTuning(), 10)

	// Assert
	assert.Len(t, connections, 4)
}

func TestWhenOnlyTwoZonesExist_PortalsLinkThemInBothDirections(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase()
	portalRoad := true
	crossroadsRule := entities.PlacementRule{Type: "Crossroads", TargetMin: 0.1, TargetMax: 0.25, Weight: 2}
	expectedConnections := []entities.Connection{
		{
			Name: "Portal-A-B", From: "Spawn-A", To: "Spawn-B",
			ConnectionType: "Portal", Road: &portalRoad,
			GuardValue: 25000, GuardWeeklyIncrement: 0.15,
			PortalPlacementRulesFrom: []entities.PlacementRule{crossroadsRule},
			PortalPlacementRulesTo:   []entities.PlacementRule{crossroadsRule},
		},
		{
			Name: "Portal-B-A", From: "Spawn-B", To: "Spawn-A",
			ConnectionType: "Portal", Road: &portalRoad,
			GuardValue: 25000, GuardWeeklyIncrement: 0.15,
			PortalPlacementRulesFrom: []entities.PlacementRule{crossroadsRule},
			PortalPlacementRulesTo:   []entities.PlacementRule{crossroadsRule},
		},
	}

	// Act
	connections := topologyBase.CreateRandomPortalConnections(
		[]string{"A", "B"}, []string{"A", "B"}, newUnitTuning(), 10)

	// Assert
	assert.ElementsMatch(t, expectedConnections, connections)
}

func TestWhenManyLabelsArePortalLinked_NoPortalLinksZoneToItself(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase()
	orderedLabels := []string{"A", "B", "C", "D", "E", "F"}
	var selfLinkedNames []string

	// Act
	connections := topologyBase.CreateRandomPortalConnections(
		[]string{"A", "B"}, orderedLabels, newUnitTuning(), 6)

	// Assert
	for _, connection := range connections {
		if connection.From == connection.To {
			selfLinkedNames = append(selfLinkedNames, connection.Name)
		}
	}
	assert.Empty(t, selfLinkedNames)
}

func TestWhenLabelsMixPlayersAndNeutrals_PortalEndpointsUseMatchingZonePrefixes(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase()
	expectedZoneNames := []string{"Spawn-A", "Spawn-B", "Neutral-C", "Neutral-D"}
	var endpointNames []string

	// Act
	connections := topologyBase.CreateRandomPortalConnections(
		[]string{"A", "B"}, []string{"A", "B", "C", "D"}, newUnitTuning(), 4)

	// Assert
	for _, connection := range connections {
		endpointNames = append(endpointNames, connection.From, connection.To)
	}
	assert.Subset(t, expectedZoneNames, endpointNames)
}
