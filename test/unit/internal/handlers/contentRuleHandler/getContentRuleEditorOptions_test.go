package contentRuleHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/content_rules"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWhenContentHasNoVariants_OffersOnlyTheFixedRules(t *testing.T) {
	t.Parallel()
	// Arrange
	content := models.SidMapping{Sid: gofakeit.Word()}
	service := &test_helpers.ContentRuleServiceMock{}
	service.On("GetVariantsForContent", content).Return([]models.VariantMapping{})
	service.On("GetDistanceDisplayNames").Return([]string{})
	handler := handlers.NewContentRuleHandler(service)

	// Act
	options := handler.GetContentRuleEditorOptions(content)

	// Assert
	assert.Equal(t, []dtos.ContentRuleKey{
		dtos.ContentRuleKeyDistanceToRoad,
		dtos.ContentRuleKeyDistanceToTown,
		dtos.ContentRuleKeyGuarded,
		dtos.ContentRuleKeySoloEncounter,
	}, ruleKeysOf(options))
}

func TestWhenContentHasVariants_AppendsTheVariantRule(t *testing.T) {
	t.Parallel()
	// Arrange
	content := models.SidMapping{Sid: gofakeit.Word()}
	service := &test_helpers.ContentRuleServiceMock{}
	service.On("GetVariantsForContent", content).Return(oneVariantMapping(1, gofakeit.Word()))
	service.On("GetDistanceDisplayNames").Return([]string{})
	handler := handlers.NewContentRuleHandler(service)

	// Act
	options := handler.GetContentRuleEditorOptions(content)

	// Assert
	assert.Equal(t, dtos.ContentRuleKeyVariant, ruleKeysOf(options)[len(options.Rules)-1])
}

func TestWhenContentHasVariants_FlattensThemIntoOptions(t *testing.T) {
	t.Parallel()
	// Arrange
	content := models.SidMapping{Sid: gofakeit.Word()}
	variantLabel := gofakeit.Word()
	service := &test_helpers.ContentRuleServiceMock{}
	service.On("GetVariantsForContent", content).Return(oneVariantMapping(7, variantLabel))
	service.On("GetDistanceDisplayNames").Return([]string{})
	handler := handlers.NewContentRuleHandler(service)

	// Act
	options := handler.GetContentRuleEditorOptions(content)

	// Assert
	assert.Equal(t, []dtos.ContentRuleVariantOptionDto{{ID: 7, Label: variantLabel}}, options.Variants)
}

func TestWhenEditorOptionsAreBuilt_ReturnsTheServicesDistanceNames(t *testing.T) {
	t.Parallel()
	// Arrange
	content := models.SidMapping{Sid: gofakeit.Word()}
	expected := []string{gofakeit.Word(), gofakeit.Word()}
	service := &test_helpers.ContentRuleServiceMock{}
	service.On("GetVariantsForContent", content).Return([]models.VariantMapping{})
	service.On("GetDistanceDisplayNames").Return(expected)
	handler := handlers.NewContentRuleHandler(service)

	// Act
	options := handler.GetContentRuleEditorOptions(content)

	// Assert
	assert.Equal(t, expected, options.Distances)
}

func TestWhenEditorOptionsAreBuilt_DescribesEachFixedRule(t *testing.T) {
	t.Parallel()
	// Arrange
	content := models.SidMapping{Sid: gofakeit.Word()}
	service := &test_helpers.ContentRuleServiceMock{}
	service.On("GetVariantsForContent", mock.Anything).Return([]models.VariantMapping{})
	service.On("GetDistanceDisplayNames").Return([]string{})
	handler := handlers.NewContentRuleHandler(service)

	// Act
	options := handler.GetContentRuleEditorOptions(content)

	// Assert
	assert.Equal(t, dtos.ContentRuleOptionDto{
		Key:         dtos.ContentRuleKeyGuarded,
		Name:        content_rules.RuleGuardedName,
		Description: content_rules.RuleGuardedDescription,
		Marker:      content_rules.RuleGuardedMarker,
		EditorKind:  dtos.ContentRuleEditorKindBoolean,
		EditorLabel: "Guarded",
	}, options.Rules[2])
}

func ruleKeysOf(options dtos.ContentRuleEditorOptionsDto) []dtos.ContentRuleKey {
	keys := make([]dtos.ContentRuleKey, 0, len(options.Rules))
	for _, rule := range options.Rules {
		keys = append(keys, rule.Key)
	}
	return keys
}

func oneVariantMapping(id int, label string) []models.VariantMapping {
	return []models.VariantMapping{{
		Variants: []data.Tuple[int, string]{{Key: id, Value: label}},
	}}
}
