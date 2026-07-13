package crossTopology_test

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

func TestWhenTwoPlayersAndFiveNeutralPlansProvided_CreatesZonePerLabel(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyCross
	playerLabels := []string{"A", "B"}
	neutralZones := neutralZone.Plans{}
	neutralZones.AddPlan("N1", neutralZone.QualityHigh, 1)
	neutralZones.AddPlan("N2", neutralZone.QualityLow, 0)
	neutralZones.AddPlan("N3", neutralZone.QualityLow, 0)
	neutralZones.AddPlan("N4", neutralZone.QualityMedium, 1)
	neutralZones.AddPlan("N5", neutralZone.QualityMedium, 1)
	tuning := models.NewGenerationTuning(configuration, 7)
	service := topology.NewCrossTopologyService()

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
	neutralZones := neutralZone.Plans{}
	neutralZones.AddPlan("N1", neutralZone.QualityHigh, 1)
	neutralZones.AddPlan("N2", neutralZone.QualityLow, 0)
	neutralZones.AddPlan("N3", neutralZone.QualityLow, 0)
	tuning := models.NewGenerationTuning(configuration, 5)
	service := topology.NewCrossTopologyService()

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
	neutralZones := neutralZone.Plans{}
	neutralZones.AddPlan("N1", neutralZone.QualityHigh, 1)
	neutralZones.AddPlan("N2", neutralZone.QualityLow, 0)
	neutralZones.AddPlan("N3", neutralZone.QualityLow, 0)
	neutralZones.AddPlan("N4", neutralZone.QualityMedium, 1)
	neutralZones.AddPlan("N5", neutralZone.QualityMedium, 1)
	tuning := models.NewGenerationTuning(configuration, 7)
	service := topology.NewCrossTopologyService()

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
			neutralZones := neutralZone.Plans{}
			tuning := models.NewGenerationTuning(configuration, len(testCase.playerLabels))
			service := topology.NewCrossTopologyService()

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
	neutralZones := neutralZone.Plans{}
	neutralZones.AddPlan("N1", neutralZone.QualityMedium, 1)
	neutralZones.AddPlan("N2", neutralZone.QualityMedium, 1)
	tuning := models.NewGenerationTuning(configuration, 4)
	service := topology.NewCrossTopologyService()

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
	neutralZones := neutralZone.Plans{}
	neutralZones.AddPlan("N1", neutralZone.QualityHigh, 1)
	neutralZones.AddPlan("N2", neutralZone.QualityLow, 0)
	neutralZones.AddPlan("N3", neutralZone.QualityMedium, 1)
	neutralZones.AddPlan("N4", neutralZone.QualityMedium, 1)
	tuning := models.NewGenerationTuning(configuration, 6)
	service := topology.NewCrossTopologyService()

	// Act
	variant := service.CreateTopologyVariant(*configuration, playerLabels, neutralZones, tuning, "")

	// Assert
	assert.NotZero(t, countPortalConnections(variant))
}
