package templateGenerator_test

import (
	"slices"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/zone_helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenNeutralZonesArePlanned_EveryPlannedTierIsRecordedOnItsZone(t *testing.T) {
	t.Parallel()
	// Arrange
	generator := test_helpers.NewTemplateGenerator(newMixedTierConfiguration())

	// Act
	generated, _ := generator.Generate()

	// Assert
	var recordedTiers []neutral_zone.Quality
	for _, zone := range generated.Variants[0].Zones {
		if zone_helpers.IsZoneNameNeutral(zone.Name) && zone.Quality != nil {
			recordedTiers = append(recordedTiers, *zone.Quality)
		}
	}
	slices.Sort(recordedTiers)
	assert.Equal(t,
		[]neutral_zone.Quality{
			neutral_zone.QualityLowest,
			neutral_zone.QualityLow,
			neutral_zone.QualityMedium,
			neutral_zone.QualityHigh,
		},
		recordedTiers)
}

func TestWhenTopologyBuildsAHubZone_RecordsItAtTheHighestTier(t *testing.T) {
	t.Parallel()
	// Arrange
	generator := test_helpers.NewTemplateGenerator(newMixedTierConfiguration())

	// Act
	generated, _ := generator.Generate()

	// Assert
	hub, ok := findZone(generated, "Hub")
	require.True(t, ok)
	require.NotNil(t, hub.Quality)
	assert.Equal(t, neutral_zone.QualityHighest, *hub.Quality)
}

func TestWhenSpawnZonesAreGenerated_RecordsNoTierForThem(t *testing.T) {
	t.Parallel()
	// Arrange
	generator := test_helpers.NewTemplateGenerator(newMixedTierConfiguration())

	// Act
	generated, _ := generator.Generate()

	// Assert
	var tieredSpawnZoneNames []string
	for _, zone := range generated.Variants[0].Zones {
		if zone_helpers.IsZoneNamePlayer(zone.Name) && zone.Quality != nil {
			tieredSpawnZoneNames = append(tieredSpawnZoneNames, zone.Name)
		}
	}
	assert.Empty(t, tieredSpawnZoneNames)
}

func TestWhenTiersAreRecorded_CoversEveryNeutralZoneOfTheVariant(t *testing.T) {
	t.Parallel()
	// Arrange
	generator := test_helpers.NewTemplateGenerator(newMixedTierConfiguration())

	// Act
	generated, _ := generator.Generate()

	// Assert
	var unrecordedZoneNames []string
	for _, zone := range generated.Variants[0].Zones {
		if !zone_helpers.IsZoneNamePlayer(zone.Name) && zone.Quality == nil {
			unrecordedZoneNames = append(unrecordedZoneNames, zone.Name)
		}
	}
	assert.Empty(t, unrecordedZoneNames)
}

func findZone(template *template_model.Template, name string) (template_model.Zone, bool) {
	for _, zone := range template.Variants[0].Zones {
		if zone.Name == name {
			return zone, true
		}
	}
	return template_model.Zone{}, false
}

// newMixedTierConfiguration asks for one neutral zone per tier on a hub
// topology, so the variant has to carry both a hub and several distinct tiers.
func newMixedTierConfiguration() *config.GeneratorConfig {
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyHubAndSpoke
	configuration.PlayerCount = 2
	configuration.ZoneConfiguration.NeutralZoneCount = 0
	configuration.ZoneConfiguration.Advanced.Enabled = true
	configuration.ZoneConfiguration.Advanced.NeutralLowestNoCastleCount = 1
	configuration.ZoneConfiguration.Advanced.NeutralLowNoCastleCount = 1
	configuration.ZoneConfiguration.Advanced.NeutralMediumNoCastleCount = 1
	configuration.ZoneConfiguration.Advanced.NeutralHighNoCastleCount = 1
	return configuration
}
