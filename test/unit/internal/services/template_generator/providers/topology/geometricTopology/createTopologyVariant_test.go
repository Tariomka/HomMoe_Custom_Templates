package geometricTopology_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenThreePlayersAndSevenNeutralPlansProvided_CreatesZonePerLabel(t *testing.T) {
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyGeometric
	playerLabels := []string{"A", "B", "C"}
	neutralZones := models.NeutralZonePlans{}
	neutralZones.AddPlan("N1", models.QualityHigh, 1)
	neutralZones.AddPlan("N2", models.QualityLow, 0)
	neutralZones.AddPlan("N3", models.QualityLow, 0)
	neutralZones.AddPlan("N4", models.QualityLow, 0)
	neutralZones.AddPlan("N5", models.QualityMedium, 1)
	neutralZones.AddPlan("N6", models.QualityMedium, 1)
	neutralZones.AddPlan("N7", models.QualityMedium, 1)
	tuning := models.NewGenerationTuning(configuration, 10)
	service := topology.NewGeometricTopologyService()

	// Act
	variant := service.CreateTopologyVariant(*configuration, playerLabels, neutralZones, tuning, "")

	// Assert
	assert.Len(t, variant.Zones, 10)
}

func TestWhenNeutralZonesExist_FirstNeutralAnchorsTheFlowerCentre(t *testing.T) {
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyGeometric
	playerLabels := []string{"A", "B"}
	neutralZones := models.NeutralZonePlans{}
	neutralZones.AddPlan("N1", models.QualityHigh, 1)
	neutralZones.AddPlan("N2", models.QualityLow, 0)
	neutralZones.AddPlan("N3", models.QualityLow, 0)
	tuning := models.NewGenerationTuning(configuration, 5)
	service := topology.NewGeometricTopologyService()

	// Act
	variant := service.CreateTopologyVariant(*configuration, playerLabels, neutralZones, tuning, "")

	// Assert
	var centreZone entities.Zone
	for _, zone := range variant.Zones {
		if zone.Name == "Neutral-N1" {
			centreZone = zone
			break
		}
	}
	require.NotNil(t, centreZone.GeneratorPosition)
	assert.Equal(t, [2]float64{0.5, 0.5}, *centreZone.GeneratorPosition)
}

func TestWhenFlowerIsBuilt_EveryConnectionReferencesExistingZones(t *testing.T) {
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyGeometric
	playerLabels := []string{"A", "B", "C"}
	neutralZones := models.NeutralZonePlans{}
	neutralZones.AddPlan("N1", models.QualityHigh, 1)
	neutralZones.AddPlan("N2", models.QualityLow, 0)
	neutralZones.AddPlan("N3", models.QualityLow, 0)
	neutralZones.AddPlan("N4", models.QualityLow, 0)
	neutralZones.AddPlan("N5", models.QualityMedium, 1)
	neutralZones.AddPlan("N6", models.QualityMedium, 1)
	neutralZones.AddPlan("N7", models.QualityMedium, 1)
	tuning := models.NewGenerationTuning(configuration, 10)
	service := topology.NewGeometricTopologyService()

	// Act
	variant := service.CreateTopologyVariant(*configuration, playerLabels, neutralZones, tuning, "")

	// Assert
	assert.Empty(t, danglingConnectionNames(variant))
}

func TestWhenNoNeutralZonesExist_FallsBackToClosedPlayerPolygon(t *testing.T) {
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyGeometric
	playerLabels := []string{"A", "B", "C"}
	neutralZones := models.NeutralZonePlans{}
	tuning := models.NewGenerationTuning(configuration, 3)
	service := topology.NewGeometricTopologyService()

	// Act
	variant := service.CreateTopologyVariant(*configuration, playerLabels, neutralZones, tuning, "")

	// Assert
	assert.Len(t, variant.Zones, 3)
}

func TestWhenPlayerConnectionsAreForbidden_NoRandomConnectionJoinsTwoSpawnZones(t *testing.T) {
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyGeometric
	configuration.NoDirectPlayerConnections = true
	playerLabels := []string{"A", "B"}
	neutralZones := models.NeutralZonePlans{}
	neutralZones.AddPlan("N1", models.QualityMedium, 1)
	neutralZones.AddPlan("N2", models.QualityMedium, 1)
	tuning := models.NewGenerationTuning(configuration, 4)
	service := topology.NewGeometricTopologyService()

	// Act
	variant := service.CreateTopologyVariant(*configuration, playerLabels, neutralZones, tuning, "")

	// Assert
	assert.Empty(t, directSpawnToSpawnNames(variant, "Rnd-"))
}

func TestWhenRandomPortalsEnabled_AddsPortalConnections(t *testing.T) {
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyGeometric
	configuration.RandomPortals = true
	playerLabels := []string{"A", "B"}
	neutralZones := models.NeutralZonePlans{}
	neutralZones.AddPlan("N1", models.QualityHigh, 1)
	neutralZones.AddPlan("N2", models.QualityLow, 0)
	neutralZones.AddPlan("N3", models.QualityLow, 0)
	neutralZones.AddPlan("N4", models.QualityMedium, 1)
	tuning := models.NewGenerationTuning(configuration, 6)
	service := topology.NewGeometricTopologyService()

	// Act
	variant := service.CreateTopologyVariant(*configuration, playerLabels, neutralZones, tuning, "")

	// Assert
	assert.NotZero(t, countPortalConnections(variant))
}
