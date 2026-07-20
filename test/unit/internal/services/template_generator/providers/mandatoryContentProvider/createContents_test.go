package mandatoryContentProvider_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers"
	"github.com/stretchr/testify/assert"
)

func TestWhenPlayerLabelsProvided_CreatesGroupPerPlayerLabel(t *testing.T) {
	t.Parallel()
	// Arrange
	provider := providers.NewMandatoryContentProvider()
	configuration := config.NewGeneratorConfig()
	configuration.SpawnRemoteFootholds = false

	// Act
	groups := provider.CreateContents(*configuration, []string{"A", "B", "C"}, nil)

	// Assert
	assert.Equal(t, []string{
		"mandatory_content_side_A",
		"mandatory_content_side_B",
		"mandatory_content_side_C",
	}, groupNames(groups))
}

func TestWhenRemoteFootholdsEnabled_PrependsFootholdItemPerCount(t *testing.T) {
	t.Parallel()
	// Arrange
	footholdSid := registry.GetMapObjectNonContentValues().RemoteFoothold
	provider := providers.NewMandatoryContentProvider()
	configuration := config.NewGeneratorConfig()
	configuration.SpawnRemoteFootholds = true
	configuration.RemoteFootholdCount = 3
	configuration.PlayerZoneMandatoryContent = []entities.MandatoryContentItem{{SID: "sawmill"}}

	// Act
	groups := provider.CreateContents(*configuration, []string{"A"}, nil)

	// Assert
	assert.Equal(t, []string{footholdSid, footholdSid, footholdSid, "sawmill"},
		itemSids(groupContent(groups, "mandatory_content_side_A")))
}

func TestWhenRemoteFootholdsDisabled_AddsNoFootholdItems(t *testing.T) {
	t.Parallel()
	// Arrange
	provider := providers.NewMandatoryContentProvider()
	configuration := config.NewGeneratorConfig()
	configuration.SpawnRemoteFootholds = false
	configuration.PlayerZoneMandatoryContent = []entities.MandatoryContentItem{{SID: "sawmill"}}

	// Act
	groups := provider.CreateContents(*configuration, []string{"A"}, nil)

	// Assert
	assert.Equal(t, []string{"sawmill"}, itemSids(groupContent(groups, "mandatory_content_side_A")))
}

func TestWhenLowTierRowsConfigured_CopiesRowsIntoLowNeutralZone(t *testing.T) {
	t.Parallel()
	// Arrange
	provider := providers.NewMandatoryContentProvider()
	configuration := config.NewGeneratorConfig()
	configuration.SpawnRemoteFootholds = false
	configuration.LowNeutralMandatoryContent = []entities.MandatoryContentItem{{SID: "low_only"}}
	plans := neutral_zone.Plans{}
	plans.AddPlan("C", neutral_zone.QualityLow, 1)

	// Act
	groups := provider.CreateContents(*configuration, nil, plans)

	// Assert
	assert.Equal(t, []string{"low_only"}, itemSids(groupContent(groups, "mandatory_content_neutral_C")))
}

func TestWhenMediumTierRowsConfigured_CopiesRowsIntoMediumNeutralZone(t *testing.T) {
	t.Parallel()
	// Arrange
	provider := providers.NewMandatoryContentProvider()
	configuration := config.NewGeneratorConfig()
	configuration.SpawnRemoteFootholds = false
	configuration.MediumNeutralMandatoryContent = []entities.MandatoryContentItem{{SID: "medium_only"}}
	plans := neutral_zone.Plans{}
	plans.AddPlan("C", neutral_zone.QualityMedium, 1)

	// Act
	groups := provider.CreateContents(*configuration, nil, plans)

	// Assert
	assert.Equal(t, []string{"medium_only"}, itemSids(groupContent(groups, "mandatory_content_neutral_C")))
}

// CreateContents must copy the configured per-tier rows into each neutral zone's
// mandatory content. The original implementation used copy() into a nil slice,
// silently dropping every row - this guards that regression.
func TestWhenHighTierRowsConfigured_CopiesRowsIntoHighNeutralZone(t *testing.T) {
	t.Parallel()
	// Arrange
	provider := providers.NewMandatoryContentProvider()
	configuration := config.NewGeneratorConfig()
	configuration.SpawnRemoteFootholds = false
	configuration.HighNeutralMandatoryContent = []entities.MandatoryContentItem{
		{SID: "university"},
		{SID: "random_item_legendary"},
	}
	plans := neutral_zone.Plans{}
	plans.AddPlan("W", neutral_zone.QualityHigh, 3)

	// Act
	groups := provider.CreateContents(*configuration, nil, plans)

	// Assert
	assert.Equal(t, []string{"university", "random_item_legendary"},
		itemSids(groupContent(groups, "mandatory_content_neutral_W")),
		"high neutral rows must reach a high-tier zone, not be dropped")
}

func TestWhenHighestTierPlanExists_CopiesHubZoneRowsIntoThatNeutralZone(t *testing.T) {
	t.Parallel()
	// Arrange
	provider := providers.NewMandatoryContentProvider()
	configuration := config.NewGeneratorConfig()
	configuration.SpawnRemoteFootholds = false
	configuration.HighNeutralMandatoryContent = []entities.MandatoryContentItem{{SID: "high_only"}}
	configuration.HubZoneMandatoryContent = []entities.MandatoryContentItem{{SID: "hub_treasure"}}
	plans := neutral_zone.Plans{}
	plans.AddPlan("V", neutral_zone.QualityHighest, 1)

	// Act
	groups := provider.CreateContents(*configuration, nil, plans)

	// Assert
	assert.Equal(t, []string{"hub_treasure"},
		itemSids(groupContent(groups, "mandatory_content_neutral_V")),
		"a Highest-quality plan must receive the hub zone rows, not the high-tier rows")
}

func TestWhenNeutralZoneHasNoCastles_StripsNearCastlePlacementRules(t *testing.T) {
	t.Parallel()
	// Arrange
	ruleTypeMainObject := registry.GetRuleTypeValues().MainObject
	provider := providers.NewMandatoryContentProvider()
	configuration := config.NewGeneratorConfig()
	configuration.SpawnRemoteFootholds = false
	configuration.MediumNeutralMandatoryContent = []entities.MandatoryContentItem{{
		SID: "treasure",
		Rules: []entities.PlacementRule{
			{Type: ruleTypeMainObject, Args: []any{"0"}},
			{Type: ruleTypeMainObject, Args: []any{"1"}},
		},
	}}
	plans := neutral_zone.Plans{}
	plans.AddPlan("C", neutral_zone.QualityMedium, 0)

	// Act
	groups := provider.CreateContents(*configuration, nil, plans)

	// Assert
	assert.Equal(t, []entities.PlacementRule{{Type: ruleTypeMainObject, Args: []any{"1"}}},
		groupContent(groups, "mandatory_content_neutral_C")[0].Rules,
		"only the near-castle (main object 0) rule must be removed")
}

func TestWhenNeutralZoneHasCastles_KeepsNearCastlePlacementRules(t *testing.T) {
	t.Parallel()
	// Arrange
	ruleTypeMainObject := registry.GetRuleTypeValues().MainObject
	provider := providers.NewMandatoryContentProvider()
	configuration := config.NewGeneratorConfig()
	configuration.SpawnRemoteFootholds = false
	configuration.MediumNeutralMandatoryContent = []entities.MandatoryContentItem{{
		SID:   "treasure",
		Rules: []entities.PlacementRule{{Type: ruleTypeMainObject, Args: []any{"0"}}},
	}}
	plans := neutral_zone.Plans{}
	plans.AddPlan("C", neutral_zone.QualityMedium, 1)

	// Act
	groups := provider.CreateContents(*configuration, nil, plans)

	// Assert
	assert.Equal(t, []entities.PlacementRule{{Type: ruleTypeMainObject, Args: []any{"0"}}},
		groupContent(groups, "mandatory_content_neutral_C")[0].Rules)
}

// Guards the cloneContentItems fix: stripping rules for a 0-castle zone must not
// corrupt the shared per-tier rows held on the configuration.
func TestWhenZeroCastleZoneStripsRules_DoesNotMutateConfiguredRows(t *testing.T) {
	t.Parallel()
	// Arrange
	ruleTypeMainObject := registry.GetRuleTypeValues().MainObject
	provider := providers.NewMandatoryContentProvider()
	configuration := config.NewGeneratorConfig()
	configuration.SpawnRemoteFootholds = false
	configuration.MediumNeutralMandatoryContent = []entities.MandatoryContentItem{{
		SID: "treasure",
		Rules: []entities.PlacementRule{
			{Type: ruleTypeMainObject, Args: []any{"0"}},
			{Type: ruleTypeMainObject, Args: []any{"1"}},
		},
	}}
	plans := neutral_zone.Plans{}
	plans.AddPlan("C", neutral_zone.QualityMedium, 0)

	// Act
	provider.CreateContents(*configuration, nil, plans)

	// Assert
	assert.Len(t, configuration.MediumNeutralMandatoryContent[0].Rules, 2,
		"configured rows must keep all their rules after generation")
}

// The hub content group is created only for the Hub & Spoke topology and only
// when the user configured hub rows, matching the parallel C# editor.
func TestWhenHubTopologyWithHubRows_EmitsHubGroupWithConfiguredRows(t *testing.T) {
	t.Parallel()
	// Arrange
	provider := providers.NewMandatoryContentProvider()
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyHubAndSpoke
	configuration.SpawnRemoteFootholds = false
	configuration.HubZoneMandatoryContent = []entities.MandatoryContentItem{{SID: "hub_treasure"}}

	// Act
	groups := provider.CreateContents(*configuration, nil, nil)

	// Assert
	assert.Equal(t, []string{"hub_treasure"}, itemSids(groupContent(groups, "mandatory_content_hub")))
}

func TestWhenHubTopologyWithoutHubRows_OmitsHubGroup(t *testing.T) {
	t.Parallel()
	// Arrange
	provider := providers.NewMandatoryContentProvider()
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyHubAndSpoke

	// Act
	groups := provider.CreateContents(*configuration, nil, nil)

	// Assert
	assert.Equal(t, 0, countMandatoryContentHubs(groups))
}

func TestWhenNonHubTopologyWithHubRows_OmitsHubGroup(t *testing.T) {
	t.Parallel()
	// Arrange
	provider := providers.NewMandatoryContentProvider()
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyRing
	configuration.HubZoneMandatoryContent = []entities.MandatoryContentItem{{SID: "hub_treasure"}}

	// Act
	groups := provider.CreateContents(*configuration, nil, nil)

	// Assert
	assert.Equal(t, 0, countMandatoryContentHubs(groups))
}
