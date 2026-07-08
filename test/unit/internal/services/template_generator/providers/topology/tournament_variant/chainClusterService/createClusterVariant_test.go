package chainClusterService_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/tournament_variant"
	"github.com/stretchr/testify/assert"
)

func newTwoNeutralPlans() models.NeutralZonePlans {
	neutralZones := models.NeutralZonePlans{}
	neutralZones.AddPlan("C", models.QualityLow, 0)
	neutralZones.AddPlan("D", models.QualityHigh, 1)
	return neutralZones
}

func TestWhenPlayerHasTwoNeutralPlans_CreatesSpawnPlusNeutralZones(t *testing.T) {
	// Arrange
	configuration := config.NewGeneratorConfig()
	neutralZones := newTwoNeutralPlans()
	tuning := models.NewGenerationTuning(configuration, 3)
	service := tournament_variant.NewChainClusterService()

	// Act
	zones, _ := service.CreateClusterVariant(*configuration, tuning, neutralZones, neutralZones, 0, "A")

	// Assert
	zoneNames := make([]string, 0, len(zones))
	for _, zone := range zones {
		zoneNames = append(zoneNames, zone.Name)
	}
	assert.Equal(t, []string{"Spawn-A", "Neutral-C", "Neutral-D"}, zoneNames)
}

func TestWhenChainIsBuilt_CreatesConnectionPerAdjacentPair(t *testing.T) {
	// Arrange
	configuration := config.NewGeneratorConfig()
	neutralZones := newTwoNeutralPlans()
	tuning := models.NewGenerationTuning(configuration, 3)
	service := tournament_variant.NewChainClusterService()

	// Act
	_, connections := service.CreateClusterVariant(*configuration, tuning, neutralZones, neutralZones, 0, "A")

	// Assert
	connectionNames := make([]string, 0, len(connections))
	for _, connection := range connections {
		connectionNames = append(connectionNames, connection.Name)
	}
	assert.Equal(t, []string{"Tourney-A-C", "Tourney-C-D"}, connectionNames)
}

func TestWhenFirstChainLinkIsBuilt_ConnectsSpawnToFirstNeutral(t *testing.T) {
	// Arrange
	configuration := config.NewGeneratorConfig()
	neutralZones := newTwoNeutralPlans()
	tuning := models.NewGenerationTuning(configuration, 3)
	service := tournament_variant.NewChainClusterService()

	// Act
	_, connections := service.CreateClusterVariant(*configuration, tuning, neutralZones, neutralZones, 0, "A")

	// Assert
	assert.Equal(t, [2]string{"Spawn-A", "Neutral-C"}, [2]string{connections[0].From, connections[0].To})
}

func TestWhenLaterChainLinkIsBuilt_ConnectsNeutralToNextNeutral(t *testing.T) {
	// Arrange
	configuration := config.NewGeneratorConfig()
	neutralZones := newTwoNeutralPlans()
	tuning := models.NewGenerationTuning(configuration, 3)
	service := tournament_variant.NewChainClusterService()

	// Act
	_, connections := service.CreateClusterVariant(*configuration, tuning, neutralZones, neutralZones, 0, "A")

	// Assert
	assert.Equal(t, [2]string{"Neutral-C", "Neutral-D"}, [2]string{connections[1].From, connections[1].To})
}

func TestWhenPlayerHasNoNeutralPlans_CreatesOnlySpawnZoneWithoutConnections(t *testing.T) {
	// Arrange
	configuration := config.NewGeneratorConfig()
	emptyPlans := models.NeutralZonePlans{}
	tuning := models.NewGenerationTuning(configuration, 1)
	service := tournament_variant.NewChainClusterService()

	// Act
	zones, connections := service.CreateClusterVariant(*configuration, tuning, emptyPlans, emptyPlans, 0, "A")

	// Assert
	assert.True(t, len(zones) == 1 && zones[0].Name == "Spawn-A" && len(connections) == 0,
		"expected a lone Spawn-A zone with no connections, got zones %v and %d connections", zones, len(connections))
}

func TestWhenSecondPlayerClusterIsBuilt_SpawnCastleBelongsToPlayerTwo(t *testing.T) {
	// Arrange
	configuration := config.NewGeneratorConfig()
	neutralZones := newTwoNeutralPlans()
	tuning := models.NewGenerationTuning(configuration, 3)
	service := tournament_variant.NewChainClusterService()

	// Act
	zones, _ := service.CreateClusterVariant(*configuration, tuning, neutralZones, neutralZones, 1, "B")

	// Assert
	assert.Equal(t, "Player2", zones[0].MainObjects[0].Spawn)
}
