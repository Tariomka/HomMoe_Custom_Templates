package contentRuleService_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/constants"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/content_rules"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenSavedRuleIsValid_RestoresRuleThatSerializesBack(t *testing.T) {
	t.Parallel()
	service := content_rules.NewContentRuleService()
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
		{"WhenRuleIsVariant_RoundTrips", utopiaVariantRule, constants.ContentIDs.DragonUtopia},
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
	saved := models.ContentRuleRowSave{Name: "gUaRdEd", IsGuarded: &isGuarded}

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
		saved models.ContentRuleRowSave
	}{
		{"WhenNameIsUnknown_ReturnsNil", models.ContentRuleRowSave{Name: "Nope"}},
		{"WhenGuardedValueIsMissing_ReturnsNil", models.ContentRuleRowSave{Name: "Guarded"}},
		{"WhenSoloEncounterValueIsMissing_ReturnsNil", models.ContentRuleRowSave{Name: "Solo Encounter"}},
		{"WhenVariantIdIsMissing_ReturnsNil", models.ContentRuleRowSave{Name: "Variant"}},
		{
			"WhenRoadDistanceNameIsUnknown_ReturnsNil",
			models.ContentRuleRowSave{Name: "Distance to road", DistanceName: "Whatever"},
		},
		{
			"WhenTownDistanceNameIsUnknown_ReturnsNil",
			models.ContentRuleRowSave{Name: "Distance to town", DistanceName: "Whatever"},
		},
		{
			"WhenVariantIdIsNotDefinedForContent_ReturnsNil",
			models.ContentRuleRowSave{Name: "Variant", VariantID: &invalidVariantID},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			// Arrange

			// Act
			restored := service.CreateRuleFromSavedRule(testCase.saved, constants.ContentIDs.DragonUtopia)

			// Assert
			assert.Nil(t, restored)
		})
	}

	t.Run("WhenContentHasNoVariants_ReturnsNil", func(t *testing.T) {
		t.Parallel()
		// Arrange
		saved := models.ContentRuleRowSave{Name: "Variant", VariantID: &someVariantID}

		// Act
		restored := service.CreateRuleFromSavedRule(saved, models.SidMapping{Sid: "x"})

		// Assert
		assert.Nil(t, restored)
	})
}
