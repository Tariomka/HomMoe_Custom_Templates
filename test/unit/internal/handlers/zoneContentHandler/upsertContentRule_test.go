package zoneContentHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenAContentRuleIsUpserted_ReturnsTheServiceRuleList(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newZoneContentHandlerFixture()
	rules := []models.ContentRuleRowSave{{Name: gofakeit.Word()}}
	rule := models.ContentRuleRowSave{Name: gofakeit.Word()}
	expected := []models.ContentRuleRowSave{{Name: gofakeit.Sentence(2)}}
	fixture.contentEditor.On("UpsertContentRule", rules, rule).Return(expected)

	// Act
	merged := fixture.handler.UpsertContentRule(rules, rule)

	// Assert
	assert.Equal(t, expected, merged)
}
