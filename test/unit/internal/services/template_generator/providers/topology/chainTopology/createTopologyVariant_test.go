package chainTopology_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
)

func TestWhenTwoPlayersAndTwoNeutralPlansProvided_CreatesZonePerLabel(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyChain
	playerLabels := []string{"A", "B"}
	neutralZones := neutral_zone.Plans{}
	neutralZones.AddPlan("N1", neutral_zone.QualityMedium, 1)
	neutralZones.AddPlan("N2", neutral_zone.QualityMedium, 1)
	tuning := test_helpers.NewGenerationTuning(configuration, 4)
	service := newChainTopologyService()

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
	neutralZones := neutral_zone.Plans{}
	neutralZones.AddPlan("N1", neutral_zone.QualityMedium, 1)
	neutralZones.AddPlan("N2", neutral_zone.QualityMedium, 1)
	tuning := test_helpers.NewGenerationTuning(configuration, 4)
	service := newChainTopologyService()

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
	neutralZones := neutral_zone.Plans{}
	neutralZones.AddPlan("N1", neutral_zone.QualityMedium, 1)
	neutralZones.AddPlan("N2", neutral_zone.QualityMedium, 1)
	tuning := test_helpers.NewGenerationTuning(configuration, 4)
	service := newChainTopologyService()

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
	neutralZones := neutral_zone.Plans{}
	tuning := test_helpers.NewGenerationTuning(configuration, 2)
	service := newChainTopologyService()

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
	neutralZones := neutral_zone.Plans{}
	neutralZones.AddPlan("N1", neutral_zone.QualityMedium, 1)
	neutralZones.AddPlan("N2", neutral_zone.QualityMedium, 1)
	neutralZones.AddPlan("N3", neutral_zone.QualityMedium, 1)
	neutralZones.AddPlan("N4", neutral_zone.QualityMedium, 1)
	tuning := test_helpers.NewGenerationTuning(configuration, 6)
	service := newChainTopologyService()

	// Act
	variant := service.CreateTopologyVariant(*configuration, playerLabels, neutralZones, tuning, "")

	// Assert
	assert.NotZero(t, countPortalConnections(variant))
}
