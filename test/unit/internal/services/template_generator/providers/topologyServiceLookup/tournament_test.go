package topologyServiceLookup_test

import (
	"strings"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
)

func TestWhenTournamentCreatorIsInvoked_BuildsTournamentVariant(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyRing
	playerLabels := []string{"A", "B"}
	neutralZones := neutral_zone.Plans{}
	neutralZones.AddPlan("C", neutral_zone.QualityMedium, 1)
	neutralZones.AddPlan("D", neutral_zone.QualityMedium, 1)
	lookup := test_helpers.NewTopologyServiceLookup(test_helpers.NewZoneFactories())

	// Act
	variant := lookup.Tournament()(
		*configuration, playerLabels, neutralZones,
		test_helpers.NewGenerationTuning(configuration, len(playerLabels)+len(neutralZones)), "")

	// Assert
	hasTournamentGuardGroup := false
	for _, connection := range variant.Connections {
		if strings.HasPrefix(connection.GuardMatchGroup, "tourney_") {
			hasTournamentGuardGroup = true
			break
		}
	}
	assert.True(t, hasTournamentGuardGroup)
}
