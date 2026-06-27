package mandatoryContent_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers"
	"github.com/stretchr/testify/assert"
)

// contentNamed returns the content items of the mandatory-content group with
// the given name.
func contentNamed(groups []entities.MandatoryContent, name string) []entities.MandatoryContentItem {
	for _, group := range groups {
		if group.Name == name {
			return group.Content
		}
	}
	return nil
}

func sids(items []entities.MandatoryContentItem) []string {
	out := make([]string, 0, len(items))
	for _, item := range items {
		out = append(out, item.SID)
	}
	return out
}

func countGroups(groups []entities.MandatoryContent, name string) int {
	count := 0
	for _, group := range groups {
		if group.Name == name {
			count++
		}
	}
	return count
}

// CreateContents must copy the configured per-tier rows into each neutral zone's
// mandatory content. The original implementation used copy() into a nil slice,
// silently dropping every row - this guards that regression.
func TestCreateContents_CopiesNeutralTierRows(t *testing.T) {
	provider := providers.NewMandatoryContentProvider()
	configuration := config.NewGeneratorConfig()
	configuration.SpawnRemoteFootholds = false
	configuration.HighNeutralMandatoryContent = []entities.MandatoryContentItem{
		{SID: "university"},
		{SID: "random_item_legendary"},
	}

	plans := models.NeutralZonePlans{}
	plans.AddPlan("W", models.QualityHigh, 3)

	groups := provider.CreateContents(*configuration, nil, plans)

	assert.Equal(t, []string{"university", "random_item_legendary"},
		sids(contentNamed(groups, "mandatory_content_neutral_W")),
		"high neutral rows must reach a high-tier zone, not be dropped")
}

// CreateContentsForZones keys content off each zone's ACTUAL quality so a zone
// re-tiered in the manual editor (Medium plan -> High zone) gets High content.
// Regression test for the manually-promoted centre zone showing no high content.
func TestCreateContentsForZones_UsesActualZoneQuality(t *testing.T) {
	provider := providers.NewMandatoryContentProvider()
	configuration := config.NewGeneratorConfig()
	configuration.SpawnRemoteFootholds = false
	configuration.MediumNeutralMandatoryContent = []entities.MandatoryContentItem{{SID: "medium_only"}}
	configuration.HighNeutralMandatoryContent = []entities.MandatoryContentItem{{SID: "high_only"}}

	// A zone the generator labelled "G" but whose pool was manually raised to a
	// high tier (t4) with three castles.
	zones := []entities.Zone{
		{
			Name:               "Neutral-G",
			GuardedContentPool: []string{"classic_template_pool_random_t4_item"},
			MainObjects: []entities.MainObject{
				{Type: "City"}, {Type: "City"}, {Type: "City"},
			},
		},
	}

	groups := provider.CreateContentsForZones(*configuration, zones)

	assert.Equal(t, []string{"high_only"},
		sids(contentNamed(groups, "mandatory_content_neutral_G")),
		"a zone manually promoted to High must get High content, not its plan tier")
}

// A castle-less neutral zone must still receive content (with near-castle rules
// stripped), confirming the clone path does not drop rows for 0-castle zones.
func TestCreateContentsForZones_ZeroCastleZoneKeepsRows(t *testing.T) {
	provider := providers.NewMandatoryContentProvider()
	configuration := config.NewGeneratorConfig()
	configuration.SpawnRemoteFootholds = false
	configuration.MediumNeutralMandatoryContent = []entities.MandatoryContentItem{{SID: "treasure"}}

	zones := []entities.Zone{
		{
			Name:               "Neutral-H",
			GuardedContentPool: []string{"classic_template_pool_random_t3_item"},
		},
	}

	groups := provider.CreateContentsForZones(*configuration, zones)

	assert.Equal(t, []string{"treasure"},
		sids(contentNamed(groups, "mandatory_content_neutral_H")))
}

// The hub content group is created only for the Hub & Spoke topology and only
// when the user configured hub rows, matching the parallel C# editor.
func TestCreateContents_HubGroupOnlyForHubTopologyWithContent(t *testing.T) {
	provider := providers.NewMandatoryContentProvider()
	hubRows := []entities.MandatoryContentItem{{SID: "hub_treasure"}}

	withContent := config.NewGeneratorConfig()
	withContent.Topology = config.TopologyHubAndSpoke
	withContent.SpawnRemoteFootholds = false
	withContent.HubZoneMandatoryContent = hubRows
	groups := provider.CreateContents(*withContent, nil, nil)
	assert.Equal(t, []string{"hub_treasure"}, sids(contentNamed(groups, "mandatory_content_hub")),
		"hub topology with hub rows must emit the hub group")

	noContent := config.NewGeneratorConfig()
	noContent.Topology = config.TopologyHubAndSpoke
	assert.Equal(t, 0, countGroups(provider.CreateContents(*noContent, nil, nil), "mandatory_content_hub"),
		"hub topology without hub rows must not emit the hub group")

	nonHub := config.NewGeneratorConfig()
	nonHub.Topology = config.TopologyRing
	nonHub.HubZoneMandatoryContent = hubRows
	assert.Equal(t, 0, countGroups(provider.CreateContents(*nonHub, nil, nil), "mandatory_content_hub"),
		"non-hub topology must not emit the hub group even with hub rows")
}

// The manual-edit path must give the hub zone its content too, emitting a single
// shared group even when several hub zones exist (tournament clusters).
func TestCreateContentsForZones_HubZoneGetsSingleHubGroup(t *testing.T) {
	provider := providers.NewMandatoryContentProvider()
	configuration := config.NewGeneratorConfig()
	configuration.SpawnRemoteFootholds = false
	configuration.HubZoneMandatoryContent = []entities.MandatoryContentItem{{SID: "hub_treasure"}}

	zones := []entities.Zone{
		{Name: "Hub", MainObjects: []entities.MainObject{{Type: "City"}}},
		{Name: "Hub-B"},
	}

	groups := provider.CreateContentsForZones(*configuration, zones)

	assert.Equal(t, 1, countGroups(groups, "mandatory_content_hub"),
		"several hub zones must still share one hub group")
	assert.Equal(t, []string{"hub_treasure"}, sids(contentNamed(groups, "mandatory_content_hub")))
}
