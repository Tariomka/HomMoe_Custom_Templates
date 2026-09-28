package zoneContentEditorService_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/content_rule_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zone_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenTheRuleTypeIsUnknown_TheCompositionIsRejected(t *testing.T) {
	t.Parallel()
	// Arrange
	service := zone_content.NewZoneContentEditorService()

	// Act
	_, ok := service.ComposeContentRule(content_rule_model.ContentRuleComposition{Name: gofakeit.Word()})

	// Assert
	assert.False(t, ok)
}

func TestWhenTheRuleTypeIsUnknown_ReturnsAnEmptyRule(t *testing.T) {
	t.Parallel()
	// Arrange
	service := zone_content.NewZoneContentEditorService()

	// Act
	rule, _ := service.ComposeContentRule(content_rule_model.ContentRuleComposition{Name: gofakeit.Word()})

	// Assert
	assert.Equal(t, editor_state_model.ContentRuleRow{}, rule)
}

func TestWhenTheDistanceIndexIsOutOfRange_TheCompositionIsRejected(t *testing.T) {
	t.Parallel()
	cases := []struct {
		scenario string
		key      content_rule_model.ContentRuleKey
		index    int
	}{
		{"DistanceToRoadNegativeIndex_IsRejected", content_rule_model.ContentRuleKeyDistanceToRoad, -1},
		{"DistanceToRoadIndexPastTheEnd_IsRejected", content_rule_model.ContentRuleKeyDistanceToRoad, 1},
		{"DistanceToTownNegativeIndex_IsRejected", content_rule_model.ContentRuleKeyDistanceToTown, -1},
		{"DistanceToTownIndexPastTheEnd_IsRejected", content_rule_model.ContentRuleKeyDistanceToTown, 1},
	}
	for _, testCase := range cases {
		t.Run(testCase.scenario, func(t *testing.T) {
			t.Parallel()
			// Arrange
			service := zone_content.NewZoneContentEditorService()

			// Act
			_, ok := service.ComposeContentRule(content_rule_model.ContentRuleComposition{
				Key:           testCase.key,
				Name:          gofakeit.Word(),
				DistanceNames: []string{gofakeit.Word()},
				DistanceIndex: testCase.index,
			})

			// Assert
			assert.False(t, ok)
		})
	}
}

func TestWhenADistanceToRoadRuleIsComposed_ItCarriesTheSelectedDistance(t *testing.T) {
	t.Parallel()
	// Arrange
	service := zone_content.NewZoneContentEditorService()
	name := gofakeit.Word()
	distance := gofakeit.Word()

	// Act
	rule, _ := service.ComposeContentRule(content_rule_model.ContentRuleComposition{
		Key:           content_rule_model.ContentRuleKeyDistanceToRoad,
		Name:          name,
		DistanceNames: []string{gofakeit.Word(), distance},
		DistanceIndex: 1,
	})

	// Assert
	assert.Equal(t, editor_state_model.ContentRuleRow{Name: name, DistanceName: distance}, rule)
}

func TestWhenADistanceToTownRuleIsComposed_ItCarriesTheSelectedDistance(t *testing.T) {
	t.Parallel()
	// Arrange
	service := zone_content.NewZoneContentEditorService()
	name := gofakeit.Word()
	distance := gofakeit.Word()

	// Act
	rule, _ := service.ComposeContentRule(content_rule_model.ContentRuleComposition{
		Key:           content_rule_model.ContentRuleKeyDistanceToTown,
		Name:          name,
		DistanceNames: []string{distance},
	})

	// Assert
	assert.Equal(t, editor_state_model.ContentRuleRow{Name: name, DistanceName: distance}, rule)
}

func TestWhenAGuardedRuleIsComposed_OnlyTheCheckboxValueIsSet(t *testing.T) {
	t.Parallel()
	cases := []struct {
		scenario string
		value    bool
	}{
		{"Checked_StoresTrue", true},
		{"Unchecked_StoresFalseRatherThanNil", false},
	}
	for _, testCase := range cases {
		t.Run(testCase.scenario, func(t *testing.T) {
			t.Parallel()
			// Arrange
			service := zone_content.NewZoneContentEditorService()
			name := gofakeit.Word()
			expected := testCase.value

			// Act
			rule, _ := service.ComposeContentRule(content_rule_model.ContentRuleComposition{
				Key:             content_rule_model.ContentRuleKeyGuarded,
				Name:            name,
				IsGuarded:       testCase.value,
				IsSoloEncounter: !testCase.value,
			})

			// Assert
			assert.Equal(t, editor_state_model.ContentRuleRow{Name: name, IsGuarded: &expected}, rule)
		})
	}
}

func TestWhenASoloEncounterRuleIsComposed_OnlyTheCheckboxValueIsSet(t *testing.T) {
	t.Parallel()
	cases := []struct {
		scenario string
		value    bool
	}{
		{"Checked_StoresTrue", true},
		{"Unchecked_StoresFalseRatherThanNil", false},
	}
	for _, testCase := range cases {
		t.Run(testCase.scenario, func(t *testing.T) {
			t.Parallel()
			// Arrange
			service := zone_content.NewZoneContentEditorService()
			name := gofakeit.Word()
			expected := testCase.value

			// Act
			rule, _ := service.ComposeContentRule(content_rule_model.ContentRuleComposition{
				Key:             content_rule_model.ContentRuleKeySoloEncounter,
				Name:            name,
				IsGuarded:       !testCase.value,
				IsSoloEncounter: testCase.value,
			})

			// Assert
			assert.Equal(t, editor_state_model.ContentRuleRow{Name: name, IsSoloEncounter: &expected}, rule)
		})
	}
}

func TestWhenTheVariantIndexIsOutOfRange_TheCompositionIsRejected(t *testing.T) {
	t.Parallel()
	cases := []struct {
		scenario string
		index    int
	}{
		{"NegativeIndex_IsRejected", -1},
		{"IndexPastTheEnd_IsRejected", 1},
	}
	for _, testCase := range cases {
		t.Run(testCase.scenario, func(t *testing.T) {
			t.Parallel()
			// Arrange
			service := zone_content.NewZoneContentEditorService()

			// Act
			_, ok := service.ComposeContentRule(content_rule_model.ContentRuleComposition{
				Key:          content_rule_model.ContentRuleKeyVariant,
				Name:         gofakeit.Word(),
				VariantIDs:   []int{gofakeit.Number(1, 500)},
				VariantIndex: testCase.index,
			})

			// Assert
			assert.False(t, ok)
		})
	}
}

func TestWhenAnIndexIsOutOfRange_ReturnsAnEmptyRule(t *testing.T) {
	t.Parallel()
	cases := []struct {
		scenario string
		key      content_rule_model.ContentRuleKey
		index    int
	}{
		{"DistanceNegativeIndex_ReturnsAnEmptyRule", content_rule_model.ContentRuleKeyDistanceToRoad, -1},
		{"DistanceIndexPastTheEnd_ReturnsAnEmptyRule", content_rule_model.ContentRuleKeyDistanceToTown, 1},
		{"VariantNegativeIndex_ReturnsAnEmptyRule", content_rule_model.ContentRuleKeyVariant, -1},
		{"VariantIndexPastTheEnd_ReturnsAnEmptyRule", content_rule_model.ContentRuleKeyVariant, 1},
	}
	for _, testCase := range cases {
		t.Run(testCase.scenario, func(t *testing.T) {
			t.Parallel()
			// Arrange
			service := zone_content.NewZoneContentEditorService()

			// Act
			rule, _ := service.ComposeContentRule(content_rule_model.ContentRuleComposition{
				Key:           testCase.key,
				Name:          gofakeit.Word(),
				DistanceNames: []string{gofakeit.Word()},
				DistanceIndex: testCase.index,
				VariantIDs:    []int{gofakeit.Number(1, 500)},
				VariantIndex:  testCase.index,
			})

			// Assert
			assert.Equal(t, editor_state_model.ContentRuleRow{}, rule)
		})
	}
}

func TestWhenAVariantRuleIsComposed_ItCarriesTheSelectedVariantId(t *testing.T) {
	t.Parallel()
	// Arrange
	service := zone_content.NewZoneContentEditorService()
	name := gofakeit.Word()
	variantID := gofakeit.Number(1, 500)

	// Act
	rule, _ := service.ComposeContentRule(content_rule_model.ContentRuleComposition{
		Key:          content_rule_model.ContentRuleKeyVariant,
		Name:         name,
		VariantIDs:   []int{gofakeit.Number(501, 900), variantID},
		VariantIndex: 1,
	})

	// Assert
	assert.Equal(t, editor_state_model.ContentRuleRow{Name: name, VariantID: &variantID}, rule)
}

func TestWhenAKnownRuleIsComposedFromValidInput_TheCompositionIsAccepted(t *testing.T) {
	t.Parallel()
	cases := []struct {
		scenario string
		key      content_rule_model.ContentRuleKey
	}{
		{"DistanceToRoad_IsAccepted", content_rule_model.ContentRuleKeyDistanceToRoad},
		{"DistanceToTown_IsAccepted", content_rule_model.ContentRuleKeyDistanceToTown},
		{"Guarded_IsAccepted", content_rule_model.ContentRuleKeyGuarded},
		{"SoloEncounter_IsAccepted", content_rule_model.ContentRuleKeySoloEncounter},
		{"Variant_IsAccepted", content_rule_model.ContentRuleKeyVariant},
	}
	for _, testCase := range cases {
		t.Run(testCase.scenario, func(t *testing.T) {
			t.Parallel()
			// Arrange
			service := zone_content.NewZoneContentEditorService()

			// Act
			_, ok := service.ComposeContentRule(content_rule_model.ContentRuleComposition{
				Key:           testCase.key,
				Name:          gofakeit.Word(),
				DistanceNames: []string{gofakeit.Word()},
				VariantIDs:    []int{gofakeit.Number(1, 500)},
			})

			// Assert
			assert.True(t, ok)
		})
	}
}
