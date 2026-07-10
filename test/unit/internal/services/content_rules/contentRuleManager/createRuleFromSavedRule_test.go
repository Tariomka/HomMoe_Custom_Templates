package contentRuleManager_test

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
	utopiaVariantID := 1
	utopiaVariantRule, err := content_rules.NewRuleVariant(&content_rules.UtopiaVariants, &utopiaVariantID)
	require.NoError(t, err)

	testCases := []struct {
		name     string
		original content_rules.ContentRule
		content  models.SidMapping
	}{
		{
			"WhenRuleIsRoadDistance_RoundTrips",
			content_rules.NewRuleDistanceToRoad(&content_rules.DistanceFar),
			models.SidMapping{Sid: "x"},
		},
		{
			"WhenRuleIsTownDistance_RoundTrips",
			content_rules.NewRuleDistanceToTown(&content_rules.DistanceNear),
			models.SidMapping{Sid: "x"},
		},
		{"WhenRuleIsGuarded_RoundTrips", content_rules.NewRuleGuarded(true), models.SidMapping{Sid: "x"}},
		{"WhenRuleIsUnguarded_RoundTrips", content_rules.NewRuleGuarded(false), models.SidMapping{Sid: "x"}},
		{"WhenRuleIsSoloEncounter_RoundTrips", content_rules.NewRuleSoloEncounter(true), models.SidMapping{Sid: "x"}},
		{"WhenRuleIsVariant_RoundTrips", utopiaVariantRule, constants.ContentIds.DragonUtopia},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			// Arrange
			saved := testCase.original.SerializeToRowSave()

			// Act
			restored := content_rules.CreateRuleFromSavedRule(saved, testCase.content)

			// Assert
			require.NotNil(t, restored)
			assert.Equal(t, saved, restored.SerializeToRowSave())
		})
	}
}

func TestWhenSavedNameDiffersOnlyByCase_RestoresRule(t *testing.T) {
	t.Parallel()
	// Arrange
	isGuarded := true
	saved := models.ContentRuleRowSave{Name: "gUaRdEd", IsGuarded: &isGuarded}

	// Act
	restored := content_rules.CreateRuleFromSavedRule(saved, models.SidMapping{Sid: "x"})

	// Assert
	assert.NotNil(t, restored)
}

func TestWhenSavedDataIsInvalid_ReturnsNil(t *testing.T) {
	t.Parallel()
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
			models.ContentRuleRowSave{Name: "Variant", VariantId: &invalidVariantID},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			// Arrange

			// Act
			restored := content_rules.CreateRuleFromSavedRule(testCase.saved, constants.ContentIds.DragonUtopia)

			// Assert
			assert.Nil(t, restored)
		})
	}

	t.Run("WhenContentHasNoVariants_ReturnsNil", func(t *testing.T) {
		t.Parallel()
		// Arrange
		saved := models.ContentRuleRowSave{Name: "Variant", VariantId: &someVariantID}

		// Act
		restored := content_rules.CreateRuleFromSavedRule(saved, models.SidMapping{Sid: "x"})

		// Assert
		assert.Nil(t, restored)
	})
}
