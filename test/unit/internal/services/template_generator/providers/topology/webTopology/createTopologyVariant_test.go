package webTopology_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology"
	"github.com/stretchr/testify/assert"
)

func TestWhenTwoPlayersAndThreeNeutralPlansProvided_CreatesZonePerLabel(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologySharedWeb
	playerLabels := []string{"A", "B"}
	neutralZones := models.NeutralZonePlans{}
	neutralZones.AddPlan("N1", models.QualityMedium, 1)
	neutralZones.AddPlan("N2", models.QualityMedium, 1)
	neutralZones.AddPlan("N3", models.QualityMedium, 1)
	tuning := models.NewGenerationTuning(configuration, 5)
	service := topology.NewSharedWebTopologyService()

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
	neutralZones := models.NeutralZonePlans{}
	neutralZones.AddPlan("N1", models.QualityMedium, 1)
	neutralZones.AddPlan("N2", models.QualityMedium, 1)
	neutralZones.AddPlan("N3", models.QualityMedium, 1)
	tuning := models.NewGenerationTuning(configuration, 5)
	service := topology.NewSharedWebTopologyService()

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
	neutralZones := models.NeutralZonePlans{}
	neutralZones.AddPlan("N1", models.QualityMedium, 1)
	neutralZones.AddPlan("N2", models.QualityMedium, 1)
	neutralZones.AddPlan("N3", models.QualityMedium, 1)
	tuning := models.NewGenerationTuning(configuration, 5)
	service := topology.NewSharedWebTopologyService()

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
	neutralZones := models.NeutralZonePlans{}
	neutralZones.AddPlan("N1", models.QualityMedium, 1)
	tuning := models.NewGenerationTuning(configuration, 3)
	service := topology.NewSharedWebTopologyService()

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
	neutralZones := models.NeutralZonePlans{}
	neutralZones.AddPlan("N1", models.QualityMedium, 1)
	neutralZones.AddPlan("N2", models.QualityMedium, 1)
	neutralZones.AddPlan("N3", models.QualityMedium, 1)
	tuning := models.NewGenerationTuning(configuration, 5)
	service := topology.NewSharedWebTopologyService()

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
	neutralZones := models.NeutralZonePlans{}
	neutralZones.AddPlan("N1", models.QualityMedium, 1)
	neutralZones.AddPlan("N2", models.QualityMedium, 1)
	neutralZones.AddPlan("N3", models.QualityMedium, 1)
	tuning := models.NewGenerationTuning(configuration, 5)
	service := topology.NewSharedWebTopologyService()

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
	neutralZones := models.NeutralZonePlans{}
	neutralZones.AddPlan("N1", models.QualityMedium, 1)
	neutralZones.AddPlan("N2", models.QualityMedium, 1)
	tuning := models.NewGenerationTuning(configuration, 4)
	service := topology.NewSharedWebTopologyService()

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
	neutralZones := models.NeutralZonePlans{}
	neutralZones.AddPlan("N1", models.QualityMedium, 1)
	neutralZones.AddPlan("N2", models.QualityMedium, 1)
	neutralZones.AddPlan("N3", models.QualityMedium, 1)
	neutralZones.AddPlan("N4", models.QualityMedium, 1)
	tuning := models.NewGenerationTuning(configuration, 6)
	service := topology.NewSharedWebTopologyService()

	// Act
	variant := service.CreateTopologyVariant(*configuration, playerLabels, neutralZones, tuning, "")

	// Assert
	assert.NotZero(t, countPortalConnections(variant))
}
