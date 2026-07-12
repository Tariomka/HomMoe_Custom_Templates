package templateGenerator_test

import (
	"fmt"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator"
	"github.com/stretchr/testify/assert"
)

// allGeneratorTopologies lists every topology the generator supports, so the
// cross-topology contract tests exercise each layout implementation.
func allGeneratorTopologies() []config.MapTopology {
	return []config.MapTopology{
		config.TopologyRing,
		config.TopologyChain,
		config.TopologyHubAndSpoke,
		config.TopologySharedWeb,
		config.TopologyRandom,
		config.TopologyCircles,
		config.TopologySquare,
		config.TopologyGeometric,
		config.TopologyCross,
		config.TopologyFractal,
	}
}

func TestWhenAnyTopologyWithVariedPlayerAndNeutralCounts_CreatesZoneForEveryPlannedZone(t *testing.T) {
	t.Parallel()
	for _, topology := range allGeneratorTopologies() {
		for _, playerCount := range []int{2, 3, 4, 8} {
			for _, neutralZoneCount := range []int{0, 1, 4} {
				subtestName := fmt.Sprintf(
					"When%sTopologyWith%dPlayersAnd%dNeutrals_CreatesZoneForEveryPlannedZone",
					topology, playerCount, neutralZoneCount)
				t.Run(subtestName, func(t *testing.T) {
					t.Parallel()
					// Arrange
					configuration := config.NewGeneratorConfig()
					configuration.Topology = topology
					configuration.PlayerCount = playerCount
					configuration.ZoneConfiguration.NeutralZoneCount = neutralZoneCount
					generator := template_generator.NewTemplateGenerator(configuration)

					// Act
					actual := generator.Generate()

					// Assert
					assert.GreaterOrEqual(t, len(actual.Variants[0].Zones), playerCount+neutralZoneCount)
				})
			}
		}
	}
}

//nolint:gocognit // This test is intentionally complex because it is a cross-topology contract test that exercises every topology implementation.
func TestWhenAnyTopologySelected_EveryZoneHasAllRequiredFields(t *testing.T) {
	t.Parallel()
	for _, topology := range allGeneratorTopologies() {
		subtestName := fmt.Sprintf("When%sTopologySelected_EveryZoneHasAllRequiredFields", topology)
		t.Run(subtestName, func(t *testing.T) {
			t.Parallel()
			// Arrange
			configuration := config.NewGeneratorConfig()
			configuration.Topology = topology
			configuration.PlayerCount = 3
			configuration.ZoneConfiguration.NeutralZoneCount = 3
			generator := template_generator.NewTemplateGenerator(configuration)

			// Act
			actual := generator.Generate()

			// Assert
			var violations []string
			for _, zone := range actual.Variants[0].Zones {
				if zone.Name == "" {
					violations = append(violations, "zone with empty name")
				}
				if zone.Layout == "" {
					violations = append(violations, zone.Name+": empty layout")
				}
				if zone.Size <= 0 || zone.Size > 2.0 {
					violations = append(violations, fmt.Sprintf("%s: size %f out of range", zone.Name, zone.Size))
				}
				if zone.GuardCutoffValue <= 0 {
					violations = append(violations, zone.Name+": non-positive guard cutoff")
				}
				if len(zone.GuardedContentPool) == 0 {
					violations = append(violations, zone.Name+": empty guarded content pool")
				}
				if len(zone.UnguardedContentPool) == 0 {
					violations = append(violations, zone.Name+": empty unguarded content pool")
				}
				if len(zone.ResourcesContentPool) == 0 {
					violations = append(violations, zone.Name+": empty resources content pool")
				}
			}
			assert.Empty(t, violations)
		})
	}
}

func TestWhenAnyTopologySelected_EveryConnectionReferencesExistingZones(t *testing.T) {
	t.Parallel()
	for _, topology := range allGeneratorTopologies() {
		subtestName := fmt.Sprintf("When%sTopologySelected_EveryConnectionReferencesExistingZones", topology)
		t.Run(subtestName, func(t *testing.T) {
			t.Parallel()
			// Arrange
			configuration := config.NewGeneratorConfig()
			configuration.Topology = topology
			configuration.PlayerCount = 3
			configuration.ZoneConfiguration.NeutralZoneCount = 2
			generator := template_generator.NewTemplateGenerator(configuration)

			// Act
			actual := generator.Generate()

			// Assert
			zoneNames := map[string]bool{}
			for _, zone := range actual.Variants[0].Zones {
				zoneNames[zone.Name] = true
			}
			var invalidReferences []string
			for _, connection := range actual.Variants[0].Connections {
				if !zoneNames[connection.From] {
					invalidReferences = append(
						invalidReferences,
						connection.Name+": unknown From zone "+connection.From)
				}
				if !zoneNames[connection.To] {
					invalidReferences = append(invalidReferences, connection.Name+": unknown To zone "+connection.To)
				}
			}
			assert.Empty(t, invalidReferences)
		})
	}
}
