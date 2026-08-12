package webTopology_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
)

func TestWhenTwoPlayersAndThreeNeutralPlansProvided_CreatesZonePerLabel(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologySharedWeb
	playerLabels := []string{"A", "B"}
	neutralZones := neutral_zone.Plans{}
	neutralZones.AddPlan("N1", neutral_zone.QualityMedium, 1)
	neutralZones.AddPlan("N2", neutral_zone.QualityMedium, 1)
	neutralZones.AddPlan("N3", neutral_zone.QualityMedium, 1)
	tuning := test_helpers.NewGenerationTuning(configuration, 5)
	service := topology.NewSharedWebTopologyService(test_helpers.NewZoneFactories())

	// Act
	variant := service.CreateTopologyVariant(*configuration, playerLabels, neutralZones, tuning, "")

	// Assert
	assert.Len(t, variant.Zones, 5)
}

func TestWhenThreeNeutralZonesFormTheRing_CreatesNeutralRingConnectionPerPair(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologySharedWeb
	playerLabels := []string{"A", "B"}
	neutralZones := neutral_zone.Plans{}
	neutralZones.AddPlan("N1", neutral_zone.QualityMedium, 1)
	neutralZones.AddPlan("N2", neutral_zone.QualityMedium, 1)
	neutralZones.AddPlan("N3", neutral_zone.QualityMedium, 1)
	tuning := test_helpers.NewGenerationTuning(configuration, 5)
	service := topology.NewSharedWebTopologyService(test_helpers.NewZoneFactories())

	// Act
	variant := service.CreateTopologyVariant(*configuration, playerLabels, neutralZones, tuning, "")

	// Assert
	assert.Equal(t, 3, countConnectionsWithPrefix(variant, "NRing-"))
}

func TestWhenPlayersAttachToTheNeutralRing_CreatesTwoWebSpokesPerPlayer(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologySharedWeb
	playerLabels := []string{"A", "B"}
	neutralZones := neutral_zone.Plans{}
	neutralZones.AddPlan("N1", neutral_zone.QualityMedium, 1)
	neutralZones.AddPlan("N2", neutral_zone.QualityMedium, 1)
	neutralZones.AddPlan("N3", neutral_zone.QualityMedium, 1)
	tuning := test_helpers.NewGenerationTuning(configuration, 5)
	service := topology.NewSharedWebTopologyService(test_helpers.NewZoneFactories())

	// Act
	variant := service.CreateTopologyVariant(*configuration, playerLabels, neutralZones, tuning, "")

	// Assert
	assert.Equal(t, 4, countConnectionsWithPrefix(variant, "Web-"))
}

func TestWhenOnlyOneNeutralZoneExists_CreatesNoNeutralRingConnections(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologySharedWeb
	playerLabels := []string{"A", "B"}
	neutralZones := neutral_zone.Plans{}
	neutralZones.AddPlan("N1", neutral_zone.QualityMedium, 1)
	tuning := test_helpers.NewGenerationTuning(configuration, 3)
	service := topology.NewSharedWebTopologyService(test_helpers.NewZoneFactories())

	// Act
	variant := service.CreateTopologyVariant(*configuration, playerLabels, neutralZones, tuning, "")

	// Assert
	assert.Zero(t, countConnectionsWithPrefix(variant, "NRing-"))
}

func TestWhenWebIsBuilt_EveryConnectionReferencesExistingZones(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologySharedWeb
	playerLabels := []string{"A", "B"}
	neutralZones := neutral_zone.Plans{}
	neutralZones.AddPlan("N1", neutral_zone.QualityMedium, 1)
	neutralZones.AddPlan("N2", neutral_zone.QualityMedium, 1)
	neutralZones.AddPlan("N3", neutral_zone.QualityMedium, 1)
	tuning := test_helpers.NewGenerationTuning(configuration, 5)
	service := topology.NewSharedWebTopologyService(test_helpers.NewZoneFactories())

	// Act
	variant := service.CreateTopologyVariant(*configuration, playerLabels, neutralZones, tuning, "")

	// Assert
	assert.Empty(t, danglingConnectionNames(variant))
}

func TestWhenCirclesTopologySelected_BalancesNeutralLabelsAcrossPlayers(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyCircles
	playerLabels := []string{"A", "B"}
	neutralZones := neutral_zone.Plans{}
	neutralZones.AddPlan("N1", neutral_zone.QualityMedium, 1)
	neutralZones.AddPlan("N2", neutral_zone.QualityMedium, 1)
	neutralZones.AddPlan("N3", neutral_zone.QualityMedium, 1)
	tuning := test_helpers.NewGenerationTuning(configuration, 5)
	service := topology.NewSharedWebTopologyService(test_helpers.NewZoneFactories())

	// Act
	variant := service.CreateTopologyVariant(*configuration, playerLabels, neutralZones, tuning, "")

	// Assert
	assert.Len(t, variant.Zones, 5)
}

func TestWhenPlayerConnectionsAreForbidden_NoDirectConnectionJoinsTwoSpawnZones(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologySharedWeb
	configuration.NoDirectPlayerConnections = true
	playerLabels := []string{"A", "B"}
	neutralZones := neutral_zone.Plans{}
	neutralZones.AddPlan("N1", neutral_zone.QualityMedium, 1)
	neutralZones.AddPlan("N2", neutral_zone.QualityMedium, 1)
	tuning := test_helpers.NewGenerationTuning(configuration, 4)
	service := topology.NewSharedWebTopologyService(test_helpers.NewZoneFactories())

	// Act
	variant := service.CreateTopologyVariant(*configuration, playerLabels, neutralZones, tuning, "")

	// Assert
	assert.Empty(t, directSpawnToSpawnNames(variant))
}

func TestWhenRandomPortalsEnabled_AddsPortalConnections(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologySharedWeb
	configuration.RandomPortals = true
	playerLabels := []string{"A", "B"}
	neutralZones := neutral_zone.Plans{}
	neutralZones.AddPlan("N1", neutral_zone.QualityMedium, 1)
	neutralZones.AddPlan("N2", neutral_zone.QualityMedium, 1)
	neutralZones.AddPlan("N3", neutral_zone.QualityMedium, 1)
	neutralZones.AddPlan("N4", neutral_zone.QualityMedium, 1)
	tuning := test_helpers.NewGenerationTuning(configuration, 6)
	service := topology.NewSharedWebTopologyService(test_helpers.NewZoneFactories())

	// Act
	variant := service.CreateTopologyVariant(*configuration, playerLabels, neutralZones, tuning, "")

	// Assert
	assert.NotZero(t, countPortalConnections(variant))
}
