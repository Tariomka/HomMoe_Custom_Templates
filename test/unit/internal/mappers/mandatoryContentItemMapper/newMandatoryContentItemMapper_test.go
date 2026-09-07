package mandatoryContentItemMapper_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/mappers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/content_rules"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWhenMapperIsCreated_ReturnsInstance(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	mapper := mappers.NewMandatoryContentItemMapper(content_rules.NewContentRuleService())

	// Assert
	assert.NotNil(t, mapper)
}

func TestWhenRowsAreMapped_UsesTheInjectedContentRuleService(t *testing.T) {
	t.Parallel()
	// Arrange
	contentRuleService := &test_helpers.ContentRuleServiceMock{}
	contentRuleService.On("RestoreRulesFromRow", mock.Anything, mock.Anything).Return(nil)
	contentRuleService.On("ApplyRulesToItem", mock.Anything, mock.Anything).Return()
	mapper := mappers.NewMandatoryContentItemMapper(contentRuleService)
	rows := []editor_state_model.ZoneContentRow{{Sid: gofakeit.LetterN(8), Count: 1}}

	// Act
	mapper.FromRows(rows)

	// Assert
	contentRuleService.AssertCalled(t, "RestoreRulesFromRow", mock.Anything, mock.Anything)
}
