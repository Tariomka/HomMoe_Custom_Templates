package topologyBase_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/base"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
)

func TestWhenFewerThanTwoLabelsExist_NoPortalsAreCreated(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())

	// Act
	connections := topologyBase.CreateRandomPortalConnections(
		[]string{"A"}, []string{"A"}, newUnitTuning(), 10, nil)

	// Assert
	assert.Empty(t, connections)
}

func TestWhenMaxCountIsBelowLabelCount_PortalCountEqualsMaxCount(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())
	orderedLabels := []string{"A", "B", "C", "D", "E"}

	// Act
	connections := topologyBase.CreateRandomPortalConnections(
		[]string{"A", "B"}, orderedLabels, newUnitTuning(), 2, nil)

	// Assert
	assert.Len(t, connections, 2)
}

func TestWhenMaxCountExceedsLabelCount_EveryLabelGetsOnePortal(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())
	orderedLabels := []string{"A", "B", "C", "D"}

	// Act
	connections := topologyBase.CreateRandomPortalConnections(
		[]string{"A", "B"}, orderedLabels, newUnitTuning(), 10, nil)

	// Assert
	assert.Len(t, connections, 4)
}

func TestWhenOnlyTwoZonesExist_PortalsLinkThemInBothDirections(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())
	portalRoad := true
	crossroadsRule := template_model.PlacementRule{Type: "Crossroads", TargetMin: 0.1, TargetMax: 0.25, Weight: 2}
	expectedConnections := []template_model.Connection{
		{
			Name: "Portal-A-B", From: "Spawn-A", To: "Spawn-B",
			ConnectionType: "Portal", Road: &portalRoad,
			GuardValue: 30000, GuardWeeklyIncrement: 0.15,
			PortalPlacementRulesFrom: []template_model.PlacementRule{crossroadsRule},
			PortalPlacementRulesTo:   []template_model.PlacementRule{crossroadsRule},
		},
		{
			Name: "Portal-B-A", From: "Spawn-B", To: "Spawn-A",
			ConnectionType: "Portal", Road: &portalRoad,
			GuardValue: 30000, GuardWeeklyIncrement: 0.15,
			PortalPlacementRulesFrom: []template_model.PlacementRule{crossroadsRule},
			PortalPlacementRulesTo:   []template_model.PlacementRule{crossroadsRule},
		},
	}

	// Act
	connections := topologyBase.CreateRandomPortalConnections(
		[]string{"A", "B"}, []string{"A", "B"}, newUnitTuning(), 10, nil)

	// Assert
	assert.ElementsMatch(t, expectedConnections, connections)
}

func TestWhenPortalJoinsTwoPlayerZones_UsesThePlayerBorderGuardValue(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())
	var guardValues []int

	// Act
	connections := topologyBase.CreateRandomPortalConnections(
		[]string{"A", "B"}, []string{"A", "B"}, newUnitTuning(), 10, nil)

	// Assert
	for _, connection := range connections {
		guardValues = append(guardValues, connection.GuardValue)
	}
	assert.Equal(t, []int{30000, 30000}, guardValues)
}

func TestWhenPortalJoinsTwoNeutralZones_UsesTheHigherTierGuardValue(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())
	neutralPlans := neutral_zone.Plans{
		{Label: "C", Quality: neutral_zone.QualityLowest},
		{Label: "D", Quality: neutral_zone.QualityHigh},
	}
	var guardValues []int

	// Act
	connections := topologyBase.CreateRandomPortalConnections(
		[]string{"A", "B"}, []string{"C", "D"}, newUnitTuning(), 10, neutralPlans)

	// Assert
	for _, connection := range connections {
		guardValues = append(guardValues, connection.GuardValue)
	}
	assert.Equal(t, []int{25000, 25000}, guardValues)
}

func TestWhenManyLabelsArePortalLinked_NoPortalLinksZoneToItself(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())
	orderedLabels := []string{"A", "B", "C", "D", "E", "F"}
	var selfLinkedNames []string

	// Act
	connections := topologyBase.CreateRandomPortalConnections(
		[]string{"A", "B"}, orderedLabels, newUnitTuning(), 6, nil)

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
	topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())
	expectedZoneNames := []string{"Spawn-A", "Spawn-B", "Neutral-C", "Neutral-D"}
	var endpointNames []string

	// Act
	connections := topologyBase.CreateRandomPortalConnections(
		[]string{"A", "B"}, []string{"A", "B", "C", "D"}, newUnitTuning(), 4, nil)

	// Assert
	for _, connection := range connections {
		endpointNames = append(endpointNames, connection.From, connection.To)
	}
	assert.Subset(t, expectedZoneNames, endpointNames)
}
