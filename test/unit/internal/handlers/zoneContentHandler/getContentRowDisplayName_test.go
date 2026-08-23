package zoneContentHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenTheRowNameIsRequested_TheContentNameAndDescriptionsAreHandedToTheService(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newZoneContentHandlerFixture()
	content := models.SidMapping{Name: gofakeit.Word()}
	rules := []models.ContentRuleRow{{Name: gofakeit.Word()}}
	description := dtos.ContentRuleDescriptionDto{Key: dtos.ContentRuleKeyVariant, Valid: true}
	fixture.contentRules.DescribeContentRuleFunc = func(
		models.SidMapping,
		models.ContentRuleRow,
	) dtos.ContentRuleDescriptionDto {
		return description
	}
	expected := gofakeit.Sentence(2)
	fixture.contentEditor.
		On("GetContentRowDisplayName", content.Name, []dtos.ContentRuleDescriptionDto{description}).
		Return(expected)

	// Act
	displayName := fixture.handler.GetContentRowDisplayName(content, rules)

	// Assert
	assert.Equal(t, expected, displayName)
}
