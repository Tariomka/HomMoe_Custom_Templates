package circlesTopology_test

import (
	"strings"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutralZone"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology"
	"github.com/stretchr/testify/assert"
)

func TestWhenTwoPlayersAndFourTieredNeutralPlansProvided_CreatesZonePerLabel(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyCircles
	playerLabels := []string{"A", "B"}
	neutralZones := neutralZone.Plans{}
	neutralZones.AddPlan("N1", neutralZone.QualityLow, 0)
	neutralZones.AddPlan("N2", neutralZone.QualityLow, 0)
	neutralZones.AddPlan("N3", neutralZone.QualityMedium, 1)
	neutralZones.AddPlan("N4", neutralZone.QualityHigh, 1)
	tuning := models.NewGenerationTuning(configuration, 6)
	service := topology.NewCirclesTopologyService()

	// Act
	variant := service.CreateTopologyVariant(*configuration, playerLabels, neutralZones, tuning, "")

	// Assert
	assert.Len(t, variant.Zones, 6)
}

func TestWhenRingsAreStamped_EveryZoneGetsGeneratorRing(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyCircles
	playerLabels := []string{"A", "B"}
	neutralZones := neutralZone.Plans{}
	neutralZones.AddPlan("N1", neutralZone.QualityLow, 0)
	neutralZones.AddPlan("N2", neutralZone.QualityMedium, 1)
	neutralZones.AddPlan("N3", neutralZone.QualityHigh, 1)
	tuning := models.NewGenerationTuning(configuration, 5)
	service := topology.NewCirclesTopologyService()

	// Act
	variant := service.CreateTopologyVariant(*configuration, playerLabels, neutralZones, tuning, "")

	// Assert
	var zonesMissingRing []string
	for _, zone := range variant.Zones {
		if zone.GeneratorRing == nil {
			zonesMissingRing = append(zonesMissingRing, zone.Name)
		}
	}
	assert.Empty(t, zonesMissingRing)
}

func TestWhenPlayerZonesSitOnTheOuterRing_TheirRingIndexIsZero(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyCircles
	playerLabels := []string{"A", "B"}
	neutralZones := neutralZone.Plans{}
	neutralZones.AddPlan("N1", neutralZone.QualityLow, 0)
	neutralZones.AddPlan("N2", neutralZone.QualityMedium, 1)
	neutralZones.AddPlan("N3", neutralZone.QualityHigh, 1)
	tuning := models.NewGenerationTuning(configuration, 5)
	service := topology.NewCirclesTopologyService()

	// Act
	variant := service.CreateTopologyVariant(*configuration, playerLabels, neutralZones, tuning, "")

	// Assert
	var spawnZonesOffOuterRing []string
	for _, zone := range variant.Zones {
		if !strings.HasPrefix(zone.Name, "Spawn-") {
			continue
		}
		if zone.GeneratorRing == nil || *zone.GeneratorRing != 0 {
			spawnZonesOffOuterRing = append(spawnZonesOffOuterRing, zone.Name)
		}
	}
	assert.Empty(t, spawnZonesOffOuterRing)
}

func TestWhenCirclesAreBuilt_EveryConnectionReferencesExistingZones(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyCircles
	playerLabels := []string{"A", "B"}
	neutralZones := neutralZone.Plans{}
	neutralZones.AddPlan("N1", neutralZone.QualityLow, 0)
	neutralZones.AddPlan("N2", neutralZone.QualityLow, 0)
	neutralZones.AddPlan("N3", neutralZone.QualityMedium, 1)
	neutralZones.AddPlan("N4", neutralZone.QualityHigh, 1)
	tuning := models.NewGenerationTuning(configuration, 6)
	service := topology.NewCirclesTopologyService()

	// Act
	variant := service.CreateTopologyVariant(*configuration, playerLabels, neutralZones, tuning, "")

	// Assert
	assert.Empty(t, danglingConnectionNames(variant))
}

func TestWhenPlayerConnectionsAreForbidden_NoRandomConnectionJoinsTwoSpawnZones(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyCircles
	configuration.NoDirectPlayerConnections = true
	playerLabels := []string{"A", "B"}
	neutralZones := neutralZone.Plans{}
	neutralZones.AddPlan("N1", neutralZone.QualityMedium, 1)
	neutralZones.AddPlan("N2", neutralZone.QualityMedium, 1)
	tuning := models.NewGenerationTuning(configuration, 4)
	service := topology.NewCirclesTopologyService()

	// Act
	variant := service.CreateTopologyVariant(*configuration, playerLabels, neutralZones, tuning, "")

	// Assert
	assert.Empty(t, directSpawnToSpawnNames(variant, "Rnd-"))
}

func TestWhenRandomPortalsEnabled_AddsPortalConnections(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyCircles
	configuration.RandomPortals = true
	playerLabels := []string{"A", "B"}
	neutralZones := neutralZone.Plans{}
	neutralZones.AddPlan("N1", neutralZone.QualityLow, 0)
	neutralZones.AddPlan("N2", neutralZone.QualityLow, 0)
	neutralZones.AddPlan("N3", neutralZone.QualityMedium, 1)
	neutralZones.AddPlan("N4", neutralZone.QualityHigh, 1)
	tuning := models.NewGenerationTuning(configuration, 6)
	service := topology.NewCirclesTopologyService()

	// Act
	variant := service.CreateTopologyVariant(*configuration, playerLabels, neutralZones, tuning, "")

	// Assert
	assert.NotZero(t, countPortalConnections(variant))
}
