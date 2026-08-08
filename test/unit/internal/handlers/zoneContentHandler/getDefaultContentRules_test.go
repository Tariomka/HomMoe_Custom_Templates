package zoneContentHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenDefaultRulesAreRequested_TheEditorOptionsAreHandedToTheService(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newZoneContentHandlerFixture()
	content := models.SidMapping{Name: gofakeit.Word()}
	options := dtos.ContentRuleEditorOptionsDto{
		Rules: []dtos.ContentRuleOptionDto{{Key: dtos.ContentRuleKeyGuarded, Name: gofakeit.Word()}},
	}
	expected := []models.ContentRuleRowSave{{Name: gofakeit.Word()}}
	fixture.contentRules.GetContentRuleEditorOptionsFunc = func(models.SidMapping) dtos.ContentRuleEditorOptionsDto {
		return options
	}
	fixture.contentEditor.On("GetDefaultContentRules", options).Return(expected)

	// Act
	rules := fixture.handler.GetDefaultContentRules(content)

	// Assert
	assert.Equal(t, expected, rules)
}
