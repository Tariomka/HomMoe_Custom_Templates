package squareTopology_test

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

func TestWhenTwoPlayersAndSixNeutralPlansProvided_CreatesZonePerLabel(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologySquare
	playerLabels := []string{"A", "B"}
	neutralZones := neutral_zone.Plans{}
	neutralZones.AddPlan("N1", neutral_zone.QualityLow, 0)
	neutralZones.AddPlan("N2", neutral_zone.QualityLow, 0)
	neutralZones.AddPlan("N3", neutral_zone.QualityMedium, 1)
	neutralZones.AddPlan("N4", neutral_zone.QualityMedium, 1)
	neutralZones.AddPlan("N5", neutral_zone.QualityHigh, 1)
	neutralZones.AddPlan("N6", neutral_zone.QualityHigh, 1)
	tuning := test_helpers.NewGenerationTuning(configuration, 8)
	service := topology.NewSquareTopologyService(test_helpers.NewZoneFactories())

	// Act
	variant := service.CreateTopologyVariant(*configuration, playerLabels, neutralZones, tuning, "")

	// Assert
	assert.Len(t, variant.Zones, 8)
}

func TestWhenSquareIsLaidOut_EveryZoneGetsPositionInsideUnitSquare(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologySquare
	playerLabels := []string{"A", "B"}
	neutralZones := neutral_zone.Plans{}
	neutralZones.AddPlan("N1", neutral_zone.QualityLow, 0)
	neutralZones.AddPlan("N2", neutral_zone.QualityMedium, 1)
	neutralZones.AddPlan("N3", neutral_zone.QualityMedium, 1)
	neutralZones.AddPlan("N4", neutral_zone.QualityHigh, 1)
	neutralZones.AddPlan("N5", neutral_zone.QualityHigh, 1)
	neutralZones.AddPlan("N6", neutral_zone.QualityHigh, 1)
	tuning := test_helpers.NewGenerationTuning(configuration, 8)
	service := topology.NewSquareTopologyService(test_helpers.NewZoneFactories())

	// Act
	variant := service.CreateTopologyVariant(*configuration, playerLabels, neutralZones, tuning, "")

	// Assert
	assert.Empty(t, zonesWithoutValidPosition(variant))
}

func TestWhenSingleInteriorNeutralExists_PlacesItAtTheSquareCenter(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologySquare
	playerLabels := []string{"A", "B"}
	neutralZones := neutral_zone.Plans{}
	neutralZones.AddPlan("N1", neutral_zone.QualityLow, 0)
	neutralZones.AddPlan("N2", neutral_zone.QualityMedium, 1)
	neutralZones.AddPlan("N3", neutral_zone.QualityHigh, 1)
	tuning := test_helpers.NewGenerationTuning(configuration, 5)
	service := topology.NewSquareTopologyService(test_helpers.NewZoneFactories())

	// Act
	variant := service.CreateTopologyVariant(*configuration, playerLabels, neutralZones, tuning, "")

	// Assert
	var interiorZone entities.Zone
	for _, zone := range variant.Zones {
		if zone.Name == "Neutral-N3" {
			interiorZone = zone
			break
		}
	}
	require.NotNil(t, interiorZone.GeneratorPosition)
	assert.Equal(t, [2]float64{0.5, 0.5}, *interiorZone.GeneratorPosition)
}

func TestWhenSquareIsBuilt_EveryConnectionReferencesExistingZones(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologySquare
	playerLabels := []string{"A", "B"}
	neutralZones := neutral_zone.Plans{}
	neutralZones.AddPlan("N1", neutral_zone.QualityLow, 0)
	neutralZones.AddPlan("N2", neutral_zone.QualityLow, 0)
	neutralZones.AddPlan("N3", neutral_zone.QualityMedium, 1)
	neutralZones.AddPlan("N4", neutral_zone.QualityMedium, 1)
	neutralZones.AddPlan("N5", neutral_zone.QualityHigh, 1)
	neutralZones.AddPlan("N6", neutral_zone.QualityHigh, 1)
	tuning := test_helpers.NewGenerationTuning(configuration, 8)
	service := topology.NewSquareTopologyService(test_helpers.NewZoneFactories())

	// Act
	variant := service.CreateTopologyVariant(*configuration, playerLabels, neutralZones, tuning, "")

	// Assert
	assert.Empty(t, danglingConnectionNames(variant))
}

func TestWhenPlayerConnectionsAreForbidden_NoRandomConnectionJoinsTwoSpawnZones(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologySquare
	configuration.NoDirectPlayerConnections = true
	playerLabels := []string{"A", "B"}
	neutralZones := neutral_zone.Plans{}
	neutralZones.AddPlan("N1", neutral_zone.QualityMedium, 1)
	neutralZones.AddPlan("N2", neutral_zone.QualityMedium, 1)
	tuning := test_helpers.NewGenerationTuning(configuration, 4)
	service := topology.NewSquareTopologyService(test_helpers.NewZoneFactories())

	// Act
	variant := service.CreateTopologyVariant(*configuration, playerLabels, neutralZones, tuning, "")

	// Assert
	assert.Empty(t, directSpawnToSpawnNames(variant, "Rnd-"))
}

func TestWhenRandomPortalsEnabled_AddsPortalConnections(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologySquare
	configuration.RandomPortals = true
	playerLabels := []string{"A", "B"}
	neutralZones := neutral_zone.Plans{}
	neutralZones.AddPlan("N1", neutral_zone.QualityLow, 0)
	neutralZones.AddPlan("N2", neutral_zone.QualityMedium, 1)
	neutralZones.AddPlan("N3", neutral_zone.QualityMedium, 1)
	neutralZones.AddPlan("N4", neutral_zone.QualityHigh, 1)
	tuning := test_helpers.NewGenerationTuning(configuration, 6)
	service := topology.NewSquareTopologyService(test_helpers.NewZoneFactories())

	// Act
	variant := service.CreateTopologyVariant(*configuration, playerLabels, neutralZones, tuning, "")

	// Assert
	assert.NotZero(t, countPortalConnections(variant))
}
