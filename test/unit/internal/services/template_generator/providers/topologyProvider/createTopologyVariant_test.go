package topologyProvider_test

import (
	"slices"
	"sort"
	"strings"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_errors"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// unsupportedTopologies are IDs no service exists for: the four retired ones,
// an unknown one and the empty one.
var unsupportedTopologies = []config.MapTopology{ //nolint:gochecknoglobals // Read-only test table.
	"Default", "HubAndSpoke", "Chain", "SharedWeb", "NotARealTopology", "",
}

func TestWhenSquareTopologySelected_CreatesZonePerLabelAndNeutralPlan(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologySquare
	playerLabels := []string{"A", "B"}
	neutralZones := twoNeutralPlans()
	tuning := buildVariantInputs(configuration, playerLabels, neutralZones)
	provider := test_helpers.NewTopologyProvider()

	// Act
	variant, err := provider.CreateTopologyVariant(*configuration, playerLabels, neutralZones, tuning, "")

	// Assert
	require.NoError(t, err)
	assert.Len(t, variant.Zones, 4)
}

func TestWhenGeometricHubTopologySelected_CreatesSingleHubZone(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyGeometricHub
	playerLabels := []string{"A", "B", "C"}
	neutralZones := neutral_zone.Plans{}
	neutralZones.AddPlan("D", neutral_zone.QualityMedium, 1)
	tuning := buildVariantInputs(configuration, playerLabels, neutralZones)
	provider := test_helpers.NewTopologyProvider()

	// Act
	variant, err := provider.CreateTopologyVariant(*configuration, playerLabels, neutralZones, tuning, "")

	// Assert
	require.NoError(t, err)
	assert.Equal(t, 1, countZonesNamed(variant, "Hub"))
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
	variant, err := provider.CreateTopologyVariant(*configuration, playerLabels, nil, tuning, "")

	// Assert
	require.NoError(t, err)
	var hubPosition *data.Vec2[float64]
	for _, zone := range variant.Zones {
		if zone.Name == "Hub" {
			hubPosition = zone.GeneratorPosition
		}
	}
	assert.NotNil(t, hubPosition)
}

func TestWhenTopologyIsUnsupported_ReturnsUnsupportedTopologyError(t *testing.T) {
	t.Parallel()
	for _, mapTopology := range unsupportedTopologies {
		for _, tournament := range []bool{false, true} {
			mode := map[bool]string{false: "Off", true: "On"}[tournament]
			name := "Topology" + string(mapTopology) + "_Tournament" + mode
			t.Run(name+"_ReturnsUnsupportedTopologyError", func(t *testing.T) {
				t.Parallel()
				// Arrange
				configuration := unsupportedConfiguration(mapTopology, tournament)
				playerLabels := []string{"A", "B"}
				tuning := buildVariantInputs(configuration, playerLabels, nil)
				provider := test_helpers.NewTopologyProvider()

				// Act
				_, err := provider.CreateTopologyVariant(*configuration, playerLabels, nil, tuning, "")

				// Assert
				assert.ErrorIs(t, err, common_errors.ErrUnsupportedTopology)
			})
		}
	}
}

func TestWhenTopologyIsUnsupported_BuildsNoVariant(t *testing.T) {
	t.Parallel()
	for _, tournament := range []bool{false, true} {
		t.Run(map[bool]string{false: "TournamentOff", true: "TournamentOn"}[tournament]+"_BuildsNoVariant",
			func(t *testing.T) {
				t.Parallel()
				// Arrange
				configuration := unsupportedConfiguration("Chain", tournament)
				playerLabels := []string{"A", "B"}
				tuning := buildVariantInputs(configuration, playerLabels, nil)
				provider := test_helpers.NewTopologyProvider()

				// Act
				variant, _ := provider.CreateTopologyVariant(*configuration, playerLabels, nil, tuning, "")

				// Assert
				assert.Equal(t, template_model.Variant{}, variant)
			})
	}
}

func TestWhenTopologyIsUnsupported_ErrorNamesTheTopology(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := unsupportedConfiguration("NotARealTopology", false)
	playerLabels := []string{"A", "B"}
	tuning := buildVariantInputs(configuration, playerLabels, nil)
	provider := test_helpers.NewTopologyProvider()

	// Act
	_, err := provider.CreateTopologyVariant(*configuration, playerLabels, nil, tuning, "")

	// Assert
	assert.EqualError(t, err, `unsupported topology: "NotARealTopology"`)
}

func TestWhenTournamentModeWithTwoPlayerLabels_CreatesTournamentVariant(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := tournamentConfiguration(config.TopologyRandom)
	playerLabels := []string{"A", "B"}
	neutralZones := twoNeutralPlans()
	tuning := buildVariantInputs(configuration, playerLabels, neutralZones)
	provider := test_helpers.NewTopologyProvider()

	// Act
	variant, err := provider.CreateTopologyVariant(*configuration, playerLabels, neutralZones, tuning, "")

	// Assert
	require.NoError(t, err)
	assert.True(t, hasGuardGroupWithPrefix(variant, "tourney_"),
		"tournament mode with 2 players must build the tournament variant")
}

// Every surviving topology that used the chain fallback now gets the balanced builder.
func TestWhenTournamentTopologyIsAnySupportedOne_BuildsBalancedClusters(t *testing.T) {
	t.Parallel()
	for _, mapTopology := range []config.MapTopology{
		config.TopologyRandom, config.TopologyCircles, config.TopologyGeometricHub, config.TopologySquare,
		config.TopologyGeometric, config.TopologyCross, config.TopologyFractal,
	} {
		t.Run(string(mapTopology)+"_BuildsBalancedClusters", func(t *testing.T) {
			t.Parallel()
			// Arrange
			configuration := tournamentConfiguration(mapTopology)
			playerLabels := []string{"A", "B"}
			neutralZones := twoNeutralPlans()
			tuning := buildVariantInputs(configuration, playerLabels, neutralZones)
			provider := test_helpers.NewTopologyProvider()

			// Act
			variant, err := provider.CreateTopologyVariant(*configuration, playerLabels, neutralZones, tuning, "")

			// Assert
			require.NoError(t, err)
			assert.True(t, hasGuardGroupWithPrefix(variant, "tourney_bal_guard_"))
		})
	}
}

func TestWhenTournamentModeWithThreePlayerLabels_UsesSelectedTopology(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := tournamentConfiguration(config.TopologyGeometricHub)
	playerLabels := []string{"A", "B", "C"}
	neutralZones := neutral_zone.Plans{}
	neutralZones.AddPlan("D", neutral_zone.QualityMedium, 1)
	tuning := buildVariantInputs(configuration, playerLabels, neutralZones)
	provider := test_helpers.NewTopologyProvider()

	// Act
	variant, err := provider.CreateTopologyVariant(*configuration, playerLabels, neutralZones, tuning, "")

	// Assert: the regular hub topology creates a single shared "Hub" zone.
	require.NoError(t, err)
	assert.Equal(t, 1, countZonesNamed(variant, "Hub"))
}

func TestWhenCalled_DoesNotMutateInputLabels(t *testing.T) {
	t.Parallel()
	// Arrange
	playerLabels := []string{"A", "B", "C", "D", "E", "F", "G", "H"}
	expectedLabels := slices.Clone(playerLabels)
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologySquare
	tuning := buildVariantInputs(configuration, playerLabels, nil)
	provider := test_helpers.NewTopologyProvider()

	// Act
	_, _ = provider.CreateTopologyVariant(*configuration, playerLabels, nil, tuning, "")

	// Assert
	assert.Equal(t, expectedLabels, playerLabels)
}

func TestWhenCalled_NamesOneSpawnZonePerPlayerLabel(t *testing.T) {
	t.Parallel()
	// Arrange
	playerLabels := []string{"A", "B", "C", "D"}
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologySquare
	tuning := buildVariantInputs(configuration, playerLabels, nil)
	provider := test_helpers.NewTopologyProvider()

	// Act
	variant, err := provider.CreateTopologyVariant(*configuration, playerLabels, nil, tuning, "")

	// Assert
	require.NoError(t, err)
	assert.Equal(t, []string{"Spawn-A", "Spawn-B", "Spawn-C", "Spawn-D"}, spawnZoneNames(variant))
}

// buildVariantInputs prepares the label/plan/tuning inputs CreateTopologyVariant
// needs for the given configuration.
func buildVariantInputs(
	configuration *config.GeneratorConfig,
	playerLabels []string,
	neutralZones neutral_zone.Plans) models.GenerationTuning {
	return test_helpers.NewGenerationTuning(configuration, len(playerLabels)+len(neutralZones))
}

func twoNeutralPlans() neutral_zone.Plans {
	neutralZones := neutral_zone.Plans{}
	neutralZones.AddPlan("C", neutral_zone.QualityMedium, 1)
	neutralZones.AddPlan("D", neutral_zone.QualityMedium, 1)
	return neutralZones
}

func tournamentConfiguration(mapTopology config.MapTopology) *config.GeneratorConfig {
	configuration := config.NewGeneratorConfig()
	configuration.Topology = mapTopology
	configuration.TournamentRules = &config.TournamentRules{
		Enabled:            true,
		FirstTournamentDay: 14,
		Interval:           7,
		PointsToWin:        2,
	}
	return configuration
}

func unsupportedConfiguration(mapTopology config.MapTopology, tournament bool) *config.GeneratorConfig {
	if tournament {
		return tournamentConfiguration(mapTopology)
	}
	configuration := config.NewGeneratorConfig()
	configuration.Topology = mapTopology
	return configuration
}

func hasGuardGroupWithPrefix(variant template_model.Variant, prefix string) bool {
	for _, connection := range variant.Connections {
		if strings.HasPrefix(connection.GuardMatchGroup, prefix) {
			return true
		}
	}
	return false
}

func countZonesNamed(variant template_model.Variant, name string) int {
	count := 0
	for _, zone := range variant.Zones {
		if zone.Name == name {
			count++
		}
	}
	return count
}

// spawnZoneNames returns the sorted names of the variant's spawn zones.
func spawnZoneNames(variant template_model.Variant) []string {
	var names []string
	for _, zone := range variant.Zones {
		if strings.HasPrefix(zone.Name, "Spawn-") {
			names = append(names, zone.Name)
		}
	}
	sort.Strings(names)
	return names
}
