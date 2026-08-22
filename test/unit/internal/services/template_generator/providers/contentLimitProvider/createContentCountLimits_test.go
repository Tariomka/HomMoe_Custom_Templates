package contentLimitProvider_test

import (
	"strings"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenDefaultConfiguration_CreatesSeventeenLimitGroups(t *testing.T) {
	t.Parallel()
	// Arrange
	provider := providers.NewContentLimitProvider()
	configuration := config.NewGeneratorConfig()

	// Act
	groups := provider.CreateContentCountLimits(*configuration)

	// Assert: base group + 0_0 group + the 15 side-pair groups (1..5 x higher).
	assert.Len(t, groups, 17)
}

func TestWhenDefaultConfiguration_NamesGroupsAfterSidePairs(t *testing.T) {
	t.Parallel()
	// Arrange
	provider := providers.NewContentLimitProvider()
	configuration := config.NewGeneratorConfig()

	// Act
	groups := provider.CreateContentCountLimits(*configuration)

	// Assert
	var actualNames []string
	for _, group := range groups {
		actualNames = append(actualNames, group.Name)
	}
	assert.Equal(t, []string{
		"content_limits_side",
		"content_limits_side_0_0",
		"content_limits_side_1_2",
		"content_limits_side_1_3",
		"content_limits_side_1_4",
		"content_limits_side_1_5",
		"content_limits_side_1_6",
		"content_limits_side_2_3",
		"content_limits_side_2_4",
		"content_limits_side_2_5",
		"content_limits_side_2_6",
		"content_limits_side_3_4",
		"content_limits_side_3_5",
		"content_limits_side_3_6",
		"content_limits_side_4_5",
		"content_limits_side_4_6",
		"content_limits_side_5_6",
	}, actualNames)
}

func TestWhenDefaultConfiguration_AllGroupsShareSameLimits(t *testing.T) {
	t.Parallel()
	// Arrange
	provider := providers.NewContentLimitProvider()
	configuration := config.NewGeneratorConfig()

	// Act
	groups := provider.CreateContentCountLimits(*configuration)

	// Assert
	assert.Equal(t, groups[0].Limits, groups[len(groups)-1].Limits)
}

func TestWhenMandatoryContentBelowDefaultCap_KeepsDefaultCap(t *testing.T) {
	t.Parallel()
	// Arrange
	fountainSid := registry.GetMapObjectHeroBuffBuildingValues().Fountain
	provider := providers.NewContentLimitProvider()
	configuration := config.NewGeneratorConfig()
	configuration.PlayerZoneMandatoryContent = []entities.MandatoryContentItem{{SID: fountainSid}}

	// Act
	groups := provider.CreateContentCountLimits(*configuration)

	// Assert: the fountain default cap is 2, one requested item must not lift it.
	assert.Equal(t, 2, maxCountFor(t, groups, fountainSid))
}

func TestWhenMandatoryContentExceedsDefaultCap_LiftsLimitToRequestedCount(t *testing.T) {
	t.Parallel()
	// Arrange
	fountainSid := registry.GetMapObjectHeroBuffBuildingValues().Fountain
	provider := providers.NewContentLimitProvider()
	configuration := config.NewGeneratorConfig()
	configuration.PlayerZoneMandatoryContent = []entities.MandatoryContentItem{
		{SID: fountainSid}, {SID: fountainSid}, {SID: fountainSid}, {SID: fountainSid},
	}

	// Act
	groups := provider.CreateContentCountLimits(*configuration)

	// Assert
	assert.Equal(t, 4, maxCountFor(t, groups, fountainSid))
}

func TestWhenRequestedSidsSpanSeveralContentLists_SumsCountsAcrossLists(t *testing.T) {
	t.Parallel()
	// Arrange
	fountainSid := registry.GetMapObjectHeroBuffBuildingValues().Fountain
	provider := providers.NewContentLimitProvider()
	configuration := config.NewGeneratorConfig()
	configuration.PlayerZoneMandatoryContent = []entities.MandatoryContentItem{
		{SID: fountainSid}, {SID: fountainSid},
	}
	configuration.LowNeutralMandatoryContent = []entities.MandatoryContentItem{
		{SID: fountainSid}, {SID: fountainSid}, {SID: fountainSid},
	}

	// Act
	groups := provider.CreateContentCountLimits(*configuration)

	// Assert
	assert.Equal(t, 5, maxCountFor(t, groups, fountainSid))
}

func TestWhenRequestedSidCaseDiffersFromLimitSid_StillLiftsLimit(t *testing.T) {
	t.Parallel()
	// Arrange
	fountainSid := registry.GetMapObjectHeroBuffBuildingValues().Fountain
	upperCaseSid := strings.ToUpper(fountainSid)
	provider := providers.NewContentLimitProvider()
	configuration := config.NewGeneratorConfig()
	configuration.PlayerZoneMandatoryContent = []entities.MandatoryContentItem{
		{SID: upperCaseSid}, {SID: upperCaseSid}, {SID: upperCaseSid},
	}

	// Act
	groups := provider.CreateContentCountLimits(*configuration)

	// Assert
	assert.Equal(t, 3, maxCountFor(t, groups, fountainSid))
}

func TestWhenMandatoryContentUsesUnlimitedSid_AddsNoNewLimitEntry(t *testing.T) {
	t.Parallel()
	// Arrange
	provider := providers.NewContentLimitProvider()
	defaultGroups := provider.CreateContentCountLimits(*config.NewGeneratorConfig())
	configuration := config.NewGeneratorConfig()
	configuration.PlayerZoneMandatoryContent = []entities.MandatoryContentItem{
		{SID: "some_sid_without_default_limit"},
	}

	// Act
	groups := provider.CreateContentCountLimits(*configuration)

	// Assert
	assert.Equal(t, defaultGroups[0].Limits, groups[0].Limits)
}

func TestWhenLowestTierRequestsMoreThanDefaultCap_LiftsLimit(t *testing.T) {
	t.Parallel()
	// Arrange
	fountainSid := registry.GetMapObjectHeroBuffBuildingValues().Fountain
	provider := providers.NewContentLimitProvider()
	configuration := config.NewGeneratorConfig()
	configuration.LowestNeutralMandatoryContent = []entities.MandatoryContentItem{
		{SID: fountainSid}, {SID: fountainSid}, {SID: fountainSid},
	}

	// Act
	groups := provider.CreateContentCountLimits(*configuration)

	// Assert
	assert.Equal(t, 3, maxCountFor(t, groups, fountainSid))
}

// maxCountFor returns the MaxCount of the limit with the given SID in the
// first limit group, and fails the test when the SID has no limit entry.
func maxCountFor(t *testing.T, groups []entities.ContentCountLimit, sid string) int {
	t.Helper()
	require.NotEmpty(t, groups)
	for _, limit := range groups[0].Limits {
		if strings.EqualFold(limit.SID, sid) {
			return limit.MaxCount
		}
	}
	t.Fatalf("no content limit found for SID %q", sid)
	return 0
}
