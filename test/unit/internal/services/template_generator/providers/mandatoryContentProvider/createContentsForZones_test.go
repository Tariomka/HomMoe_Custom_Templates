package mandatoryContentProvider_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers"
	"github.com/stretchr/testify/assert"
)

// CreateContentsForZones keys content off each zone's ACTUAL quality so a zone
// re-tiered in the manual editor (Medium plan -> High zone) gets High content.
// Regression test for the manually-promoted center zone showing no high content.
func TestWhenZoneManuallyPromotedToHighTier_UsesHighTierRows(t *testing.T) {
	t.Parallel()
	// Arrange
	provider := providers.NewMandatoryContentProvider()
	configuration := config.NewGeneratorConfig()
	configuration.SpawnRemoteFootholds = false
	configuration.MediumNeutralMandatoryContent = []entities.MandatoryContentItem{{SID: "medium_only"}}
	configuration.HighNeutralMandatoryContent = []entities.MandatoryContentItem{{SID: "high_only"}}
	// A zone the generator labelled "G" but whose profile was manually raised to
	// the high tier (treasure layout, t4 pool) with three castles.
	zones := []entities.Zone{{
		Name:               "Neutral-G",
		Layout:             registry.GetLayoutValues().TreasureZone,
		GuardedContentPool: []string{"classic_template_pool_random_t4_item"},
		MainObjects: []entities.MainObject{
			{Type: "City"}, {Type: "City"}, {Type: "City"},
		},
	}}

	// Act
	groups := provider.CreateContentsForZones(*configuration, zones)

	// Assert
	assert.Equal(t, []string{"high_only"},
		itemSids(groupContent(groups, "mandatory_content_neutral_G")),
		"a zone manually promoted to High must get High content, not its plan tier")
}

func TestWhenZonePoolIsLowTier_UsesLowTierRows(t *testing.T) {
	t.Parallel()
	// Arrange
	provider := providers.NewMandatoryContentProvider()
	configuration := config.NewGeneratorConfig()
	configuration.SpawnRemoteFootholds = false
	configuration.LowNeutralMandatoryContent = []entities.MandatoryContentItem{{SID: "low_only"}}
	configuration.MediumNeutralMandatoryContent = []entities.MandatoryContentItem{{SID: "medium_only"}}
	zones := []entities.Zone{{
		Name:               "Neutral-C",
		Layout:             registry.GetLayoutValues().Sides,
		GuardedContentPool: []string{"classic_template_pool_random_t2_item"},
		MainObjects:        []entities.MainObject{{Type: "City"}},
	}}

	// Act
	groups := provider.CreateContentsForZones(*configuration, zones)

	// Assert
	assert.Equal(t, []string{"low_only"}, itemSids(groupContent(groups, "mandatory_content_neutral_C")))
}

func TestWhenZonePoolIsLowestTier_UsesLowestTierRows(t *testing.T) {
	t.Parallel()
	// Arrange
	provider := providers.NewMandatoryContentProvider()
	configuration := config.NewGeneratorConfig()
	configuration.SpawnRemoteFootholds = false
	configuration.LowestNeutralMandatoryContent = []entities.MandatoryContentItem{{SID: "lowest_only"}}
	configuration.LowNeutralMandatoryContent = []entities.MandatoryContentItem{{SID: "low_only"}}
	zones := []entities.Zone{{
		Name:               "Neutral-C",
		Layout:             registry.GetLayoutValues().Sides,
		GuardedContentPool: []string{"classic_template_pool_random_t1_item"},
		MainObjects:        []entities.MainObject{{Type: "City"}},
	}}

	// Act
	groups := provider.CreateContentsForZones(*configuration, zones)

	// Assert
	assert.Equal(t, []string{"lowest_only"}, itemSids(groupContent(groups, "mandatory_content_neutral_C")))
}

func TestWhenZoneHasNoRecognizableQuality_CreatesAnEmptyGroup(t *testing.T) {
	t.Parallel()
	// Arrange
	provider := providers.NewMandatoryContentProvider()
	configuration := config.NewGeneratorConfig()
	configuration.SpawnRemoteFootholds = false
	configuration.MediumNeutralMandatoryContent = []entities.MandatoryContentItem{{SID: "medium_only"}}
	zones := []entities.Zone{{
		Name:        "Neutral-C",
		MainObjects: []entities.MainObject{{Type: "City"}},
	}}

	// Act
	groups := provider.CreateContentsForZones(*configuration, zones)

	// Assert
	assert.Empty(t, itemSids(groupContent(groups, "mandatory_content_neutral_C")))
}

// A castle-less neutral zone must still receive content (with near-castle rules
// stripped), confirming the clone path does not drop rows for 0-castle zones.
func TestWhenZoneHasNoCastles_KeepsConfiguredRows(t *testing.T) {
	t.Parallel()
	// Arrange
	provider := providers.NewMandatoryContentProvider()
	configuration := config.NewGeneratorConfig()
	configuration.SpawnRemoteFootholds = false
	configuration.MediumNeutralMandatoryContent = []entities.MandatoryContentItem{{SID: "treasure"}}
	zones := []entities.Zone{{
		Name:               "Neutral-H",
		Layout:             registry.GetLayoutValues().TreasureZone,
		GuardedContentPool: []string{"classic_template_pool_random_t3_item"},
	}}

	// Act
	groups := provider.CreateContentsForZones(*configuration, zones)

	// Assert
	assert.Equal(t, []string{"treasure"}, itemSids(groupContent(groups, "mandatory_content_neutral_H")))
}

func TestWhenSpawnZoneProvided_CreatesPlayerGroupNamedAfterZoneSuffix(t *testing.T) {
	t.Parallel()
	// Arrange
	provider := providers.NewMandatoryContentProvider()
	configuration := config.NewGeneratorConfig()
	configuration.SpawnRemoteFootholds = false
	configuration.PlayerZoneMandatoryContent = []entities.MandatoryContentItem{{SID: "sawmill"}}
	zones := []entities.Zone{{Name: "Spawn-B"}}

	// Act
	groups := provider.CreateContentsForZones(*configuration, zones)

	// Assert
	assert.Equal(t, []string{"sawmill"}, itemSids(groupContent(groups, "mandatory_content_side_B")))
}

func TestWhenZoneNameIsUnrecognized_SkipsZone(t *testing.T) {
	t.Parallel()
	// Arrange
	provider := providers.NewMandatoryContentProvider()
	configuration := config.NewGeneratorConfig()
	configuration.MediumNeutralMandatoryContent = []entities.MandatoryContentItem{{SID: "treasure"}}
	zones := []entities.Zone{{Name: "SomethingElse"}}

	// Act
	groups := provider.CreateContentsForZones(*configuration, zones)

	// Assert
	assert.Empty(t, groups)
}

// The manual-edit path must give the hub zone its content too, emitting a single
// shared group even when several hub zones exist (tournament clusters).
func TestWhenMultipleHubZonesProvided_EmitsSingleHubGroup(t *testing.T) {
	t.Parallel()
	// Arrange
	provider := providers.NewMandatoryContentProvider()
	configuration := config.NewGeneratorConfig()
	configuration.SpawnRemoteFootholds = false
	configuration.HubZoneMandatoryContent = []entities.MandatoryContentItem{{SID: "hub_treasure"}}
	zones := []entities.Zone{
		{Name: "Hub", MainObjects: []entities.MainObject{{Type: "City"}}},
		{Name: "Hub-B"},
	}

	// Act
	groups := provider.CreateContentsForZones(*configuration, zones)

	// Assert
	assert.Equal(t, 1, countMandatoryContentHubs(groups),
		"several hub zones must still share one hub group")
}

func TestWhenMultipleHubZonesProvided_SharedHubGroupContainsConfiguredRows(t *testing.T) {
	t.Parallel()
	// Arrange
	provider := providers.NewMandatoryContentProvider()
	configuration := config.NewGeneratorConfig()
	configuration.SpawnRemoteFootholds = false
	configuration.HubZoneMandatoryContent = []entities.MandatoryContentItem{{SID: "hub_treasure"}}
	zones := []entities.Zone{
		{Name: "Hub", MainObjects: []entities.MainObject{{Type: "City"}}},
		{Name: "Hub-B"},
	}

	// Act
	groups := provider.CreateContentsForZones(*configuration, zones)

	// Assert
	assert.Equal(t, []string{"hub_treasure"}, itemSids(groupContent(groups, "mandatory_content_hub")))
}

func TestWhenHubZoneProvidedWithoutHubRows_OmitsHubGroup(t *testing.T) {
	t.Parallel()
	// Arrange
	provider := providers.NewMandatoryContentProvider()
	configuration := config.NewGeneratorConfig()
	zones := []entities.Zone{{Name: "Hub", MainObjects: []entities.MainObject{{Type: "City"}}}}

	// Act
	groups := provider.CreateContentsForZones(*configuration, zones)

	// Assert
	assert.Equal(t, 0, countMandatoryContentHubs(groups))
}
