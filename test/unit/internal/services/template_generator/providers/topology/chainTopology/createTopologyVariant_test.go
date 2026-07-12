package chainTopology_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutralZone"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology"
	"github.com/stretchr/testify/assert"
)

func TestWhenTwoPlayersAndTwoNeutralPlansProvided_CreatesZonePerLabel(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyChain
	playerLabels := []string{"A", "B"}
	neutralZones := neutralZone.Plans{}
	neutralZones.AddPlan("N1", neutralZone.QualityMedium, 1)
	neutralZones.AddPlan("N2", neutralZone.QualityMedium, 1)
	tuning := models.NewGenerationTuning(configuration, 4)
	service := topology.NewChainTopologyService()

	// Act
	variant := service.CreateTopologyVariant(*configuration, playerLabels, neutralZones, tuning, "")

	// Assert
	assert.Len(t, variant.Zones, 4)
}

func TestWhenFourLabelsFormTheChain_CreatesConnectionPerAdjacentPair(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyChain
	playerLabels := []string{"A", "B"}
	neutralZones := neutralZone.Plans{}
	neutralZones.AddPlan("N1", neutralZone.QualityMedium, 1)
	neutralZones.AddPlan("N2", neutralZone.QualityMedium, 1)
	tuning := models.NewGenerationTuning(configuration, 4)
	service := topology.NewChainTopologyService()

	// Act
	variant := service.CreateTopologyVariant(*configuration, playerLabels, neutralZones, tuning, "")

	// Assert
	assert.Equal(t, 3, countConnectionsWithPrefix(variant, "Chain-"))
}

func TestWhenChainIsBuilt_EveryConnectionReferencesExistingZones(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyChain
	playerLabels := []string{"A", "B"}
	neutralZones := neutralZone.Plans{}
	neutralZones.AddPlan("N1", neutralZone.QualityMedium, 1)
	neutralZones.AddPlan("N2", neutralZone.QualityMedium, 1)
	tuning := models.NewGenerationTuning(configuration, 4)
	service := topology.NewChainTopologyService()

	// Act
	variant := service.CreateTopologyVariant(*configuration, playerLabels, neutralZones, tuning, "")

	// Assert
	assert.Empty(t, danglingConnectionNames(variant))
}

func TestWhenIsolatedPlayersAreAdjacentInTheChain_SkipsTheirChainConnection(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyChain
	configuration.NoDirectPlayerConnections = true
	playerLabels := []string{"A", "B"}
	neutralZones := neutralZone.Plans{}
	tuning := models.NewGenerationTuning(configuration, 2)
	service := topology.NewChainTopologyService()

	// Act
	variant := service.CreateTopologyVariant(*configuration, playerLabels, neutralZones, tuning, "")

	// Assert
	assert.Zero(t, countConnectionsWithPrefix(variant, "Chain-"))
}

func TestWhenRandomPortalsEnabled_AddsPortalConnections(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyChain
	configuration.RandomPortals = true
	playerLabels := []string{"A", "B"}
	neutralZones := neutralZone.Plans{}
	neutralZones.AddPlan("N1", neutralZone.QualityMedium, 1)
	neutralZones.AddPlan("N2", neutralZone.QualityMedium, 1)
	neutralZones.AddPlan("N3", neutralZone.QualityMedium, 1)
	neutralZones.AddPlan("N4", neutralZone.QualityMedium, 1)
	tuning := models.NewGenerationTuning(configuration, 6)
	service := topology.NewChainTopologyService()

	// Act
	variant := service.CreateTopologyVariant(*configuration, playerLabels, neutralZones, tuning, "")

	// Assert
	assert.NotZero(t, countPortalConnections(variant))
}
