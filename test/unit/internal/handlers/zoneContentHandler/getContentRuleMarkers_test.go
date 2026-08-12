package zoneContentHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenMarkersAreRequested_TheDescribedRulesAreHandedToTheService(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newZoneContentHandlerFixture()
	content := models.SidMapping{Name: gofakeit.Word()}
	rules := []models.ContentRuleRowSave{{Name: gofakeit.Word()}, {Name: gofakeit.Sentence(2)}}
	fixture.contentRules.DescribeContentRuleFunc = func(
		_ models.SidMapping,
		savedRule models.ContentRuleRowSave,
	) dtos.ContentRuleDescriptionDto {
		return dtos.ContentRuleDescriptionDto{Marker: savedRule.Name, Valid: true}
	}
	expected := gofakeit.Word()
	fixture.contentEditor.On("GetContentRuleMarkers", []dtos.ContentRuleDescriptionDto{
		{Marker: rules[0].Name, Valid: true},
		{Marker: rules[1].Name, Valid: true},
	}).Return(expected)

	// Act
	markers := fixture.handler.GetContentRuleMarkers(content, rules)

	// Assert
	assert.Equal(t, expected, markers)
}
