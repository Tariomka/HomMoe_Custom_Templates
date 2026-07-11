package hubTopology_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenOuterLabelsSurroundTheHub_CreatesSingleHubZone(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyHubAndSpoke
	playerLabels := []string{"A", "B"}
	neutralZones := models.NeutralZonePlans{}
	neutralZones.AddPlan("N1", models.QualityMedium, 1)
	neutralZones.AddPlan("N2", models.QualityMedium, 1)
	tuning := models.NewGenerationTuning(configuration, 5)
	service := topology.NewHubTopologyService()

	// Act
	variant := service.CreateTopologyVariant(*configuration, playerLabels, neutralZones, tuning, false)

	// Assert
	hubZoneCount := 0
	for _, zone := range variant.Zones {
		if zone.Name == "Hub" {
			hubZoneCount++
		}
	}
	assert.Equal(t, 1, hubZoneCount)
}

func TestWhenTwoPlayersAndTwoNeutralPlansProvided_CreatesHubPlusOuterZones(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyHubAndSpoke
	playerLabels := []string{"A", "B"}
	neutralZones := models.NeutralZonePlans{}
	neutralZones.AddPlan("N1", models.QualityMedium, 1)
	neutralZones.AddPlan("N2", models.QualityMedium, 1)
	tuning := models.NewGenerationTuning(configuration, 5)
	service := topology.NewHubTopologyService()

	// Act
	variant := service.CreateTopologyVariant(*configuration, playerLabels, neutralZones, tuning, false)

	// Assert
	assert.Len(t, variant.Zones, 5)
}

func TestWhenHubAndSpokesAreBuilt_EveryConnectionReferencesExistingZones(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyHubAndSpoke
	playerLabels := []string{"A", "B"}
	neutralZones := models.NeutralZonePlans{}
	neutralZones.AddPlan("N1", models.QualityMedium, 1)
	neutralZones.AddPlan("N2", models.QualityMedium, 1)
	tuning := models.NewGenerationTuning(configuration, 5)
	service := topology.NewHubTopologyService()

	// Act
	variant := service.CreateTopologyVariant(*configuration, playerLabels, neutralZones, tuning, false)

	// Assert
	assert.Empty(t, danglingConnectionNames(variant))
}

func TestWhenHubMandatoryContentConfigured_HubZoneReferencesHubContentName(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyHubAndSpoke
	configuration.HubZoneMandatoryContent = []entities.MandatoryContentItem{{}}
	playerLabels := []string{"A", "B"}
	neutralZones := models.NeutralZonePlans{}
	neutralZones.AddPlan("N1", models.QualityMedium, 1)
	tuning := models.NewGenerationTuning(configuration, 4)
	service := topology.NewHubTopologyService()

	// Act
	variant := service.CreateTopologyVariant(*configuration, playerLabels, neutralZones, tuning, false)

	// Assert
	var hubZone entities.Zone
	for _, zone := range variant.Zones {
		if zone.Name == "Hub" {
			hubZone = zone
			break
		}
	}
	require.Equal(t, "Hub", hubZone.Name)
	assert.Contains(t, hubZone.MandatoryContent, "mandatory_content_hub")
}

func TestWhenIsolatedPlayersAreAdjacentOuterLabels_SkipsTheirPseudoConnection(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyHubAndSpoke
	configuration.NoDirectPlayerConnections = true
	playerLabels := []string{"A", "B"}
	neutralZones := models.NeutralZonePlans{}
	neutralZones.AddPlan("N1", models.QualityMedium, 1)
	tuning := models.NewGenerationTuning(configuration, 4)
	service := topology.NewHubTopologyService()

	// Act
	variant := service.CreateTopologyVariant(*configuration, playerLabels, neutralZones, tuning, false)

	// Assert
	assert.NotContains(t, connectionNames(variant), "Pseudo-A-B")
}

func TestWhenCirclesTopologyCannotHonorSeparation_BalancesOuterLabelsWithoutSeparation(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyCircles
	configuration.MinNeutralZonesBetweenPlayers = 0
	playerLabels := []string{"A", "B"}
	neutralZones := models.NeutralZonePlans{}
	neutralZones.AddPlan("N1", models.QualityMedium, 1)
	neutralZones.AddPlan("N2", models.QualityMedium, 1)
	tuning := models.NewGenerationTuning(configuration, 5)
	service := topology.NewHubTopologyService()

	// Act
	variant := service.CreateTopologyVariant(*configuration, playerLabels, neutralZones, tuning, false)

	// Assert
	assert.Len(t, variant.Zones, 5)
}

func TestWhenCirclesTopologyHonorsNeutralSeparation_CreatesHubPlusOuterZones(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyCircles
	configuration.MinNeutralZonesBetweenPlayers = 1
	configuration.PlayerCount = 2
	configuration.ZoneConfiguration.NeutralZoneCount = 2
	playerLabels := []string{"A", "B"}
	neutralZones := models.NeutralZonePlans{}
	neutralZones.AddPlan("N1", models.QualityMedium, 1)
	neutralZones.AddPlan("N2", models.QualityMedium, 1)
	tuning := models.NewGenerationTuning(configuration, 5)
	service := topology.NewHubTopologyService()

	// Act
	variant := service.CreateTopologyVariant(*configuration, playerLabels, neutralZones, tuning, false)

	// Assert
	assert.Len(t, variant.Zones, 5)
}

func TestWhenRandomPortalsEnabled_AddsPortalConnections(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyHubAndSpoke
	configuration.RandomPortals = true
	playerLabels := []string{"A", "B"}
	neutralZones := models.NeutralZonePlans{}
	neutralZones.AddPlan("N1", models.QualityMedium, 1)
	neutralZones.AddPlan("N2", models.QualityMedium, 1)
	neutralZones.AddPlan("N3", models.QualityMedium, 1)
	neutralZones.AddPlan("N4", models.QualityMedium, 1)
	tuning := models.NewGenerationTuning(configuration, 7)
	service := topology.NewHubTopologyService()

	// Act
	variant := service.CreateTopologyVariant(*configuration, playerLabels, neutralZones, tuning, false)

	// Assert
	assert.NotZero(t, countPortalConnections(variant))
}
