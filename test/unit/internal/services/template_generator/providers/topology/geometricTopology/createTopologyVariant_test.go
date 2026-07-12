package geometricTopology_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutralZone"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenThreePlayersAndSevenNeutralPlansProvided_CreatesZonePerLabel(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyGeometric
	playerLabels := []string{"A", "B", "C"}
	neutralZones := neutralZone.Plans{}
	neutralZones.AddPlan("N1", neutralZone.QualityHigh, 1)
	neutralZones.AddPlan("N2", neutralZone.QualityLow, 0)
	neutralZones.AddPlan("N3", neutralZone.QualityLow, 0)
	neutralZones.AddPlan("N4", neutralZone.QualityLow, 0)
	neutralZones.AddPlan("N5", neutralZone.QualityMedium, 1)
	neutralZones.AddPlan("N6", neutralZone.QualityMedium, 1)
	neutralZones.AddPlan("N7", neutralZone.QualityMedium, 1)
	tuning := models.NewGenerationTuning(configuration, 10)
	service := topology.NewGeometricTopologyService()

	// Act
	variant := service.CreateTopologyVariant(*configuration, playerLabels, neutralZones, tuning, "")

	// Assert
	assert.Len(t, variant.Zones, 10)
}

func TestWhenNeutralZonesExist_FirstNeutralAnchorsTheFlowerCentre(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyGeometric
	playerLabels := []string{"A", "B"}
	neutralZones := neutralZone.Plans{}
	neutralZones.AddPlan("N1", neutralZone.QualityHigh, 1)
	neutralZones.AddPlan("N2", neutralZone.QualityLow, 0)
	neutralZones.AddPlan("N3", neutralZone.QualityLow, 0)
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
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyGeometric
	playerLabels := []string{"A", "B", "C"}
	neutralZones := neutralZone.Plans{}
	neutralZones.AddPlan("N1", neutralZone.QualityHigh, 1)
	neutralZones.AddPlan("N2", neutralZone.QualityLow, 0)
	neutralZones.AddPlan("N3", neutralZone.QualityLow, 0)
	neutralZones.AddPlan("N4", neutralZone.QualityLow, 0)
	neutralZones.AddPlan("N5", neutralZone.QualityMedium, 1)
	neutralZones.AddPlan("N6", neutralZone.QualityMedium, 1)
	neutralZones.AddPlan("N7", neutralZone.QualityMedium, 1)
	tuning := models.NewGenerationTuning(configuration, 10)
	service := topology.NewGeometricTopologyService()

	// Act
	variant := service.CreateTopologyVariant(*configuration, playerLabels, neutralZones, tuning, "")

	// Assert
	assert.Empty(t, danglingConnectionNames(variant))
}

func TestWhenNoNeutralZonesExist_FallsBackToClosedPlayerPolygon(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyGeometric
	playerLabels := []string{"A", "B", "C"}
	neutralZones := neutralZone.Plans{}
	tuning := models.NewGenerationTuning(configuration, 3)
	service := topology.NewGeometricTopologyService()

	// Act
	variant := service.CreateTopologyVariant(*configuration, playerLabels, neutralZones, tuning, "")

	// Assert
	assert.Len(t, variant.Zones, 3)
}

func TestWhenPlayerConnectionsAreForbidden_NoRandomConnectionJoinsTwoSpawnZones(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyGeometric
	configuration.NoDirectPlayerConnections = true
	playerLabels := []string{"A", "B"}
	neutralZones := neutralZone.Plans{}
	neutralZones.AddPlan("N1", neutralZone.QualityMedium, 1)
	neutralZones.AddPlan("N2", neutralZone.QualityMedium, 1)
	tuning := models.NewGenerationTuning(configuration, 4)
	service := topology.NewGeometricTopologyService()

	// Act
	variant := service.CreateTopologyVariant(*configuration, playerLabels, neutralZones, tuning, "")

	// Assert
	assert.Empty(t, directSpawnToSpawnNames(variant, "Rnd-"))
}

func TestWhenRandomPortalsEnabled_AddsPortalConnections(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyGeometric
	configuration.RandomPortals = true
	playerLabels := []string{"A", "B"}
	neutralZones := neutralZone.Plans{}
	neutralZones.AddPlan("N1", neutralZone.QualityHigh, 1)
	neutralZones.AddPlan("N2", neutralZone.QualityLow, 0)
	neutralZones.AddPlan("N3", neutralZone.QualityLow, 0)
	neutralZones.AddPlan("N4", neutralZone.QualityMedium, 1)
	tuning := models.NewGenerationTuning(configuration, 6)
	service := topology.NewGeometricTopologyService()

	// Act
	variant := service.CreateTopologyVariant(*configuration, playerLabels, neutralZones, tuning, "")

	// Assert
	assert.NotZero(t, countPortalConnections(variant))
}
