package hubClusterService_test

import (
	"strings"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/tournament_variant"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
)

func newTwoNeutralPlans() neutral_zone.Plans {
	neutralZones := neutral_zone.Plans{}
	neutralZones.AddPlan("C", neutral_zone.QualityLow, 0)
	neutralZones.AddPlan("D", neutral_zone.QualityHigh, 1)
	return neutralZones
}

func TestWhenPlayerHasTwoNeutralPlans_CreatesHubSpawnAndNeutralZones(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	neutralZones := newTwoNeutralPlans()
	tuning := test_helpers.NewGenerationTuning(configuration, 4)
	service := tournament_variant.NewHubClusterService(test_helpers.NewZoneFactories())

	// Act
	zones, _ := service.CreateClusterVariant(*configuration, tuning, neutralZones, neutralZones, 0, "A")

	// Assert
	zoneNames := make([]string, 0, len(zones))
	for _, zone := range zones {
		zoneNames = append(zoneNames, zone.Name)
	}
	assert.Equal(t, []string{"Hub-A", "Spawn-A", "Neutral-C", "Neutral-D"}, zoneNames)
}

func TestWhenClusterIsBuilt_HubZoneIsNamedAfterPlayerLabel(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	neutralZones := newTwoNeutralPlans()
	tuning := test_helpers.NewGenerationTuning(configuration, 4)
	service := tournament_variant.NewHubClusterService(test_helpers.NewZoneFactories())

	// Act
	zones, _ := service.CreateClusterVariant(*configuration, tuning, neutralZones, neutralZones, 1, "B")

	// Assert
	assert.Equal(t, "Hub-B", zones[0].Name)
}

func TestWhenASpokeConnectionTouchesTheHub_ItIsGuardedAtTheHighestTier(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	neutralZones := newTwoNeutralPlans()
	tuning := test_helpers.NewGenerationTuning(configuration, 4)
	tuning.BorderGuardStrengthMultiplier = 1.0
	service := tournament_variant.NewHubClusterService(test_helpers.NewZoneFactories())
	var hubGuardValues []int

	// Act
	_, connections := service.CreateClusterVariant(*configuration, tuning, neutralZones, neutralZones, 0, "A")

	// Assert
	for _, connection := range connections {
		if connection.From == "Hub-A" || connection.To == "Hub-A" {
			hubGuardValues = append(hubGuardValues, connection.GuardValue)
		}
	}
	assert.Equal(t, []int{35000, 35000, 35000}, hubGuardValues)
}

func TestWhenSpokesAreBuilt_CreatesSpokeConnectionPerSpokeZone(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	neutralZones := newTwoNeutralPlans()
	tuning := test_helpers.NewGenerationTuning(configuration, 4)
	service := tournament_variant.NewHubClusterService(test_helpers.NewZoneFactories())

	// Act
	_, connections := service.CreateClusterVariant(*configuration, tuning, neutralZones, neutralZones, 0, "A")

	// Assert
	spokeCount := 0
	for _, connection := range connections {
		if strings.HasPrefix(connection.Name, "THubSpoke-") {
			spokeCount++
		}
	}
	assert.Equal(t, 3, spokeCount)
}

func TestWhenSpokeConnectionsAreBuilt_EverySpokeStartsAtHubZone(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	neutralZones := newTwoNeutralPlans()
	tuning := test_helpers.NewGenerationTuning(configuration, 4)
	service := tournament_variant.NewHubClusterService(test_helpers.NewZoneFactories())

	// Act
	_, connections := service.CreateClusterVariant(*configuration, tuning, neutralZones, neutralZones, 0, "A")

	// Assert
	fromHub := 0
	spokeCount := 0
	for _, connection := range connections {
		if strings.HasPrefix(connection.Name, "THubSpoke-") {
			spokeCount++
			if connection.From == "Hub-A" {
				fromHub++
			}
		}
	}
	assert.Equal(t, spokeCount, fromHub)
}

func TestWhenProximityRingIsBuilt_CreatesProximityConnectionPerSpokePair(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	neutralZones := newTwoNeutralPlans()
	tuning := test_helpers.NewGenerationTuning(configuration, 4)
	service := tournament_variant.NewHubClusterService(test_helpers.NewZoneFactories())

	// Act
	_, connections := service.CreateClusterVariant(*configuration, tuning, neutralZones, neutralZones, 0, "A")

	// Assert
	proximityRingCount := 0
	for _, connection := range connections {
		if strings.HasPrefix(connection.Name, "THubRing-") && connection.ConnectionType == "Proximity" {
			proximityRingCount++
		}
	}
	assert.Equal(t, 3, proximityRingCount)
}

func TestWhenHubMandatoryContentIsConfigured_HubZoneReferencesHubContentGroup(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.HubZoneMandatoryContent = []entities.MandatoryContentItem{{}}
	neutralZones := newTwoNeutralPlans()
	tuning := test_helpers.NewGenerationTuning(configuration, 4)
	service := tournament_variant.NewHubClusterService(test_helpers.NewZoneFactories())

	// Act
	zones, _ := service.CreateClusterVariant(*configuration, tuning, neutralZones, neutralZones, 0, "A")

	// Assert
	assert.Contains(t, []string(zones[0].MandatoryContent), "mandatory_content_hub")
}

func TestWhenHubMandatoryContentIsEmpty_HubZoneHasNoMandatoryContent(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.HubZoneMandatoryContent = nil
	neutralZones := newTwoNeutralPlans()
	tuning := test_helpers.NewGenerationTuning(configuration, 4)
	service := tournament_variant.NewHubClusterService(test_helpers.NewZoneFactories())

	// Act
	zones, _ := service.CreateClusterVariant(*configuration, tuning, neutralZones, neutralZones, 0, "A")

	// Assert
	assert.Empty(t, zones[0].MandatoryContent)
}
