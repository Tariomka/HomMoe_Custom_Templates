package ringClusterService_test

import (
	"strings"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutralZone"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/tournament_variant"
	"github.com/stretchr/testify/assert"
)

func newThreeNeutralPlans() neutralZone.Plans {
	neutralZones := neutralZone.Plans{}
	neutralZones.AddPlan("C", neutralZone.QualityLow, 0)
	neutralZones.AddPlan("D", neutralZone.QualityMedium, 1)
	neutralZones.AddPlan("E", neutralZone.QualityHigh, 1)
	return neutralZones
}

func TestWhenPlayerHasThreeNeutralPlans_CreatesSpawnPlusNeutralZones(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	neutralZones := newThreeNeutralPlans()
	tuning := models.NewGenerationTuning(configuration, 4)
	service := tournament_variant.NewRingClusterService()

	// Act
	zones, _ := service.CreateClusterVariant(*configuration, tuning, neutralZones, neutralZones, 0, "A")

	// Assert
	assert.Len(t, zones, 4)
}

func TestWhenRingIsBuilt_FirstZoneIsPlayerSpawn(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	neutralZones := newThreeNeutralPlans()
	tuning := models.NewGenerationTuning(configuration, 4)
	service := tournament_variant.NewRingClusterService()

	// Act
	zones, _ := service.CreateClusterVariant(*configuration, tuning, neutralZones, neutralZones, 0, "A")

	// Assert
	assert.Equal(t, "Spawn-A", zones[0].Name)
}

func TestWhenRingIsBuilt_CreatesConnectionPerRingSegment(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	neutralZones := newThreeNeutralPlans()
	tuning := models.NewGenerationTuning(configuration, 4)
	service := tournament_variant.NewRingClusterService()

	// Act
	_, connections := service.CreateClusterVariant(*configuration, tuning, neutralZones, neutralZones, 0, "A")

	// Assert
	assert.Len(t, connections, 4)
}

func TestWhenRingIsBuilt_EveryConnectionNameCarriesRingPrefix(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	neutralZones := newThreeNeutralPlans()
	tuning := models.NewGenerationTuning(configuration, 4)
	service := tournament_variant.NewRingClusterService()

	// Act
	_, connections := service.CreateClusterVariant(*configuration, tuning, neutralZones, neutralZones, 0, "A")

	// Assert
	prefixed := 0
	for _, connection := range connections {
		if strings.HasPrefix(connection.Name, "TRing-") {
			prefixed++
		}
	}
	assert.Equal(t, len(connections), prefixed)
}

func TestWhenRingIsBuilt_EveryConnectionReferencesExistingZones(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	neutralZones := newThreeNeutralPlans()
	tuning := models.NewGenerationTuning(configuration, 4)
	service := tournament_variant.NewRingClusterService()

	// Act
	zones, connections := service.CreateClusterVariant(*configuration, tuning, neutralZones, neutralZones, 0, "A")

	// Assert
	zoneNames := map[string]bool{}
	for _, zone := range zones {
		zoneNames[zone.Name] = true
	}
	var dangling []string
	for _, connection := range connections {
		if !zoneNames[connection.From] || !zoneNames[connection.To] {
			dangling = append(dangling, connection.Name)
		}
	}
	assert.Empty(t, dangling)
}

func TestWhenPlayerHasNoNeutralPlans_CreatesLoneSpawnZoneWithoutConnections(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	emptyPlans := neutralZone.Plans{}
	tuning := models.NewGenerationTuning(configuration, 1)
	service := tournament_variant.NewRingClusterService()

	// Act
	zones, connections := service.CreateClusterVariant(*configuration, tuning, emptyPlans, emptyPlans, 0, "A")

	// Assert
	assert.True(t, len(zones) == 1 && zones[0].Name == "Spawn-A" && len(connections) == 0,
		"expected a lone Spawn-A zone with no connections, got zones %v and %d connections", zones, len(connections))
}

func TestWhenSecondPlayerClusterIsBuilt_SpawnCastleBelongsToPlayerTwo(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	neutralZones := newThreeNeutralPlans()
	tuning := models.NewGenerationTuning(configuration, 4)
	service := tournament_variant.NewRingClusterService()

	// Act
	zones, _ := service.CreateClusterVariant(*configuration, tuning, neutralZones, neutralZones, 1, "B")

	// Assert
	assert.Equal(t, "Player2", zones[0].MainObjects[0].Spawn)
}
