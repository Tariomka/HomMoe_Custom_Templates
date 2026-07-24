package geometricTopology_test

import (
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
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
	neutralZones := neutral_zone.Plans{}
	neutralZones.AddPlan("N1", neutral_zone.QualityHigh, 1)
	neutralZones.AddPlan("N2", neutral_zone.QualityLow, 0)
	neutralZones.AddPlan("N3", neutral_zone.QualityLow, 0)
	neutralZones.AddPlan("N4", neutral_zone.QualityLow, 0)
	neutralZones.AddPlan("N5", neutral_zone.QualityMedium, 1)
	neutralZones.AddPlan("N6", neutral_zone.QualityMedium, 1)
	neutralZones.AddPlan("N7", neutral_zone.QualityMedium, 1)
	tuning := test_helpers.NewGenerationTuning(configuration, 10)
	service := topology.NewGeometricTopologyService()

	// Act
	variant := service.CreateTopologyVariant(*configuration, playerLabels, neutralZones, tuning, "")

	// Assert
	assert.Len(t, variant.Zones, 10)
}

func TestWhenNeutralZonesExist_FirstNeutralAnchorsTheFlowerCenter(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyGeometric
	playerLabels := []string{"A", "B"}
	neutralZones := neutral_zone.Plans{}
	neutralZones.AddPlan("N1", neutral_zone.QualityHigh, 1)
	neutralZones.AddPlan("N2", neutral_zone.QualityLow, 0)
	neutralZones.AddPlan("N3", neutral_zone.QualityLow, 0)
	tuning := test_helpers.NewGenerationTuning(configuration, 5)
	service := topology.NewGeometricTopologyService()

	// Act
	variant := service.CreateTopologyVariant(*configuration, playerLabels, neutralZones, tuning, "")

	// Assert
	var centerZone entities.Zone
	for _, zone := range variant.Zones {
		if zone.Name == "Neutral-N1" {
			centerZone = zone
			break
		}
	}
	require.NotNil(t, centerZone.GeneratorPosition)
	assert.Equal(t, [2]float64{0.5, 0.5}, *centerZone.GeneratorPosition)
}

func TestWhenFlowerIsBuilt_EveryConnectionReferencesExistingZones(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyGeometric
	playerLabels := []string{"A", "B", "C"}
	neutralZones := neutral_zone.Plans{}
	neutralZones.AddPlan("N1", neutral_zone.QualityHigh, 1)
	neutralZones.AddPlan("N2", neutral_zone.QualityLow, 0)
	neutralZones.AddPlan("N3", neutral_zone.QualityLow, 0)
	neutralZones.AddPlan("N4", neutral_zone.QualityLow, 0)
	neutralZones.AddPlan("N5", neutral_zone.QualityMedium, 1)
	neutralZones.AddPlan("N6", neutral_zone.QualityMedium, 1)
	neutralZones.AddPlan("N7", neutral_zone.QualityMedium, 1)
	tuning := test_helpers.NewGenerationTuning(configuration, 10)
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
	neutralZones := neutral_zone.Plans{}
	tuning := test_helpers.NewGenerationTuning(configuration, 3)
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
	neutralZones := neutral_zone.Plans{}
	neutralZones.AddPlan("N1", neutral_zone.QualityMedium, 1)
	neutralZones.AddPlan("N2", neutral_zone.QualityMedium, 1)
	tuning := test_helpers.NewGenerationTuning(configuration, 4)
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
	neutralZones := neutral_zone.Plans{}
	neutralZones.AddPlan("N1", neutral_zone.QualityHigh, 1)
	neutralZones.AddPlan("N2", neutral_zone.QualityLow, 0)
	neutralZones.AddPlan("N3", neutral_zone.QualityLow, 0)
	neutralZones.AddPlan("N4", neutral_zone.QualityMedium, 1)
	tuning := test_helpers.NewGenerationTuning(configuration, 6)
	service := topology.NewGeometricTopologyService()

	// Act
	variant := service.CreateTopologyVariant(*configuration, playerLabels, neutralZones, tuning, "")

	// Assert
	assert.NotZero(t, countPortalConnections(variant))
}
