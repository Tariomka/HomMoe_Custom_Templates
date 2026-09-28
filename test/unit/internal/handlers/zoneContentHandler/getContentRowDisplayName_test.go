package zoneContentHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/content_rule_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenTheRowNameIsRequested_TheContentNameAndDescriptionsAreHandedToTheService(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newZoneContentHandlerFixture()
	content := models.SidMapping{Name: gofakeit.Word()}
	rules := []editor_state_model.ContentRuleRow{{Name: gofakeit.Word()}}
	description := content_rule_model.ContentRuleDescription{
		Key:          content_rule_model.ContentRuleKeyVariant,
		VariantLabel: gofakeit.Word(),
		Valid:        true,
	}
	fixture.contentRules.DescribeContentRuleFunc = func(
		models.SidMapping,
		editor_state_model.ContentRuleRow,
	) dtos.ContentRuleDescriptionDto {
		return dtos.ContentRuleDescriptionDto{ContentRuleDescription: description}
	}
	expected := gofakeit.Sentence(2)
	fixture.contentEditor.
		On("GetContentRowDisplayName", content.Name, []content_rule_model.ContentRuleDescription{description}).
		Return(expected)

	// Act
	displayName := fixture.handler.GetContentRowDisplayName(content, rules)

	// Assert
	assert.Equal(t, expected, displayName)
}
