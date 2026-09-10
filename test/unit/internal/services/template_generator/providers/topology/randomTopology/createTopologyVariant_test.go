package randomTopology_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
)

func TestWhenTwoPlayersAndFourNeutralPlansProvided_CreatesZonePerLabel(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyRandom
	playerLabels := []string{"A", "B"}
	neutralZones := neutral_zone.Plans{}
	neutralZones.AddPlan("N1", neutral_zone.QualityMedium, 1)
	neutralZones.AddPlan("N2", neutral_zone.QualityMedium, 1)
	neutralZones.AddPlan("N3", neutral_zone.QualityMedium, 1)
	neutralZones.AddPlan("N4", neutral_zone.QualityMedium, 1)
	tuning := test_helpers.NewGenerationTuning(configuration, 6)
	service := topology.NewRandomTopologyService(test_helpers.NewZoneFactories())

	// Act
	variant := service.CreateTopologyVariant(*configuration, playerLabels, neutralZones, tuning, "")

	// Assert
	assert.Len(t, variant.Zones, 6)
}

func TestWhenPositionsAreRolled_EveryZoneGetsPositionInsideUnitSquare(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyRandom
	playerLabels := []string{"A", "B"}
	neutralZones := neutral_zone.Plans{}
	neutralZones.AddPlan("N1", neutral_zone.QualityMedium, 1)
	neutralZones.AddPlan("N2", neutral_zone.QualityMedium, 1)
	neutralZones.AddPlan("N3", neutral_zone.QualityMedium, 1)
	neutralZones.AddPlan("N4", neutral_zone.QualityMedium, 1)
	tuning := test_helpers.NewGenerationTuning(configuration, 6)
	service := topology.NewRandomTopologyService(test_helpers.NewZoneFactories())

	// Act
	variant := service.CreateTopologyVariant(*configuration, playerLabels, neutralZones, tuning, "")

	// Assert
	assert.Empty(t, zonesWithoutValidPosition(variant))
}

func TestWhenTriangulationIsBuilt_CreatesRandomPrefixedConnections(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyRandom
	playerLabels := []string{"A", "B"}
	neutralZones := neutral_zone.Plans{}
	neutralZones.AddPlan("N1", neutral_zone.QualityMedium, 1)
	neutralZones.AddPlan("N2", neutral_zone.QualityMedium, 1)
	neutralZones.AddPlan("N3", neutral_zone.QualityMedium, 1)
	neutralZones.AddPlan("N4", neutral_zone.QualityMedium, 1)
	tuning := test_helpers.NewGenerationTuning(configuration, 6)
	service := topology.NewRandomTopologyService(test_helpers.NewZoneFactories())

	// Act
	variant := service.CreateTopologyVariant(*configuration, playerLabels, neutralZones, tuning, "")

	// Assert
	assert.NotZero(t, countConnectionsWithPrefix(variant, "Rnd-"))
}

func TestWhenTriangulationIsBuilt_EveryConnectionReferencesExistingZones(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyRandom
	playerLabels := []string{"A", "B"}
	neutralZones := neutral_zone.Plans{}
	neutralZones.AddPlan("N1", neutral_zone.QualityMedium, 1)
	neutralZones.AddPlan("N2", neutral_zone.QualityMedium, 1)
	neutralZones.AddPlan("N3", neutral_zone.QualityMedium, 1)
	neutralZones.AddPlan("N4", neutral_zone.QualityMedium, 1)
	tuning := test_helpers.NewGenerationTuning(configuration, 6)
	service := topology.NewRandomTopologyService(test_helpers.NewZoneFactories())

	// Act
	variant := service.CreateTopologyVariant(*configuration, playerLabels, neutralZones, tuning, "")

	// Assert
	assert.Empty(t, danglingConnectionNames(variant))
}

func TestWhenPlayerConnectionsAreForbidden_NoRandomConnectionJoinsTwoSpawnZones(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyRandom
	configuration.NoDirectPlayerConnections = true
	playerLabels := []string{"A", "B"}
	neutralZones := neutral_zone.Plans{}
	neutralZones.AddPlan("N1", neutral_zone.QualityMedium, 1)
	neutralZones.AddPlan("N2", neutral_zone.QualityMedium, 1)
	tuning := test_helpers.NewGenerationTuning(configuration, 4)
	service := topology.NewRandomTopologyService(test_helpers.NewZoneFactories())

	// Act
	variant := service.CreateTopologyVariant(*configuration, playerLabels, neutralZones, tuning, "")

	// Assert
	assert.Empty(t, directSpawnToSpawnNames(variant, "Rnd-"))
}

func TestWhenRoadlessRandomIsolationHasNoNeutral_FallbackIsKeptWithoutConnectionRoads(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyRandom
	configuration.GenerateRoads = false
	configuration.NoDirectPlayerConnections = true
	playerLabels := []string{"A", "B"}
	tuning := test_helpers.NewGenerationTuning(configuration, 2)
	service := topology.NewRandomTopologyService(test_helpers.NewZoneFactories())

	// Act
	variant := service.CreateTopologyVariant(*configuration, playerLabels, nil, tuning, "")

	// Assert
	assert.Equal(t, []string{"Fallback-A-B"}, connectionNames(variant.Connections))
}

func TestWhenRoadlessRandomIsolationHasNoNeutral_DoesNotMaterializeFallbackConnectionRoads(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyRandom
	configuration.GenerateRoads = false
	configuration.NoDirectPlayerConnections = true
	playerLabels := []string{"A", "B"}
	tuning := test_helpers.NewGenerationTuning(configuration, 2)
	service := topology.NewRandomTopologyService(test_helpers.NewZoneFactories())

	// Act
	variant := service.CreateTopologyVariant(*configuration, playerLabels, nil, tuning, "")

	// Assert
	assert.Empty(t, connectionRoadTargets(variant.Zones))
}

func TestWhenRandomPortalsEnabled_AddsPortalConnections(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyRandom
	configuration.RandomPortals = true
	playerLabels := []string{"A", "B"}
	neutralZones := neutral_zone.Plans{}
	neutralZones.AddPlan("N1", neutral_zone.QualityMedium, 1)
	neutralZones.AddPlan("N2", neutral_zone.QualityMedium, 1)
	neutralZones.AddPlan("N3", neutral_zone.QualityMedium, 1)
	neutralZones.AddPlan("N4", neutral_zone.QualityMedium, 1)
	tuning := test_helpers.NewGenerationTuning(configuration, 6)
	service := topology.NewRandomTopologyService(test_helpers.NewZoneFactories())

	// Act
	variant := service.CreateTopologyVariant(*configuration, playerLabels, neutralZones, tuning, "")

	// Assert
	assert.NotZero(t, countPortalConnections(variant))
}

func connectionNames(connections []template_model.Connection) []string {
	names := make([]string, len(connections))
	for index, connection := range connections {
		names[index] = connection.Name
	}
	return names
}

func connectionRoadTargets(zones []template_model.Zone) []string {
	var connectionNames []string
	for _, zone := range zones {
		for _, road := range zone.Roads {
			if road.From.Type == "Connection" && len(road.From.Args) > 0 {
				connectionNames = append(connectionNames, road.From.Args[0])
			}
			if road.To.Type == "Connection" && len(road.To.Args) > 0 {
				connectionNames = append(connectionNames, road.To.Args[0])
			}
		}
	}
	return connectionNames
}
