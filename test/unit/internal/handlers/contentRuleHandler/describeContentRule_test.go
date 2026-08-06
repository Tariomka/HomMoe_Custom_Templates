package contentRuleHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/content_rules"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWhenSavedRuleCannotBeRestored_ReturnsAnInvalidDescription(t *testing.T) {
	t.Parallel()
	// Arrange
	savedRule := models.ContentRuleRowSave{Name: gofakeit.Word()}
	service := &test_helpers.ContentRuleServiceMock{}
	service.On("CreateRuleFromSavedRule", mock.Anything, mock.Anything).Return(nil)
	handler := handlers.NewContentRuleHandler(service)

	// Act
	description := handler.DescribeContentRule(models.SidMapping{}, savedRule)

	// Assert
	assert.Equal(t, dtos.ContentRuleDescriptionDto{
		DisplayText: savedRule.Name,
		SavedRule:   savedRule,
	}, description)
}

func TestWhenSavedRuleIsRestored_ReturnsTheRulesDisplayText(t *testing.T) {
	t.Parallel()
	// Arrange
	expected := gofakeit.Sentence(3)
	service := &test_helpers.ContentRuleServiceMock{}
	service.On("CreateRuleFromSavedRule", mock.Anything, mock.Anything).
		Return(ruleShowing(expected, gofakeit.Letter()))
	handler := handlers.NewContentRuleHandler(service)

	// Act
	description := handler.DescribeContentRule(models.SidMapping{}, models.ContentRuleRowSave{})

	// Assert
	assert.Equal(t, expected, description.DisplayText)
}

func TestWhenSavedRuleIsRestored_ReturnsTheRulesMarker(t *testing.T) {
	t.Parallel()
	// Arrange
	expected := gofakeit.Letter()
	service := &test_helpers.ContentRuleServiceMock{}
	service.On("CreateRuleFromSavedRule", mock.Anything, mock.Anything).
		Return(ruleShowing(gofakeit.Sentence(3), expected))
	handler := handlers.NewContentRuleHandler(service)

	// Act
	description := handler.DescribeContentRule(models.SidMapping{}, models.ContentRuleRowSave{})

	// Assert
	assert.Equal(t, expected, description.Marker)
}

func TestWhenSavedRuleIsRestored_MarksTheDescriptionValid(t *testing.T) {
	t.Parallel()
	// Arrange
	service := &test_helpers.ContentRuleServiceMock{}
	service.On("CreateRuleFromSavedRule", mock.Anything, mock.Anything).
		Return(ruleShowing(gofakeit.Sentence(3), gofakeit.Letter()))
	handler := handlers.NewContentRuleHandler(service)

	// Act
	description := handler.DescribeContentRule(models.SidMapping{}, models.ContentRuleRowSave{})

	// Assert
	assert.True(t, description.Valid)
}

func TestWhenSavedRuleNamesAKnownRule_ReturnsItsKey(t *testing.T) {
	t.Parallel()
	// Arrange
	service := &test_helpers.ContentRuleServiceMock{}
	service.On("CreateRuleFromSavedRule", mock.Anything, mock.Anything).Return(nil)
	handler := handlers.NewContentRuleHandler(service)

	// Act
	description := handler.DescribeContentRule(
		models.SidMapping{},
		models.ContentRuleRowSave{Name: content_rules.RuleGuardedName})

	// Assert
	assert.Equal(t, dtos.ContentRuleKeyGuarded, description.Key)
}

func TestWhenSavedRuleNamesAnUnknownRule_ReturnsAnEmptyKey(t *testing.T) {
	t.Parallel()
	// Arrange
	service := &test_helpers.ContentRuleServiceMock{}
	service.On("CreateRuleFromSavedRule", mock.Anything, mock.Anything).Return(nil)
	handler := handlers.NewContentRuleHandler(service)

	// Act
	description := handler.DescribeContentRule(
		models.SidMapping{},
		models.ContentRuleRowSave{Name: gofakeit.UUID()})

	// Assert
	assert.Equal(t, dtos.ContentRuleKey(""), description.Key)
}

func TestWhenSavedRuleSelectsAKnownVariant_ReturnsItsLabel(t *testing.T) {
	t.Parallel()
	// Arrange
	variantID := 7
	expectedLabel := gofakeit.Word()
	content := models.SidMapping{Sid: gofakeit.Word()}
	service := &test_helpers.ContentRuleServiceMock{}
	service.On("CreateRuleFromSavedRule", mock.Anything, mock.Anything).
		Return(ruleShowing(gofakeit.Sentence(3), gofakeit.Letter()))
	service.On("GetVariantForContentByID", content, variantID).
		Return(oneVariantMapping(variantID, expectedLabel)[0], true)
	handler := handlers.NewContentRuleHandler(service)

	// Act
	description := handler.DescribeContentRule(content, models.ContentRuleRowSave{VariantID: &variantID})

	// Assert
	assert.Equal(t, expectedLabel, description.VariantLabel)
}

func TestWhenSavedRuleSelectsAnUnknownVariant_ReturnsNoVariantLabel(t *testing.T) {
	t.Parallel()
	// Arrange
	variantID := gofakeit.IntRange(1, 100)
	content := models.SidMapping{Sid: gofakeit.Word()}
	service := &test_helpers.ContentRuleServiceMock{}
	service.On("CreateRuleFromSavedRule", mock.Anything, mock.Anything).
		Return(ruleShowing(gofakeit.Sentence(3), gofakeit.Letter()))
	service.On("GetVariantForContentByID", content, variantID).Return(models.VariantMapping{}, false)
	handler := handlers.NewContentRuleHandler(service)

	// Act
	description := handler.DescribeContentRule(content, models.ContentRuleRowSave{VariantID: &variantID})

	// Assert
	assert.Empty(t, description.VariantLabel)
}

// ruleShowing returns a content rule mock with the given display text and marker.
func ruleShowing(displayText, marker string) *test_helpers.ContentRuleMock {
	rule := &test_helpers.ContentRuleMock{}
	rule.On("DisplayText").Return(displayText)
	rule.On("Marker").Return(marker)
	return rule
}
