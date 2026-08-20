package topologyProvider_test

import (
	"slices"
	"sort"
	"strings"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
)

// buildVariantInputs prepares the label/plan/tuning inputs CreateTopologyVariant
// needs for the given configuration.
func buildVariantInputs(
	configuration *config.GeneratorConfig,
	playerLabels []string,
	neutralZones neutral_zone.Plans) models.GenerationTuning {
	return test_helpers.NewGenerationTuning(configuration, len(playerLabels)+len(neutralZones))
}

// spawnZoneNames returns the sorted names of the variant's spawn zones.
func spawnZoneNames(variant entities.Variant) []string {
	var names []string
	for _, zone := range variant.Zones {
		if strings.HasPrefix(zone.Name, "Spawn-") {
			names = append(names, zone.Name)
		}
	}
	sort.Strings(names)
	return names
}

func TestWhenRingTopologySelected_CreatesZonePerLabelAndNeutralPlan(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyRing
	playerLabels := []string{"A", "B"}
	neutralZones := neutral_zone.Plans{}
	neutralZones.AddPlan("C", neutral_zone.QualityMedium, 1)
	neutralZones.AddPlan("D", neutral_zone.QualityMedium, 1)
	tuning := buildVariantInputs(configuration, playerLabels, neutralZones)
	provider := test_helpers.NewTopologyProvider()

	// Act
	variant := provider.CreateTopologyVariant(*configuration, playerLabels, neutralZones, tuning, "")

	// Assert
	assert.Len(t, variant.Zones, 4)
}

func TestWhenHubAndSpokeTopologySelected_CreatesHubZone(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyHubAndSpoke
	playerLabels := []string{"A", "B", "C"}
	neutralZones := neutral_zone.Plans{}
	neutralZones.AddPlan("D", neutral_zone.QualityMedium, 1)
	tuning := buildVariantInputs(configuration, playerLabels, neutralZones)
	provider := test_helpers.NewTopologyProvider()

	// Act
	variant := provider.CreateTopologyVariant(*configuration, playerLabels, neutralZones, tuning, "")

	// Assert
	hubZoneCount := 0
	for _, zone := range variant.Zones {
		if zone.Name == "Hub" {
			hubZoneCount++
		}
	}
	assert.Equal(t, 1, hubZoneCount)
}

func TestWhenGeometricHubTopologySelected_CreatesPositionedHubZone(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyGeometricHub
	playerLabels := []string{"A", "B"}
	tuning := buildVariantInputs(configuration, playerLabels, nil)
	provider := test_helpers.NewTopologyProvider()

	// Act
	variant := provider.CreateTopologyVariant(*configuration, playerLabels, nil, tuning, "")

	// Assert
	var hubPosition *[2]float64
	for _, zone := range variant.Zones {
		if zone.Name == "Hub" {
			hubPosition = zone.GeneratorPosition
		}
	}
	assert.NotNil(t, hubPosition)
}

func TestWhenTopologyIsUnknown_UsesRingProvider(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.MapTopology("Unknown")
	playerLabels := []string{"A", "B"}
	neutralZones := neutral_zone.Plans{}
	neutralZones.AddPlan("C", neutral_zone.QualityMedium, 1)
	neutralZones.AddPlan("D", neutral_zone.QualityMedium, 1)
	tuning := buildVariantInputs(configuration, playerLabels, neutralZones)
	provider := test_helpers.NewTopologyProvider()

	// Act
	variant := provider.CreateTopologyVariant(*configuration, playerLabels, neutralZones, tuning, "")

	// Assert
	assert.Len(t, variant.Connections, 4)
}

func TestWhenTournamentModeWithTwoPlayerLabels_CreatesTournamentVariant(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyRing
	configuration.TournamentRules = &config.TournamentRules{
		Enabled:            true,
		FirstTournamentDay: 14,
		Interval:           7,
		PointsToWin:        2,
	}
	playerLabels := []string{"A", "B"}
	neutralZones := neutral_zone.Plans{}
	neutralZones.AddPlan("C", neutral_zone.QualityMedium, 1)
	neutralZones.AddPlan("D", neutral_zone.QualityMedium, 1)
	tuning := buildVariantInputs(configuration, playerLabels, neutralZones)
	provider := test_helpers.NewTopologyProvider()

	// Act
	variant := provider.CreateTopologyVariant(*configuration, playerLabels, neutralZones, tuning, "")

	// Assert
	hasTournamentGuardGroup := false
	for _, connection := range variant.Connections {
		if strings.HasPrefix(connection.GuardMatchGroup, "tourney_") {
			hasTournamentGuardGroup = true
			break
		}
	}
	assert.True(t, hasTournamentGuardGroup, "tournament mode with 2 players must build the tournament variant")
}

func TestWhenTournamentModeWithThreePlayerLabels_UsesSelectedTopology(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyHubAndSpoke
	configuration.TournamentRules = &config.TournamentRules{
		Enabled:            true,
		FirstTournamentDay: 14,
		Interval:           7,
		PointsToWin:        2,
	}
	playerLabels := []string{"A", "B", "C"}
	neutralZones := neutral_zone.Plans{}
	neutralZones.AddPlan("D", neutral_zone.QualityMedium, 1)
	tuning := buildVariantInputs(configuration, playerLabels, neutralZones)
	provider := test_helpers.NewTopologyProvider()

	// Act
	variant := provider.CreateTopologyVariant(*configuration, playerLabels, neutralZones, tuning, "")

	// Assert: the regular hub topology creates a single shared "Hub" zone,
	// unlike the tournament variant's per-player "Hub-X" clusters.
	hubZoneCount := 0
	for _, zone := range variant.Zones {
		if zone.Name == "Hub" {
			hubZoneCount++
		}
	}
	assert.Equal(t, 1, hubZoneCount)
}

func TestWhenCalled_DoesNotMutateInputLabels(t *testing.T) {
	t.Parallel()
	// Arrange
	playerLabels := []string{"A", "B", "C", "D", "E", "F", "G", "H"}
	expectedLabels := slices.Clone(playerLabels)
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyRing
	tuning := buildVariantInputs(configuration, playerLabels, nil)
	provider := test_helpers.NewTopologyProvider()

	// Act
	provider.CreateTopologyVariant(*configuration, playerLabels, nil, tuning, "")

	// Assert
	assert.Equal(t, expectedLabels, playerLabels)
}

func TestWhenCalled_NamesOneSpawnZonePerPlayerLabel(t *testing.T) {
	t.Parallel()
	// Arrange
	playerLabels := []string{"A", "B", "C", "D"}
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyRing
	tuning := buildVariantInputs(configuration, playerLabels, nil)
	provider := test_helpers.NewTopologyProvider()

	// Act
	variant := provider.CreateTopologyVariant(*configuration, playerLabels, nil, tuning, "")

	// Assert
	assert.Equal(t, []string{"Spawn-A", "Spawn-B", "Spawn-C", "Spawn-D"}, spawnZoneNames(variant))
}
