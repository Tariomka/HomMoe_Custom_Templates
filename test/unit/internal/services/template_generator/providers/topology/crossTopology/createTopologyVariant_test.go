package crossTopology_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenTwoPlayersAndFiveNeutralPlansProvided_CreatesZonePerLabel(t *testing.T) {
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyCross
	playerLabels := []string{"A", "B"}
	neutralZones := models.NeutralZonePlans{}
	neutralZones.AddPlan("N1", models.QualityHigh, 1)
	neutralZones.AddPlan("N2", models.QualityLow, 0)
	neutralZones.AddPlan("N3", models.QualityLow, 0)
	neutralZones.AddPlan("N4", models.QualityMedium, 1)
	neutralZones.AddPlan("N5", models.QualityMedium, 1)
	tuning := models.NewGenerationTuning(configuration, 7)
	service := topology.NewCrossTopologyService()

	// Act
	variant := service.CreateTopologyVariant(*configuration, playerLabels, neutralZones, tuning, "")

	// Assert
	assert.Len(t, variant.Zones, 7)
}

func TestWhenNeutralZonesExist_FirstNeutralAnchorsTheCrossCentre(t *testing.T) {
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyCross
	playerLabels := []string{"A", "B"}
	neutralZones := models.NeutralZonePlans{}
	neutralZones.AddPlan("N1", models.QualityHigh, 1)
	neutralZones.AddPlan("N2", models.QualityLow, 0)
	neutralZones.AddPlan("N3", models.QualityLow, 0)
	tuning := models.NewGenerationTuning(configuration, 5)
	service := topology.NewCrossTopologyService()

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

func TestWhenCrossIsBuilt_EveryConnectionReferencesExistingZones(t *testing.T) {
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyCross
	playerLabels := []string{"A", "B"}
	neutralZones := models.NeutralZonePlans{}
	neutralZones.AddPlan("N1", models.QualityHigh, 1)
	neutralZones.AddPlan("N2", models.QualityLow, 0)
	neutralZones.AddPlan("N3", models.QualityLow, 0)
	neutralZones.AddPlan("N4", models.QualityMedium, 1)
	neutralZones.AddPlan("N5", models.QualityMedium, 1)
	tuning := models.NewGenerationTuning(configuration, 7)
	service := topology.NewCrossTopologyService()

	// Act
	variant := service.CreateTopologyVariant(*configuration, playerLabels, neutralZones, tuning, "")

	// Assert
	assert.Empty(t, danglingConnectionNames(variant))
}

func TestWhenNoNeutralZonesExist_JoinsPlayerTipsInARing(t *testing.T) {
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
			// Arrange
			configuration := config.NewGeneratorConfig()
			configuration.Topology = config.TopologyCross
			configuration.PlayerCount = len(testCase.playerLabels)
			neutralZones := models.NeutralZonePlans{}
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
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyCross
	configuration.NoDirectPlayerConnections = true
	playerLabels := []string{"A", "B"}
	neutralZones := models.NeutralZonePlans{}
	neutralZones.AddPlan("N1", models.QualityMedium, 1)
	neutralZones.AddPlan("N2", models.QualityMedium, 1)
	tuning := models.NewGenerationTuning(configuration, 4)
	service := topology.NewCrossTopologyService()

	// Act
	variant := service.CreateTopologyVariant(*configuration, playerLabels, neutralZones, tuning, "")

	// Assert
	assert.Empty(t, spawnToSpawnNamesWithPrefix(variant, "Rnd-"))
}

func TestWhenRandomPortalsEnabled_AddsPortalConnections(t *testing.T) {
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyCross
	configuration.RandomPortals = true
	playerLabels := []string{"A", "B"}
	neutralZones := models.NeutralZonePlans{}
	neutralZones.AddPlan("N1", models.QualityHigh, 1)
	neutralZones.AddPlan("N2", models.QualityLow, 0)
	neutralZones.AddPlan("N3", models.QualityMedium, 1)
	neutralZones.AddPlan("N4", models.QualityMedium, 1)
	tuning := models.NewGenerationTuning(configuration, 6)
	service := topology.NewCrossTopologyService()

	// Act
	variant := service.CreateTopologyVariant(*configuration, playerLabels, neutralZones, tuning, "")

	// Assert
	assert.NotZero(t, countPortalConnections(variant))
}
