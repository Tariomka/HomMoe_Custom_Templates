package fractalTopology_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology"
	"github.com/stretchr/testify/assert"
)

// addTieredNeutralPlans adds two low, two medium and two high quality plans so
// every fractal tier band receives zones.
func addTieredNeutralPlans(neutralZones *neutral_zone.Plans) {
	neutralZones.AddPlan("N1", neutral_zone.QualityLow, 0)
	neutralZones.AddPlan("N2", neutral_zone.QualityLow, 0)
	neutralZones.AddPlan("N3", neutral_zone.QualityMedium, 1)
	neutralZones.AddPlan("N4", neutral_zone.QualityMedium, 1)
	neutralZones.AddPlan("N5", neutral_zone.QualityHigh, 1)
	neutralZones.AddPlan("N6", neutral_zone.QualityHigh, 1)
}

func TestWhenTwoPlayersAndSixTieredNeutralPlansProvided_CreatesZonePerLabel(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyFractal
	playerLabels := []string{"A", "B"}
	neutralZones := neutral_zone.Plans{}
	addTieredNeutralPlans(&neutralZones)
	tuning := models.NewGenerationTuning(configuration, 8)
	service := topology.NewFractalTopologyService()

	// Act
	variant := service.CreateTopologyVariant(*configuration, playerLabels, neutralZones, tuning, "")

	// Assert
	assert.Len(t, variant.Zones, 8)
}

func TestWhenEveryTierIsPopulated_NoRandomConnectionJoinsTwoSpawnZones(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyFractal
	playerLabels := []string{"A", "B"}
	neutralZones := neutral_zone.Plans{}
	addTieredNeutralPlans(&neutralZones)
	tuning := models.NewGenerationTuning(configuration, 8)
	service := topology.NewFractalTopologyService()

	// Act
	variant := service.CreateTopologyVariant(*configuration, playerLabels, neutralZones, tuning, "")

	// Assert
	assert.Empty(t, randomSpawnToSpawnNames(variant))
}

func TestWhenFractalsAreLaidOut_EveryZoneGetsPositionInsideUnitSquare(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyFractal
	playerLabels := []string{"A", "B"}
	neutralZones := neutral_zone.Plans{}
	addTieredNeutralPlans(&neutralZones)
	tuning := models.NewGenerationTuning(configuration, 8)
	service := topology.NewFractalTopologyService()

	// Act
	variant := service.CreateTopologyVariant(*configuration, playerLabels, neutralZones, tuning, "")

	// Assert
	assert.Empty(t, zonesWithoutValidPosition(variant))
}

func TestWhenFractalsAreBuilt_EveryConnectionReferencesExistingZones(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyFractal
	playerLabels := []string{"A", "B"}
	neutralZones := neutral_zone.Plans{}
	addTieredNeutralPlans(&neutralZones)
	tuning := models.NewGenerationTuning(configuration, 8)
	service := topology.NewFractalTopologyService()

	// Act
	variant := service.CreateTopologyVariant(*configuration, playerLabels, neutralZones, tuning, "")

	// Assert
	assert.Empty(t, danglingConnectionNames(variant))
}

func TestWhenFewerNeutralZonesThanPlayers_CreatesZonePerLabel(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyFractal
	configuration.PlayerCount = 3
	playerLabels := []string{"A", "B", "C"}
	neutralZones := neutral_zone.Plans{}
	neutralZones.AddPlan("N1", neutral_zone.QualityLow, 0)
	tuning := models.NewGenerationTuning(configuration, 4)
	service := topology.NewFractalTopologyService()

	// Act
	variant := service.CreateTopologyVariant(*configuration, playerLabels, neutralZones, tuning, "")

	// Assert
	assert.Len(t, variant.Zones, 4)
}

func TestWhenFewerNeutralZonesThanPlayers_EveryConnectionReferencesExistingZones(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyFractal
	configuration.PlayerCount = 3
	playerLabels := []string{"A", "B", "C"}
	neutralZones := neutral_zone.Plans{}
	neutralZones.AddPlan("N1", neutral_zone.QualityLow, 0)
	tuning := models.NewGenerationTuning(configuration, 4)
	service := topology.NewFractalTopologyService()

	// Act
	variant := service.CreateTopologyVariant(*configuration, playerLabels, neutralZones, tuning, "")

	// Assert
	assert.Empty(t, danglingConnectionNames(variant))
}

func TestWhenNoNeutralZonesExist_CreatesOnlyPlayerZones(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyFractal
	playerLabels := []string{"A", "B"}
	neutralZones := neutral_zone.Plans{}
	tuning := models.NewGenerationTuning(configuration, 2)
	service := topology.NewFractalTopologyService()

	// Act
	variant := service.CreateTopologyVariant(*configuration, playerLabels, neutralZones, tuning, "")

	// Assert
	assert.Len(t, variant.Zones, 2)
}

func TestWhenPlayerConnectionsAreForbidden_NoRandomConnectionJoinsTwoSpawnZones(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyFractal
	configuration.NoDirectPlayerConnections = true
	playerLabels := []string{"A", "B"}
	neutralZones := neutral_zone.Plans{}
	addTieredNeutralPlans(&neutralZones)
	tuning := models.NewGenerationTuning(configuration, 8)
	service := topology.NewFractalTopologyService()

	// Act
	variant := service.CreateTopologyVariant(*configuration, playerLabels, neutralZones, tuning, "")

	// Assert
	assert.Empty(t, randomSpawnToSpawnNames(variant))
}

func TestWhenRandomPortalsEnabled_AddsPortalConnections(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyFractal
	configuration.RandomPortals = true
	playerLabels := []string{"A", "B"}
	neutralZones := neutral_zone.Plans{}
	addTieredNeutralPlans(&neutralZones)
	tuning := models.NewGenerationTuning(configuration, 8)
	service := topology.NewFractalTopologyService()

	// Act
	variant := service.CreateTopologyVariant(*configuration, playerLabels, neutralZones, tuning, "")

	// Assert
	assert.NotZero(t, countPortalConnections(variant))
}
