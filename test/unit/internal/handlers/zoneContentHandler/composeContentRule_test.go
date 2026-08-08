package zoneContentHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenAContentRuleIsComposed_ReturnsTheServiceResult(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newZoneContentHandlerFixture()
	request := dtos.ContentRuleCompositionRequestDto{
		Option: dtos.ContentRuleOptionDto{Key: dtos.ContentRuleKeyGuarded, Name: gofakeit.Word()},
	}
	expected := dtos.ContentRuleCompositionResultDto{
		Rule:  models.ContentRuleRowSave{Name: gofakeit.Word()},
		Valid: true,
	}
	fixture.contentEditor.On("ComposeContentRule", request).Return(expected)

	// Act
	result := fixture.handler.ComposeContentRule(request)

	// Assert
	assert.Equal(t, expected, result)
}
