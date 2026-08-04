package crossTopology_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenTwoPlayersAndFiveNeutralPlansProvided_CreatesZonePerLabel(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyCross
	playerLabels := []string{"A", "B"}
	neutralZones := neutral_zone.Plans{}
	neutralZones.AddPlan("N1", neutral_zone.QualityHigh, 1)
	neutralZones.AddPlan("N2", neutral_zone.QualityLow, 0)
	neutralZones.AddPlan("N3", neutral_zone.QualityLow, 0)
	neutralZones.AddPlan("N4", neutral_zone.QualityMedium, 1)
	neutralZones.AddPlan("N5", neutral_zone.QualityMedium, 1)
	tuning := test_helpers.NewGenerationTuning(configuration, 7)
	service := topology.NewCrossTopologyService(test_helpers.NewZoneFactories())

	// Act
	variant := service.CreateTopologyVariant(*configuration, playerLabels, neutralZones, tuning, "")

	// Assert
	assert.Len(t, variant.Zones, 7)
}

func TestWhenNeutralZonesExist_FirstNeutralAnchorsTheCrossCenter(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyCross
	playerLabels := []string{"A", "B"}
	neutralZones := neutral_zone.Plans{}
	neutralZones.AddPlan("N1", neutral_zone.QualityHigh, 1)
	neutralZones.AddPlan("N2", neutral_zone.QualityLow, 0)
	neutralZones.AddPlan("N3", neutral_zone.QualityLow, 0)
	tuning := test_helpers.NewGenerationTuning(configuration, 5)
	service := topology.NewCrossTopologyService(test_helpers.NewZoneFactories())

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

func TestWhenCrossIsBuilt_EveryConnectionReferencesExistingZones(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyCross
	playerLabels := []string{"A", "B"}
	neutralZones := neutral_zone.Plans{}
	neutralZones.AddPlan("N1", neutral_zone.QualityHigh, 1)
	neutralZones.AddPlan("N2", neutral_zone.QualityLow, 0)
	neutralZones.AddPlan("N3", neutral_zone.QualityLow, 0)
	neutralZones.AddPlan("N4", neutral_zone.QualityMedium, 1)
	neutralZones.AddPlan("N5", neutral_zone.QualityMedium, 1)
	tuning := test_helpers.NewGenerationTuning(configuration, 7)
	service := topology.NewCrossTopologyService(test_helpers.NewZoneFactories())

	// Act
	variant := service.CreateTopologyVariant(*configuration, playerLabels, neutralZones, tuning, "")

	// Assert
	assert.Empty(t, danglingConnectionNames(variant))
}

func TestWhenNoNeutralZonesExist_JoinsPlayerTipsInARing(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name              string
		playerLabels      []string
		expectedRingEdges int
	}{
		{
			name:              "WhenTwoPlayersHaveNoNeutrals_CreatesSingleTipEdge",
			playerLabels:      []string{"A", "B"},
			expectedRingEdges: 1,
		},
		{
			name:              "WhenThreePlayersHaveNoNeutrals_CreatesTipEdgePerPlayer",
			playerLabels:      []string{"A", "B", "C"},
			expectedRingEdges: 3,
		},
		{
			name:              "WhenFourPlayersHaveNoNeutrals_CreatesTipEdgePerPlayer",
			playerLabels:      []string{"A", "B", "C", "D"},
			expectedRingEdges: 4,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			// Arrange
			configuration := config.NewGeneratorConfig()
			configuration.Topology = config.TopologyCross
			configuration.PlayerCount = len(testCase.playerLabels)
			neutralZones := neutral_zone.Plans{}
			tuning := test_helpers.NewGenerationTuning(configuration, len(testCase.playerLabels))
			service := topology.NewCrossTopologyService(test_helpers.NewZoneFactories())

			// Act
			variant := service.CreateTopologyVariant(*configuration, testCase.playerLabels, neutralZones, tuning, "")

			// Assert
			assert.Len(t, spawnToSpawnNamesWithPrefix(variant, "Rnd-"), testCase.expectedRingEdges,
				"player tips of a %d-player cross must close into a ring", len(testCase.playerLabels))
		})
	}
}

func TestWhenPlayerConnectionsAreForbidden_NoRandomConnectionJoinsTwoSpawnZones(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyCross
	configuration.NoDirectPlayerConnections = true
	playerLabels := []string{"A", "B"}
	neutralZones := neutral_zone.Plans{}
	neutralZones.AddPlan("N1", neutral_zone.QualityMedium, 1)
	neutralZones.AddPlan("N2", neutral_zone.QualityMedium, 1)
	tuning := test_helpers.NewGenerationTuning(configuration, 4)
	service := topology.NewCrossTopologyService(test_helpers.NewZoneFactories())

	// Act
	variant := service.CreateTopologyVariant(*configuration, playerLabels, neutralZones, tuning, "")

	// Assert
	assert.Empty(t, spawnToSpawnNamesWithPrefix(variant, "Rnd-"))
}

func TestWhenRandomPortalsEnabled_AddsPortalConnections(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyCross
	configuration.RandomPortals = true
	playerLabels := []string{"A", "B"}
	neutralZones := neutral_zone.Plans{}
	neutralZones.AddPlan("N1", neutral_zone.QualityHigh, 1)
	neutralZones.AddPlan("N2", neutral_zone.QualityLow, 0)
	neutralZones.AddPlan("N3", neutral_zone.QualityMedium, 1)
	neutralZones.AddPlan("N4", neutral_zone.QualityMedium, 1)
	tuning := test_helpers.NewGenerationTuning(configuration, 6)
	service := topology.NewCrossTopologyService(test_helpers.NewZoneFactories())

	// Act
	variant := service.CreateTopologyVariant(*configuration, playerLabels, neutralZones, tuning, "")

	// Assert
	assert.NotZero(t, countPortalConnections(variant))
}
