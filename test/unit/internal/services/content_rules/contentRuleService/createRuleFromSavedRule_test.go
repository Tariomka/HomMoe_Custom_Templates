package contentRuleService_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/content_rules"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The SID is duplicated here on purpose: internal-layer tests must not import the GUI catalogue.
const dragonUtopiaSid = "dragon_utopia"

func TestWhenSavedRuleIsValid_RestoresRuleThatSerializesBack(t *testing.T) {
	t.Parallel()
	service := content_rules.NewContentRuleService()
	dragonUtopia := models.SidMapping{Sid: dragonUtopiaSid}
	far := models.DistancePreset{Name: "Far", Min: 0.5, Max: 0.75}
	near := models.DistancePreset{Name: "Near", Min: 0.1, Max: 0.25}
	utopiaVariantID := 1
	defaultMapping := content_rules.NewVariantMappingCatalog().GetDefaultMapping()
	utopiaVariantRule, err := content_rules.NewRuleVariant(&defaultMapping, &utopiaVariantID)
	require.NoError(t, err)

	testCases := []struct {
		name     string
		original content_rules.IContentRule
		content  models.SidMapping
	}{
		{"WhenRuleIsRoadDistance_RoundTrips", content_rules.NewRuleDistanceToRoad(&far), models.SidMapping{Sid: "x"}},
		{"WhenRuleIsTownDistance_RoundTrips", content_rules.NewRuleDistanceToTown(&near), models.SidMapping{Sid: "x"}},
		{"WhenRuleIsGuarded_RoundTrips", content_rules.NewRuleGuarded(true), models.SidMapping{Sid: "x"}},
		{"WhenRuleIsUnguarded_RoundTrips", content_rules.NewRuleGuarded(false), models.SidMapping{Sid: "x"}},
		{"WhenRuleIsSoloEncounter_RoundTrips", content_rules.NewRuleSoloEncounter(true), models.SidMapping{Sid: "x"}},
		{"WhenRuleIsVariant_RoundTrips", utopiaVariantRule, dragonUtopia},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			// Arrange
			saved := testCase.original.SerializeToRowSave()

			// Act
			restored := service.CreateRuleFromSavedRule(saved, testCase.content)

			// Assert
			require.NotNil(t, restored)
			assert.Equal(t, saved, restored.SerializeToRowSave())
		})
	}
}

func TestWhenSavedNameDiffersOnlyByCase_RestoresRule(t *testing.T) {
	t.Parallel()
	// Arrange
	service := content_rules.NewContentRuleService()
	isGuarded := true
	saved := models.ContentRuleRow{Name: "gUaRdEd", IsGuarded: &isGuarded}

	// Act
	restored := service.CreateRuleFromSavedRule(saved, models.SidMapping{Sid: "x"})

	// Assert
	assert.NotNil(t, restored)
}

func TestWhenSavedDataIsInvalid_ReturnsNil(t *testing.T) {
	t.Parallel()
	service := content_rules.NewContentRuleService()
	invalidVariantID := 99
	someVariantID := 0
	testCases := []struct {
		name  string
		saved models.ContentRuleRow
	}{
		{"WhenNameIsUnknown_ReturnsNil", models.ContentRuleRow{Name: "Nope"}},
		{"WhenGuardedValueIsMissing_ReturnsNil", models.ContentRuleRow{Name: "Guarded"}},
		{"WhenSoloEncounterValueIsMissing_ReturnsNil", models.ContentRuleRow{Name: "Solo Encounter"}},
		{"WhenVariantIdIsMissing_ReturnsNil", models.ContentRuleRow{Name: "Variant"}},
		{
			"WhenRoadDistanceNameIsUnknown_ReturnsNil",
			models.ContentRuleRow{Name: "Distance to road", DistanceName: "Whatever"},
		},
		{
			"WhenTownDistanceNameIsUnknown_ReturnsNil",
			models.ContentRuleRow{Name: "Distance to town", DistanceName: "Whatever"},
		},
		{
			"WhenVariantIdIsNotDefinedForContent_ReturnsNil",
			models.ContentRuleRow{Name: "Variant", VariantID: &invalidVariantID},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			// Arrange

			// Act
			restored := service.CreateRuleFromSavedRule(testCase.saved, models.SidMapping{Sid: dragonUtopiaSid})

			// Assert
			assert.Nil(t, restored)
		})
	}

	t.Run("WhenContentHasNoVariants_ReturnsNil", func(t *testing.T) {
		t.Parallel()
		// Arrange
		saved := models.ContentRuleRow{Name: "Variant", VariantID: &someVariantID}

		// Act
		restored := service.CreateRuleFromSavedRule(saved, models.SidMapping{Sid: "x"})

		// Assert
		assert.Nil(t, restored)
	})
}
