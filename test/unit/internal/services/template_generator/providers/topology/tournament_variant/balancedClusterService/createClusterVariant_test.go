package balancedClusterService_test

import (
	"strings"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/tournament_variant"
	"github.com/stretchr/testify/assert"
)

func newFourNeutralPlans() models.NeutralZonePlans {
	neutralZones := models.NeutralZonePlans{}
	neutralZones.AddPlan("C", models.QualityLow, 0)
	neutralZones.AddPlan("D", models.QualityMedium, 1)
	neutralZones.AddPlan("E", models.QualityMedium, 1)
	neutralZones.AddPlan("F", models.QualityHigh, 1)
	return neutralZones
}

func TestWhenPlayerHasFourNeutralPlans_CreatesSpawnPlusNeutralZones(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	neutralZones := newFourNeutralPlans()
	tuning := models.NewGenerationTuning(configuration, 5)
	service := tournament_variant.NewBalancedClusterService()

	// Act
	zones, _ := service.CreateClusterVariant(*configuration, tuning, neutralZones, neutralZones, 0, "A")

	// Assert
	assert.Len(t, zones, 5)
}

func TestWhenClusterIsBuilt_EveryZoneReceivesGeneratorAndManualPositions(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	neutralZones := newFourNeutralPlans()
	tuning := models.NewGenerationTuning(configuration, 5)
	service := tournament_variant.NewBalancedClusterService()

	// Act
	zones, _ := service.CreateClusterVariant(*configuration, tuning, neutralZones, neutralZones, 0, "A")

	// Assert
	positioned := 0
	for _, zone := range zones {
		if zone.GeneratorPosition != nil && zone.ManualPosition != nil && zone.GeneratorRing != nil {
			positioned++
		}
	}
	assert.Equal(t, len(zones), positioned)
}

func TestWhenFirstPlayerClusterIsBuilt_AllZonePositionsStayInLeftHalf(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	neutralZones := newFourNeutralPlans()
	tuning := models.NewGenerationTuning(configuration, 5)
	service := tournament_variant.NewBalancedClusterService()

	// Act
	zones, _ := service.CreateClusterVariant(*configuration, tuning, neutralZones, neutralZones, 0, "A")

	// Assert
	inLeftHalf := 0
	for _, zone := range zones {
		if zone.GeneratorPosition[0] < 0.5 {
			inLeftHalf++
		}
	}
	assert.Equal(t, len(zones), inLeftHalf)
}

func TestWhenSecondPlayerClusterIsBuilt_AllZonePositionsStayInRightHalf(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	neutralZones := newFourNeutralPlans()
	tuning := models.NewGenerationTuning(configuration, 5)
	service := tournament_variant.NewBalancedClusterService()

	// Act
	zones, _ := service.CreateClusterVariant(*configuration, tuning, neutralZones, neutralZones, 1, "B")

	// Assert
	inRightHalf := 0
	for _, zone := range zones {
		if zone.GeneratorPosition[0] > 0.5 {
			inRightHalf++
		}
	}
	assert.Equal(t, len(zones), inRightHalf)
}

func TestWhenConnectionsAreBuilt_BalancedConnectionsCarryClusterPrefix(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	neutralZones := newFourNeutralPlans()
	tuning := models.NewGenerationTuning(configuration, 5)
	service := tournament_variant.NewBalancedClusterService()

	// Act
	_, connections := service.CreateClusterVariant(*configuration, tuning, neutralZones, neutralZones, 0, "A")

	// Assert
	assert.NotZero(t, countWithPrefix(connections, "TBal-"))
}

func TestWhenClusterIsBuilt_EveryConnectionReferencesExistingZones(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	neutralZones := newFourNeutralPlans()
	tuning := models.NewGenerationTuning(configuration, 5)
	service := tournament_variant.NewBalancedClusterService()

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

func TestWhenClusterIsBuilt_AllZonesFormSingleConnectedComponent(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	neutralZones := newFourNeutralPlans()
	tuning := models.NewGenerationTuning(configuration, 5)
	service := tournament_variant.NewBalancedClusterService()

	// Act
	zones, connections := service.CreateClusterVariant(*configuration, tuning, neutralZones, neutralZones, 0, "A")

	// Assert
	parent := map[string]string{}
	var find func(name string) string
	find = func(name string) string {
		if parent[name] != name {
			parent[name] = find(parent[name])
		}
		return parent[name]
	}
	for _, zone := range zones {
		parent[zone.Name] = zone.Name
	}
	for _, connection := range connections {
		if _, ok := parent[connection.From]; !ok {
			continue
		}
		if _, ok := parent[connection.To]; !ok {
			continue
		}
		rootFrom, rootTo := find(connection.From), find(connection.To)
		if rootFrom != rootTo {
			parent[rootFrom] = rootTo
		}
	}
	components := map[string]bool{}
	for _, zone := range zones {
		components[find(zone.Name)] = true
	}
	assert.Len(t, components, 1)
}

func countWithPrefix(connections []entities.Connection, prefix string) int {
	count := 0
	for _, connection := range connections {
		if strings.HasPrefix(connection.Name, prefix) {
			count++
		}
	}
	return count
}
